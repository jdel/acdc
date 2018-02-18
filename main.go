package main

import (
	"log"

	"github.com/google/gops/agent"
	"github.com/jdel/acdc/cmd"
	"github.com/jdel/acdc/util"
)

func main() {
	if err := agent.Listen(&agent.Options{}); err != nil {
		log.Fatal(err)
	}

	go util.ProcessQueue()

	cmd.Execute()
}
