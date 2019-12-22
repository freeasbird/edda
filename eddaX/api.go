package eddaX

// 解绑


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
//			id, err := InsertApp(collection, c.Request.Body)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
//			return
//		case "GET":
//			// one
//			if id != "" {
//				data, err := FindOneApp(collection, id)
//				if err != nil {
//					c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//					return
//				}
//				c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
//				return
//			}
//			// many
//			data, err := FindAllApp(collection, (page-1)*size, size)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
//			return
//		case "DELETE":
//			err := Delete(collection, id)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
//			return
//		case "POST":
//			msg, err := ResolveSerial(c.PostForm("code"))
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
//			id, err := InsertNode(collection, c.Request.Body)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
//			return
//		case "GET":
//			// one
//			if id != "" {
//				data, err := FindOneNode(collection, id)
//				if err != nil {
//					c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//					return
//				}
//				c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
//				return
//			}
//			// many
//			data, err := FindNode(collection, bson.D{}, (page-1)*size, size)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
//			return
//		case "DELETE":
//			err := Delete(collection, id)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": nil})
//			return
//		case "POST":
//			err := Update(collection, id, c.Request.Body)
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
//			id, err := InsertCustomer(collection, c.Request.Body)
//			if err != nil {
//				c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//				return
//			}
//			c.JSON(200, gin.H{"code": 200, "msg": "success", "data": id})
//			return
//		case "GET":
//			// one
//			if id != "" {
//				data, err := FindOneCustomer(collection, id)
//				if err != nil {
//					c.JSON(200, gin.H{"code": 404, "msg": err.Error(), "data": nil})
//					return
//				}
//				c.JSON(200, gin.H{"code": 200, "msg": "success", "data": data})
//				return
//			}
//			// many
//			data, err := FindAllCustomer(collection, (page-1)*size, size)
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
//			err := Delete(collection, id)
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
//			err := Update(collection, id, c.Request.Body)
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
