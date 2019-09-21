package endersa

import (
	"encoding/base64"
	"fmt"
	"odin.ren/endecrypt"
	"testing"
)

var (
	text = []byte("hello world.")
)

func BenchmarkPriKeyEncrypt(b *testing.B) {
	// 加密
	result, err := PriKeyEncrypt(text, endecrypt.PirvatekeyAuth)
	if err != nil {
		panic(err)
	}
	// nRmbAgLEsFSZzieUekELhA==
	fmt.Println(base64.StdEncoding.EncodeToString(result))
	// 解密
	oriData, err := PublicDecrypt(result, endecrypt.PubkeyAuth)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(oriData))
}

func BenchmarkPriKeyEncryptAuth1024(b *testing.B) {
	// 加密
	result, err := PriKeyEncrypt(text, endecrypt.PirvatekeyAuth)
	if err != nil {
		panic(err)
	}
	// nRmbAgLEsFSZzieUekELhA==
	fmt.Println(base64.StdEncoding.EncodeToString(result))
	// 解密
	oriData, err := PublicDecrypt(result, endecrypt.PubkeyAuth)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(oriData))
}

func BenchmarkPriPubDecryptAuth1024(b *testing.B) {
	// 加密
	result, err := PriKeyEncrypt(text, endecrypt.PirvatekeyAuth)
	if err != nil {
		panic(err)
	}
	// nRmbAgLEsFSZzieUekELhA==
	fmt.Println(base64.StdEncoding.EncodeToString(result))
	// 解密
	oriData, err := PublicDecrypt(result, endecrypt.PubkeyAuth)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(oriData))
}
