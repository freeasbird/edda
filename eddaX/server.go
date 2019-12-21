package eddaX

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"

	"runtime"
	"strings"

	"github.com/offer365/edda/asset"
	"github.com/offer365/edda/log"
)

var (
	AssetPath  string
	User       = "admin"
	gs         *grpc.Server
	ListenAddr string
)

func args() {
	flag.StringVar(&ListenAddr, "l", ":19527", "listen addr.")
	flag.Parse()
}

func init() {
	args()
}

func Main() {
	RestoreAsset()
	Run(ListenAddr)
}

func Run(addr string) {
	var err error
	gs, err = AuthGRpcServer()
	if err != nil {
		log.Sugar.Fatal(err)
		return
	}
	RegisterAuthorizationServer(gs, AuthServer)
	ws := route()
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			gs.ServeHTTP(w, r) // grpc server
		} else {
			ws.ServeHTTP(w, r) // gin web server
		}
		return
	})
	listener, err := NewTlsListen([]byte(Cfg.GRpcServerCrt), []byte(Cfg.GRpcServerKey), []byte(Cfg.GRpcCaCrt), addr)
	if err != nil {
		log.Sugar.Fatal(err)
		return
	}
	err = http.Serve(listener, handle)
	if err != nil {
		log.Sugar.Fatal(err)
		return
	}
}

func NewTlsListen(crt, key, ca []byte, addr string) (net.Listener, error) {
	certificate, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Sugar.Fatal(err)
		return nil, err
	}
	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		err = errors.New("failed to append ca certs")
		log.Sugar.Fatal(err)
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
func route() http.Handler {
	gin.SetMode(gin.ReleaseMode) // 生产模式
	r := gin.New()
	r.Use(gin.Recovery()) // Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	r.LoadHTMLGlob(AssetPath + "views/*")

	// api 路由组
	api := r.Group("/edda/api/v1")

	// 授权码
	api.POST("/license", LicenseAPI)
	// 应用
	api.Any("/app/*id", AppAPI)
	// 生成密文
	// api.GET("/cipher/:lid", CipherAPI)
	api.GET("/untied/:app/:id", UntiedApi)
	// 解析
	api.POST("/server/:do", ServerAPI)

	// r.Use(SimpleSession)
	r.Static("/static", AssetPath+"static")
	r.Any("", func(c *gin.Context) {
		c.Request.URL.Path = "/index"
		r.HandleContext(c)
	})

	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "首页",
		})
	})

	r.StaticFile("/favicon.ico", AssetPath+"static/favicon.ico")
	return r
}

// 释放静态资源
func RestoreAsset() {
	// 解压 静态文件的位置
	if runtime.GOOS == "linux" {
		AssetPath = "/usr/share/.asset/.temp/"
	} else {
		AssetPath = "./"
	}
	// go get -u github.com/jteeuwen/go-bindata/...
	// 重新生成静态资源在项目的根目录下 go-bindata -o=asset/asset.go -pkg=asset views/... static/...
	dirs := []string{"views", "static"}
	for _, dir := range dirs {
		if err := asset.RestoreAssets(AssetPath, dir); err != nil {
			log.Sugar.Error("RestoreAssets error. ", err.Error())
			_ = os.RemoveAll(filepath.Join(AssetPath, dir))
			continue
		}
	}
}
