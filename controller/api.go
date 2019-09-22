package controller

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/offer365/edda/config"
	"github.com/offer365/edda/logic"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
)

var (
	User     = "admin"
	secrets  = gin.H{}
	Accounts gin.Accounts
)

func Secrets() {
	for user, pwd := range config.CFG.Users {
		secrets[user] = pwd
	}
}

func CountAPI(c *gin.Context) {
	var (
		collection string
	)
	collection = c.Param("coll")
	collection = strings.Trim(collection, "/")
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		num, err := logic.Count(collection)
		if err != nil {
			c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "count": num})
			return
		}
		c.JSON(200, gin.H{"code": 200, "msg": "success", "count": num})
		return
	}
}

// 客户
func CustomerAPI(c *gin.Context) {
	var (
		id         string
		collection = "customers"
	)

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
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
			if user != User {
				c.JSON(200, gin.H{"code": 404, "msg": "No permission.", "data": nil})
				return
			}
			id, err := logic.InsertCustomer(collection, c.Request.Body)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
			return
		case "GET":
			// one
			if id != "" {
				data, err := logic.FindOneCustomer(collection, id)
				if err != nil {
					c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
					return
				}
				c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
				return
			}
			// many
			data, err := logic.FindAllCustomer(collection, (page-1)*size, size)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
			return
		case "DELETE":
			if user != User {
				c.JSON(200, gin.H{"code": 404, "msg": "No permission.", "data": nil})
				return
			}
			err := logic.Delete(collection, id)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
			return
		case "POST":
			if user != User {
				c.JSON(200, gin.H{"code": 404, "msg": "No permission.", "data": nil})
				return
			}
			err := logic.Update(collection, id, c.Request.Body)
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
}

// 应用
func AppAPI(c *gin.Context) {
	var (
		id         string
		collection = "apps"
	)

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
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
			id, err := logic.InsertApp(collection, c.Request.Body)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
			return
		case "GET":
			// one
			if id != "" {
				data, err := logic.FindOneApp(collection, id)
				if err != nil {
					c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
					return
				}
				c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
				return
			}
			// many
			data, err := logic.FindAllApp(collection, (page-1)*size, size)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
			return
		case "DELETE":
			err := logic.Delete(collection, id)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
			return
		case "POST":
			err := logic.Update(collection, id, c.Request.Body)
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

}

// 序列号
func SerialAPI(c *gin.Context) {
	var (
		id         string
		collection = "projects"
	)

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
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
			id, err := logic.InsertApp(collection, c.Request.Body)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
			return
		case "GET":
			// one
			if id != "" {
				data, err := logic.FindOneApp(collection, id)
				if err != nil {
					c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
					return
				}
				c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
				return
			}
			// many
			data, err := logic.FindAllApp(collection, (page-1)*size, size)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
			return
		case "DELETE":
			err := logic.Delete(collection, id)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
			return
		case "POST":
			msg, err := logic.ResolveSerial(collection, c.Request.Body)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": msg})
			return
		default:
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "Method error.",
			})
		}
	}
}

// license
func LicenseAPI(c *gin.Context) {
	var (
		id         string
		collection = "licenses"
	)

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
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
			id, err := logic.InsertLicense(collection, c.Request.Body)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
			return
		case "GET":
			// one
			if id != "" {
				data, err := logic.FindOneLicense(collection, id)
				if err != nil {
					c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
					return
				}
				c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
				return
			}
			// many
			data, err := logic.FindLicense(collection, bson.D{}, (page-1)*size, size)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
			return
		case "DELETE":
			err := logic.Delete(collection, id)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
			return
		case "POST":
			//err := logic.UpdateProduct(collection, id, c.Request.Body)
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
}

func NodeAPI(c *gin.Context) {
	var (
		id         string
		collection = "copyrights"
	)

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
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
			id, err := logic.InsertNode(collection, c.Request.Body)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
			return
		case "GET":
			// one
			if id != "" {
				data, err := logic.FindOneNode(collection, id)
				if err != nil {
					c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
					return
				}
				c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
				return
			}
			// many
			data, err := logic.FindNode(collection, bson.D{}, (page-1)*size, size)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
			return
		case "DELETE":
			err := logic.Delete(collection, id)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
			return
		case "POST":
			err := logic.Update(collection, id, c.Request.Body)
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
}

// web登录Api
func LoginAPI(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		c.JSON(200, gin.H{"cookie": base64.StdEncoding.EncodeToString([]byte(user))})
	}
}

func HelpAPI(c *gin.Context) {
	c.JSON(200, gin.H{"code": 200, "msg": []string{
		"获取二维码图片：",
		"curl -X GET http://127.0.0.1:9999/odin/api/v1/server/qr-code",
		"其他问题请联系技术人员。",
	}})
}
