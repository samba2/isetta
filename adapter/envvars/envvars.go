package envvars

import (
	"fmt"
	"os"
	"strings"
	log "org.samba/isetta/simplelogger"
)

const exportEnvironmentVariables = `
export HTTPS_PROXY=http://%[1]v:%[2]v
export HTTP_PROXY=http://%[1]v:%[2]v
export https_proxy=http://%[1]v:%[2]v
export http_proxy=http://%[1]v:%[2]v
export NO_PROXY=localhost,127.0.0.1,%[1]v
`

const unsetEnvironmentVariables = `
unset HTTPS_PROXY
unset HTTP_PROXY
unset https_proxy
unset http_proxy
unset NO_PROXY
`

type ConsoleEnvVarPrinter struct{
	WindowsIp string
	PxProxyPort int
}

func (c *ConsoleEnvVarPrinter) PrintExportCommands() {
	fmt.Println(c.buildPrintExportCommands())
}

func (c *ConsoleEnvVarPrinter) buildPrintExportCommands() string {
	trimmed := strings.Trim(exportEnvironmentVariables, "\n") 
	return fmt.Sprintf(trimmed, c.WindowsIp, c.PxProxyPort)
}

func (c *ConsoleEnvVarPrinter) PrintUnsetCommands() {
	fmt.Print(strings.Trim(unsetEnvironmentVariables, "\n"))
}

func (c *ConsoleEnvVarPrinter) WarnIfProxyVarSet() {
	if c.areHttpEnvVarsSet() {
		log.Logger.Warn("This shell still has one ore more http(s)_proxy environment variables set. You are directly connected, don't forget to unset them.")
	}
}

func (c *ConsoleEnvVarPrinter) areHttpEnvVarsSet() bool {
	proxyEnvVars := []string{
		"http_proxy",
		"https_proxy",
		"HTTP_PROXY",
		"HTTPS_PROXY",
	}

	for _, envVar := range proxyEnvVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}
	return false
}