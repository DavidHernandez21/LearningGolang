package keyboardinterrup

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Listening() {
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	log.Printf("received signal: %v...cancelling the channel", s)

}
