package config

import (
	"bytes"
	"fmt"

	"github.com/3th1nk/cidr"
	"github.com/spf13/viper"

	"org.samba/isetta/helper"
	log "org.samba/isetta/simplelogger"
)

var defaults = map[string]string{
	"general.internet_access_test_url": "https://www.google.com/",
	"general.log_level":                "info",
	"network.wsl_to_windows_subnet":    "169.254.254.0/24",
	"network.px_proxy_port":            "3128",
	"dns.public_server":                "8.8.8.8",
}

type Config struct {
	General General
	Network Network
	Dns     Dns
}

type General struct {
	InternetAccessTestUrl string `mapstructure:"internet_access_test_url" validate:"url"`
	LogLevel              string `mapstructure:"log_level" validate:"alpha"`
}

type Network struct {
	WslToWindowsSubnet string `mapstructure:"wsl_to_windows_subnet" validate:"cidrv4"`
	PxProxyPort        int    `mapstructure:"px_proxy_port" validate:"min=1,max=65535"`
	P2p                P2p
	NoProxy   []string `mapstructure:"no_proxy"`
}

type P2p struct {
	WindowsIp  string
	LinuxIp    string
	SubnetMask string
}

type Dns struct {
	InternalServer string `mapstructure:"internal_server" validate:"required,ip4_addr"`
	PublicServer   string `mapstructure:"public_server" validate:"ip4_addr"`
}

func init() {
	viper.SetConfigName(".isetta")
	viper.SetConfigType("toml")
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
}

func FromConfigFile(configPath string, validLogLevels []string) Config {
	conf := readConfigFromFile(configPath)
	validateAnDetermineIps(&conf, validLogLevels)
	return conf
}

func readConfigFromFile(configPath string) Config {
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Logger.Info("Error reading config file from %v, error was: %v", configPath, err)
	}
	return unmarshalConfig()
}

func GetProxyUrl(conf Config) string {
	return fmt.Sprintf("http://%v:%v", conf.Network.P2p.WindowsIp, conf.Network.PxProxyPort)
}

// for testing
func FromByteBuffer(buffer *bytes.Buffer, validLogLevels []string) Config {
	conf := readConfigFromBuffer(buffer)
	validateAnDetermineIps(&conf, validLogLevels)
	return conf
}

func readConfigFromBuffer(buffer *bytes.Buffer) Config {
	err := viper.ReadConfig(buffer)
	helper.AssertNoError(err, "error loading config")
	return unmarshalConfig()
}

func validateAnDetermineIps(conf *Config, validLogLevels []string) {
	myValidator := NewValidator(conf, validLogLevels)
	err := myValidator.DoValidate()
	helper.AssertNoError2(err)

	err = determineP2pAddresses(conf)
	helper.AssertNoError2(err)
}

func unmarshalConfig() Config {
	var conf Config
	err := viper.Unmarshal(&conf)
	helper.AssertNoError(err, "error parsing config")
	return conf
}

func determineP2pAddresses(conf *Config) error {
	subnet := parseCidr(conf.Network.WslToWindowsSubnet)
	conf.Network.P2p.SubnetMask = subnet.Mask()
	conf.Network.P2p.WindowsIp, conf.Network.P2p.LinuxIp = determineFirstTwoIps(subnet)
	return nil
}

func parseCidr(subnet string) *cidr.CIDR {
	c, err := cidr.ParseCIDR(subnet)
	helper.AssertNoError(err, "error parsing wsl_to_windows_subnet %v.", subnet)
	return c
}

func determineFirstTwoIps(subnet *cidr.CIDR) (string, string) {
	firstIp := ""
	secondIp := ""

	cnt := 0
	// iterate all IPs in subnet
	// skip 1st (network address)
	subnet.ForEachIP(func(ip string) error {
		if cnt == 1 {
			firstIp = ip
		} else if cnt == 2 {
			secondIp = ip
		}
		cnt++
		return nil
	})
	return firstIp, secondIp
}
