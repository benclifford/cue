package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"github.com/rs/xid"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var optVerbose = getopt.BoolLong("verbose", 'V', "", "output verbose progress information")

func main() {
	var optExtra = getopt.StringLong("docker-args", 'D', "", "add extra docker command-line arguments")

	getopt.Parse()

	logInfo("parsed options\n")

	var environmentName string

	var cmdlineArgs = getopt.Args()

	if len(cmdlineArgs) < 1 {
		logInfo("usage: cue <image name> [cmd]\n")
		os.Exit(75)
	}

	environmentName = cmdlineArgs[0]

	logInfo("environment: %s\n", environmentName)

	var imageId string
	imageId = resolveNameToImage(environmentName)

	logInfo("image ID: %s\n", imageId)

	// this directory needs to be somewhere that will be
	// mounted inside the container, so that it will be
	// accessible there. (BUG/FEATUREREQ: these files can
	// be prepared in a local tmp and copied into the container
	// which might also help in making the container
	// restartable, if such a mode is desired - allowing the
	// shared temp directory to be cleared out but the container
	// to still start up properly)

	var sharedTmpDir string = getHomeDir() + "/tmp/cue" // TODO
	logInfo("shared temporary directory: %s\n", sharedTmpDir)

	logInfo("creating temporary directory\n")
	err := os.MkdirAll(sharedTmpDir, 0755)
	exitOnError("when creating temporary directory", 76, err)

	userName := getUsername()
	id := userName + "-" + xid.New().String()

	rootFilename, rootFile := createSharedScript(sharedTmpDir, "rootfile-"+id)
	defer rootFile.Close()

	_, err = rootFile.WriteString("#!/bin/bash\n")
	exitOnError("writing to rootFile", 68, err)

	if *optVerbose {
		_, err = rootFile.WriteString("echo cue: root: starting initialisation\n")
		exitOnError("writing to rootFile", 68, err)
	}

	userFilename, userFile := createSharedScript(sharedTmpDir, "userfile-"+id)
	defer userFile.Close()

	_, err = userFile.WriteString("#!/bin/bash\n")
	exitOnError("writing to userFile", 73, err)

	if *optVerbose {
		_, err = userFile.WriteString("echo cue: user: starting initialisation\n")
		exitOnError("writing to userFile", 73, err)
	}

	extraArgs := []string{}

	// Create user
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

	sanitisedEnvironmentName := strings.Replace(environmentName, ".", "-", -1)
	// set container name and hostname
	uniquifier := getUniquifier(sharedTmpDir)
	hostname := sanitisedEnvironmentName + "-" + uniquifier
	containerName := "cue." + userName + "." + sanitisedEnvironmentName + "." + uniquifier
	extraArgs = append(extraArgs, "--name", containerName, "--hostname", hostname)

	// If $DISPLAY is set to :0, mount
	// /tmp/.X11-unix/X0 into the container.
	// This could be generalised to arbitrary $DISPLAY values
	// with more effort.

	display := os.Getenv("DISPLAY")
	if display == ":0" {
		logInfo("mounting X server\n")
		extraArgs = append(extraArgs, "-v", "/tmp/.X11-unix/X0:/tmp/.X11-unix/X0")
		_, err = userFile.WriteString("export DISPLAY=:0\n")
		exitOnError("writing to userFile", 73, err)
	}

	// Mount the SSH agent in the container if it exists
	sshAgent, sshAgentPresent := os.LookupEnv("SSH_AUTH_SOCK")
	if sshAgentPresent {
		logInfo("mounting SSH agent socket\n")
		extraArgs = append(extraArgs, "-v", sshAgent+":"+sshAgent)
		_, err = userFile.WriteString("export SSH_AUTH_SOCK=" + sshAgent + "\n")
		exitOnError("writing to userFile", 73, err)
	}

	// Handle docker extra args
	ex := *optExtra
	if ex != "" {
		axs := strings.Split(*optExtra, " ")
		logInfo("docker extra args: %d >%s<\n", len(axs), ex)
		extraArgs = append(extraArgs, axs...)
	}

	// After everything else is done, run a shell
	// or eventually a passed-in command.

	// TODO: the choice of shell is interesting here.
	// Should it be the user's default shell, which is not necessarily
	// installed?

	cmdFilename, cmdFile := createSharedScript(sharedTmpDir, "cmdfile-"+id)

	if len(cmdlineArgs) == 1 {
		_, err = cmdFile.WriteString("/bin/bash\n")
		exitOnError("writing user shell to cmdFile", 73, err)
	} else {
		// TODO: there will be some string escaping bug here
		// one day, but string escaping in shell frustrates me
		// too much to deal with at time of writing.
		for _, element := range cmdlineArgs[1:] {
			_, err = cmdFile.WriteString(element)
			exitOnError("writing user command to cmdFile", 73, err)
			_, err = cmdFile.WriteString(" ")
			exitOnError("writing user command to cmdFile", 73, err)
		}

		exitOnError("writing user command to userFile", 73, err)
		_, err = userFile.WriteString("\n")
		exitOnError("writing user command newline to userFile", 73, err)
	}

	// this tests for file existence (-f) rather than executability (-x)
	// because if there is a non-executable /cue.shell, I want to get
	// an execution error rather than a silent ignore.
	_, err = userFile.WriteString("if [ -f /cue.shell ] ; then /cue.shell " + cmdFilename + " ; else " + cmdFilename + " ; fi\n")
	exitOnError("writing cmdFile invocation to userFile", 73, err)

	err = rootFile.Close()
	exitOnError("closing rootFile", 69, err)

	err = userFile.Close()
	exitOnError("closing userFile", 71, err)

	err = cmdFile.Close()
	exitOnError("closing cmdFile", 78, err)

	exitStatus := runImage(imageId, rootFilename, extraArgs)

	logInfo("done\n")
	os.Exit(exitStatus)
}

func runImage(imageId string, rootFile string, dockerArgs []string) int {
	attributes := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	// TODO: get docker from the path

	homeDir := getHomeDir()

	argsPre := []string{"docker", "run", "--rm", "-ti", "-v", homeDir + ":" + homeDir}

	argsPost := []string{imageId, rootFile}
	args := append(argsPre, append(dockerArgs, argsPost...)...)

	process, err := os.StartProcess("/usr/bin/docker", args, &attributes)
	exitOnError("running Docker container", 65, err)

	status, err := process.Wait()
	exitOnError("waiting for Docker container process", 66, err)

	logInfo("runImage: docker container process finished with status: %s\n", status)

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
		logInfo("resolveNameToImage: environment directory exists - using docker build\n")
		cmd := "docker"

		username := getUsername()

		tagname := "cue/" + username + "/" + environment
		args := []string{"build", "--quiet", "--tag", tagname, environmentPath}
		output, err := exec.Command(cmd, args...).CombinedOutput()
		exitOnError("running Docker build", 64, err)

		logInfo("resolveNameToImage: successful output from docker build:\n%s\n", output)
		return strings.TrimSpace(string(output))
	} else {
		logInfo("resolveNameToImage: environment directory does not exist - using name as raw docker image identifier\n")
		return strings.TrimSpace(environment)
	}
}

func exitOnError(message string, exitCode int, err error) {
	if err != nil {
		logError("%s: %s\n", message, err)
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

func createSharedScript(sharedTmpDir string, filenameBase string) (string, *os.File) {
	var filename string = sharedTmpDir + "/" + filenameBase
	file, err := os.Create(filename)
	exitOnError("when opening file "+filename, 67, err)
	err = file.Chmod(0755)
	exitOnError("chmod'ing file", 69, err)
	return filename, file
}

// I'd like this to return a simple small id that is
// unique wrt other cue containers with the same user name
// and environment name. There is no requirement for
// sequencing. There is no requirement for the return value
// to be a number, although it is in this implementation.
func getUniquifier(tmpdir string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	return strconv.Itoa(rand.Intn(10000))
}

func logInfo(format string, a ...interface{}) (n int, err error) {
	if *optVerbose {
		return fmt.Printf("cue: "+format, a...)
	} else {
		return 0, nil
	}
}

func logError(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf("cue: ERROR: "+format, a...)
}
