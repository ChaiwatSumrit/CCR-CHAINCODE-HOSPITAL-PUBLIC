package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"os"
)

func GetPrivateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err.Error)
	}
	return privateKey
}

func GetPublicKeyFromPrivateKey(privateKey *rsa.PrivateKey) rsa.PublicKey {
	publicKey := privateKey.PublicKey
	return publicKey
}
func SaveByKey(fileName string, key interface{}) {

	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)

	}
	defer outFile.Close()
	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(2)

	}

}

func SaveByPEMKey(fileName string, key *rsa.PrivateKey) {

	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)

	}
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(2)

	}
}

func SaveByPublicPEMKey(fileName string, pubKey rsa.PublicKey) {

	asn1Bytes, err := asn1.Marshal(pubKey)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)

	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(2)

	}
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(3)

	}
}

func Sign(hash [32]byte, privateKey *rsa.PrivateKey) (s []byte, err error) {

	s, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return
	}
	return s, nil
}

func Verify(publicKey rsa.PublicKey, msg, s []byte) (err error) {

	hashv := sha256.Sum256(msg)
	return rsa.VerifyPKCS1v15(&publicKey, crypto.SHA256, hashv[:], s)
}
