package app

import (
	"fmt"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	ver "github.com/vault-thirteen/auxie/Versioneer"
)

func showIntro(v *ver.Versioneer, serviceName string) {
	v.ShowIntroText(serviceName)
	v.ShowComponentsInfoText()
	fmt.Println()
}

func waitForQuitSignalFromOS(serverMustBeStopped *chan bool) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range osSignals {
			switch sig {
			case syscall.SIGINT,
				syscall.SIGTERM:
				log.Println(c.MsgQuitSignalIsReceived, sig)
				*serverMustBeStopped <- true
			}
		}
	}()
}
