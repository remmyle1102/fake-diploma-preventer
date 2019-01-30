package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	blockchainServer = make(chan []Block)
	t := time.Now()
	genesisBlock := Block{0, "", 0, "", "", "", t.String()}
	spew.Dump(generateBlock)
	Blockchain = append(Blockchain, genesisBlock)

	//start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close() //close the server when no longer need it.

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}
