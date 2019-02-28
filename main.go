package main

import (
	"log"
	"net"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	genesisBlock := createGenesisBlock()
	if len(Blockchain) != 0 {
		spew.Dump(generateBlock)
	} else {
		spew.Dump(genesisBlock)
	}
	Blockchain = append(Blockchain, genesisBlock)

	//start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close() //close the server when no longer need it.

	go func() {
		for candidate := range candidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
		}
	}()

	go func() {
		for {
			pickWinner()
		}
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}
