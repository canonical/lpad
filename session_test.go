package lpad_test

import (
	"http"
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
	"os"
)

var _ = Suite(&SessionS{})
var _ = Suite(&SessionI{})

type SessionS struct {
	HTTPSuite
}

type SessionI struct{
	SuiteI
}

type dummyAuth struct {
	loginBaseURL string
	loginErr os.Error
	signReq *http.Request
	signErr os.Error
}

func (a *dummyAuth) Login(baseURL string) os.Error {
	a.loginBaseURL = baseURL
	return a.loginErr
}

func (a *dummyAuth) Sign(r *http.Request) os.Error {
	a.signReq = r
	return a.signErr
}

func (s *SessionS) TestLogin(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)

	auth := &dummyAuth{}
	root, err := lpad.Login(testServer.URL, auth)
	c.Assert(err, IsNil)
	c.Assert(auth.loginBaseURL, Equals, testServer.URL)

	c.Assert(root.BaseURL(), Equals, testServer.URL)
	c.Assert(root.URL(), Equals, testServer.URL)
	c.Assert(len(root.Map()), Equals, 0)

	err = root.Get(nil)
	c.Assert(err, IsNil)
	c.Assert(root.Map()["ok"], Equals, true)

	c.Assert(auth.signReq, NotNil)
	c.Assert(auth.signReq.URL.String(), Equals, testServer.URL)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.RawPath, Equals, "/")
}

var lpadAuth = &lpad.OAuth{
	Token: "SfVJpl7pJgSLJX9cm0wj",
	TokenSecret: "CXJGg1t5gTdjDqtFG0HNBFQn8WLWq8QQ3B2sHh9NmgLxQ6kGl9m123gQLZpDF8HFxQzk8HV78c9sGHQb",
}

func (s *SessionI) TestLogin(c *C) {
	root, err := lpad.Login(lpad.Production, lpadAuth)
	me, err := root.Me()
	c.Assert(err, IsNil)
	c.Assert(me.DisplayName(), Equals, "Lpad Test User")
}

