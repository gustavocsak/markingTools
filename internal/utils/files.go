package utils

import (
	"fmt"
	"log"
	"os"
)

func Log(msgType string, msg string) {
	prefix := fmt.Sprintf("[%s]", msgType)
	log.Printf("%10s %s\n", prefix, msg)
}

func CheckExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
