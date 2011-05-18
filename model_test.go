package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
)

var _ = Suite(&ModelS{})
var _ = Suite(&ModelI{})

type ModelS struct {
	HTTPSuite
}

type ModelI struct {
	SuiteI
}

type M map[string]interface{}

func (s *ModelS) TestRootMe(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"display_name": "Joe"}`)

	root := lpad.Root{lpad.NewResource(nil, testServer.URL, "", nil)}

	me, err := root.Me()
	c.Assert(err, IsNil)
	c.Assert(me.DisplayName(), Equals, "Joe")

	req := testServer.WaitRequest()
	c.Assert(req.URL.Path, Equals, "/people/+me")
}

func (s *ModelS) TestPerson(c *C) {
	m := M{
		"display_name": "Joe",
	}
	person := lpad.Person{lpad.NewResource(nil, "", "", m)}
	c.Assert(person.DisplayName(), Equals, "Joe")
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
