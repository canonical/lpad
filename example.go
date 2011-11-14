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
	root, err := lpad.Login(lpad.Production, &lpad.ConsoleOAuth{})
	check(err)
	me, err := root.Me()
	check(err)
	fmt.Println(me.DisplayName())

	v, err := root.GetLocation("/bugs/123456")
	fmt.Printf("%#v\n", v.Map())

	//check(err)
	//
	//	list, err := root.FindTeams("ensemble")
	//	check(err)
	//
	//	fmt.Printf("Found %d members.\n", list.TotalSize())
	//
	//	i := 0
	//	err = list.For(func(t lpad.Team) os.Error {
	//		fmt.Printf("%s\n", t.DisplayName())
	//		i++
	//		return nil
	//	})
	//	check(err)
	//
	//	fmt.Printf("Had %d results, iterated over %d.\n", list.TotalSize(), i)

	//me, err := root.Me()
	//check(err)
	//fmt.Println(me.DisplayName())
	//fmt.Printf("me.M: %#v\n", me.Map()["is_team"])

	//nicks, err := me.IRCNicks()
	//check(err)
	//for _, nick := range nicks {
	//	println(nick.Network(), "=", nick.Nick())
	//	if nick.Network() == "irc.freenode.net" {
	//		nick.SetNick("newer-freenode-nick")
	//		err := nick.Patch()
	//		check(err)
	//	}
	//}

	//langs, err := me.GetLink("irc_nicknames_collection_link")
	//check(err)
	//println(langs.ListSize())
	//check(err)
	//err = langs.ListIter(func(v *lpad.Value) os.Error {
	//	fmt.Printf("Entry: %#v\n", v.Map())
	//	return nil
	//})
	//check(err)
}
