package solo

import (
	"github.com/deroproject/derosuite/block"
	"github.com/deroproject/derosuite/blockchain"
	"github.com/deroproject/derosuite/blockchain/mempool"
	"github.com/deroproject/derosuite/consensus"
	"github.com/deroproject/derosuite/crypto"
	"github.com/deroproject/derosuite/globals"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

var logger *log.Entry

func init() {
	consensus.Reg("solo",NewSolo)
}

type Solo struct {
	quitCh chan bool

	chain    *blockchain.Blockchain
	enable   bool
	pending  bool
	interval int64
}

func NewSolo(cfg consensus.Consensus_object, chain *blockchain.Blockchain) (consensus.Consensus, error) {
	logger = globals.Logger.WithFields(log.Fields{"com": "CONSENSUS"})
	solo := &Solo{
		quitCh:   make(chan bool, 5),
		chain:    chain,
		enable:   false,
		pending:  false,
		interval: cfg.GetInterval(),
	}
	return solo,nil
}

func (solo *Solo) Start() {
	solo.enable = true
	go solo.eventLoop()
}

func (solo *Solo) eventLoop() {
	// solo 共识打包时间可配置
	t := time.NewTimer(time.Second * time.Duration(solo.interval)).C
	for {
		select {
		case <-solo.quitCh:
			return
		case <-t:
			solo.CreateBlock()
		}
	}
}

func (solo *Solo) Stop() {
	solo.quitCh <- true
}

func (solo *Solo) CreateBlock() {
	needSleep := true
	if !solo.enable {
		logger.Fatal("Please use start mining command to start create block.")
		return
	}
	for {
		// 是否在暂停状态
		if solo.pending {
			return
		}

		if needSleep {
			time.Sleep(time.Duration(solo.interval) * time.Second)
		}

		// txs := solo.chain.Mempool.Mempool_List_TX_SortedInfo()
		txs := makeTestTxs(10)
		if len(txs) == 0 {
			needSleep = true
			continue
		}

		needSleep = false

		//check dup

		var newblock *block.Block

		for i := range txs {
			newblock.Tx_hashes = append(newblock.Tx_hashes, txs[i].Hash)
		}

		newblock.Nonce = rand.New(globals.NewCryptoRandSource()).Uint32()

		// 更新区块
		solo.chain.Store_BL(nil, newblock)
	}
	return
}

func makeTestTxs(num int) []mempool.TX_Sorting_struct {
	var txs []mempool.TX_Sorting_struct
	for i := 0; i < num; i++ {
		tx := mempool.TX_Sorting_struct{FeesPerByte: 1, Hash: crypto.HashHexToHash(time.Now().String()), Size: 1}
		txs = append(txs, tx)
	}
	return txs
}
