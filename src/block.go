package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// init difficuty: 2
const DIFFICULTY int = 2

type Block struct {
	Index      int
	Timestamp  string
	Content    string
	Hash       string
	PreHash    string
	Difficulty int
	Nonce      string
}

// 计算hash
func CalculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.Content + block.PreHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

// 生成新的区块
func GenerateBlock(oldBlock Block, Content string) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Content = Content
	newBlock.PreHash = oldBlock.Hash
	newBlock.Hash = CalculateHash(newBlock)
	newBlock.Difficulty = DIFFICULTY + oldBlock.Index

	// 计算
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(CalculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(CalculateHash(newBlock), "continu computed...")
			continue
		} else {
			fmt.Println(CalculateHash(newBlock), "牛逼了💯, your right!")
			newBlock.Hash = CalculateHash(newBlock)
			break
		}
	}
	return newBlock, nil
}

// 验证区块
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	// 新区块中的prehash应该是前个区块的hash
	if oldBlock.Hash != newBlock.PreHash {
		return false
	}

	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// 验证区块是否符合要求
func isHashValid(hash string, difficulty int) bool {
	/* hash的前缀0的个数 */
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}
