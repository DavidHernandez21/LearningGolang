package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/google/uuid"
)

func main() {
	s := "Make the computer speak"
	uuidWithHyphen := fmt.Sprintf("%v.wav", uuid.New().String())

	cmd := exec.Command("espeak", "-w", uuidWithHyphen, s)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
