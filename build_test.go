package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
)

func (s *ModelS) TestBuild(c *C) {
	m := M{
		"title":                           "thetitle",
		"arch_tag":                        "armel",
		"buildstate":                      "Failed to build",
		"build_log_url":                   "http://logurl",
		"upload_log_url":                  "http://uploadurl",
		"web_link":                        "http://page",
		"datecreated":                     "2011-10-10T00:00:00",
		"datebuilt":                       "2011-10-10T00:00:10",
		"current_source_publication_link": testServer.URL + "/current_source_publication_link",
	}
	build := &lpad.Build{lpad.NewValue(nil, "", "", m)}
	c.Assert(build.Title(), Equals, "thetitle")
	c.Assert(build.Arch(), Equals, "armel")
	c.Assert(build.State(), Equals, lpad.BuildState("Failed to build"))
	c.Assert(build.BuildLogURL(), Equals, "http://logurl")
	c.Assert(build.UploadLogURL(), Equals, "http://uploadurl")
	c.Assert(build.WebPage(), Equals, "http://page")
	c.Assert(build.DateCreated(), Equals, "2011-10-10T00:00:00")
	c.Assert(build.DateBuilt(), Equals, "2011-10-10T00:00:10")

	testServer.PrepareResponse(200, jsonType, `{"source_package_name": "packagename"}`)
	ph, err := build.PubHistory()
	c.Assert(err, IsNil)
	c.Assert(ph.PackageName(), Equals, "packagename")

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/current_source_publication_link")
}

func (s *ModelS) TestBuildRetry(c *C) {
	testServer.PrepareResponse(200, jsonType, "{}")
	build := &lpad.Build{lpad.NewValue(nil, testServer.URL, testServer.URL, nil)}
	err := build.Retry()
	c.Assert(err, IsNil)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.Form["ws.op"], Equals, []string{"retry"})
}

func (s *ModelS) TestPubHistory(c *C) {
	m := M{
		"source_package_name":    "pkgname",
		"source_package_version": "pkgversion",
		"component_name":         "main",
		"distro_series_link":     testServer.URL + "/distro_series_link",
		"archive_link":           testServer.URL + "/archive_link",
	}
	ph := &lpad.PubHistory{lpad.NewValue(nil, "", "", m)}
	c.Assert(ph.PackageName(), Equals, "pkgname")
	c.Assert(ph.PackageVersion(), Equals, "pkgversion")
	c.Assert(ph.Component(), Equals, "main")

	testServer.PrepareResponse(200, jsonType, `{"name": "archivename"}`)
	archive, err := ph.Archive()
	c.Assert(err, IsNil)
	c.Assert(archive.Name(), Equals, "archivename")

	testServer.PrepareResponse(200, jsonType, `{"name": "seriesname"}`)
	series, err := ph.DistroSeries()
	c.Assert(err, IsNil)
	c.Assert(series.Name(), Equals, "seriesname")
}
