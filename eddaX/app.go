package eddaX

import (
	"encoding/json"
	"io"
	"io/ioutil"

)

type APP struct {
	ID int `bson:"id" json:"id"`
	App
}

var apps []*APP

func FindOneApp(id int) (instance *APP) {
	if id < len(apps) {
		return apps[id]
	}
	return
}

func FindAllApp() (instances []*APP) {
	return apps
}

func InsertApp(body io.Reader) (id string, err error) {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	instance := new(APP)
	err = json.Unmarshal(byt, instance)
	if err != nil {
		return
	}
	instance.ID=len(apps)
	apps = append(apps, instance)
	return
}

func DeleteApp(i int)  {
	if i < len(apps){
		apps = append(apps[:i], apps[i+1:]...)
	}
}

func UpdateApp(i int,body io.Reader)  {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	instance := new(APP)
	err = json.Unmarshal(byt, instance)
	if err != nil {
		return
	}
	apps[i]=instance
	return
}