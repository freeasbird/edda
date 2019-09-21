package endersa

// https://github.com/wenzhenxi/gorsa

// 公钥加密
func PublicEncrypt(data []byte, publicKey string) ([]byte, error) {

	gRsa := RSASecurity{}
	gRsa.SetPublicKey(publicKey)

	rsaData, err := gRsa.PubKeyENCTYPT(data)
	if err != nil {
		return nil, err
	}
	// base64.StdEncoding.EncodeToString(rsaData), nil
	return rsaData, nil
}

// 私钥加密
func PriKeyEncrypt(data []byte, privateKey string) ([]byte, error) {

	gRsa := RSASecurity{}
	gRsa.SetPrivateKey(privateKey)

	rsaData, err := gRsa.PriKeyENCTYPT(data)
	if err != nil {
		return nil, err
	}

	// base64.StdEncoding.EncodeToString(rsaData), nil
	return rsaData, nil
}

// 公钥解密
func PublicDecrypt(data []byte, publicKey string) ([]byte, error) {

	//dataByt, _ := base64.StdEncoding.DecodeString(data)

	gRsa := RSASecurity{}
	gRsa.SetPublicKey(publicKey)

	rsaData, err := gRsa.PubKeyDECRYPT(data)
	if err != nil {
		return nil, err
	}

	return rsaData, nil

}

// 私钥解密
func PriKeyDecrypt(data []byte, privateKey string) ([]byte, error) {

	//dataByt, _ := base64.StdEncoding.DecodeString(data)

	gRsa := RSASecurity{}
	gRsa.SetPrivateKey(privateKey)

	rsaData, err := gRsa.PriKeyDECRYPT(data)
	if err != nil {
		return nil, err
	}

	return rsaData, nil
}
