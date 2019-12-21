package eddaX

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
)

type Result struct {
	SerialNum string          `json:"serial_num"`
	Apps      map[string]*App `json:"apps"`
}

func GenAuth(body io.Reader) (code string, err error) {
	var authresp *AuthResp
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	result := new(Result)
	err = json.Unmarshal(byt, result)
	if err != nil {
		return
	}
	cipher := &Cipher{
		Code: result.SerialNum,
	}
	ar := &AuthReq{
		Cipher: cipher,
		Apps:   result.Apps,
	}
	if authresp, err = AuthServer.Authorized(context.Background(), ar); err != nil {
		return
	}
	return authresp.Cipher.Code, err

}
