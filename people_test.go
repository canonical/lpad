package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
	"os"
)

func (s *ModelS) TestRootMe(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"display_name": "Joe"}`)

	root := lpad.Root{lpad.NewResource(nil, testServer.URL, "", nil)}

	me, err := root.Me()
	c.Assert(err, IsNil)
	c.Assert(me.DisplayName(), Equals, "Joe")

	req := testServer.WaitRequest()
	c.Assert(req.URL.Path, Equals, "/people/+me")
}

func (s *ModelS) TestRootPerson(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"display_name": "Joe"}`)

	root := lpad.Root{lpad.NewResource(nil, testServer.URL, "", nil)}

	person, err := root.Person("joe")
	c.Assert(err, IsNil)
	c.Assert(person.DisplayName(), Equals, "Joe")

	req := testServer.WaitRequest()
	c.Assert(req.URL.Path, Equals, "/people/~joe")
}

func (s *ModelS) TestRootTeam(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"display_name": "Ensemble", "is_team": true}`)

	root := lpad.Root{lpad.NewResource(nil, testServer.URL, "", nil)}

	team, err := root.Team("ensemble")
	c.Assert(err, IsNil)
	c.Assert(team.DisplayName(), Equals, "Ensemble")

	req := testServer.WaitRequest()
	c.Assert(req.URL.Path, Equals, "/people/~ensemble")
}

func (s *ModelS) TestRootFindMembers(c *C) {
	data := `{
		"total_size": 2,
		"start": 0,
		"entries": [{
			"self_link": "http://self0",
			"display_name": "Name0",
			"is_team": false
		}, {
			"self_link": "http://self1",
			"display_name": "Name1",
			"is_team": true
		}]
	}`
	testServer.PrepareResponse(200, jsonType, data)
	root := lpad.Root{lpad.NewResource(nil, testServer.URL, "", nil)}
	list, err := root.FindMembers("someuser")
	c.Assert(err, IsNil)
	c.Assert(list.TotalSize(), Equals, 2)

	names := []string{}
	list.For(func(r lpad.Resource) os.Error {
		if r.BoolField("is_team") {
			t := r.(lpad.Team)
			names = append(names, t.DisplayName())
		} else {
			p := r.(lpad.Person)
			names = append(names, p.DisplayName())
		}
		return nil
	})
	c.Assert(names, Equals, []string{"Name0", "Name1"})

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/people")
	c.Assert(req.Form["ws.op"], Equals, []string{"find"})
	c.Assert(req.Form["text"], Equals, []string{"someuser"})
}

func (s *ModelS) TestRootFindPeople(c *C) {
	data := `{
		"total_size": 2,
		"start": 0,
		"entries": [{
			"self_link": "http://self0",
			"display_name": "Name0"
		}, {
			"self_link": "http://self1",
			"display_name": "Name1"
		}]
	}`
	testServer.PrepareResponse(200, jsonType, data)
	root := lpad.Root{lpad.NewResource(nil, testServer.URL, "", nil)}
	list, err := root.FindPeople("someuser")
	c.Assert(err, IsNil)
	c.Assert(list.TotalSize(), Equals, 2)

	names := []string{}
	list.For(func(p lpad.Person) os.Error {
		names = append(names, p.DisplayName())
		return nil
	})
	c.Assert(names, Equals, []string{"Name0", "Name1"})

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/people")
	c.Assert(req.Form["ws.op"], Equals, []string{"findPerson"})
	c.Assert(req.Form["text"], Equals, []string{"someuser"})
}

func (s *ModelS) TestRootFindTeams(c *C) {
	data := `{
		"total_size": 2,
		"start": 0,
		"entries": [{
			"self_link": "http://self0",
			"display_name": "Name0",
			"is_team": true
		}, {
			"self_link": "http://self1",
			"display_name": "Name1",
			"is_team": true
		}]
	}`
	testServer.PrepareResponse(200, jsonType, data)
	root := lpad.Root{lpad.NewResource(nil, testServer.URL, "", nil)}
	list, err := root.FindTeams("someuser")
	c.Assert(err, IsNil)
	c.Assert(list.TotalSize(), Equals, 2)

	names := []string{}
	list.For(func(t lpad.Team) os.Error {
		names = append(names, t.DisplayName())
		return nil
	})
	c.Assert(names, Equals, []string{"Name0", "Name1"})

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/people")
	c.Assert(req.Form["ws.op"], Equals, []string{"findTeam"})
	c.Assert(req.Form["text"], Equals, []string{"someuser"})
}

func (s *ModelS) TestPerson(c *C) {
	m := M{
		"display_name": "Joe",
	}
	person := lpad.Person{lpad.NewResource(nil, "", "", m)}
	c.Assert(person.DisplayName(), Equals, "Joe")
	person.SetDisplayName("Name")
	c.Assert(person.DisplayName(), Equals, "Name")
}

func (s *ModelS) TestTeam(c *C) {
	m := M{
		"name": "myteam",
		"display_name": "My Team",
	}
	team := lpad.Team{lpad.NewResource(nil, "", "", m)}
	c.Assert(team.Name(), Equals, "myteam")
	team.SetName("ateam")
	c.Assert(team.Name(), Equals, "ateam")
	c.Assert(team.DisplayName(), Equals, "My Team")
	team.SetDisplayName("A Team")
	c.Assert(team.DisplayName(), Equals, "A Team")
}

func (s *ModelS) TestIRCNick(c *C) {
	m := M{
		"resource_type_link": "https://api.launchpad.net/1.0/#irc_id",
		"self_link":          "https://api.launchpad.net/1.0/~lpad-test/+ircnick/28983",
		"person_link":        "https://api.launchpad.net/1.0/~lpad-test",
		"web_link":           "https://api.launchpad.net/~lpad-test/+ircnick/28983",
		"nickname":           "canonical-nick",
		"network":            "irc.canonical.com",
		"http_etag":          "\"the-etag\"",
	}
	nick := lpad.IRCNick{lpad.NewResource(nil, "", "", m)}
	c.Assert(nick.Nick(), Equals, "canonical-nick")
	c.Assert(nick.Network(), Equals, "irc.canonical.com")
}

func (s *ModelS) TestIRCNickChange(c *C) {
	nick := lpad.IRCNick{lpad.NewResource(nil, "", "", nil)}
	nick.SetNick("mynick")
	nick.SetNetwork("mynetwork")
	c.Assert(nick.Nick(), Equals, "mynick")
	c.Assert(nick.Network(), Equals, "mynetwork")
}

func (s *ModelS) TestPersonNicks(c *C) {
	m := M{
		"irc_nicknames_collection_link": testServer.URL + "/~lpad-test/irc_nicknames",
	}
	data := `{
		"total_size": 2,
		"start": 0,
		"entries": [{
			"resource_type_link": "https://api.launchpad.net/1.0/#irc_id",
			"network": "irc.canonical.com",
			"person_link": "https://api.launchpad.net/1.0/~lpad-test",
			"web_link": "https://api.launchpad.net/~lpad-test/+ircnick/28983",
			"http_etag": "\"the-etag1\"",
			"self_link": "https://api.launchpad.net/1.0/~lpad-test/+ircnick/28983",
			"nickname": "canonical-nick"
		}, {
			"resource_type_link": "https://api.launchpad.net/1.0/#irc_id",
			"network": "irc.freenode.net",
			"person_link": "https://api.launchpad.net/1.0/~lpad-test",
			"web_link": "https://api.launchpad.net/~lpad-test/+ircnick/28982",
			"http_etag": "\"the-etag2\"",
			"self_link": "https://api.launchpad.net/1.0/~lpad-test/+ircnick/28982",
			"nickname": "freenode-nick"
		}],
		"resource_type_link": "https://api.launchpad.net/1.0/#irc_id-page-resource"
	}`
	testServer.PrepareResponse(200, jsonType, data)
	person := lpad.Person{lpad.NewResource(nil, "", "", m)}
	nicks, err := person.IRCNicks()
	c.Assert(err, IsNil)
	c.Assert(len(nicks), Equals, 2)
	c.Assert(nicks[0].Nick(), Equals, "canonical-nick") 
	c.Assert(nicks[1].Nick(), Equals, "freenode-nick") 
}
