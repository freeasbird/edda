package endeaesrsa

import (
	"encoding/base64"
	"github.com/offer365/endecrypt/endeaes"
	"github.com/offer365/endecrypt/endersa"
)

// 公钥加密
func PubEncrypt(src []byte, rsaKey, aesKey string) (string, error) {
	// RSA公钥加密
	code, err := endersa.PublicEncrypt(src, rsaKey)
	if err != nil {
		return err.Error(), err
	}
	// AES加密
	result, err := endeaes.AesCbcEncrypt(code, []byte(aesKey))
	if err != nil {
		return err.Error(), err
	}

	return base64.StdEncoding.EncodeToString(result), err
}

// 公钥解密
func PubDecrypt(src, rsaKey, aesKey string) (string, error) {
	byt, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return err.Error(), err
	}
	// aes 解密
	tmp, err := endeaes.AesCbcDecrypt(byt, []byte(aesKey))
	if err != nil {
		return err.Error(), err
	}
	// rsa 公钥解密
	tmp, err = endersa.PublicDecrypt(tmp, rsaKey)
	if err != nil {
		return err.Error(), err
	}
	return string(tmp), nil
}

// 私钥加密
func PriEncrypt(src []byte, rsaKey, aesKey string) (string, error) {
	str, err := endersa.PriKeyEncrypt(src, rsaKey)
	if err != nil {
		return err.Error(), err
	}
	// AES加密
	result, err := endeaes.AesCbcEncrypt([]byte(str), []byte(aesKey))
	if err != nil {
		return err.Error(), err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

// 私钥解密
func PirDecrypt(src, rsaKey, aesKey string) (string, error) {
	byt, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return err.Error(), err
	}
	// aes 解密
	tmp, err := endeaes.AesCbcDecrypt(byt, []byte(aesKey))
	if err != nil {
		return err.Error(), err
	}
	// rsa 私钥解密
	tmp, err = endersa.PriKeyDecrypt(tmp, rsaKey)
	if err != nil {
		return err.Error(), err
	}
	return string(tmp), nil
}
