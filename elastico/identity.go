/*
 * Copyright © 2018 Lynn <lynn9388@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package elastico

import (
	"bytes"
	"log"
	"math"
	"strconv"

	"github.com/lynn9388/pox/pow"
)

const (
	random     = "lynn9388"
	difficulty = 5
	shardNum   = 4
)

// NewIDProof returns a new proof for identity with PoW.
func NewIDProof(addr string, pk []byte) *IDProof {
	var data bytes.Buffer
	data.WriteString(random)
	data.WriteString(addr)
	data.Write(pk)
	nonce := pow.GetNonce(data.Bytes(), difficulty)

	return &IDProof{
		Addr:  addr,
		PK:    pk,
		Nonce: nonce,
	}
}

// toByte converts the identity proof to slice of byte without nonce field.
func (p *IDProof) toByte() []byte {
	var data bytes.Buffer
	data.WriteString(random)
	data.WriteString(p.Addr)
	data.Write(p.PK)
	return data.Bytes()
}

// Verify verifies if the identity proof is correct.
func (p *IDProof) Verify() bool {
	return pow.Fulfill(p.toByte(), p.Nonce, difficulty)
}

// GetCommitteeNo generates the committee number based on the identity proof.
func (p *IDProof) GetCommitteeNo() int {
	hash := pow.Hash(p.toByte(), p.Nonce)
	bitNum := int(math.Ceil(math.Log2(shardNum)))
	lastBytes := hash[len(hash)-bitNum:]
	value, err := strconv.ParseUint(lastBytes, 16, 4*bitNum)
	if err != nil {
		log.Fatal(err)
	}
	return int(value) % shardNum
}
