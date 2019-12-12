package proto

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/offer365/example/endecrypt"
	cores "github.com/offer365/example/grpc/core/server"
	"github.com/offer365/example/tools"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)


var (
	salt      = []byte(hashSalt)
	Auth      AuthorizationServer
)

func init() {
	Auth = NewAuthServer()
}

func NewAuthServer() AuthorizationServer {
	return &auth{}
}

type auth struct {}

type untied struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Date  int64  `json:"date"`
}

// 解析序列号
func (a *auth) Resolved(ctx context.Context, cipher *Cipher) (sn *SerialNum, err error) {
	var (
		byt []byte
	)
	if byt, err = base64.StdEncoding.DecodeString(cipher.Code); err != nil {
		return
	}
	// 私钥解密
	if byt, err = endecrypt.Decrypt(endecrypt.Pri1AesRsa2048, byt); err != nil {
		return
	}
	if byt == nil || len(byt) == 0 {
		err = errors.New("decrypt error ")
		return
	}
	sn = new(SerialNum)
	sn.Nodes = make(map[string]*Node, 0)

	if err = json.Unmarshal(byt, sn); err != nil {
		return
	}
	now := time.Now().Unix()
	if len(sn.Nodes) == 0 || sn.Sid == "" || sn.Date > now || (now-sn.Date) > 60*60*24 {
		err = errors.New("serial time error ")
		return
	}
	return
}

// 生成授全信息
func (a *auth) Authorized(ctx context.Context, req *AuthReq) (resp *AuthResp, err error) {
	var (
		sn  *SerialNum
		lic *License
	)
	if len(req.Apps) == 0 {
		err = errors.New("app length error ")
		return
	}
	if sn, err = a.Resolved(ctx, req.Cipher); err != nil {
		return
	}
	lic = new(License)
	lic.Lid = uuid.Must(uuid.NewV4()).String()
	lic.Sid = sn.Sid
	lic.Devices = make(map[string]string)
	for name, node := range sn.Nodes {
		lic.Devices[name] = node.Attrs.Hwmd5
	}
	lic.Apps = make(map[string]*App, 0)
	for _, app := range req.Apps {
		lic.Apps[app.Key] = app
	}
	lic.Update = time.Now().Unix()
	lic.Generate = time.Now().Unix()
	resp = new(AuthResp)
	resp.Cipher = new(Cipher)
	if resp.Cipher.Code, err = a.lic2Str(lic); err != nil {
		return
	}
	resp.Lic = lic
	return
}

// 生成解绑码
func (a *auth) Untied(ctx context.Context, req *UntiedReq) (cipher *Cipher, err error) {
	var (
		byt []byte
	)
	if req.App == "" || req.Id == "" {
		err = errors.New("app or id error ")
		return
	}
	untie := &untied{
		Key:   tools.Sha256sum([]byte(req.App), salt),
		Value: tools.Sha256sum([]byte(req.Id), salt),
		Date:  time.Now().Unix(),
	}
	if byt, err = json.Marshal(untie); err != nil {
		return
	}
	if byt, err = endecrypt.Encrypt(endecrypt.Pri1AesRsa2048, byt); err != nil {
		return
	}
	cipher = &Cipher{Code: base64.StdEncoding.EncodeToString(byt)}
	return
}

// 解析注销码
func (a *auth) Cleared(ctx context.Context, cipher *Cipher) (clear *Clear, err error) {
	var (
		byt []byte
	)
	if byt, err = base64.StdEncoding.DecodeString(cipher.Code); err != nil {
		return
	}

	clear = new(Clear)
	if err = json.Unmarshal(byt, clear); err != nil {
		return
	}
	now := time.Now().Unix()
	if clear.Date > now || (now-clear.Date) > 60*60*24 {
		err = errors.New("clear time error")
		return
	}
	if len(clear.Lic.Apps) != 0 {
		err = errors.New("clear apps error")
		return
	}
	lic, err := a.str2lic(clear.Cipher.Code)
	if err != nil || lic.Lid != clear.Lic.Lid || lic.Sid != clear.Lic.Sid || lic.Generate != clear.Lic.Generate {
		err = errors.New("clear license error")
		return
	}
	return
}

// 反序列化license
func (a *auth) str2lic(cipher string) (lic *License, err error) {
	var (
		byt []byte
	)
	if byt, err = base64.StdEncoding.DecodeString(cipher); err != nil {
		return
	}
	if byt == nil || len(byt) == 0 {
		return
	}
	lic = new(License)
	if byt, err = endecrypt.Decrypt(endecrypt.Pub1AesRsa2048, byt); err != nil {
		return
	}
	if err = json.Unmarshal(byt, lic); err != nil {
		return
	}
	return
}

func (a *auth) lic2Str(lic interface{}) (cipher string, err error) {
	var (
		byt []byte
	)
	// 生成密文
	if byt, err = json.Marshal(lic); err != nil {
		return
	}
	// 私钥加密
	if byt, err = endecrypt.Encrypt(endecrypt.Pri1AesRsa2048, byt); err != nil {
		return
	}
	return base64.StdEncoding.EncodeToString(byt), err
}

func AuthGRpcServer() (*grpc.Server, error) {
	// Token认证
	auth := func(ctx context.Context) error {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Errorf(codes.Unauthenticated, "missing credentials")
		}

		var user string
		var pwd string

		if val, ok := md["user"]; ok {
			user = val[0]
		}
		if val, ok := md["password"]; ok {
			pwd = val[0]
		}

		if user != _username || pwd != _password {
			return status.Errorf(codes.Unauthenticated, "invalid token")
		}

		return nil
	}

	// 一元拦截器
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}

	// 实例化grpc Server
	return cores.NewRpcServer(
		cores.WithServerOption(grpc.UnaryInterceptor(interceptor)),
		cores.WithCert([]byte(Server_crt)),
		cores.WithKey([]byte(Server_key)),
		cores.WithCa([]byte(Ca_crt)),
	)
}

// Authentication 自定义认证
// 要实现对每个gRPC方法进行认证，需要实现grpc.PerRPCCredentials接口
// type Authentication struct {
//	User     string
//	Password string
// }
//
// func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
//	return map[string]string{"user": a.User, "password": a.Password}, nil
// }
// func (a *Authentication) RequireTransportSecurity() bool {
//	return true
// }
