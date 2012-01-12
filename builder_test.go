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

func (s *ModelS) TestBuilderList(c *C) {
	data := `{
		"total_size": 2,
		"start": 0,
        "entries": [
            {"name":"builder1", "builderok":true},
            {"name":"builder2", "builderok":true}
            ]
        }`

	testServer.PrepareResponse(200, jsonType, data)
	root := &lpad.Root{lpad.NewValue(nil, testServer.URL, "", nil)}

	builders, err := root.Builders()
	c.Assert(err, IsNil)
	builders.For(func(builder *lpad.Builder) error {
		c.Assert(builder.BuilderOK(), Equals, true)
		return nil
	})

}
