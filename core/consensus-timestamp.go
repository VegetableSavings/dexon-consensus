// Copyright 2018 The dexon-consensus-core Authors
// This file is part of the dexon-consensus-core library.
//
// The dexon-consensus-core library is free software: you can redistribute it
// and/or modify it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// The dexon-consensus-core library is distributed in the hope that it will be
// useful, but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser
// General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the dexon-consensus-core library. If not, see
// <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	"time"

	"github.com/dexon-foundation/dexon-consensus-core/core/types"
)

// consensusTimestamp is for Concensus Timestamp Algorithm.
type consensusTimestamp struct {
	chainTimestamps []time.Time
}

var (
	// ErrTimestampNotIncrease would be reported if the timestamp is not strickly
	// increasing on the same chain.
	ErrTimestampNotIncrease = errors.New("timestamp is not increasing")
)

// newConsensusTimestamp create timestamper object.
func newConsensusTimestamp() *consensusTimestamp {
	return &consensusTimestamp{}
}

// ProcessBlocks is the entry function.
func (ct *consensusTimestamp) processBlocks(blocks []*types.Block) (err error) {
	for _, block := range blocks {
		if !block.IsGenesis() {
			block.Notary.Timestamp, err = getMedianTime(ct.chainTimestamps)
			if err != nil {
				return
			}
		}

		for uint32(len(ct.chainTimestamps)) <= block.Position.ChainID {
			ct.chainTimestamps = append(ct.chainTimestamps, time.Time{})
		}

		if !block.Timestamp.After(ct.chainTimestamps[block.Position.ChainID]) {
			return ErrTimestampNotIncrease
		}

		ct.chainTimestamps[block.Position.ChainID] = block.Timestamp
	}
	return
}
