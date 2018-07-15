package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// 存储区块链数据
var Blockchain []Block

// 全网的公链的block
var blockchainServer chan []Block

// lock
var mutex = &sync.Mutex{}

// listen global block port
var port string = ":9000"

// 接受post过来参数
type Message struct {
	Content string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	blockchainServer = make(chan []Block)

	t := time.Now()
	// 创世纪块
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), "first block", CalculateHash(genesisBlock), "", DIFFICULTY, ""}
	spew.Dump(genesisBlock)
	mutex.Lock()
	Blockchain = append(Blockchain, genesisBlock)
	mutex.Unlock()

	// listen 9000 port
	server, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		connection, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(connection)
	}
}

// handleConnection
func handleConnection(connection net.Conn) {
	defer connection.Close()

	io.WriteString(connection, "please write any contents:")

	scanner := bufio.NewScanner(connection)

	go func() {
		for scanner.Scan() {
			content := scanner.Text()

			newBlock, err := GenerateBlock(Blockchain[len(Blockchain)-1], string(content))
			if err != nil {
				log.Println(err)
				continue
			}

			// 取最长的链为有效链
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}

			blockchainServer <- Blockchain
			io.WriteString(connection, "\n please write any contents:")
		}
	}()

	// The timer is broadcast
	go func() {
		for {
			time.Sleep(40 * time.Second)
			mutex.Lock()
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Unlock()
			io.WriteString(connection, string(output))
		}
	}()

	for _ = range blockchainServer {
		spew.Dump(Blockchain)
	}
}

//以最长的链为有效链
func replaceChain(newBlocks []Block) {
	mutex.Lock()
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
	mutex.Unlock()
}
