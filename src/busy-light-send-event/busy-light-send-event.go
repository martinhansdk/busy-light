package main

import (
	"flag"
	"fmt"
	"encoding/json"
	"net"
)

type EventMessage struct {
	Event* string
	Cmdline* string
	Shellid* string
	Pwd* string
	Exitcode* int
}

func main() {
	msg := EventMessage{
		Event : flag.String("event", "", "The event to send - either \"start\" or \"stop\"."),
		Cmdline : flag.String("cmdline", "", "The commandline to send."),
		Shellid: flag.String("shellid", "", "The shell id to send."),
		Pwd : flag.String("pwd", "", "The pwd to send."),
		Exitcode : flag.Int("exitcode", 0, "The exit code to send."),
	}

	server := flag.String("server", "127.0.0.1:5050", "The IP address and port of the busy-light-server to send the update to in the form server:port")
	
	conn, err := net.Dial("udp", *server)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	
	flag.Parse()
	
	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error:", err)
	}
	conn.Write(b)

	defer conn.Close()
}
