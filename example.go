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
	fmt.Printf("me.M: %#v\n", me.Map())

	nicks, err := me.IRCNicks()
	check(err)
	for _, nick := range nicks {
		println(nick.Network(), "=", nick.Nick())
	}

	//langs, err := me.GetLink("irc_nicknames_collection_link")
	//check(err)
	//println(langs.ListSize())
	//check(err)
	//err = langs.ListIter(func(r lpad.Resource) os.Error {
	//	fmt.Printf("Entry: %#v\n", r.Map())
	//	return nil
	//})
	//check(err)
}
