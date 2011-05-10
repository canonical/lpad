package main

import (
	"fmt"
	"launchpad.net/lpad"
	"os"
)

func check(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	root, err := lpad.Login(lpad.Production, &lpad.ConsoleOAuth{})
	check(err)
	me, err := root.Me()
	check(err)
	fmt.Println(me.DisplayName())
	//fmt.Printf("me.M: %#v\n", (*lpad.Resource)(me).M)
}
