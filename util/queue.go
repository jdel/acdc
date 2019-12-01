package util // import jdel.org/acdc/util

import (
	"os/exec"
)

// QueueItem represents a queue item
type QueueItem struct {
	Key    string
	Ticket int
	Cmd    *exec.Cmd
}

// InputQueue takes *exec.Cmd to process them
var InputQueue = make(chan QueueItem, 100)

// OutputLog collects outputs from InputQueue
var OutputLog = make(map[string]map[int]string)

// TicketNumber distributes ticket numbers
var TicketNumber = make(chan int, 100)

// NextTicket gets the next ticket in line
var NextTicket func() int

func init() {
	NextTicket = newTicket()
}

func newTicket() func() int {
	i := 0
	return func() int {
		i++
		return i
	}

}

// ProcessQueue processes *exec.Cmd
func ProcessQueue() {
	for {
		item := <-InputQueue
		if OutputLog[item.Key] == nil {
			OutputLog[item.Key] = map[int]string{}
		}
		o, err := item.Cmd.CombinedOutput()
		if err != nil {
			OutputLog[item.Key][item.Ticket] = err.Error()
		}
		OutputLog[item.Key][item.Ticket] = string(o)
	}
}
