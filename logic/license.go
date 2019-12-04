package logic

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"

	pb "github.com/offer365/edda/proto"
)

type Result struct {
	SerialNum string             `json:"serial_num"`
	Apps      map[string]*pb.App `json:"apps"`
}

func GenAuth(body io.Reader) (code string,err error)  {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	result := new(Result)
	err = json.Unmarshal(byt, result)
	if err != nil {
		return
	}
	cipher:=&pb.Cipher{
		Code:                 result.SerialNum,
	}
	ar:=&pb.AuthReq{
		Cipher:               cipher,
		Apps:                 result.Apps,
	}
	authresp,err:=pb.Auth.Authorized(context.Background(),ar)
	return  authresp.Cipher.Code,err

}
