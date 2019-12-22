package eddaX

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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


func UntiedApi(c *gin.Context) {
	var (
		app, id string
	)

	app = c.Param("app")
	id = c.Param("id")
	req := UntiedReq{
		App: app,
		Id:  id,
	}
	cipher, err := AuthServer.Untied(context.TODO(), &req)
	if err != nil {
		c.JSON(401, map[string]string{"code": "error"})
		return
	}
	c.JSON(200, map[string]string{"code": cipher.Code})
}

// 应用
func AppAPI(c *gin.Context) {
	var (
		id string
	)

	id = c.Param("id")
	id = strings.Trim(id, "/")
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		page = 1
	}
	if page <= 0 {
		page = 1
	}
	size, err := strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 64)
	if err != nil {
		size = 10
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	switch c.Request.Method {
	case "PUT":
		id, err := InsertApp(c.Request.Body)
		if err != nil {
			c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
		return
	case "GET":
		// one
		if id != "" {
			_id, err := strconv.Atoi(id)
			data := FindOneApp(_id)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
			return
		}
		// many
		data := FindAllApp()
		if err != nil {
			c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
		return
	case "DELETE":
		_id, err := strconv.Atoi(id)
		DeleteApp(_id)
		if err != nil {
			c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
		return
	case "POST":
		_id, err := strconv.Atoi(id)
		UpdateApp(_id, c.Request.Body)
		if err != nil {
			c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
		return
	default:
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Method error.",
		})
	}

}

func LicenseAPI(c *gin.Context) {
	if code, err := GenAuth(c.Request.Body); err == nil {
		c.JSON(200, gin.H{"code": 200, "data": code})
	}
}

func ServerAPI(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {

		return
	}
	ctx := context.TODO()
	switch c.Param("do") {
	case "resolved":
		req := new(Cipher)
		err = json.Unmarshal(data, req)
		if err != nil {
			return
		}

		resp, err := AuthServer.Resolved(ctx, req)
		c.JSON(200, gin.H{"serial": resp, "msg": err})
		return
	case "authorized":
		req := new(AuthReq)
		err = json.Unmarshal(data, req)
		if err != nil {
			return
		}

		resp, err := AuthServer.Authorized(ctx, req)
		c.JSON(200, gin.H{"auth": resp, "msg": err})
		return
	case "untied":
		req := new(UntiedReq)
		err = json.Unmarshal(data, req)
		if err != nil {
			return
		}

		resp, err := AuthServer.Untied(ctx, req)
		c.JSON(200, gin.H{"cipher": resp, "msg": err})
		return
	case "cleared":
		req := new(Cipher)
		err = json.Unmarshal(data, req)
		if err != nil {
			return
		}

		resp, err := AuthServer.Cleared(ctx, req)
		c.JSON(200, gin.H{"clear": resp, "msg": err})
		return
	default:
		c.JSON(404, nil)
	}
}