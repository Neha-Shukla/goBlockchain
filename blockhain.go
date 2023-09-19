package main

import (
	"fmt"
	"log"
	"time"
)

type Block struct{
	nonce int;
	previousHash string;
	timestamp int64;
	transactions []string;
}

type Blockchain struct{
	transactionPool []string;
	chain []*Block;
}

func (b *Block) Print(){
	fmt.Printf("Nonce          %d\n",b.nonce)
	fmt.Printf("Timestamp      %d\n",b.timestamp)
	fmt.Printf("Previous_Hash  %s\n",b.previousHash)
	fmt.Printf("Transactions   %d\n",b.transactions)
}

func (bc *Blockchain) Print(){
	fmt.Println(bc)
}


func init(){
	log.SetPrefix("Blockchain: ")
}

func newBlockchain() *Blockchain{
bc:=new(Blockchain)
bc.createBlock(0,"init hash")
return bc
}

func newBlock(nonce int,previousHash string) *Block{
b:=new(Block)
b.timestamp=time.Now().UnixNano()
b.previousHash=previousHash
b.nonce=nonce
return b
}

func (bc *Blockchain) createBlock(nonce int, previousHash string) *Block{
	b:=newBlock(nonce,previousHash)
	bc.chain = append(bc.chain,b)
	return b
}

func main(){
	bc:=newBlockchain()
	bc.createBlock(1,"hash1")
	bc.createBlock(2,"hash2")
	bc.createBlock(3,"hash3")
	bc.Print()
}