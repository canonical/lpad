package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
)

func (s *ModelS) TestBlueprint(c *C) {
	m := M{
		"name":       "thename",
		"title":      "Title",
		"summary":    "Summary",
		"whiteboard": "Whiteboard",
		"web_link":   "http://page",
	}
	project := &lpad.Blueprint{lpad.NewValue(nil, "", "", m)}
	c.Assert(project.Name(), Equals, "thename")
	c.Assert(project.Title(), Equals, "Title")
	c.Assert(project.Summary(), Equals, "Summary")
	c.Assert(project.Whiteboard(), Equals, "Whiteboard")
	c.Assert(project.WebPage(), Equals, "http://page")
	project.SetName("newname")
	project.SetTitle("New Title")
	project.SetSummary("New summary")
	project.SetWhiteboard("New whiteboard")
	c.Assert(project.Name(), Equals, "newname")
	c.Assert(project.Title(), Equals, "New Title")
	c.Assert(project.Summary(), Equals, "New summary")
	c.Assert(project.Whiteboard(), Equals, "New whiteboard")

	//testServer.PrepareResponse(200, jsonType, `{"name": "seriesname"}`)
	//series, err := project.FocusSeries()
	//c.Assert(err, IsNil)
	//c.Assert(series.Name(), Equals, "seriesname")

	//req := testServer.WaitRequest()
	//c.Assert(req.Method, Equals, "GET")
	//c.Assert(req.URL.Path, Equals, "/focus_link")
}

func (s *ModelS) TestBlueprintLinkBranch(c *C) {
	testServer.PrepareResponse(200, jsonType, `{}`)
	bp := &lpad.Blueprint{lpad.NewValue(nil, "", testServer.URL+"/project/+spec/the-bp", nil)}
	branch := lpad.Branch{lpad.NewValue(nil, testServer.URL, testServer.URL+"~joe/ensemble/some-branch", nil)}

	err := bp.LinkBranch(branch)
	c.Assert(err, IsNil)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/project/+spec/the-bp")
	c.Assert(req.Form["ws.op"], Equals, []string{"linkBranch"})
	c.Assert(req.Form["branch"], Equals, []string{branch.AbsLoc()})
}
