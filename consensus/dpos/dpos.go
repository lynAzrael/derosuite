package dpos

import (
	"github.com/deroproject/derosuite/blockchain"
	"github.com/deroproject/derosuite/consensus"
)

type Dpos struct {
	quitCh chan bool

	chain *blockchain.Blockchain

	enable   bool
	pending  bool
	interval int64
}

func NewDpos(cfg consensus.Consensus_object) *Dpos {
	dpos := &Dpos{
		quitCh:   make(chan bool, 5),
		enable:   false,
		pending:  false,
		interval: cfg.GetInterval(),
	}
	return dpos
}

func (dpos *Dpos) Start()  {

}

func (dpos *Dpos) Stop() {

}
