package cmd

import (
	"log"
	"storage/internal/controller"
)

func Start() {
	if err := controller.StartServer(); err != nil {
		log.Fatal(err)
	}
}
