package core

import (
	"github.com/boltdb/bolt"
	"bubba_coin/utils"
	"errors"
)

const (
	dbFile        = "bubba.db"
	blockBucket   = "blocks"
	versionBucket = "versions"
	hashComplexity = 20
)

type Blockchain struct {
	versions       map[int]int
	currentVersion int
	lastBlock      []byte
	db             *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bci *BlockchainIterator) Next() *Block {
	var block *Block

	err := bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		if b == nil {
			return errors.New("block bucket was not found during iteration")
		}
		block = Deserialize(b.Get(bci.currentHash))
		bci.currentHash = block.Header.PreviousBlockHash

		return nil
	})
	if err != nil {
		panic(err)
	}

	return block
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.lastBlock, bc.db}

	return bci
}

func (bc *Blockchain) CloseDb() {
	bc.db.Close()
}

func (bc *Blockchain) AddBlock(data string) error {
	newBlock := NewBlock(data, bc.lastBlock, bc.versions[bc.currentVersion])
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		if b == nil {
			return errors.New("block bucket was not found during update")
		}
		err := b.Put(newBlock.Header.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), newBlock.Header.Hash)
		if err != nil {
			return err
		}
		bc.lastBlock = newBlock.Header.Hash
		return nil
	})

	return err
}

func (bc *Blockchain) ChangeVersion(newMiningComplexity int) error {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(versionBucket))
		if v == nil {
			return errors.New("version bucket was not found during update")
		}
		bc.currentVersion++
		bc.versions[bc.currentVersion] = newMiningComplexity
		err := v.Put([]byte("v"), utils.SerializeMap(bc.versions))
		if err != nil {
			return err
		}
		err = v.Put([]byte("lv"), utils.SerializeInt(bc.currentVersion))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func NewBlockchain() *Blockchain {
	var (
		last        []byte
		lastVersion int
		versions    map[int]int
	)
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		v := tx.Bucket([]byte(versionBucket))

		if b == nil {
			genesisBlock := NewBlock("Genesis block", []byte{}, hashComplexity)
			b, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				return err
			}
			err = b.Put(genesisBlock.Header.Hash, genesisBlock.Serialize())
			if err != nil {
				return err
			}
			err = b.Put([]byte("l"), genesisBlock.Header.Hash)
			if err != nil {
				return err
			}

			last = genesisBlock.Header.Hash
		} else {
			last = b.Get([]byte("l"))
		}

		if v == nil {
			versions = make(map[int]int)
			lastVersion = 0
			versions[lastVersion] = hashComplexity
			v, err = tx.CreateBucket([]byte(versionBucket))
			if err != nil {
				return err
			}
			err = v.Put([]byte("v"), utils.SerializeMap(versions))
			if err != nil {
				return err
			}
			err = v.Put([]byte("lv"), utils.SerializeInt(lastVersion))
			if err != nil {
				return err
			}
		} else {
			versions = utils.DeserializeMap(v.Get([]byte("v")))
			lastVersion = utils.DeserializeInt(v.Get([]byte("lv")))
			if versions[lastVersion] != hashComplexity {
				lastVersion++
				versions[lastVersion] = hashComplexity
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
	bc := Blockchain{lastBlock: last, db: db, versions: versions, currentVersion: lastVersion}

	return &bc
}
