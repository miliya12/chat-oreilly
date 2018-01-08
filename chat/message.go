package main

import "time"

// message expresses one message
type message struct {
	Name    string
	Message string
	When    time.Time
}
