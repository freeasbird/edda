package logic

import (
	"context"
	"flag"

	pb "github.com/offer365/edda/proto"
	corec "github.com/offer365/example/grpc/core/client"
	"google.golang.org/grpc"
)


var (
	auth      *Authentication
	_username = "C205v406x68f5IM7"
	_password = "c9bJ3v7FQ11681EP"
	cli       pb.AuthorizationClient
	ListenAddr  string
)

func args() {
	flag.StringVar(&ListenAddr, "l", ":19527", "listen addr.")
	flag.Parse()
}

func init() {
	args()
	auth = &Authentication{
		User:     _username,
		Password: _password,
	}
	gRpcClient()
}

func gRpcClient() {
	var (
		conn *grpc.ClientConn
		err  error
	)

	conn, err = corec.NewRpcClient(
		corec.WithAddr(ListenAddr),
		corec.WithDialOption(grpc.WithPerRPCCredentials(auth)),
		corec.WithServerName("server.io"),
		corec.WithCert([]byte(pb.Client_crt)),
		corec.WithKey([]byte(pb.Client_key)),
		corec.WithCa([]byte(pb.Ca_crt)),
	)
	if err != nil {
		return
	}
	cli = pb.NewAuthorizationClient(conn)
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
