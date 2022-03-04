package btc

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blockBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// tip 存储最后一个块的hash
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) Close() {
	bc.db.Close()
}

func (bc *Blockchain) AddBlock(data string) {
	// var lastHash []byte
	//
	// err := bc.db.View(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte(blockBucket))
	// 	lastHash = b.Get([]byte("l"))
	//
	// 	return nil
	// })
	//
	// if err != nil {
	// 	log.Panic(err)
	// }
	//
	// newBlock := NewBlock(data, lastHash)
	// err = bc.db.Update(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte(blockBucket))
	// 	err := b.Put(newBlock.Hash, newBlock.Seralize())
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	//
	// 	err = b.Put([]byte("l"), newBlock.Hash)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	//
	// 	bc.tip = newBlock.Hash
	// 	return nil
	// })
	//
}

func dbExist() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// 创建一个有创世快的新连
func NewBlockchain(data string) *Blockchain {
	if dbExist() == false {
		fmt.Println("No existing blockchain found, Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		tip = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{
		tip: tip,
		db:  db,
	}
	return &bc
}

// 创建一个新的区块链数据库， address 用来接收挖出的创世快的奖励
func CreateBlockChain(address string) *Blockchain {
	if dbExist() {
		fmt.Println("Blockchain already exists")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(blockBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Seralize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}

		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{
		tip: tip,
		db:  db,
	}
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		encodeBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodeBlock)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash
	return block
}
