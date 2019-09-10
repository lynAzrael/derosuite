package dpos

import (
	"github.com/deroproject/derosuite/blockchain"
	"github.com/deroproject/derosuite/consensus"
	"github.com/deroproject/derosuite/structures"
)

func init() {
	consensus.Reg("dpos", NewDpos)
}

type Dpos struct {
	quitCh chan bool

	chain *blockchain.Blockchain

	enable   bool
	pending  bool
	interval int64
}

func NewDpos(cfg structures.Consensus_object, chain *blockchain.Blockchain) (consensus.Consensus, error) {
	dpos := &Dpos{
		quitCh:   make(chan bool, 5),
		chain:    chain,
		enable:   false,
		pending:  false,
		interval: cfg.GetInterval(),
	}
	return dpos, nil
}

func (dpos *Dpos) Start() {

}

func (dpos *Dpos) Stop() {

}
