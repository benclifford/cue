package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"
        "github.com/rs/xid"
)

var optVerbose = getopt.BoolLong("verbose", 'V', "", "output verbose progress information")

func main() {
	logInfo("cue: starting\n")
	var optPrivileged = getopt.BoolLong("privileged", 'P', "", "add assorted extra privileges")

	var optExtra = getopt.StringLong("docker-args", 'D', "", "add extra docker command-line arguments")

	getopt.Parse()

	logInfo("cue: parsed options\n")

	var environmentName string

	var cmdlineArgs = getopt.Args()

	if len(cmdlineArgs) < 1 {
		logInfo("cue: usage: cue <image name> [cmd]\n")
		os.Exit(75)
	}

	environmentName = cmdlineArgs[0]

	logInfo("cue: environment: %s\n", environmentName)

	var imageId string
	imageId = resolveNameToImage(environmentName)

	logInfo("cue: image ID: %s\n", imageId)

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

	var sharedTmpDir string = getHomeDir() + "/tmp/cue" // TODO
	logInfo("cue: shared temporary directory: %s\n", sharedTmpDir)

	logInfo("cue: creating temporary directory\n")
	err := os.MkdirAll(sharedTmpDir, 0755)
	exitOnError("when creating temporary directory", 76, err)

	userName := getUsername()
        id := userName + "-" + xid.New().String()

	var rootFilename string = sharedTmpDir + "/rootfile-" + id

	rootFile, err := os.Create(rootFilename)
	exitOnError("when opening rootfile", 67, err)

	defer rootFile.Close()

	_, err = rootFile.WriteString("#!/bin/bash\n")
	exitOnError("writing to rootFile", 68, err)

	if *optVerbose {
		_, err = rootFile.WriteString("echo cue: root: starting initialisation\n")
		exitOnError("writing to rootFile", 68, err)
	}

	err = rootFile.Chmod(0755)
	exitOnError("chmod'ing rootFile", 69, err)

	// TODO: factor with root file creation (eg. createSharedScript())
	var userFilename string = sharedTmpDir + "/userfile-" + id
	userFile, err := os.Create(userFilename)
	exitOnError("when opening userFile", 70, err)

	defer userFile.Close()

	_, err = userFile.WriteString("#!/bin/bash\n")
	exitOnError("writing to userFile", 73, err)

	if *optVerbose {
		_, err = userFile.WriteString("echo cue: user: starting initialisation\n")
		exitOnError("writing to userFile", 73, err)
	}

	err = userFile.Chmod(0755)
	exitOnError("chmod'ing userFile", 72, err)

	extraArgs := []string{}

	// Create user
	// TODO: read from current user information

	uid := getUid()

	if *optVerbose {
		_, err = rootFile.WriteString("echo cue: root: creating local user\n")
		exitOnError("writing to rootFile", 68, err)
	}

	// Create user and setup sudo - this varies by the nature of the underlying
	// tooling, which for now can be `debian` style (by default), or `redhat`
	// style, set by the CUE_USERMODE environment variable inside the container

	_, err = rootFile.WriteString(`# switch based on user creation mode
if [ "$CUE_USERMODE" == "debian" ] || [ "$CUE_USERMODE" == "" ] ; then
  useradd ` + userName + ` --uid=` + uid + ` --shell /bin/sh > /dev/null
  echo '%sudo   ALL=(ALL:ALL) NOPASSWD: ALL' > /etc/sudoers
  adduser root sudo > /dev/null
  adduser ` + userName + ` sudo > /dev/null
elif [ "$CUE_USERMODE" == "redhat" ] ; then
  useradd ` + userName + ` --uid=` + uid + ` --no-create-home --shell /bin/sh > /dev/null
  groupadd --system sudo
  gpasswd -a root sudo > /dev/null
  gpasswd -a ` + userName + ` sudo > /dev/null
  echo '%sudo   ALL=(ALL:ALL) NOPASSWD: ALL' > /etc/sudoers
else
  echo UNKNOWN USER MODE $CUE_USERMODE
  exit 1
fi
`)
	exitOnError("writing to rootFile", 68, err)

	// Run user shell (TODO: run user command)
	_, err = userFile.WriteString("#!/bin/bash\n")
	exitOnError("writing to userFile", 73, err)

	workdir, err := os.Getwd()
	exitOnError("getting workdir", 74, err)

	_, err = userFile.WriteString("cd " + workdir + "\n")
	exitOnError("writing to userFile", 73, err)

	// Run user level initialisation

	if *optVerbose {
		_, err = rootFile.WriteString("echo cue: root: running user level\n")
		exitOnError("writing to rootFile", 68, err)
	}

	_, err = rootFile.WriteString("sudo -u " + userName + " -i " + userFilename + "\n")
	exitOnError("writing to rootFile", 68, err)

	// TODO: if $DISPLAY is set to :0, mount
	// /tmp/.X11-unix/X0 into the container.
	// This could be generalised to arbitrary $DISPLAY values
	// with more effort.

	display := os.Getenv("DISPLAY")
	if display == ":0" {
		logInfo("cue: mounting X server\n")
		extraArgs = append(extraArgs, "-v", "/tmp/.X11-unix/X0:/tmp/.X11-unix/X0")
		_, err = userFile.WriteString("export DISPLAY=:0\n")
		exitOnError("writing to userFile", 73, err)
	}

	// Handle docker extra args
	ex := *optExtra
	if ex != "" {
		axs := strings.Split(*optExtra, " ")
		logInfo("cue: docker extra args: %d >%s<\n", len(axs), ex)
		extraArgs = append(extraArgs, axs...)
	}

	// After everything else is done, run a shell
	// or eventually a passed-in command.

	// TODO: the choice of shell is interesting here.
	// Should it be the user's default shell, which is not necessarily
	// installed?

	if len(cmdlineArgs) == 1 {
		_, err = userFile.WriteString("/bin/bash\n")
		exitOnError("writing user shell to userFile", 73, err)
	} else {
		// TODO: there will be some string escaping bug here
		// one day, but string escaping in shell frustrates me
		// too much to deal with at time of writing.
		for _, element := range cmdlineArgs[1:] {
			_, err = userFile.WriteString(element)
			exitOnError("writing user command to userFile", 73, err)
			_, err = userFile.WriteString(" ")
			exitOnError("writing user command to userFile", 73, err)
		}
		exitOnError("writing user command to userFile", 73, err)
		_, err = userFile.WriteString("\n")
		exitOnError("writing user command newline to userFile", 73, err)
	}

	err = rootFile.Close()
	exitOnError("closing rootFile", 69, err)

	err = userFile.Close()
	exitOnError("closing userFile", 71, err)

	exitStatus := runImage(imageId, rootFilename, extraArgs, *optPrivileged)

	logInfo("cue: done\n")
	os.Exit(exitStatus)
}

func runImage(imageId string, rootFile string, dockerArgs []string, privileged bool) int {
	attributes := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	// TODO: get docker from the path

	homeDir := getHomeDir()

	argsPre := []string{"docker", "run", "--rm", "-ti", "-v", homeDir + ":" + homeDir}

	argsPrivileged := []string{}

	if privileged {
		// volume mount /dev/bus/usb: if not, a /dev/bus/usb is created that appears to be a replicate of the container start time /dev/bus/usb, not the "live" version with new devices.
		// forward port 8080 for MSE development
		fmt.Print("cue: runImage: adding extra privileges\n")
		argsPrivileged = []string{"-v", "/dev/bus/usb:/dev/bus/usb", "--privileged", "-p", "8080"}

	}

	argsPre = append(argsPre, argsPrivileged...)

	argsPost := []string{imageId, rootFile}
	args := append(argsPre, append(dockerArgs, argsPost...)...)

	process, err := os.StartProcess("/usr/bin/docker", args, &attributes)
	exitOnError("running Docker container", 65, err)

	status, err := process.Wait()
	exitOnError("waiting for Docker container process", 66, err)

	logInfo("cue: runImage: docker container process finished with status: %s\n", status)

	ws := status.Sys().(syscall.WaitStatus)

	return ws.ExitStatus()
}

// given a name, resolves it to a docker image ID

// If a directory <environment>/Dockerfile exists, it
// will be build and used.
// (in this case, docker build will be run every time,
// but usually will // be fast (as the output will be cached),
// and when it isn't fast, it is because the environment needed
// rebuilding anyway.

// Otherwise, the environment name will be passed directly
// to docker to be resolved - for example as a container ID
// or docker image tag.

func resolveNameToImage(environment string) string {

	dockerfileLibrary := getHomeDir() + "/src/cue/dockerfiles/"
	environmentPath := dockerfileLibrary + "/" + environment

	if stat, err := os.Stat(environmentPath); err == nil && stat.IsDir() {
		logInfo("cue: resolveNameToImage: environment directory exists - using docker build\n")
		cmd := "docker"

		username := getUsername()

		tagname := "cue/" + username + "/" + environment
		args := []string{"build", "--quiet", "--tag", tagname, environmentPath}
		output, err := exec.Command(cmd, args...).CombinedOutput()
		exitOnError("running Docker build", 64, err)

		logInfo("cue: resolveNameToImage: successful output from docker build:\n%s\n", output)
		return strings.TrimSpace(string(output))
	} else {
		logInfo("cue: resolveNameToImage: environment directory does not exist - using name as raw docker image identifier\n")
		return strings.TrimSpace(environment)
	}
}

func exitOnError(message string, exitCode int, err error) {
	if err != nil {
		logError("cue: ERROR: %s: %s\n", message, err)
		os.Exit(exitCode)
	}
}

func getHomeDir() string {
	usr, err := user.Current()
	exitOnError("Getting current user info", 77, err)
	return usr.HomeDir
}

func getUsername() string {
	usr, err := user.Current()
	exitOnError("Getting current user info", 77, err)
	return usr.Username
}

func getUid() string {
	usr, err := user.Current()
	exitOnError("Getting current user info", 77, err)
	return usr.Uid
}

func logInfo(format string, a ...interface{}) (n int, err error) {
	if *optVerbose {
		return fmt.Printf(format, a...)
	} else {
		return 0, nil
	}
}

func logError(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}
