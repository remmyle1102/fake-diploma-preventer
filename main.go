package main

import (
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{Difficulty: difficulty, Timestamp: t.String(), Index: 0, Nonce: "", Hash: calculateHash(genesisBlock)}

		spew.Dump(genesisBlock)

		Blockchain = append(Blockchain, genesisBlock)
	}()
	log.Fatal(run())
}
