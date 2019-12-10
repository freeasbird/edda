package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/offer365/edda/logic"
	pb "github.com/offer365/edda/proto"

	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"google.golang.org/grpc"

	"os"
	"path/filepath"
	"runtime"

	"github.com/offer365/edda/asset"
	"github.com/offer365/edda/controller"

	log "github.com/sirupsen/logrus"
)

var (
	debug     bool
	AssetPath string
)

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
			log.Error("RestoreAssets error. ", err.Error())
			_ = os.RemoveAll(filepath.Join(AssetPath, dir))
			continue
		}
	}
}

func init() {
	RestoreAsset()
	log.SetFormatter(&log.JSONFormatter{})
	debug = true
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

}

var (
	gs      *grpc.Server
	secrets = gin.H{"admin": nil}
)

func main() {
	Run(logic.ListenAddr)
}

func Run(addr string) {
	var err error
	gs, err = pb.AuthGRpcServer()
	if err != nil {
		log.Fatal(err)
		return
	}
	pb.RegisterAuthorizationServer(gs, pb.Auth)
	ws := route()
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
		log.Fatal(err)
		return
	}
	err = http.Serve(listener, handle)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func NewTlsListen(crt, key, ca []byte, addr string) (net.Listener, error) {
	certificate, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		err = errors.New("failed to append ca certs")
		log.Fatal(err)
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
	api.POST("/license", controller.LicenseAPI)
	// 应用
	api.Any("/app/*id", controller.AppAPI)
	// 生成密文
	// api.GET("/cipher/:lid", controller.CipherAPI)
	api.GET("/untied/:app/:id", controller.UntiedApi)
	// 解析
	api.POST("/server/:do", controller.ServerAPI)

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


