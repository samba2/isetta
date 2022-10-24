package gsudo

import (
	_ "embed"
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"org.samba/isetta/helper"
	log "org.samba/isetta/simplelogger"
)

//go:embed resources/gsudo.exe
var gsudoBinary []byte

type Gsudo struct {
	windowsTempDirPath    string
	windowsTempDirWslPath string
	gsudoWslPath          string
	gsudoWindowsPath      string
}

func (gsudo *Gsudo) Init() {
	gsudo.setupPaths()
	gsudo.copyGsudoBinary()
	gsudo.preflightCheck()
}

func (gsudo *Gsudo) setupPaths() {
	gsudo.windowsTempDirPath = getWindowsTempDir()
	gsudo.windowsTempDirWslPath = windowsPathToWsl(gsudo.windowsTempDirPath)
	gsudo.gsudoWslPath = path.Join(gsudo.windowsTempDirWslPath, "gsudo-isetta.exe")
	log.Logger.Trace("gsudo WSL path: %v", gsudo.gsudoWslPath)
	gsudo.gsudoWindowsPath = gsudo.windowsTempDirPath + "\\gsudo-isetta.exe" // concat since no Windows join available
	log.Logger.Trace("gsudo Windows path: %v", gsudo.gsudoWindowsPath)
}

func (gsudo *Gsudo) copyGsudoBinary() {
	log.Logger.Debug("Making gsudo available at %v", gsudo.gsudoWslPath)
	configFunc := func() bool {
		err := os.WriteFile(gsudo.gsudoWslPath, gsudoBinary, 0775)
		return err == nil
	}

	err := helper.Retry(helper.RetryParams{
		Description: "Writing gsudo binary",
		Attempts:    10,
		Sleep:       1 * time.Second,
		Func:        configFunc,
	})
	helper.AssertNoError2(err)
}

func (gsudo *Gsudo) preflightCheck() {
	fullCommand := []string{gsudo.gsudoWslPath, "--help"}
	fullCommandStr := strings.Join(fullCommand, " ")
	log.Logger.Trace("Preflight check. Executing command '%v'", fullCommandStr)
	out, err := exec.Command(fullCommand[0], fullCommand[1:]...).CombinedOutput()
	helper.AssertNoError(err, "Preflight check failed. Failed to run: %v, output was: %v", fullCommandStr, string(out))
	log.Logger.Trace("Preflight check was successful")
}

// should be called via 'defer' to cleanup the binary
func (gsudo *Gsudo) Cleanup() {
	log.Logger.Trace("Resetting cache")
	gsudo.run("--reset-timestamp")
	log.Logger.Trace("Removing gsudo binary from %v", gsudo.gsudoWslPath)
	os.Remove(gsudo.gsudoWslPath)
}

// executed elevated Windows command. By default it is expected that the command executes with
// exit code 0 (success). Error checking can also be deactivating by passing 'false' as second argument
func (gsudo *Gsudo) RunElevated(command string, checkError ...bool) string {
	gsudo.tryActivateCache()
	return gsudo.run(command, checkError...)
}

func (gsudo *Gsudo) tryActivateCache() {
	log.Logger.Trace("Trying activate gsudo cache")
	statusOutput := gsudo.run("status")

	if gsudo.isCacheActive(statusOutput) {
		log.Logger.Trace("Credential cache is active. Won't start a new session.")
	} else {
		log.Logger.Debug("Credential cache not active, starting it (will prompt for admin credentials)")
		gsudo.run("cache on --pid 0 --duration 00:00:30")
		gsudo.waitForCacheActive()
	}
}

func (gsudo *Gsudo) waitForCacheActive()  {
	checkFunc := func() bool {
		statusOutput := gsudo.run("status")
		return gsudo.isCacheActive(statusOutput)
	}

	err := helper.Retry(helper.RetryParams{
		Description: "Credential cache started",
		Attempts:    10,
		Sleep:       250 * time.Millisecond,
		Func:        checkFunc,
	})
	helper.AssertNoError2(err)	
}

// Execute gsudo
// after some experiments, this is how it worked reliable:
// - binary must reside somewhere under c:\
// - executed via cmd.exe
// - gsudo path + passed in command need to be single string (otherwise multi args don't work)
// - the command needs to run inside an existing Windows dir (prevent cmd.exe warnings)
//
// 'checkError' is an optional bool which controls if execution errors should
// stop program flow
func (gsudo *Gsudo) run(command string, checkError ...bool) string {
	checkError2, err := isCheckError(checkError)
	helper.AssertNoError2(err)
	cmdCommand := gsudo.gsudoWindowsPath + " " + command
	fullCommand := []string{"cmd.exe", "/c", cmdCommand}
	log.Logger.Trace("Executing command '%v'", strings.Join(fullCommand, " "))

	cmd := exec.Command(fullCommand[0], fullCommand[1:]...)
	cmd.Dir = gsudo.windowsTempDirWslPath // prevent warnings, see comment above
	outBytes, err := cmd.CombinedOutput()
	out := string(outBytes)

	if checkError2 && err != nil {
		log.Logger.Error("Error running: %v, error was: %v, output was: %v", fullCommand, err, out)
	}

	log.Logger.Trace("Output was: %v", out)
	return out
}

// determine default argument
func isCheckError(checkError []bool) (bool, error) {
	myCheckError := true

	if len(checkError) == 1 {
		myCheckError = checkError[0]
	} else if len(checkError) > 1 {
		return false, errors.New("'checkError' is only allowed to have one argument")
	}
	return myCheckError, nil
}

func (gsudo *Gsudo) isCacheActive(statusOutput string) bool {
	searchString := "Available for this process: True"
	log.Logger.Trace("Looking for search string '%v' in output", searchString)
	return strings.Contains(statusOutput, searchString)
}

func getWindowsTempDir() string {
	cmd := exec.Command("cmd.exe", "/c", "echo %TEMP%")
	cmd.Dir = "/mnt/c/"
	out, err := cmd.CombinedOutput()
	helper.AssertNoError(err, "failed to determine Windows temp dir. Output was: %v", string(out))

	windowsTempDir := strings.TrimRight(string(out), "\r\n")
	return windowsTempDir
}

func windowsPathToWsl(p string) string {
	out, err := exec.Command("wslpath", "-u", p).Output()
	helper.AssertNoError2(err)
	return strings.Trim(string(out), "\n\r")
}
