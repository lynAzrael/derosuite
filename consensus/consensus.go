package consensus

import (
	"encoding/json"
	"github.com/deroproject/derosuite/blockchain"
	"github.com/deroproject/derosuite/globals"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync/atomic"
)

type RegFunc func(cfg Consensus_object, chain *blockchain.Blockchain) (Consensus, error)

var RegMap = make(map[string]RegFunc)

type Consensus interface {
	Start()
	Stop()
}

var loggerpool *log.Entry

type Consensus_object struct {
	name      string
	iinterval int64
}

func Init_Consensus(params map[string]interface{}, chain *blockchain.Blockchain) (*Consensus, error) {
	var consensus Consensus

	loggerpool = globals.Logger.WithFields(log.Fields{"com": "CONSENSUS"}) // all components must use this logger
	loggerpool.Infof("Consensus started")
	atomic.AddUint32(&globals.Subsystem_Active, 1) // increment subsystem

	// get config info
	consensus_file := filepath.Join(globals.GetDataDirectory(), "consensus.json")

	file, err := os.Open(consensus_file)
	if err != nil {
		loggerpool.Warnf("Error opening consensus data file %s err %s", consensus_file, err)
	} else {
		defer file.Close()

		var cfg Consensus_object
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&cfg)
		if err != nil {
			loggerpool.Warnf("Error unmarshalling consensus data err %s", err)
		} else { // successfully unmarshalled data, add it to consensus
			loggerpool.Debugf("Will try to init consensus：", cfg.name)
			load := Load(cfg.name)
			if load == nil {
				loggerpool.Fatalf("Consensus driver not found: %s", cfg.name)
			}

			consensus, err = load(cfg, chain)
			if err != nil {
				return nil, err
			}
		}
	}
	return &consensus, nil
}

func Reg(name string, regFunc RegFunc) {
	if regFunc == nil {
		panic("Consensus: Register driver is nil")
	}
	if _, dup := RegMap[name]; dup {
		panic("Consensus: Register called twice for driver " + name)
	}
	RegMap[name] = regFunc
}

func Load(name string) (RegFunc) {
	if reg, ok := RegMap[name]; ok {
		return reg
	}
	return nil
}

func (cfg *Consensus_object) GetName() string {
	return cfg.name
}

func (cfg *Consensus_object) GetInterval() int64 {
	return cfg.iinterval
}
