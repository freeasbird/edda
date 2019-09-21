package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Node struct {
	ID primitive.ObjectID `bson:"_id" json:"-"`
	*Attr
	*Hardware
}

type Attr struct {
	Name  string `bson:"name" json:"name"`
	IP    string `bson:"ip" json:"ip"`
	Start int64  `bson:"start" json:"start"` // 启动时间
	HwMd5 string `bson:"md5" json:"md5"`
	Now   int64  `bson:"now" json:"now"`
}

type Hardware struct {
	Host     *Host      `bson:"host" json:"host"`
	Product  *Product   `bson:"product"`
	Board    *Board     `bson:"board"`
	Chassis  *Chassis   `bson:"chassis"`
	Bios     *Bios      `bson:"bios"`
	Cpu      *Cpu       `bson:"cpu" json:"cpu"`
	Mem      *Mem       `bson:"mem" json:"mem"`
	Networks []*Network `bson:"networks" json:"networks"`
}

type Host struct {
	Machineid    string `bson:"machineid" json:"machineid"` // 设备id
	Hypervisor   string `bson:"hypervisor" json:"hypervisor"`
	Architecture string `bson:"architecture" json:"architecture"` // 架构
}

type Product struct {
	Name    string `bson:"name" json:"name"`
	Vendor  string `bson:"vendor" json:"vendor"`
	Version string `bson:"version" json:"version"`
	Serial  string `bson:"serial" json:"serial"`
}

type Board struct {
	Name     string `bson:"name" json:"name"`
	Vendor   string `bson:"vendor" json:"vendor"`
	Version  string `bson:"version" json:"version"`
	Serial   string `bson:"serial" json:"serial"`
	Assettag string `bson:"assettag" json:"assettag"`
}

type Chassis struct {
	Type     uint   `bson:"type" json:"type"`
	Vendor   string `bson:"vendor" json:"vendor"`
	Version  string `bson:"version" json:"version"`
	Serial   string `bson:"vSerial" json:"serial"`
	Assettag string `bson:"assettag" json:"assettag"`
}

type Bios struct {
	Vendor string `bson:"vendor" json:"vendor"`
}

type Cpu struct {
	Vendor  string `bson:"vendor" json:"vendor"`
	Model   string `bson:"model" json:"model"`
	Speed   uint   `bson:"speed" json:"speed"`
	Cache   uint   `bson:"cache" json:"cache"`
	Cpus    uint   `bson:"cpus" json:"cpus"`
	Cores   uint   `bson:"cores" json:"cores"`
	Threads uint   `bson:"threads" json:"threads"`
}

// 内存
type Mem struct {
	Type  string `bson:"type" json:"type"`   // type
	Speed uint   `bson:"speed" json:"speed"` // 速率
}

type Storage struct {
	Driver string `bson:"driver" json:"driver"`
	Vendor string `bson:"vendor" json:"vendor"`
	Model  string `bson:"model" json:"model"`
	Serial string `bson:"serial" json:"serial"`
}

type Network struct {
	Driver     string `bson:"driver" json:"driver"`
	Macaddress string `bson:"macaddress" json:"macaddress"`
	Speed      uint   `bson:"speed" json:"speed"`
}
