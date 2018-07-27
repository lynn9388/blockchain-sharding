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
	"crypto/rsa"

	"crypto/sha256"

	"crypto"

	"log"
)

func newKey() (*rsa.PrivateKey, *rsa.PublicKey) {
	sk, err := rsa.GenerateKey(rand.Reader, 2018)
	if err != nil {
		log.Fatal(err)
	}

	return sk, &sk.PublicKey
}

func encrypt(pk *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pk, plaintext, []byte(""))
}

func decrypt(sk *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, sk, ciphertext, []byte(""))
}

func hash(msg []byte) []byte {
	sha := sha256.New()
	sha.Write(msg)
	return sha.Sum(nil)
}

func signature(sk *rsa.PrivateKey, hash []byte) ([]byte, error) {
	return rsa.SignPSS(rand.Reader, sk, crypto.SHA256, hash, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto})
}

func verify(pk *rsa.PublicKey, hash []byte, sig []byte) error {
	return rsa.VerifyPSS(pk, crypto.SHA256, hash, sig, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto})
}
