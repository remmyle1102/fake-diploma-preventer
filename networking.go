package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// blockchainServer handles incoming concurrent Blocks
var blockchainServer chan []Block

// handle incomming connection
func handleConn(conn net.Conn) {
	defer conn.Close()
	var studentID string
	var studentName string
	var eduInstitution string
	var grade int
	io.WriteString(conn, "Fake Diploma Preventer\n")

	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()

	// validator address
	var address string

	// allow user to allocate number of tokens to stake
	// the greater the number of tokens, the greater chance to forging a new block
	io.WriteString(conn, "\nEnter token balance:")
	balanceScanner := bufio.NewScanner(conn)
	for balanceScanner.Scan() {
		balance, err := strconv.Atoi(balanceScanner.Text())
		if err != nil {
			log.Printf("%v is not a number: %v", balanceScanner.Text(), err)
			return
		}
		t := time.Now()
		address = calculateHash(t.String())
		validators[address] = balance
		fmt.Println(validators)
		break
	}

	scanner := bufio.NewScanner(conn)

	go func() {
		for scanner.Scan() {
			io.WriteString(conn, "Enter student's ID")
			for scanner.Scan() {
				studentID = scanner.Text()
				break
			}
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
					delete(validators, address)
					/*
						This is why Proof of Stake is secure,
						If you do something wrong and get caught, you will lose your entire staked token balance.
					*/
					conn.Close()
				}
				grade = studentGrade
				break
			}
			io.WriteString(conn, "Enter Edu Institution's Name\n")
			for scanner.Scan() {
				eduInstitution = scanner.Text()
				break
			}
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], studentName, studentID, grade, eduInstitution, address)
			if err != nil {
				log.Println(err)
				continue
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				// send the created block to the candidateBlocks channel for futher processing
				candidateBlocks <- newBlock
			}
			io.WriteString(conn, "Enter new student\n")
		}
	}()

	//simulate receiving broadcast
	for {
		time.Sleep(time.Minute)
		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, string(output)+"\n")
	}
}

func pickWinner() {
	time.Sleep(30 * time.Second)
	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	lotteryPool := []string{}
	if len(temp) > 0 {
	OUTER:
		for _, block := range temp {
			// if already in lottery pool, skip
			for _, node := range lotteryPool {
				if block.Validator == node {
					continue OUTER
				}
			}

			// lock list of validators to prevent data race
			mutex.Lock()
			setValidators := validators
			mutex.Unlock()

			k, ok := setValidators[block.Validator]
			if ok {
				for i := 0; i < k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		//randomly pick a winner from lottery pool
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]

		// add block of winner to blockchain and let all the other nodes know
		for _, block := range temp {
			if block.Validator == lotteryWinner {
				mutex.Lock()
				Blockchain = append(Blockchain, block)
				mutex.Unlock()
				for range validators {
					announcements <- "\nwinning validator: " + lotteryWinner + "\n"
				}
				break
			}
		}
	}
	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()
}
