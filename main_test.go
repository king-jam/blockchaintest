package main

import (
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	_ "net/http/pprof"
	"os"
	"sync"
	"testing"
	"time"
)

var data = []byte("mytest")

func benchmarkBlockchain(datasize int64, b *testing.B) {
	file, err := ioutil.TempFile("/mnt/test/", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	Blockchain := make([]Block, 1000000)
	enc := gob.NewEncoder(file)
	var mutex = &sync.Mutex{}
	go func() {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{0, t, []byte{}, calculateHash(genesisBlock), ""}
		//spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		err := enc.Encode(Blockchain)
		if err != nil {
			os.Remove(file.Name())
			log.Fatalf("Broke: %s", err.Error())
		}
		file.Sync()
		mutex.Unlock()
	}()

	time.Sleep(1 * time.Second)

	data := make([]byte, datasize)
	_, err = rand.Read(data)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	b.ResetTimer()
	b.SetBytes(datasize)
	for i := 0; i < b.N; i++ {
		mutex.Lock()
		newBlock := generateBlock(Blockchain[len(Blockchain)-1], data)
		mutex.Unlock()

		if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
			Blockchain = append(Blockchain, newBlock)
			err := enc.Encode(Blockchain)
			if err != nil {
				os.Remove(file.Name())
			}
			file.Sync()
			//spew.Dump(Blockchain)
		}
	}
}

func Benchmark1K(b *testing.B)   { benchmarkBlockchain(1000, b) }
func Benchmark4K(b *testing.B)   { benchmarkBlockchain(4000, b) }
func Benchmark8K(b *testing.B)   { benchmarkBlockchain(8000, b) }
func Benchmark16K(b *testing.B)  { benchmarkBlockchain(16000, b) }
func Benchmark32K(b *testing.B)  { benchmarkBlockchain(32000, b) }
func Benchmark64K(b *testing.B)  { benchmarkBlockchain(64000, b) }
func Benchmark128K(b *testing.B) { benchmarkBlockchain(128000, b) }
func Benchmark256K(b *testing.B) { benchmarkBlockchain(256000, b) }
func Benchmark512K(b *testing.B) { benchmarkBlockchain(512000, b) }
func Benchmark1M(b *testing.B)   { benchmarkBlockchain(1000000, b) }
func Benchmark2M(b *testing.B)   { benchmarkBlockchain(2000000, b) }
func Benchmark4M(b *testing.B)   { benchmarkBlockchain(4000000, b) }
func Benchmark8M(b *testing.B)   { benchmarkBlockchain(8000000, b) }
func Benchmark16M(b *testing.B)  { benchmarkBlockchain(16000000, b) }
func Benchmark32M(b *testing.B)  { benchmarkBlockchain(32000000, b) }
func Benchmark64M(b *testing.B)  { benchmarkBlockchain(64000000, b) }
