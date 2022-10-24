package windows

import (
	"fmt"
	"time"

	"org.samba/isetta/gsudo"
	"org.samba/isetta/helper"
)


type WindowsConfigurerImpl struct {
	WindowsIp string
	SubnetMask string
	PxProxyPort int
	Gsudo *gsudo.Gsudo	
}

func (w *WindowsConfigurerImpl) Init() {
	w.Gsudo.Init()
}

func (w *WindowsConfigurerImpl) Cleanup() {
	w.Gsudo.Cleanup()
}

func (w *WindowsConfigurerImpl) AddP2pAddress(successChecker func() bool) error {
	cmd := fmt.Sprintf("netsh interface ip add address \"vEthernet (WSL)\" %v %v", w.WindowsIp, w.SubnetMask)
	w.Gsudo.RunElevated(cmd, false)
	
	// letting config change settle
	return helper.Retry(helper.RetryParams{
		Description: "Setting Windows P2P address",
		Attempts:    10,
		Sleep:       100 * time.Millisecond,
		Func:        successChecker,
	})
}

func (w *WindowsConfigurerImpl) SetPortProxy(successChecker func() bool) error {
	// re-run config when last config attempt had issues
	configFunc := func () bool  {
		w.resetPortProxy()
		w.addPortProxy()
		return successChecker()
	}
	
	return helper.Retry(helper.RetryParams{
		Description: "Setting Windows portproxy",
		Attempts:    10,
		Sleep:       200 * time.Millisecond,
		Func:        configFunc,
	})
}

func (w *WindowsConfigurerImpl) resetPortProxy() {
	w.Gsudo.RunElevated("netsh interface portproxy reset")
}

func (w *WindowsConfigurerImpl) addPortProxy() {
	cmd := fmt.Sprintf("netsh interface portproxy add v4tov4 listenaddress=%[1]v listenport=%[2]v connectaddress=127.0.0.1 connectport=%[2]v", w.WindowsIp, w.PxProxyPort)
	w.Gsudo.RunElevated(cmd)
}
