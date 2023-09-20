package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY=3
	MINING_SENDER = "THE BLOCKCHAIN"
	MINING_REWARD =1.0
)
type Block struct{
	timestamp int64;
	nonce int;
	previousHash [32]byte;
	transactions []*Transaction;
}

type Blockchain struct{
	transactionPool []*Transaction;
	chain []*Block; 
	blockchainAddess string
}

type Transaction struct{
	senderBlockchainAddress string
	receiverBlockchainAddress string
	value float32
}

func (b *Block) Print(){
	fmt.Printf("Nonce          %d\n",b.nonce)
	fmt.Printf("Timestamp      %d\n",b.timestamp)
	fmt.Printf("Previous_Hash  %x\n",b.previousHash)
	fmt.Printf("%s\n",strings.Repeat("-",40))
	for _,tx:=range b.transactions{
		tx.Print()
	}
}

func (b *Block) Hash() [32]byte{
	m:=b.MarshalJson()
	// fmt.Println("json is",string(m))
	hash:=sha256.Sum256([]byte(m))
	// fmt.Println("Hash",hash)
	return hash
}

func(b *Block) MarshalJson() []byte{
	m,_:=json.Marshal(struct{
		Timestamp int64 `json:"timestamp"`;
		Nonce int `json:"nonce"`;
		PreviousHash [32]byte `json:"previousHash"`;
		Transactions []*Transaction `json:"transactions"`;
	}{
		Timestamp:b.timestamp ,
		Nonce:b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
	return m
}

func (t *Transaction) Print(){
	fmt.Printf("%s\n",strings.Repeat("-",40))
	fmt.Printf("Sender    %s\n",t.senderBlockchainAddress)
	fmt.Printf("Receiver  %s\n",t.receiverBlockchainAddress)
	fmt.Printf("Value  	  %.001f\n",t.value)
}

func(t *Transaction) MarshalJson() []byte{
	m,_:=json.Marshal(struct{
		Sender string `json:"sender_blockchain_address"`;
		Receiver string `json:"receiver_blockchain_address"`;
		Value float32 `json:"value"`;
	}{
		Sender:t.senderBlockchainAddress ,
		Receiver:t.receiverBlockchainAddress,
		Value: t.value,
	})
	return m
}

func (bc *Blockchain) Print(){

	fmt.Printf("%s\n",strings.Repeat("*",25))

	for key,chain:=range bc.chain{
		fmt.Printf("Chain %s %d %s \n",strings.Repeat("=",25),key,strings.Repeat("=",25))
		chain.Print()
	}
	fmt.Printf("%s\n",strings.Repeat("*",25))
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block{
	b:=newBlock(nonce,previousHash,bc.transactionPool)
	bc.chain = append(bc.chain,b)
	bc.transactionPool=[]*Transaction{}
	return b
}

func (bc *Blockchain) LastBlock() *Block{
	block:=bc.chain[len(bc.chain)-1]
	return block
}

func (bc *Blockchain)AddTransaction(sender string, receiver string, value float32){
	t:=NewTransaction(sender,receiver,value)
	bc.transactionPool = append(bc.transactionPool,t)
}

func (bc *Blockchain) CopyTransaction() []*Transaction{
	transactions:=make([]*Transaction,0)
	for _,t:=range bc.transactionPool{
		transactions=append(transactions,NewTransaction(t.senderBlockchainAddress,t.receiverBlockchainAddress,t.value))
	}
	return transactions;
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte , transactions []*Transaction, difficulty int) bool{
	zeroes:=strings.Repeat("0", difficulty)
	guessBlock:=Block{
		0,
		nonce,
		previousHash,
		transactions,
	}

	guessHashStr:=fmt.Sprintf("%x\n",guessBlock.Hash())
	fmt.Println(guessHashStr)
	return guessHashStr[:difficulty] ==zeroes 
}

func (bc *Blockchain) Mining() bool{
	bc.AddTransaction(MINING_SENDER,bc.blockchainAddess,MINING_REWARD)
	nonce:=bc.proofOfWork()
	previousHash:=bc.LastBlock().Hash()
	bc.CreateBlock(nonce,previousHash)
	log.Print("action=mining, status=success")
	return true
}

func (bc *Blockchain) CalculateAmount(blockchainAddress string) float32{
	var amount float32=0.0
	for _,b:=range(bc.chain){
		for _,t:=range(b.transactions){
			if t.senderBlockchainAddress==blockchainAddress{
				amount-=t.value
			} else if t.receiverBlockchainAddress==blockchainAddress{
				amount+=t.value
			}
		}
	}
	return amount
}

func (bc *Blockchain) proofOfWork() int{
	transactions:=bc.CopyTransaction()
	previousHash:=bc.LastBlock().Hash()
	nonce:=0
	for !bc.ValidProof(nonce,previousHash,transactions, MINING_DIFFICULTY){
		nonce+=1
	}
	return nonce
}

func newBlockchain(blockchainAddress string) *Blockchain{
	b:=&Block{nonce:1}
	bc:=new(Blockchain)
	bc.blockchainAddess=blockchainAddress
	bc.CreateBlock(0,b.Hash())
	return bc
}

func newBlock(nonce int,previousHash [32]byte, transactions []*Transaction) *Block{
	b:=new(Block)
	b.timestamp=time.Now().UnixNano()
	b.previousHash=previousHash
	b.nonce=nonce
	b.transactions=transactions
	return b
}

func NewTransaction(Sender string,Receiver string,Value float32 ) *Transaction{
	return &Transaction{Sender,Receiver,Value }
}

func init(){
	log.SetPrefix("Blockchain: ")
}

func main(){
	myBlockchainAddress:="my_blockchain_address"
	bc:=newBlockchain(myBlockchainAddress)
	bc.AddTransaction("A","B",1)
	bc.Mining()
	bc.AddTransaction("C","D",1.1)
	bc.AddTransaction("B","D",0.1)
	bc.AddTransaction("X","Y",2.1)
	bc.Mining()
	bc.Print()

	fmt.Printf("Balance of A is %.1f\n",bc.CalculateAmount("A"))
	fmt.Printf("Balance of B is %.1f\n",bc.CalculateAmount("B"))
	fmt.Printf("Balance of C is %.1f\n",bc.CalculateAmount("C"))
	fmt.Printf("Balance of D is %.1f\n",bc.CalculateAmount("D"))
	fmt.Printf("Balance of X is %.1f\n",bc.CalculateAmount("X"))
	fmt.Printf("Balance of Y is %.1f\n",bc.CalculateAmount("Y"))
	fmt.Printf("Balance of Miner is %.1f\n",bc.CalculateAmount(myBlockchainAddress))
}