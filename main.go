package main

import (
	"log"

	"gitlab.seakoi.net/engineer/backend/be-template/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Println("application run failed:", err.Error())
	}
}
