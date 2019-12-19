package main

import (
	"github.com/offer365/edda/eddaX"
	"github.com/offer365/edda/utils"
	"github.com/offer365/example/endecrypt/endeaes"
	"github.com/offer365/example/endecrypt/endeaesrsa"
	"github.com/offer365/example/endecrypt/endeaesrsaecc"
	"github.com/offer365/example/endecrypt/endersa"
)

func main()  {
	cfg:=&eddaX.Config{
		GRpcServerCrt:  server_crt,
		GRpcServerKey:  server_key,
		GRpcClientCrt:  client_crt,
		GRpcClientKey:  client_key,
		GRpcCaCrt:      ca_crt,
		GRpcUser:       grpcUser,
		GRpcPwd:        grpcPwd,
		GRpcServerName: server_name,
		GRpcListen:     "19527",
		LicenseEncrypt:             licenseEncrypt1,
		LicenseDecrypt:             licenseDecrypt1,
		SerialEncrypt:              serialEncrypt1,
		SerialDecrypt:              serialDecrypt1,
		UntiedEncrypt:              untiedEncrypt1,
		UntiedDecrypt:              untiedDecrypt1,
		TokenHash:                  HashFunc1,
	}
	eddaX.Start(cfg)
}


// odin & edda

// Pub Encrypt Rsa2048 Aes256
func licenseEncrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PubEncrypt(src, []byte(_rsa2048pub1), []byte(_aes256key1))
}
// Pri Decrypt Rsa2048 Aes256
func licenseDecrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PriDecrypt(src, []byte(_rsa2048pri1), []byte(_aes256key1))
}

// Pub Encrypt Ecc256 Rsa204 8Aes256
func licenseEncrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PubEncrypt(src, []byte(_eccpub1), []byte(_rsa2048pub1), []byte(_aes256key1))
}

// Pri Decrypt Ecc25 6Rsa2048 Aes256
func licenseDecrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PriDecrypt(src, []byte(_eccpri1), []byte(_rsa2048pri1), []byte(_aes256key1))
}

// Pub Encrypt Rsa2048 Aes256
func serialEncrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PubEncrypt(src, []byte(_rsa2048pub2), []byte(_aes256key2))
}
// Pri Decrypt Rsa2048 Aes256
func serialDecrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PriDecrypt(src, []byte(_rsa2048pri2), []byte(_aes256key2))
}

// Pub Encrypt Ecc256 Rsa204 8Aes256
func serialEncrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PubEncrypt(src, []byte(_eccpub2), []byte(_rsa2048pub2), []byte(_aes256key2))
}

// Pri Decrypt Ecc25 6Rsa2048 Aes256
func serialDecrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PriDecrypt(src, []byte(_eccpri2), []byte(_rsa2048pri2), []byte(_aes256key2))
}


// Pub Encrypt Rsa2048 Aes256
func untiedEncrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PubEncrypt(src, []byte(_rsa2048pub3), []byte(_aes256key3))
}
// Pri Decrypt Rsa2048 Aes256
func untiedDecrypt1(src []byte) ([]byte, error) {
	return endeaesrsa.PriDecrypt(src, []byte(_rsa2048pri3), []byte(_aes256key3))
}

// Pub Encrypt Ecc256 Rsa204 8Aes256
func untiedEncrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PubEncrypt(src, []byte(_eccpub3), []byte(_rsa2048pub3), []byte(_aes256key3))
}

// Pri Decrypt Ecc25 6Rsa2048 Aes256
func untiedDecrypt2(src []byte) ([]byte, error) {
	return endeaesrsaecc.PriDecrypt(src, []byte(_eccpri3), []byte(_rsa2048pri3), []byte(_aes256key3))
}

func PriDecryptRsa2048(src []byte) ([]byte, error) {
	return endersa.PriDecrypt(src, []byte(_rsa4096pri1))
}

func Aes256key1(src []byte) ([]byte, error) {
	return endeaes.AesCbcEncrypt(src, []byte(_aes256key4))
}

func Aes256key2(src []byte) ([]byte, error) {
	return endeaes.AesCbcEncrypt(src, []byte(_aes256key4))
}

func HashFunc1(src []byte) string {
	return utils.Sha256Hex(src, []byte(storeHashSalt))
}

func HashFunc2(src []byte) string {
	return utils.Sha256Hex(src, []byte(storeHashSalt))
}