package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
)

func (s *ModelS) TestBuild(c *C) {
	m := M{
		"title":          "thetitle",
		"arch_tag":       "armel",
		"buildstate":     "Failed to build",
		"build_log_url":  "http://logurl",
		"upload_log_url": "http://uploadurl",
		"web_link":       "http://page",
		"self_link":      "http://apipage",
		"datecreated":    "2011-10-10T00:00:00",
		"datebuilt":      "2011-10-10T00:00:10",
	}
	build := &lpad.Build{lpad.NewValue(nil, "", "", m)}
	c.Assert(build.Title(), Equals, "thetitle")
	c.Assert(build.ArchTag(), Equals, "armel")
	c.Assert(build.State(), Equals, "Failed to build")
	c.Assert(build.BuildLogURL(), Equals, "http://logurl")
	c.Assert(build.UploadLogURL(), Equals, "http://uploadurl")
	c.Assert(build.WebPage(), Equals, "http://page")
	c.Assert(build.SelfLink(), Equals, "http://apipage")
	c.Assert(build.DateCreated(), Equals, "2011-10-10T00:00:00")
	c.Assert(build.DateBuilt(), Equals, "2011-10-10T00:00:10")

}

func (s *ModelS) TestSPPH(c *C) {
	m := M{
		"source_package_name":    "pkgname",
		"source_package_version": "pkgversion",
		"component_name":         "main",
		"distro_series_link":     testServer.URL + "/distro_series_link",
		"archive_link":           testServer.URL + "/archive_link",
	}
	spph := &lpad.SPPH{lpad.NewValue(nil, "", "", m)}
	c.Assert(spph.PackageName(), Equals, "pkgname")
	c.Assert(spph.PackageVersion(), Equals, "pkgversion")
	c.Assert(spph.Component(), Equals, "main")

	testServer.PrepareResponse(200, jsonType, `{"name": "archivename"}`)
	archive, err := spph.Archive()
	c.Assert(err, IsNil)
	c.Assert(archive.Name(), Equals, "archivename")

	testServer.PrepareResponse(200, jsonType, `{"name": "seriesname"}`)
	series, err := spph.DistroSeries()
	c.Assert(err, IsNil)
	c.Assert(series.Name(), Equals, "seriesname")
}
