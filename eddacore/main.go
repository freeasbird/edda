package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	pb "github.com/offer365/edda/eddacore/proto"
	"go.etcd.io/etcd/pkg/logutil"
	"go.uber.org/zap"

	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	sugar   *zap.SugaredLogger
	gs      *grpc.Server
	secrets = gin.H{"admin": nil}
	addr    string
)

func args() {
	flag.StringVar(&addr, "a", ":19527", "listen addr.")
	flag.Parse()
}

func init() {
	args()
	lg, _ := zap.NewProduction()
	defer lg.Sync()
	cfg := logutil.DefaultZapLoggerConfig
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	lg, _ = cfg.Build()
	sugar = lg.Sugar()
}

func main() {
	Run(addr)
}

func Run(addr string) {
	var err error
	gs, err = pb.AuthGRpcServer()
	if err != nil {
		sugar.Fatal(err)
		return
	}
	pb.RegisterAuthorizationServer(gs, pb.Auth)
	ws := ginServer()
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			gs.ServeHTTP(w, r) // grpc server
		} else {
			ws.ServeHTTP(w, r) // gin web server
		}
		return
	})
	listener, err := NewTlsListen([]byte(pb.Server_crt), []byte(pb.Server_key), []byte(pb.Ca_crt), addr)
	if err != nil {
		sugar.Fatal(err)
		return
	}
	err = http.Serve(listener, handle)
	if err != nil {
		sugar.Fatal(err)
		return
	}
}

func NewTlsListen(crt, key, ca []byte, addr string) (net.Listener, error) {
	certificate, err := tls.X509KeyPair(crt, key)
	if err != nil {
		sugar.Fatal(err)
		return nil, err
	}
	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		err = errors.New("failed to append ca certs")
		sugar.Fatal(err)
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		ClientAuth:         tls.NoClientCert, // NOTE: 这是可选的!
		ClientCAs:          certPool,
		InsecureSkipVerify: true,
		Rand:               rand.Reader,
		Time:               time.Now,
		NextProtos:         []string{"http/1.1", http2.NextProtoTLS},
	}
	return tls.Listen("tcp", addr, tlsConfig)
}

// gin 路由
func ginServer() http.Handler {
	gin.SetMode(gin.ReleaseMode) // 生产模式
	r := gin.New()
	r.Use(gin.Recovery()) //Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	api := r.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		"admin": "66666",
	}))

	// 解析
	api.POST("/server/:do", handler)
	return r
}

func handler(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {

			return
		}
		ctx := context.TODO()
		switch c.Param("do") {
		case "resolved":
			req := new(pb.Cipher)
			err = json.Unmarshal(data, req)
			if err != nil {
				return
			}

			resp, err := pb.Auth.Resolved(ctx, req)
			c.JSON(200, gin.H{"serial": resp, "msg": err})
			return
		case "authorized":
			req := new(pb.AuthReq)
			err = json.Unmarshal(data, req)
			if err != nil {
				return
			}

			resp, err := pb.Auth.Authorized(ctx, req)
			c.JSON(200, gin.H{"auth": resp, "msg": err})
			return
		case "untied":
			req := new(pb.UntiedReq)
			err = json.Unmarshal(data, req)
			if err != nil {
				return
			}

			resp, err := pb.Auth.Untied(ctx, req)
			c.JSON(200, gin.H{"cipher": resp, "msg": err})
			return
		case "cleared":
			req := new(pb.Cipher)
			err = json.Unmarshal(data, req)
			if err != nil {
				return
			}

			resp, err := pb.Auth.Cleared(ctx, req)
			c.JSON(200, gin.H{"clear": resp, "msg": err})
			return
		default:
			c.JSON(404, nil)
		}
	}
}
