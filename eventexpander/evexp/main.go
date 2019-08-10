package main

import (
	"log"
	"os"

	"github.com/graham/rrule/eventexpander/evexp/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
