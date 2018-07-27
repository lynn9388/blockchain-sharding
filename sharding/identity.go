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

package sharding

import (
	"crypto/rand"
	"net"

	"crypto/rsa"

	"errors"

	"github.com/bwesterb/go-pow"
)

type identity struct {
	request string
	data    []byte
	proof   string
	id      []byte
}

const (
	idByteLength = 1
)

func newIdentity(epochRandom []byte, ip net.IP, pk *rsa.PublicKey) *identity {
	var i identity

	random := make([]byte, 32)
	rand.Read(random)
	i.request = pow.NewRequest(getDifficulty(), random)

	i.data = epochRandom
	i.data = append(i.data, []byte(ip)...)
	i.data = append(i.data, pk.N.String()...)

	i.proof, _ = pow.Fulfil(i.request, i.data)

	i.id = hash(append(i.data, i.proof...))
	i.id = i.id[len(i.id)-idByteLength:]
	return &i
}

func checkIdentity(i *identity) error {
	ok, _ := pow.Check(i.request, i.proof, i.data)
	if !ok {
		return errors.New("proof is not fulfilled")
	}

	hash := hash(append(i.data, i.proof...))
	if string(i.id) != string(hash[len(hash)-idByteLength:]) {
		return errors.New("id is not fulfilled")
	}
	return nil
}

func getDifficulty() uint32 {
	return 5
}
