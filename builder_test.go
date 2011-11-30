package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
)

func (s *ModelS) TestBuilder(c *C) {
	m := M{
		"name":        "thename",
		"title":       "Title",
		"active":      true,
		"builderok":   true,
		"virtualized": "false",
		"vm_host":     "foobar",
		"web_link":    "http://page",
	}
	builder := &lpad.Builder{lpad.NewValue(nil, "", "", m)}
	c.Assert(builder.Name(), Equals, "thename")
	c.Assert(builder.Title(), Equals, "Title")
	c.Assert(builder.Active(), Equals, true)
	c.Assert(builder.BuilderOK(), Equals, true)
	c.Assert(builder.Virtualized(), Equals, false)
	c.Assert(builder.VMHost(), Equals, "foobar")
	c.Assert(builder.WebPage(), Equals, "http://page")
}
