package main

import (
	"fmt"
	"launchpad.net/lpad"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//root, err := lpad.Login(lpad.Production, &lpad.ConsoleOAuth{})
	//check(err)
	root := lpad.Anonymous(lpad.Production)

	v, err := root.Location("/bugs/123456").Get(nil)
	check(err)
	fmt.Printf("%#v\n", v.Map())
}
