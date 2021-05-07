// Package main provides ...
package main

import (
	"io/ioutil"

	"github.com/eiblog/migrate/migrate"

	"gopkg.in/yaml.v2"
)

type fromTo struct {
	From migrate.Store `yaml:"from"`
	To   migrate.Store `yaml:"to"`
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

}
