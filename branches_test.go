package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
)

func (s *ModelS) TestBranch(c *C) {
	m := M{
		"bzr_identity": "lp:~joe/ensemble",
		"unique_name":  "lp:~joe/ensemble/some-branch",
		"web_link": "http://page",
	}
	branch := lpad.Branch{lpad.NewValue(nil, "", "", m)}
	c.Assert(branch.Id(), Equals, "lp:~joe/ensemble")
	c.Assert(branch.UniqueName(), Equals, "lp:~joe/ensemble/some-branch")
	c.Assert(branch.WebPage(), Equals, "http://page")
}

func (s *ModelS) TestRootBranch(c *C) {
	data := `{"unique_name": "lp:branch"}`
	testServer.PrepareResponse(200, jsonType, data)

	root := lpad.Root{lpad.NewValue(nil, testServer.URL, "", nil)}

	branch, err := root.Branch("lp:~joe/project/branch-name")
	c.Assert(err, IsNil)
	c.Assert(branch.UniqueName(), Equals, "lp:branch")

	req := testServer.WaitRequest()
	c.Assert(req.URL.Path, Equals, "/branches")
	c.Assert(req.Form["ws.op"], Equals, []string{"getByUrl"})
	c.Assert(req.Form["url"], Equals, []string{"lp:~joe/project/branch-name"})
}

func (s *ModelS) TestMergeProposal(c *C) {
	m := M{
		"description": "Description",
		"commit_message": "Commit message",
		"queue_status": "Needs review",
		"address": "some@email.com",
		"web_link": "http://page",
		"prerequisite_branch_link": testServer.URL + "/prereq_link",
		"target_branch_link": testServer.URL + "/target_link",
		"source_branch_link": testServer.URL + "/source_link",
	}
	mp := lpad.MergeProposal{lpad.NewValue(nil, "", "", m)}
	c.Assert(mp.Description(), Equals, "Description")
	c.Assert(mp.CommitMessage(), Equals, "Commit message")
	c.Assert(mp.Status(), Equals, "Needs review")
	c.Assert(mp.Email(), Equals, "some@email.com")
	c.Assert(mp.WebPage(), Equals, "http://page")

	testServer.PrepareResponse(200, jsonType, `{"unique_name": "branch1"}`)
	testServer.PrepareResponse(200, jsonType, `{"unique_name": "branch2"}`)
	testServer.PrepareResponse(200, jsonType, `{"unique_name": "branch3"}`)

	b1, err := mp.Target()
	c.Assert(err, IsNil)
	c.Assert(b1.UniqueName(), Equals, "branch1")

	b2, err := mp.PreReq()
	c.Assert(err, IsNil)
	c.Assert(b2.UniqueName(), Equals, "branch2")

	b3, err := mp.Source()
	c.Assert(err, IsNil)
	c.Assert(b3.UniqueName(), Equals, "branch3")

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/target_link")

	req = testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/prereq_link")

	req = testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/source_link")
}

func (s *ModelS) TestBranchProposeMerge(c *C) {
	data := `{"description": "Description"}`
	testServer.PrepareResponse(200, jsonType, data)
	branch := lpad.Branch{lpad.NewValue(nil, testServer.URL, testServer.URL + "/~joe/ensemble/some-branch", nil)}
	target := lpad.Branch{lpad.NewValue(nil, testServer.URL, testServer.URL + "/~ensemble/ensemble/trunk", nil)}

	stub := lpad.MergeStub{
		Description: "Description",
		CommitMessage: "Commit message",
		NeedsReview: true,
		Target: target,
	}

	mp, err := branch.ProposeMerge(&stub)
	c.Assert(err, IsNil)
	c.Assert(mp.Description(), Equals, "Description")

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/~joe/ensemble/some-branch")
	c.Assert(req.Form["commit_message"], Equals, []string{"Commit message"})
	c.Assert(req.Form["initial_comment"], Equals, []string{"Description"})
	c.Assert(req.Form["needs_review"], Equals, []string{"true"})
	c.Assert(req.Form["target_branch"], Equals, []string{target.AbsLoc()})
}

func (s *ModelS) TestBranchProposeMergePreReq(c *C) {
	data := `{"description": "Description"}`
	testServer.PrepareResponse(200, jsonType, data)
	branch := lpad.Branch{lpad.NewValue(nil, testServer.URL, testServer.URL + "/~joe/ensemble/some-branch", nil)}
	target := lpad.Branch{lpad.NewValue(nil, testServer.URL, testServer.URL + "~ensemble/ensemble/trunk", nil)}
	prereq := lpad.Branch{lpad.NewValue(nil, testServer.URL, testServer.URL + "~ensemble/ensemble/prereq", nil)}

	stub := lpad.MergeStub{
		Target: target,
		PreReq: prereq,
	}

	mp, err := branch.ProposeMerge(&stub)
	c.Assert(err, IsNil)
	c.Assert(mp.Description(), Equals, "Description")

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/~joe/ensemble/some-branch")
	c.Assert(req.Form["commit_message"], Equals, []string{})
	c.Assert(req.Form["initial_comment"], Equals, []string{})
	c.Assert(req.Form["needs_review"], Equals, []string{})
	c.Assert(req.Form["target_branch"], Equals, []string{target.AbsLoc()})
	c.Assert(req.Form["prerequisite_branch"], Equals, []string{prereq.AbsLoc()})
}
