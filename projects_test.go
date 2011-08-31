package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
)

func (s *ModelS) TestProject(c *C) {
	m := M{
		"name": "thename",
		"display_name": "Display Name",
		"title": "Title",
		"summary": "Summary",
		"description": "Description",
	}
	project := lpad.Project{lpad.NewResource(nil, "", "", m)}
	c.Assert(project.Name(), Equals, "thename")
	c.Assert(project.DisplayName(), Equals, "Display Name")
	c.Assert(project.Title(), Equals, "Title")
	c.Assert(project.Summary(), Equals, "Summary")
	c.Assert(project.Description(), Equals, "Description")
	project.SetName("newname")
	project.SetDisplayName("New Display Name")
	project.SetTitle("New Title")
	project.SetSummary("New summary")
	project.SetDescription("New description")
	c.Assert(project.Name(), Equals, "newname")
	c.Assert(project.DisplayName(), Equals, "New Display Name")
	c.Assert(project.Title(), Equals, "New Title")
	c.Assert(project.Summary(), Equals, "New summary")
	c.Assert(project.Description(), Equals, "New description")
}

func (s *ModelS) TestRootProject(c *C) {
	data := `{
		"name": "Name",
		"title": "Title",
		"description": "Description"
	}`
	testServer.PrepareResponse(200, jsonType, data)
	root := lpad.Root{lpad.NewResource(nil, testServer.URL, "", nil)}
	project, err := root.Project("myproj")
	c.Assert(err, IsNil)
	c.Assert(project.Name(), Equals, "Name")

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/myproj")
}

