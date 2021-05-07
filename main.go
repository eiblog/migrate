// Package main provides ...
package main

import (
	"io/ioutil"

	"github.com/eiblog/migrate/store"
	v1 "github.com/eiblog/migrate/v1"

	"gopkg.in/yaml.v2"
)

type fromTo struct {
	From store.Store `yaml:"from"`
	To   store.Store `yaml:"to"`
}

func main() {
	data, err := ioutil.ReadFile("./app.yml")
	if err != nil {
		panic(err)
	}
	var info fromTo
	err = yaml.Unmarshal(data, &info)
	if err != nil {
		panic(err)
	}
	// 判断条件
	// v2 -> v1
	if info.From.Version > info.To.Version {
		panic("can not downgrade version")
	}
	if info.From.Version == "v1" &&
		info.From.Driver != "mongodb" {
		panic("v1 nonsupport driver :" + info.From.Driver)
	}
	// find from
	from := store.Migrate[info.From.Driver]
	if from == nil {
		panic("unsupported driver: " + info.From.Driver)
	}
	// find to
	to := store.Migrate[info.To.Driver]
	if to == nil {
		panic("unsupported driver: " + info.From.Driver)
	}
	// migrate data
	blog, err := from.LoadEiBlog(info.From)
	if err != nil {
		panic(err)
	}
	switch {
	case info.From.Version == "v1" && info.To.Version == "v1":
		err = to.StoreEiBlog(info.To, blog)
	case info.From.Version == "v1" && info.To.Version == "v2":
		data := store.V1ToV2(blog.(*v1.EiBlog))
		err = to.StoreEiBlog(info.To, data)
	case info.From.Version == "v2" && info.To.Version == "v2":
		err = to.StoreEiBlog(info.To, blog)
	}
	if err != nil {
		panic(err)
	}
}
