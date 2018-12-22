package main

import (
	"encoding/json"
	"fmt"
	"github.com/tarm/serial"
	"net"
	"log"
	"path/filepath"
	"regexp"
)

const MAX_LEN = 4096

var regexes = [...]*regexp.Regexp {
	regexp.MustCompile("forge"),
	regexp.MustCompile("p4 sync"),
}

var event2pattern = map[string]byte{
    "start":  'f',
    "stop":   '+',
}


type EventMessage struct {
	Event string
	Cmdline string
	Shellid string
	Pwd string
	Exitcode int
}

func main() {
	var serport *serial.Port;

	serialPorts, err := filepath.Glob("/dev/ttyUSB*")
	
	if err != nil {
		log.Fatal(err)
	}
	for i := range serialPorts {
		config := &serial.Config{Name: serialPorts[i], Baud: 9600}
		serport, err = serial.OpenPort(config)
		
		if err != nil {
			// ok, try the next port
			log.Printf("Could not open port %s: %s", serialPorts[i], err)
		} else {
			log.Printf("Opened port %s", serialPorts[i])
			break
		}
	}

	
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

			filterMessage(eventMessage, serport);
		}
	}
}

func filterMessage(eventMessage EventMessage, serport *serial.Port) {
	if(eventMessage.Event != "start" && eventMessage.Event != "stop") {
		return
	}
	
	for _, re := range regexes {
		if(re.MatchString(eventMessage.Cmdline)) {
			processEvent(eventMessage, serport)
		}
	}
}

func processEvent(eventMessage EventMessage, serport *serial.Port) {
	fmt.Printf("%s : %s : %s\n", eventMessage.Event, eventMessage.Shellid, eventMessage.Cmdline)
	serport.Write([]byte{event2pattern[eventMessage.Event]});
}
