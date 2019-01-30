package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// blockchainServer handles incoming concurrent Blocks
var blockchainServer chan []Block

// handle incomming connection
func handleConn(conn net.Conn) {
	defer conn.Close()
	io.WriteString(conn, "Fake Diploma Preventer")

	scanner := bufio.NewScanner(conn)
	go func() {
		var studentName string
		var eduInstitution string
		var grade int
		for scanner.Scan() {
			io.WriteString(conn, "Enter student's name\n")
			studentName = scanner.Text()
			io.WriteString(conn, "Enter grade \n")
			inputGrade := scanner.Text()
			grade, _ = strconv.Atoi(inputGrade)
			io.WriteString(conn, "Enter Edu Institution's Name\n")
			eduInstitution = scanner.Text()
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], studentName, grade, eduInstitution)
			if err != nil {
				log.Println(err)
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}
			blockchainServer <- Blockchain // throw new Blockchain into the created channel
			io.WriteString(conn, "\n Enter a new Student")
		}
	}()

	//BroadCasting the blockchain

	go func() {
		for {
			time.Sleep(30 * time.Second)
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()

	for range blockchainServer {
		spew.Dump(Blockchain)
	}
}
