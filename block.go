package main

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// Blockchain is a series of validated Blocks
var Blockchain []Block

// Block store data that will be written to the blockchain
type Block struct {
	Index          int // position of the data record in the blockchain
	StudentID      string
	StudentName    string
	Grade          int
	EduInstitution string // name of the educational institution
	Hash           string // is a SHA256 identifier representing data record
	PrevBlockHash  string // is the previous block hash record in the chain
	Timestamp      string // the time the data is written
}

var mutex = &sync.Mutex{}

func calculateHash(block Block) string {
	record := string(block.Index) + string(block.Grade) + block.EduInstitution + block.PrevBlockHash + block.Timestamp + block.StudentName
	hashFunc := sha256.New()
	hashFunc.Write([]byte(record))
	hashed := hashFunc.Sum(nil)
	return hex.EncodeToString(hashed)
}

func createGenesisBlock() Block {
	genesisBlock := Block{}
	t := time.Now()
	genesisBlock = Block{Index: 0, StudentID: "", StudentName: "", EduInstitution: "", Hash: calculateHash(genesisBlock), Timestamp: t.String()}
	return genesisBlock
}

func generateBlock(prevBlock Block, studentID string, studentName string, grade int, eduInstitution string) (Block, error) {
	var newBlock Block
	t := time.Now()

	newBlock.Index = prevBlock.Index + 1
	newBlock.StudentID = studentID
	newBlock.StudentName = studentName
	newBlock.Grade = grade
	newBlock.EduInstitution = eduInstitution
	newBlock.PrevBlockHash = prevBlock.Hash
	newBlock.Timestamp = t.String()
	newBlock.Hash = calculateHash(newBlock)

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
