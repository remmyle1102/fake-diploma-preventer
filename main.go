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
		genesisBlock := Block{0, 0, "", "", "", t.String()}
		spew.Dump(generateBlock)
		Blockchain = append(Blockchain, genesisBlock)
	}()
	log.Fatal(run())
}
