package wflags

import (
	"flag"
	"fmt"
	"log"
)

var (
	Environment string
)

func ProcessFlags() {
	flag.StringVar(&Environment, "env", "", "application environment")

	validateFlags()
}

func validateFlags() {
	flag.Parse()

	fmt.Println("Process Flags Environment", Environment)

	if Environment != "development" && Environment != "production" {
		log.Fatalln("[ERROR] You must specify valid Environment")
	}
}
