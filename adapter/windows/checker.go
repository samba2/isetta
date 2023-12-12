package windows

import (
	_ "embed"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"org.samba/isetta/helper"
	log "org.samba/isetta/simplelogger"
)

type WindowsCheckerImpl struct {
	PxProxyPort int
}

func (WindowsCheckerImpl) IsPingable(host string) bool {
	log.Logger.Trace("Checking if reachable inside Windows")
	cmd := fmt.Sprintf("ping.exe -n 2 -w 100 %v > $null; $LASTEXITCODE", host)
	exitCode := runInPowerShell(cmd)
	return exitCode == "0"
}

func (w *WindowsCheckerImpl) IsPxProxyRunning() bool {
	pxProxyPort := fmt.Sprint(w.PxProxyPort)
	return isPortOpenOnWindows(pxProxyPort)
}

func (WindowsCheckerImpl) IsRunningOnWsl2() bool {
	log.Logger.Trace("Checking if running inside WSL 2")
	resultUtf16 := runInPowerShell("wsl.exe --list --verbose")

	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	result, err := decoder.String(string(resultUtf16))
	helper.AssertNoError(err, "Error decoding cmd output")

	return parseListOutput(result)
}

// find the WSL2 version in the output
// (?m) is for multiline match
// example:
// * Ubuntu    Running         2 
func parseListOutput(output string) bool {
	versionRegexLine := "(?m)^* .+ 2\\s*$"
	match, err := regexp.MatchString(versionRegexLine, output)
	helper.AssertNoError(err, "Error executing regex")
	return match
}

func isPortOpenOnWindows(port string) bool {
	command := fmt.Sprintf("Test-NetConnection -ComputerName 127.0.0.1 -Port %v -InformationLevel Quiet", port)
	result := runInPowerShell(command)
	if result == "True" {
		return true
	} else {
		return false
	}
}

func runInPowerShell(command string) string {
	log.Logger.Trace("Running in Powershell: %v", command)
	result, err := exec.Command("powershell.exe", "-NoProfile", "-Command", command).CombinedOutput()
	helper.AssertNoError(err, "Error executing Powershell command: %v", command)

	return strings.TrimSuffix(string(result), "\r\n")
}
