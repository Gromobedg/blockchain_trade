package logerror

import (
	"log"
	"syscall"
)

func Log(err interface{}) {
	log.Printf("ERROR: %v\n", err)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}