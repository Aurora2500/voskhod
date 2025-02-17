package main

import (
	"log"
	"voskhod/protocol"
)

func main() {
	url := "gemini://geminiprotocol.net/"
	db, err := protocol.InitCertsDB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	response, err := protocol.FetchUrl(url, db)
	if err != nil {
		log.Fatalln(err.Error())
	}

	println(response)
}
