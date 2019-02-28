package main

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// Block store data that will be written to the blockchain
type Block struct {
	Index          int // position of the data record in the blockchain
	StudentName    string
	StudentID      string
	Grade          int
	EduInstitution string // name of the educational institution
	Hash           string // is a SHA256 identifier representing data record
	PrevBlockHash  string // is the previous block hash record in the chain
	Timestamp      string // the time the data is written
	Validator      string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

var tempBlocks []Block

// candidateBlocks, a chnnel of blocks, handles incoming blocks for validation
var candidateBlocks = make(chan Block)

// announcements broadcasts winning validator to all nodes
var announcements = make(chan string)

// validators keeps track of open validators and balances
var validators = make(map[string]int)

// mutex prevents data races and make sure blocks aren't generated at the same time
var mutex = &sync.Mutex{}

func calculateHash(s string) string {
	hashFunc := sha256.New()
	hashFunc.Write([]byte(s))
	hashed := hashFunc.Sum(nil)
	return hex.EncodeToString(hashed)
}

func calculateBlockHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.StudentName + block.StudentName + string(block.EduInstitution) + string(block.Grade) + block.PrevBlockHash
	return calculateHash(record)
}

func createGenesisBlock() Block {
	genesisBlock := Block{}
	t := time.Now()
	genesisBlock = Block{Index: 0, StudentID: "", StudentName: "", EduInstitution: "", Hash: calculateBlockHash(genesisBlock), Timestamp: t.String()}
	return genesisBlock
}

func generateBlock(prevBlock Block, studentName string, studentID string, grade int, eduInstitution string, address string) (Block, error) {
	var newBlock Block
	t := time.Now()

	newBlock.Index = prevBlock.Index + 1
	newBlock.StudentName = studentName
	newBlock.StudentID = studentID
	newBlock.Grade = grade
	newBlock.EduInstitution = eduInstitution
	newBlock.PrevBlockHash = prevBlock.Hash
	newBlock.Timestamp = t.String()
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Validator = address
	return newBlock, nil
}

func isBlockValid(newBlock, prevBlock Block) bool {
	if prevBlock.Index+1 != newBlock.Index {
		return false
	}
	if prevBlock.Hash != newBlock.PrevBlockHash {
		return false
	}
	if calculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
