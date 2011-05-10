package lpad_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
)

var _ = Suite(&ResS{})
var _ = Suite(&ResI{})

type ResS struct {
	HTTPSuite
}

type ResI struct{
	SuiteI
}

var jsonType = map[string]string{
	"Content-Type": "application/json",
}

func (s *ResS) TestMapInit(c *C) {
	r := lpad.NewResource(nil, "", "", nil)
	m := r.Map()
	c.Assert(m, NotNil)
	m["a"] = 1
	c.Assert(r.Map()["a"], Equals, 1)
}

func (s *ResS) TestGet(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"a": 1, "b": [1, 2]}`)
	r := lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	err := r.Get(nil)
	c.Assert(err, IsNil)
	c.Assert(r.Map()["a"], Equals, float64(1))
	c.Assert(r.Map()["b"], Equals, []interface{}{float64(1), float64(2)})

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/myresource")
	c.Assert(req.Header.Get("Accept"), Equals, "application/json")
}

func (s *ResS) TestGetWithParams(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)
	r := lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	err := r.Get(lpad.Params{"k": "v"})
	c.Assert(err, IsNil)
	c.Assert(r.Map()["ok"], Equals, true)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/myresource")
	c.Assert(req.URL.RawQuery, Equals, "k=v")
}

func (s *ResS) TestGetSign(c *C) {
	oauth := &lpad.OAuth{Token: "mytoken", TokenSecret: "mytokensecret"}
	session := lpad.NewSession(oauth)

	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)
	r := lpad.NewResource(session, "", testServer.URL + "/myresource", nil)
	err := r.Get(nil)
	c.Assert(err, IsNil)
	c.Assert(r.Map()["ok"], Equals, true)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/myresource")
	c.Assert(req.Header["Authorization"], NotNil)
	c.Assert(req.Header["Authorization"][0], Matches, "OAuth.*")
}

func (s *ResS) TestGetRedirect(c *C) {
	headers := map[string]string{
		"Location": testServer.URL + "/myotherresource",
	}
	testServer.PrepareResponse(303, headers, "")
	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)
	r := lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	err := r.Get(nil)
	c.Assert(err, IsNil)
	c.Assert(r.URL(), Equals, testServer.URL + "/myotherresource")
	c.Assert(r.Map()["ok"], Equals, true)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/myresource")
}

func (s *ResS) TestGetNonJSONContent(c *C) {
	headers := map[string]string{
		"Content-Type": "text/plain",
	}
	testServer.PrepareResponse(200, headers, "NOT JSON")
	r := lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	err := r.Get(nil)
	c.Assert(err, Matches, "Non-JSON content-type: text/plain.*")
}

func (s *ResS) TestGetError(c *C) {
	testServer.PrepareResponse(500, jsonType, `{"what": "ever"}`)
	r := lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	err := r.Get(nil)
	c.Assert(err, Matches, `Server returned 500 and body: {"what": "ever"}`)

	testServer.PrepareResponse(404, jsonType, "")
	r = lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	err = r.Get(nil)
	c.Assert(err, Matches, `Server returned 404 and no body.`)
}

func (s *ResS) TestGetRedirectWithoutLocation(c *C) {
	headers := map[string]string{
		"Content-Type": "application/json", // Should be ignored.
	}
	testServer.PrepareResponse(303, headers, `{"ok": true}`)
	r := lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	err := r.Get(nil)
	c.Assert(err, Matches, "Got redirection status 303 without a Location")
}

func (s *ResS) TestPost(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)
	r := lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	other, err := r.Post(nil)
	c.Assert(err, IsNil)
	c.Assert(len(r.Map()), Equals, 0)
	c.Assert(len(other.Map()), Equals, 0)
	c.Assert(other.URL(), Equals, r.URL())

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/myresource")
}

func (s *ResS) TestPostWithParams(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)
	r := lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	_, err := r.Post(lpad.Params{"k": "v"})
	c.Assert(err, IsNil)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/myresource")
	c.Assert(req.Form["k"], Equals, []string{"v"})
}

func (s *ResS) TestPostCreation(c *C) {
	headers := map[string]string{
		"Location": "http://example.com",
		"Content-Type": "application/json", // Should be ignored.
	}
	testServer.PrepareResponse(201, headers, `{"ok": true}`)
	r := lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	other, err := r.Post(nil)
	c.Assert(err, IsNil)
	c.Assert(len(r.Map()), Equals, 0)
	c.Assert(len(other.Map()), Equals, 0)
	c.Assert(other.URL(), Equals, "http://example.com")

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/myresource")
}

func (s *ResS) TestPostSign(c *C) {
	oauth := &lpad.OAuth{Token: "mytoken", TokenSecret: "mytokensecret"}
	session := lpad.NewSession(oauth)

	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)
	r := lpad.NewResource(session, "", testServer.URL + "/myresource", nil)
	other, err := r.Post(nil)
	c.Assert(err, IsNil)
	c.Assert(len(r.Map()), Equals, 0)
	c.Assert(len(other.Map()), Equals, 0)
	c.Assert(other.URL(), Equals, r.URL())
	c.Assert(other.Session(), Equals, r.Session())

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/myresource")
	c.Assert(req.Header["Authorization"], NotNil)
	c.Assert(req.Header["Authorization"][0], Matches, "OAuth.*")
}

type locationTest struct {
	BaseURL, URL, Location, Result string
}

var locationTests = []locationTest{
	{"http://e.c/base/", "http://e.c/base/more/foo", "bar", "http://e.c/base/more/foo/bar"},
	{"http://e.c/base/", "http://e.c/base/more/foo", "../bar", "http://e.c/base/more/bar"},
	{"http://e.c/base/", "http://e.c/base/more/foo", "/bar", "http://e.c/base/bar"},
	{"http://e.c/base", "http://e.c/base/more/foo", "/bar", "http://e.c/base/bar"},
}

func (s *ResS) TestLocation(c *C) {
	oauth := &lpad.OAuth{Token: "mytoken", TokenSecret: "mytokensecret"}
	session := lpad.NewSession(oauth)

	for _, test := range locationTests {
		r1 := lpad.NewResource(session, test.BaseURL, test.URL, nil)
		r2 := r1.Location(test.Location)
		c.Assert(r2.URL(), Equals, test.Result)
		c.Assert(r2.BaseURL(), Equals, test.BaseURL)
		c.Assert(r2.Session(), Equals, session)
	}
}

func (s *ResS) TestLink(c *C) {
	oauth := &lpad.OAuth{Token: "mytoken", TokenSecret: "mytokensecret"}
	session := lpad.NewSession(oauth)

	for _, test := range locationTests {
		m := map[string]interface{}{"some_link": test.Location}
		r1 := lpad.NewResource(session, test.BaseURL, test.URL, m)
		r2 := r1.Link("some_link")
		c.Assert(r2.URL(), Equals, test.Result)
		c.Assert(r2.BaseURL(), Equals, test.BaseURL)
		c.Assert(r2.Session(), Equals, session)
	}
}
