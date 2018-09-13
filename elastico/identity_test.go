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
	"testing"

	"github.com/lynn9388/pox/pow"

	"github.com/lynn9388/blockchain-sharding/crypto"
)

func TestIDProof_Verify(t *testing.T) {
	sk, err := crypto.NewKey()
	if err != nil {
		t.Fatal(err)
	}

	proof := NewIDProof("localhost:9388", sk.D.Bytes())
	if proof.Verify() != true {
		t.Fail()
	}

	proof.Addr = "localhost:9488"
	if proof.Verify() == true {
		t.Fail()
	}
}

func TestIDProof_GetCommitteeNo(t *testing.T) {
	sk, err := crypto.NewKey()
	if err != nil {
		t.Fatal(err)
	}

	proof := NewIDProof("localhost:9388", sk.D.Bytes())
	hash := pow.Hash(proof.toByte(), proof.Nonce)
	no := proof.GetCommitteeNo()
	t.Logf("committee number: %v -> %v", hash, no)
}
