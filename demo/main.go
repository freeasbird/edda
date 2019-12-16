package main

import (
	"context"
	"flag"

	corec "github.com/offer365/example/grpc/core/client"
	"github.com/offer365/edda/eddaX"
	"google.golang.org/grpc"
)

var (
	auth      *Authentication
	_username = "C205v406x68f5IM7"
	_password = "c9bJ3v7FQ11681EP"
	cli       eddaX.AuthorizationClient
	ListenAddr  string
	Cfg *eddaX.Config
)


func args() {
	flag.StringVar(&ListenAddr, "l", ":19527", "listen addr.")
	flag.Parse()
}

func main() {
	args()
	auth = &Authentication{
		User:     _username,
		Password: _password,
	}

	Cfg=&eddaX.Config{
		GRpcServerCrt:  "",
		GRpcServerKey:  "",
		GRpcClientCrt:  "",
		GRpcClientKey:  "",
		GRpcCaCrt:      "",
		GRpcUser:       "",
		GRpcPwd:        "",
		GRpcServerName: "",
		GRpcListen:     "",
		RestfulPwd:     "",
		LicenseEncrypt: nil,
		LicenseDecrypt: nil,
		SerialEncrypt:  nil,
		SerialDecrypt:  nil,
		UntiedEncrypt:  nil,
		UntiedDecrypt:  nil,
		TokenHash:      nil,
	}
	gRpcClient()
	// cli.Authorized()
}

func gRpcClient() {
	var (
		conn *grpc.ClientConn
		err  error
	)

	conn, err = corec.NewRpcClient(
		corec.WithAddr(ListenAddr),
		corec.WithDialOption(grpc.WithPerRPCCredentials(auth)),
		corec.WithServerName(Cfg.GRpcServerName),
		corec.WithCert([]byte(Cfg.GRpcClientCrt)),
		corec.WithKey([]byte(Cfg.GRpcClientKey)),
		corec.WithCa([]byte(Cfg.GRpcCaCrt)),
	)
	if err != nil {
		return
	}
	cli = eddaX.NewAuthorizationClient(conn)
}

type Authentication struct {
	User     string
	Password string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}
