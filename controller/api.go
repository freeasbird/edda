package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/offer365/edda/logic"
	pb "github.com/offer365/edda/proto"
)

var (
	User     = "admin"
	secrets  = gin.H{}
	Accounts gin.Accounts
	salt     = []byte("build857484914")
)


// 解绑
func UntiedApi(c *gin.Context) {
	var (
		app, id string
	)

	app = c.Param("app")
	id = c.Param("id")
	req := pb.UntiedReq{
		App: app,
		Id:  id,
	}
	cipher, err := pb.Auth.Untied(context.TODO(), &req)
	if err != nil {
		c.JSON(401, map[string]string{"code": "error"})
		return
	}
	c.JSON(200, map[string]string{"code": cipher.Code})
}

// 应用
func AppAPI(c *gin.Context) {
	var (
		id         string
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
		id, err := logic.InsertApp(c.Request.Body)
		if err != nil {
			c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
		return
	case "GET":
		// one
		if id != "" {
			_id,err:=strconv.Atoi(id)
			data := logic.FindOneApp(_id)
			if err != nil {
				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
				return
			}
			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
			return
		}
		// many
		data:= logic.FindAllApp()
		if err != nil {
			c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
		return
	case "DELETE":
		_id,err:=strconv.Atoi(id)
		logic.DeleteApp(_id)
		if err != nil {
			c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
			return
		}
		c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
		return
	case "POST":
		_id,err:=strconv.Atoi(id)
		logic.UpdateApp(_id, c.Request.Body)
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

func LicenseAPI(c *gin.Context)  {
	if code,err:=logic.GenAuth(c.Request.Body);err==nil{
		c.JSON(200,gin.H{"code":200,"data":code})
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


// 序列号
// func SerialAPI(c *gin.Context) {
//	var (
//		id         string
//		collection = "projects"
//	)
//
//	user := c.MustGet(gin.AuthUserKey).(string)
//	if _, ok := secrets[user]; ok {
//		id = c.Param("id")
//		id = strings.Trim(id, "/")
//		page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
//		if err != nil {
//			page = 1
//		}
//		if page <= 0 {
//			page = 1
//		}
//		size, err := strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 64)
//		if err != nil {
//			size = 10
//		}
//		if size <= 0 || size > 100 {
//			size = 10
//		}
//
//		switch c.Request.Method {
//		case "PUT":
//			id, err := logic.InsertApp(collection, c.Request.Body)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
//			return
//		case "GET":
//			// one
//			if id != "" {
//				data, err := logic.FindOneApp(collection, id)
//				if err != nil {
//					c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//					return
//				}
//				c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
//				return
//			}
//			// many
//			data, err := logic.FindAllApp(collection, (page-1)*size, size)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
//			return
//		case "DELETE":
//			err := logic.Delete(collection, id)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
//			return
//		case "POST":
//			msg, err := logic.ResolveSerial(c.PostForm("code"))
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": msg})
//			return
//		default:
//			c.JSON(200, gin.H{
//				"code": 1,
//				"msg":  "Method error.",
//			})
//		}
//	}
// }

// func NodeAPI(c *gin.Context) {
//	var (
//		id         string
//		collection = "copyrights"
//	)
//
//	user := c.MustGet(gin.AuthUserKey).(string)
//	if _, ok := secrets[user]; ok {
//		id = c.Param("id")
//		id = strings.Trim(id, "/")
//		page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
//		if err != nil {
//			page = 1
//		}
//		if page <= 0 {
//			page = 1
//		}
//		size, err := strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 64)
//		if err != nil {
//			size = 10
//		}
//		if size <= 0 || size > 100 {
//			size = 10
//		}
//		switch c.Request.Method {
//		case "PUT":
//			id, err := logic.InsertNode(collection, c.Request.Body)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
//			return
//		case "GET":
//			// one
//			if id != "" {
//				data, err := logic.FindOneNode(collection, id)
//				if err != nil {
//					c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//					return
//				}
//				c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
//				return
//			}
//			// many
//			data, err := logic.FindNode(collection, bson.D{}, (page-1)*size, size)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
//			return
//		case "DELETE":
//			err := logic.Delete(collection, id)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
//			return
//		case "POST":
//			err := logic.Update(collection, id, c.Request.Body)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
//			return
//		default:
//			c.JSON(200, gin.H{
//				"code": 1,
//				"msg":  "Method error.",
//			})
//		}
//	}
// }

// // 客户
// func CustomerAPI(c *gin.Context) {
//	var (
//		id         string
//		collection = "customers"
//	)
//
//	user := c.MustGet(gin.AuthUserKey).(string)
//	if _, ok := secrets[user]; ok {
//		id = c.Param("id")
//		id = strings.Trim(id, "/")
//		page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
//		if err != nil {
//			page = 1
//		}
//		if page <= 0 {
//			page = 1
//		}
//		size, err := strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 64)
//		if err != nil {
//			size = 10
//		}
//		if size <= 0 || size > 100 {
//			size = 10
//		}
//		switch c.Request.Method {
//		case "PUT":
//			if user != User {
//				c.JSON(200, gin.H{"code": 404, "msg": "No permission.", "data": nil})
//				return
//			}
//			id, err := logic.InsertCustomer(collection, c.Request.Body)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
//			return
//		case "GET":
//			// one
//			if id != "" {
//				data, err := logic.FindOneCustomer(collection, id)
//				if err != nil {
//					c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//					return
//				}
//				c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
//				return
//			}
//			// many
//			data, err := logic.FindAllCustomer(collection, (page-1)*size, size)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
//			return
//		case "DELETE":
//			if user != User {
//				c.JSON(200, gin.H{"code": 404, "msg": "No permission.", "data": nil})
//				return
//			}
//			err := logic.Delete(collection, id)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
//			return
//		case "POST":
//			if user != User {
//				c.JSON(200, gin.H{"code": 404, "msg": "No permission.", "data": nil})
//				return
//			}
//			err := logic.Update(collection, id, c.Request.Body)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
//			return
//		default:
//			c.JSON(200, gin.H{
//				"code": 1,
//				"msg":  "Method error.",
//			})
//		}
//	}
// }
