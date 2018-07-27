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
	"testing"
)

func TestDecrypt(t *testing.T) {
	sk, pk := newKey()
	msg := []byte("Hello, world!")
	ciphertext, _ := encrypt(pk, msg)
	plaintext, _ := decrypt(sk, ciphertext)
	if string(plaintext) != string(msg) {
		t.FailNow()
	}
}

func TestVerify(t *testing.T) {
	sk, pk := newKey()
	hash := hash([]byte("Hello, world!"))
	sig, _ := signature(sk, hash)
	err := verify(pk, hash, sig)
	if err != nil {
		t.FailNow()
	}
}
