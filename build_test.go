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
		"self_link":      "http://apipage",
	}
	build := &lpad.Build{lpad.NewValue(nil, "", "", m)}
	c.Assert(build.Title(), Equals, "thetitle")
	c.Assert(build.ArchTag(), Equals, "armel")
	c.Assert(build.BuildLogURL(), Equals, "http://logurl")
	c.Assert(build.WebPage(), Equals, "http://page")
	c.Assert(build.SelfLink(), Equals, "http://apipage")
}

func (s *ModelS) TestSPPH(c *C) {
	m := M{
		"source_package_name":         "pkgname",
		"source_package_version":      "pkgversion",
		"component_name":             "main",
	}
	spph := &lpad.SPPH{lpad.NewValue(nil, "", "", m)}
	c.Assert(spph.PackageName(), Equals, "pkgname")
	c.Assert(spph.PackageVersion(), Equals, "pkgversion")
	c.Assert(spph.Component(), Equals, "main")
}
