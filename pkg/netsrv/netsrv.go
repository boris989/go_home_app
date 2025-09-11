package netsrv

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func StartServer(port string, handler func(word string) []string) {
	listener, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Сервер запущен на порту %s\n", port)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleRequest(conn, handler)
	}
}

func handleRequest(conn net.Conn, handler func(word string) []string) {
	defer conn.Close()
	defer fmt.Println("Connection Closed")

	r := bufio.NewReader(conn)
	for {
		msg, _, err := r.ReadLine()
		if err != nil {
			return
		}

		word := string(msg)

		word = strings.TrimSpace(word)

		res := handler(word)

		resMsg := []byte("")

		for _, line := range res {
			resMsg = append(resMsg, []byte(line)...)
			resMsg = append(resMsg, []byte("\n")...)
		}

		resMsg = append(resMsg, []byte("Поиск завершен\n")...)

		_, err = conn.Write(resMsg)
		if err != nil {
			return
		}
	}
}
