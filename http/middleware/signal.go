package middleware

import (
	"fmt"
	"os"
	"os/signal"
)

var (
	signalChannel = make(chan os.Signal)
)

func Signal() {
	go func() {
		signal.Notify(signalChannel, os.Interrupt, os.Kill)
		<-signalChannel

		exit()
	}()
}

func exit() {
	defer os.Exit(0)

	fmt.Println("perform cleanup...")

	fmt.Println("complete!")
}
