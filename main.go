package main

import (
	"log"

	"github.com/chaihaobo/be-template/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Println("application run failed:", err.Error())
	}
}
