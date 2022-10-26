package envvars

import (
	"fmt"
	"os"
	"strings"
	log "org.samba/isetta/simplelogger"
)

const exportHttpProxyVariables = `
export HTTPS_PROXY=http://%[1]v:%[2]v
export HTTP_PROXY=http://%[1]v:%[2]v
export https_proxy=http://%[1]v:%[2]v
export http_proxy=http://%[1]v:%[2]v
`

const exportNoProxyVariable = `
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
	trimmedHttpVars := strings.Trim(exportHttpProxyVariables, "\n")
	httpVars := fmt.Sprintf(trimmedHttpVars, c.WindowsIp, c.PxProxyPort)
	
	trimmedNoProxyVar := strings.Trim(exportNoProxyVariable, "\n") 
	noProxyVar := fmt.Sprintf(trimmedNoProxyVar, c.WindowsIp)
	appendNoProxyEnvVar(&noProxyVar)
	
	return fmt.Sprintf("%v\n%v", httpVars, noProxyVar)
}

func appendNoProxyEnvVar(noProxyVar *string) {
	noProxyFromEnv := os.Getenv("NO_PROXY")
	if noProxyFromEnv != "" {
		*noProxyVar = fmt.Sprintf("%v,%v", *noProxyVar, noProxyFromEnv)
	}
}

func (c *ConsoleEnvVarPrinter) PrintUnsetCommands() {
	fmt.Println(strings.Trim(unsetEnvironmentVariables, "\n"))
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