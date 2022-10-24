package core

import "org.samba/isetta/mocks"

// shared between the tests of this package :-/
var mockWinChecker *mocks.WindowsChecker
var mockWinConfigurer *mocks.WindowsConfigurer
var mockHttpChecker *mocks.HttpChecker
var mockDnsConfigurer *mocks.DnsConfigurer
var mockEnvVarPrinter *mocks.EnvVarPrinter
var mockLinuxPinger *mocks.LinuxPinger
var mockLinuxConfigurer *mocks.LinuxConfigurer
