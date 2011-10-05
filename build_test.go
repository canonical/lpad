package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
)

func (s *ModelS) TestBuild(c *C) {
	m := M{
		"title":         "thetitle",
		"arch_tag":      "armel",
		"build_log_url": "http://logurl",
		"web_link":      "http://page",
	}
	build := lpad.Build{lpad.NewValue(nil, "", "", m)}
	c.Assert(build.Title(), Equals, "thetitle")
	c.Assert(build.ArchTag(), Equals, "armel")
	c.Assert(build.BuildLogURL(), Equals, "http://logurl")
	c.Assert(build.WebPage(), Equals, "http://page")
}
