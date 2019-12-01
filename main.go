package main // import jdel.org/acdc/main

import (
	"log"

	"github.com/google/gops/agent"
	"jdel.org/acdc/cmd"
	"jdel.org/acdc/util"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}

	go util.ProcessQueue()

	cmd.Execute()
}
