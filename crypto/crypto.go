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

// Package crypto provides some cryptography help functions.
package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// NewKey returns a new RSA private key.
func NewKey() (*rsa.PrivateKey, error) {
	sk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return sk, nil
}

// Encrypt encrypts the given plaintext with RSA public key.
func Encrypt(pk *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pk, plaintext, []byte(""))
}

// Decrypt decrypts the given ciphertext with RSA private key.
func Decrypt(sk *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, sk, ciphertext, []byte(""))
}

// Hash returns the hash value of the given message.
func Hash(msg []byte) []byte {
	hash := sha256.Sum256(msg)
	return hash[:]
}

// Sign returns the signature of the data with RSA private key.
func Sign(sk *rsa.PrivateKey, hash []byte) ([]byte, error) {
	return rsa.SignPSS(rand.Reader, sk, crypto.SHA256, hash, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto})
}

// Verify verifies the signature is correct or the data is not modified.
func Verify(pk *rsa.PublicKey, hash []byte, sig []byte) error {
	return rsa.VerifyPSS(pk, crypto.SHA256, hash, sig, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto})
}
