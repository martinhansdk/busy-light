package main

import (
	"encoding/json"
	"fmt" 
	"net"
	"regexp"
)

const MAX_LEN = 4096

var regexes = [...]*regexp.Regexp {
	regexp.MustCompile("forge"),
	regexp.MustCompile("p4 sync"),
}

type EventMessage struct {
	Event string
	Cmdline string
	Shellid string
	Pwd string
	Exitcode int
}

func main() {
	msg := make([]byte, MAX_LEN)
	addr := net.UDPAddr{
		Port: 5050,
		IP: net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		len,_,err := ser.ReadFromUDP(msg)
		
		if(len < MAX_LEN) {
			if err !=  nil {
				fmt.Printf("Some error  %v", err)
				continue
			}

			var eventMessage EventMessage
			err := json.Unmarshal(msg[0:len], &eventMessage)
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Printf("%+v\n", eventMessage)

			filterMessage(eventMessage);
		}
	}
}

func filterMessage(eventMessage EventMessage) {
	if(eventMessage.Event != "start" && eventMessage.Event != "stop") {
		return
	}
	
	for _, re := range regexes {
		if(re.MatchString(eventMessage.Cmdline)) {
			processEvent(eventMessage)
		}
	}
}

func processEvent(eventMessage EventMessage) {
	fmt.Printf("%s : %s : %s\n", eventMessage.Event, eventMessage.Shellid, eventMessage.Cmdline)
}
