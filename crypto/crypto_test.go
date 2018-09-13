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

package crypto

import "testing"

func TestDecrypt(t *testing.T) {
	sk, err := NewKey()
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("lynn9388")
	ciphertext, err := Encrypt(&sk.PublicKey, msg)
	if err != nil {
		t.Error(err)
	}

	plaintext, err := Decrypt(sk, ciphertext)
	if err != nil {
		t.Error(err)
	}

	if string(plaintext) != string(msg) {
		t.Error("failed to decrypt ciphertext")
	}
}

func TestVerify(t *testing.T) {
	sk, err := NewKey()
	if err != nil {
		t.Fatal(err)
	}

	hash := Hash([]byte("lynn9388"))
	sig, err := Sign(sk, hash)
	if err != nil {
		t.Error(err)
	}

	if err = Verify(&sk.PublicKey, hash, sig); err != nil {
		t.Error(err)
	}
}
