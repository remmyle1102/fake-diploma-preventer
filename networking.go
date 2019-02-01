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
	io.WriteString(conn, "Fake Diploma Preventer\n")

	scanner := bufio.NewScanner(conn)
	go func() {
		for scanner.Scan() {
			var studentName string
			var eduInstitution string
			var grade int
			io.WriteString(conn, "Enter student's name\n")
			for scanner.Scan() {
				studentName = scanner.Text()
				break
			}
			io.WriteString(conn, "Enter grade \n")
			for scanner.Scan() {
				studentGrade, err := strconv.Atoi(scanner.Text())
				if err != nil {
					log.Printf("%v is not a number: %v ", scanner.Text(), err)
					continue
				}
				grade = studentGrade
				break
			}
			io.WriteString(conn, "Enter Edu Institution's Name\n")
			for scanner.Scan() {
				eduInstitution = scanner.Text()
				break
			}
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], studentName, grade, eduInstitution)
			if err != nil {
				log.Println(err)
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}
			blockchainServer <- Blockchain // throw new Blockchain into the created channel
			io.WriteString(conn, "Enter a new Student\n")
		}
	}()

	//BroadCasting the blockchain

	go func() {
		for {
			time.Sleep(10 * time.Second)
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
