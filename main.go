package main

import (
	"crypto/tls"
	"io"
	"log"
	"strings"
)

func main() {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", "geminiquickst.art:1965", config)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close()

	conn.Write([]byte("gemini://geminiquickst.art/\r\n"))
	sb := strings.Builder{}
	io.Copy(&sb, conn)
	println(sb.String())
}
