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

const defaultNoProxyHosts = "localhost,127.0.0.1"

const unsetEnvironmentVariables = `
unset HTTPS_PROXY
unset HTTP_PROXY
unset https_proxy
unset http_proxy
unset NO_PROXY
unset no_proxy
`

type ConsoleEnvVarPrinter struct{
	WindowsIp string
	PxProxyPort int
	NoProxy []string
}

func (c *ConsoleEnvVarPrinter) PrintExportCommands() {
	fmt.Println(c.buildPrintExportCommands())
}

func (c *ConsoleEnvVarPrinter) buildPrintExportCommands() string {
	trimmedHttpVars := strings.Trim(exportHttpProxyVariables, "\n")
	httpVars := fmt.Sprintf(trimmedHttpVars, c.WindowsIp, c.PxProxyPort)
	
	upperNoProxyVar := c.buildNoProxyEnvVar("NO_PROXY")
	lowerNoProxyVar := c.buildNoProxyEnvVar("no_proxy")

	return fmt.Sprintf("%v\n%v\n%v", httpVars, upperNoProxyVar, lowerNoProxyVar)
}

func (c *ConsoleEnvVarPrinter) buildNoProxyEnvVar(envVarName string) string {
	out := fmt.Sprintf("export %v=%v,%v", envVarName, defaultNoProxyHosts, c.WindowsIp)
	out += appendEnvVarIfSet(envVarName)
	out += c.appendNoProxyConfigIfSet(envVarName)
	return out
}

func appendEnvVarIfSet(key string) string {
	value := os.Getenv(key)
	if value != "" {
		return fmt.Sprintf(",%v", value)
	}
	return ""
}

func (c *ConsoleEnvVarPrinter) appendNoProxyConfigIfSet(envVarName string) string {
	if len(c.NoProxy) > 0 {
		noProxyList := strings.Join(c.NoProxy, ",")
		return fmt.Sprintf(",%v", noProxyList)
	}
	return ""
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