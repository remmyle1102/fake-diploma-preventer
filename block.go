package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const difficulty = 1

// Block store data that will be written to the blockchain
type Block struct {
	Index          int // position of the data record in the blockchain
	StudentName    string
	Grade          int
	EduInstitution string // name of the educational institution
	Hash           string // is a SHA256 identifier representing data record
	PrevBlockHash  string // is the previous block hash record in the chain
	Timestamp      string // the time the data is written
	Difficulty     int
	Nonce          string
}

func calculateHash(block Block) string {
	record := string(block.Index) + string(block.Grade) + block.EduInstitution + block.PrevBlockHash + block.Timestamp + block.StudentName + block.Nonce
	hashFunc := sha256.New()
	hashFunc.Write([]byte(record))
	hashed := hashFunc.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(prevBlock Block, studentName string, grade int, eduInstitution string) (Block, error) {
	var newBlock Block
	t := time.Now()

	newBlock.Index = prevBlock.Index + 1
	newBlock.StudentName = studentName
	newBlock.Grade = grade
	newBlock.EduInstitution = eduInstitution
	newBlock.PrevBlockHash = prevBlock.Hash
	newBlock.Timestamp = t.String()
	newBlock.Difficulty = difficulty

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), " Keep Working")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), " Block created, Good job!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}
	}

	return newBlock, nil
}

func isBlockValid(newBlock, prevBlock Block) bool {
	if prevBlock.Index+1 != newBlock.Index {
		return false
	}
	if prevBlock.Hash != newBlock.PrevBlockHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
