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

	if err != nil {
		fmt.Printf("cue: ERROR: when opening rootFile %s: %s", rootFile, err)
		os.Exit(67)
	}

	defer rootFile.Close()

	_, err = rootFile.WriteString("#!/bin/bash\necho cue: root: starting initialisation\n")
	if err != nil {
		fmt.Printf("cue: ERROR: writing to rootFile: %s", err)
		os.Exit(68)
	}

	err = rootFile.Chmod(0755)

	if err != nil {
		fmt.Printf("cue: ERROR: chmod rootFile: %s", err)
		os.Exit(69)
	}

	// TODO: factor with root file creation (eg. createSharedScript())
	var userFilename string = sharedTmpDir + "/userfile" // TODO
	userFile, err := os.Create(userFilename)
	if err != nil {
		fmt.Printf("cue: ERROR: when opening userFile %s: %s", userFile, err)
		os.Exit(70)
	}
	defer userFile.Close()

	_, err = userFile.WriteString("#!/bin/bash\necho cue: user: starting initialisation\n")
	if err != nil {
		fmt.Printf("cue: ERROR: writing to userFile: %s", err)
		os.Exit(68)
	}

	err = userFile.Chmod(0755)
	if err != nil {
		fmt.Printf("cue: ERROR: chmod rootFile: %s", err)
		os.Exit(72)
	}

	// Create user
	// TODO: read from current user information

	userName := "benc"
	uid := "1000"

	_, err = rootFile.WriteString("echo cue: root: creating local user\n")
	if err != nil {
		fmt.Printf("cue: ERROR: writing to rootFile: %s", err)
		os.Exit(68)
	}

	_, err = rootFile.WriteString("useradd " + userName + " --uid=" + uid + "\n")
	if err != nil {
		fmt.Printf("cue: ERROR: writing to rootFile: %s", err)
		os.Exit(68)
	}

	// Diddle sudo

	_, err = rootFile.WriteString("echo '%sudo   ALL=(ALL:ALL) NOPASSWD: ALL' > /etc/sudoers\n")
	if err != nil {
		fmt.Printf("cue: ERROR: writing to rootFile: %s", err)
		os.Exit(68)
	}

	_, err = rootFile.WriteString("adduser root sudo\n")
	if err != nil {
		fmt.Printf("cue: ERROR: writing to rootFile: %s", err)
		os.Exit(68)
	}
	_, err = rootFile.WriteString("adduser " + userName + " sudo\n")
	if err != nil {
		fmt.Printf("cue: ERROR: writing to rootFile: %s", err)
		os.Exit(68)
	}

	// Run user shell (TODO: run user command)
	_, err = userFile.WriteString("/bin/bash\n")
	if err != nil {
		fmt.Printf("cue: ERROR: writing to userFile: %s", err)
		os.Exit(68)
	}

	// Run user level initialisation
	_, err = rootFile.WriteString("echo cue: root: running user level\n")
	if err != nil {
		fmt.Printf("cue: ERROR: writing to rootFile: %s", err)
		os.Exit(68)
	}
	_, err = rootFile.WriteString("sudo -u " + userName + " -i " + userFilename + "\n")
	if err != nil {
		fmt.Printf("cue: ERROR: writing to rootFile: %s", err)
		os.Exit(68)
	}

	err = rootFile.Close()
	if err != nil {
		fmt.Printf("cue: ERROR: closing rootFile: %s", err)
		os.Exit(69)
	}

	err = userFile.Close()
	if err != nil {
		fmt.Printf("cue: ERROR: closing userFile: %s", err)
		os.Exit(71)
	}

	runImage(imageId, rootFilename)

	fmt.Printf("cue: done\n")
	// TODO: return container exit code
}

func runImage(imageId string, rootFile string) {
	attributes := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	// TODO: get docker from the path
	process, err := os.StartProcess("/usr/bin/docker", []string{"docker", "run", "--rm", "-ti", "-v", "/home/benc:/home/benc", imageId, rootFile}, &attributes)

	if err != nil {
		fmt.Printf("cue: runImage: ERROR: running docker container: %s\n", err)
		os.Exit(65)
	}

	status, err := process.Wait()
	if err != nil {
		fmt.Printf("cue: runImage: ERROR: docker container process wait returned: %s\nstatus: %s", err, status)
		os.Exit(66)
	}

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
	if err != nil {
		fmt.Printf("cue: resolveNameToImage: ERROR: error resolving name to image, running docker build: %s\n%s", err, output)
		os.Exit(64)
	}

	fmt.Printf("cue: resolveNameToImage: successful output from docker build:\n%s\n", output)

	return strings.TrimSpace(string(output))
}
