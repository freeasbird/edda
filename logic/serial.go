package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/offer365/edda/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"time"
)

func FindOneSerial(coll string, id string) (instances []*model.SerialNum, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}
	instance := new(model.SerialNum)
	err = device.FindOne(coll, filter, instance)
	return []*model.SerialNum{instance}, err
}

func FindAllSerial(coll string, skip, limit int64) (instances []*model.SerialNum, err error) {
	instances = make([]*model.SerialNum, 0)
	fu := func(cursor *mongo.Cursor) (err error) {
		// 遍历结果集
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		for cursor.Next(ctx) {
			instance := new(model.SerialNum)
			if err = cursor.Decode(instance); err == nil { // 反序列化bson到对象
				instances = append(instances, instance)
			}
		}
		return
	}
	err = device.Find(coll, make(map[string]string), fu, skip, limit, 1)
	return
}

func InsertSerial(coll string, body io.Reader) (id string, err error) {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	instance := new(model.SerialNum)
	err = json.Unmarshal(byt, instance)
	if err != nil {
		return
	}
	instance.ID = primitive.NewObjectID()
	return device.Insert(coll, instance)
}

func ResolveSerial(coll string, body io.Reader) (msg []string, err error) {
	var (
		filter   bson.D
		ins      []*model.Node
		byt      []byte
		licenses []*model.License
		sn       *model.SerialNum
	)
	msg = make([]string, 0)
	ins = make([]*model.Node, 0)
	byt, err = ioutil.ReadAll(body)
	data := make(map[string]string, 1)
	err = json.Unmarshal(byt, &data)
	sn, err = model.Decrypt(data["code"])
	byt, err = json.Marshal(sn)
	for _, n := range sn.Nodes {
		// 硬件md5
		filter = bson.D{{"attr.md5", n.HwMd5}}
		nodes, _ := FindNode("nodes", filter, 0, 0)
		if len(nodes) == 0 {
			msg = append(msg, "未在数据库中找到该设备:"+n.Name+"。")
			// 设备id
			filter = bson.D{{"hardware.host.machineid", n.Host.Machineid}}
			machine, _ := FindNode("nodes", filter, 0, 0)
			if len(machine) > 0 {
				for _, mh := range machine {
					msg = append(msg, "在数据库中找到与该设备一致的设备id:"+mh.Host.Machineid+"。")
				}
			}
			nodes = append(nodes, machine...)
			// 产品序列号
			filter = bson.D{{"hardware.product.serial", n.Product.Serial}}
			serial, _ := FindNode("nodes", filter, 0, 0)
			if len(serial) > 0 {
				for _, sr := range serial {
					msg = append(msg, "在数据库中找到与该设备一致的产品序列号:"+sr.Product.Serial+"。")
				}
			}
			nodes = append(nodes, serial...)
			for _, nw := range n.Hardware.Networks {
				//  mac
				filter = bson.D{{"hardware.networks.macaddress", nw.Macaddress}}
				mac, _ := FindNode("nodes", filter, 0, 0)
				if len(mac) > 0 {
					for _, m := range mac {
						for _, ma := range m.Networks {
							msg = append(msg, "在数据库中找到与该设备一致的MAC地址:"+ma.Macaddress+"。")
						}
					}
				}
				nodes = append(nodes, mac...)
			}
		}
		ins = append(ins, nodes...)
	}
	devices := make(map[primitive.ObjectID]*model.Node)
	for _, in := range ins {
		devices[in.ID] = in
	}
	for _, dev := range devices {
		filter = bson.D{{"devices.md5", dev.Attr.HwMd5}}
		licenses, err = FindLicense("licenses", filter, 0, 0)
		licenses = append(licenses, licenses...)
	}

	licenseM := make(map[primitive.ObjectID]*model.License)
	for _, lic := range licenses {
		licenseM[lic.ID] = lic
	}
	msg = append(msg, "在数据库中找到相似的设备:"+fmt.Sprintf("%d个。", len(devices)))
	customers := make(map[string]int)
	for _, lic := range licenseM {
		key := fmt.Sprintf("{相关客户: %s,相关项目: %s}", lic.Customer, lic.Project)
		customers[key] += 1
	}
	for custom, n := range customers {
		msg = append(msg, fmt.Sprintf("%s*%d", custom, n))
	}

	return
}
