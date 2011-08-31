package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
)

func (s *ModelS) TestBug(c *C) {
	m := M{
		"id": 123456.0,
		"title": "Title",
		"description": "Description",
		"tags": "a b c",
		"private": true,
		"security_related": true,
	}
	bug := lpad.Bug{lpad.NewResource(nil, "", "", m)}
	c.Assert(bug.Id(), Equals, 123456)
	c.Assert(bug.Title(), Equals, "Title")
	c.Assert(bug.Description(), Equals, "Description")
	c.Assert(bug.Tags(), Equals, []string{"a", "b", "c"})
	c.Assert(bug.Private(), Equals, true)
	c.Assert(bug.SecurityRelated(), Equals, true)
	bug.SetTitle("New title")
	bug.SetDescription("New description")
	bug.SetTags([]string{"new", "tags"})
	bug.SetPrivate(false)
	bug.SetSecurityRelated(false)
	c.Assert(bug.Title(), Equals, "New title")
	c.Assert(bug.Description(), Equals, "New description")
	c.Assert(bug.Tags(), Equals, []string{"new", "tags"})
	c.Assert(bug.Private(), Equals, false)
	c.Assert(bug.SecurityRelated(), Equals, false)
}

func (s *ModelS) TestRootCreateBug(c *C) {
	data := `{
		"id": 123456,
		"title": "Title",
		"description": "Description",
		"private": true,
		"security_related": true,
		"tags": "a b c"
	}`
	testServer.PrepareResponse(200, jsonType, data)
	root := lpad.Root{lpad.NewResource(nil, testServer.URL, "", nil)}
	stub := lpad.BugStub{
		Title: "Title",
		Description: "Description.",
		Private: true,
		SecurityRelated: true,
		Tags: []string{"a", "b", "c"},
		Target: lpad.NewResource(nil, "", "http://target", nil),
	}
	bug, err := root.CreateBug(&stub)
	c.Assert(err, IsNil)
	c.Assert(bug.Title(), Equals, "Title")

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/bugs")
	c.Assert(req.Form["ws.op"], Equals, []string{"createBug"})
	c.Assert(req.Form["title"], Equals, []string{"Title"})
	c.Assert(req.Form["description"], Equals, []string{"Description."})
	c.Assert(req.Form["private"], Equals, []string{"true"})
	c.Assert(req.Form["security_related"], Equals, []string{"true"})
	c.Assert(req.Form["tags"], Equals, []string{"a b c"})
	c.Assert(req.Form["target"], Equals, []string{"http://target"})
}
