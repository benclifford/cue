package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Printf("cue: starting\n")

	var environmentName string

	if len(os.Args) < 2 {
		fmt.Printf("cue: usage: cue <image name>\n")
		os.Exit(75)
	}

	environmentName = os.Args[1]

	fmt.Printf("cue: environment: %s\n", environmentName)

	var imageId string
	imageId = resolveNameToImage(environmentName)

	fmt.Printf("cue: image ID: %s\n", imageId)

	// TODO: prep docker commandline, root and user prep files in
	//       thematic stages

	// this directory needs to be somewhere that will be
	// mounted inside the container, so that it will be
	// accessible there. (BUG/FEATUREREQ: these files can
	// be prepared in a local tmp and copied into the container
	// which might also help in making the container
	// restartable, if such a mode is desired - allowing the
	// shared temp directory to be cleared out but the container
	// to still start up properly)
	var sharedTmpDir string = "/home/benc/tmp/cue" // TODO
	fmt.Printf("cue: shared temporary directory: %s\n", sharedTmpDir)

	var rootFilename string = sharedTmpDir + "/rootfile" // TODO

	rootFile, err := os.Create(rootFilename)
	exitOnError("when opening rootfile", 67, err)

	defer rootFile.Close()

	_, err = rootFile.WriteString("#!/bin/bash\necho cue: root: starting initialisation\n")
	exitOnError("writing to rootFile", 68, err)

	err = rootFile.Chmod(0755)
	exitOnError("chmod'ing rootFile", 69, err)

	// TODO: factor with root file creation (eg. createSharedScript())
	var userFilename string = sharedTmpDir + "/userfile" // TODO
	userFile, err := os.Create(userFilename)
	exitOnError("when opening userFile", 70, err)

	defer userFile.Close()

	_, err = userFile.WriteString("#!/bin/bash\necho cue: user: starting initialisation\n")
	exitOnError("writing to userFile", 73, err)

	err = userFile.Chmod(0755)
	exitOnError("chmod'ing userFile", 72, err)

	extraArgs := []string{}

	// Create user
	// TODO: read from current user information

	userName := "benc"
	uid := "1000"

	_, err = rootFile.WriteString("echo cue: root: creating local user\n")
	exitOnError("writing to rootFile", 68, err)

	_, err = rootFile.WriteString("useradd " + userName + " --uid=" + uid + "\n")
	exitOnError("writing to rootFile", 68, err)

	// Diddle sudo

	_, err = rootFile.WriteString("echo '%sudo   ALL=(ALL:ALL) NOPASSWD: ALL' > /etc/sudoers\n")
	exitOnError("writing to rootFile", 68, err)

	_, err = rootFile.WriteString("adduser root sudo\n")
	exitOnError("writing to rootFile", 68, err)

	_, err = rootFile.WriteString("adduser " + userName + " sudo\n")
	exitOnError("writing to rootFile", 68, err)

	// Run user shell (TODO: run user command)
	_, err = userFile.WriteString("#!/bin/bash\n")
	exitOnError("writing to userFile", 73, err)

	workdir, err := os.Getwd()
	exitOnError("getting workdir", 74, err)

	_, err = userFile.WriteString("cd " + workdir + "\n")
	exitOnError("writing to userFile", 73, err)

	// Run user level initialisation
	_, err = rootFile.WriteString("echo cue: root: running user level\n")
	exitOnError("writing to rootFile", 68, err)

	_, err = rootFile.WriteString("sudo -u " + userName + " -i " + userFilename + "\n")
	exitOnError("writing to rootFile", 68, err)

	// TODO: if $DISPLAY is set to :0, mount
	// /tmp/.X11-unix/X0 into the container.
	// This could be generalised to arbitrary $DISPLAY values
	// with more effort.

	display := os.Getenv("DISPLAY")
	if display == ":0" {
		fmt.Printf("cue: mounting X server\n")
		extraArgs = append(extraArgs, "-v", "/tmp/.X11-unix/X0:/tmp/.X11-unix/X0")
		_, err = userFile.WriteString("export DISPLAY=:0\n")
		exitOnError("writing to userFile", 73, err)
	}

	// After everything else is done, run a shell
	// or TODO eventually a passed-in command
	_, err = userFile.WriteString("/bin/bash\n")
	exitOnError("writing to userFile", 73, err)

	err = rootFile.Close()
	exitOnError("closing rootFile", 69, err)

	err = userFile.Close()
	exitOnError("closing userFile", 71, err)

	runImage(imageId, rootFilename, extraArgs)

	fmt.Printf("cue: done\n")
	// TODO: return container exit code
}

func runImage(imageId string, rootFile string, dockerArgs []string) {
	attributes := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	// TODO: get docker from the path

	argsPre := []string{"docker", "run", "--rm", "-ti", "-v", "/home/benc:/home/benc"}
	argsPost := []string{imageId, rootFile}
	args := append(argsPre, append(dockerArgs, argsPost...)...)

	process, err := os.StartProcess("/usr/bin/docker", args, &attributes)
	exitOnError("running Docker container", 65, err)

	status, err := process.Wait()
	exitOnError("waiting for Docker container process", 66, err)

	fmt.Printf("cue: runImage: docker container process finished with status: %s\n", status)

}

// given a name, resolves it to a docker image ID
// perhaps one day doing something fancy but for now
// using a directory containing <environment>/Dockerfile
// so basically wraps a call to docker build.
// docker build will be run every time, but usually will
// be fast (as the output will be cached), and when it
// isn't fast, it is because the environment needed
// rebuilding anyway.
func resolveNameToImage(environment string) string {
	cmd := "docker"
	args := []string{"build", "--quiet", "/home/benc/src/cue/dockerfiles/" + environment}
	output, err := exec.Command(cmd, args...).CombinedOutput()
	exitOnError("running Docker build", 64, err)

	fmt.Printf("cue: resolveNameToImage: successful output from docker build:\n%s\n", output)

	return strings.TrimSpace(string(output))
}

func exitOnError(message string, exitCode int, err error) {
	if err != nil {
		fmt.Printf("cue: ERROR: %s: %s\n", message, err)
		os.Exit(exitCode)
	}
}
