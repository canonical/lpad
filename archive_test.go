package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
)

func (s *ModelS) TestArchive(c *C) {
	m := M{
		"name":        "thename",
		"displayname": "The Name",
		"description": "The Description",
		"self_link":   "http://apipage",
		"web_link":    "http://page",
	}
	archive := &lpad.Archive{lpad.NewValue(nil, "", "", m)}
	c.Assert(archive.Name(), Equals, "thename")
	c.Assert(archive.DisplayName(), Equals, "The Name")
	c.Assert(archive.Description(), Equals, "The Description")
	c.Assert(archive.SelfLink(), Equals, "http://apipage")
	c.Assert(archive.WebPage(), Equals, "http://page")
}
