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

type ModelI struct{
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

func (s *ModelS) TestPersonDisplayName(c *C) {
	m := M{
		"display_name": "Joe",
	}
	person := lpad.Person{lpad.NewResource(nil, "", "", m)}
	c.Assert(person.DisplayName(), Equals, "Joe")

	person = lpad.Person{lpad.NewResource(nil, "", "", nil)}
	c.Assert(person.DisplayName(), Equals, "")
}
