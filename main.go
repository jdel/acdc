package main

import (
	"log"

	"github.com/google/gops/agent"
	"github.com/jdel/acdc/cmd"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	cmd.Execute()
}
