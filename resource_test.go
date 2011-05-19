package lpad_test

import (
	"fmt"
	"http"
	"json"
	. "launchpad.net/gocheck"
	"launchpad.net/lpad"
	"os"
	"strconv"
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

func (s *ResS) TestFieldMethods(c *C) {
	m := M{
		"n": nil,
		"s": "string",
		"f": 42.1,
		"b": true,
	}
	r := lpad.NewResource(nil, "", "", m)
	c.Assert(r.StringField("s"), Equals, "string")
	c.Assert(r.StringField("n"), Equals, "")
	c.Assert(r.StringField("x"), Equals, "")
	c.Assert(r.IntField("f"), Equals, 42)
	c.Assert(r.IntField("n"), Equals, 0)
	c.Assert(r.IntField("x"), Equals, 0)
	c.Assert(r.FloatField("f"), Equals, 42.1)
	c.Assert(r.FloatField("n"), Equals, 0.0)
	c.Assert(r.FloatField("x"), Equals, 0.0)
	c.Assert(r.BoolField("b"), Equals, true)
	c.Assert(r.BoolField("n"), Equals, false)
	c.Assert(r.BoolField("x"), Equals, false)
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

func (s *ResS) TestGetWithParamsMerging(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)
	r := lpad.NewResource(nil, "", testServer.URL + "/myresource?k2=v2", nil)
	err := r.Get(lpad.Params{"k1": "v1"})
	c.Assert(err, IsNil)
	c.Assert(r.Map()["ok"], Equals, true)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/myresource")
	params, err := http.ParseQuery(req.URL.RawQuery)
	c.Assert(err, IsNil)
	c.Assert(params["k1"], Equals, []string{"v1"})
	c.Assert(params["k2"], Equals, []string{"v2"})
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
	c.Assert(r.Map(), Equals, map[string]interface{}{})
	c.Assert(other.Map()["ok"], Equals, true)
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
		"Location": testServer.URL + "/newresource",
		"Content-Type": "application/json", // Should be ignored.
	}
	testServer.PrepareResponse(201, headers, `{"ok": false}`)
	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)

	r := lpad.NewResource(nil, testServer.URL, testServer.URL + "/myresource", nil)
	other, err := r.Post(nil)
	c.Assert(err, IsNil)
	c.Assert(len(r.Map()), Equals, 0)
	c.Assert(other.BaseURL(), Equals, testServer.URL)
	c.Assert(other.URL(), Equals, testServer.URL + "/newresource")
	c.Assert(other.Map()["ok"], Equals, true)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/myresource")

	req = testServer.WaitRequest()
	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/newresource")
	c.Assert(len(req.Form), Equals, 0)
}

func (s *ResS) TestPostSign(c *C) {
	oauth := &lpad.OAuth{Token: "mytoken", TokenSecret: "mytokensecret"}
	session := lpad.NewSession(oauth)

	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)
	r := lpad.NewResource(session, "", testServer.URL + "/myresource", nil)
	other, err := r.Post(nil)
	c.Assert(err, IsNil)
	c.Assert(len(r.Map()), Equals, 0)
	c.Assert(other.Map()["ok"], Equals, true)
	c.Assert(other.URL(), Equals, r.URL())
	c.Assert(other.Session(), Equals, r.Session())

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/myresource")
	c.Assert(req.Header["Authorization"], NotNil)
	c.Assert(req.Header["Authorization"][0], Matches, "OAuth.*")
}

func (s *ResS) TestPatch(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"a": 1, "b": 2}`)
	testServer.PrepareResponse(200, nil, "")

	r := lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	err := r.Get(nil)
	c.Assert(err, IsNil)

	r.SetField("a", 3)
	r.SetField("c", "string")
	c.Assert(r.Map()["a"], Equals, 3.0)
	c.Assert(r.Map()["b"], Equals, 2.0)
	c.Assert(r.Map()["c"], Equals, "string")

	err = r.Patch()
	c.Assert(err, IsNil)

	req1 := testServer.WaitRequest()
	c.Assert(req1.Method, Equals, "GET")
	c.Assert(req1.URL.Path, Equals, "/myresource")

	req2 := testServer.WaitRequest()
	c.Assert(req2.Method, Equals, "PATCH")
	c.Assert(req2.URL.Path, Equals, "/myresource")
	c.Assert(req2.Header.Get("Accept"), Equals, "application/json")
	c.Assert(req2.Header.Get("Content-Type"), Equals, "application/json")

	var m M
	err = json.Unmarshal([]byte(body(req2)), &m)
	c.Assert(err, IsNil)
	c.Assert(m, Equals, M{"a": 3.0, "c": "string"})
}

func (s *ResS) TestPatchWithContent(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"a": 1, "b": 2}`)
	testServer.PrepareResponse(209, jsonType, `{"new": "content"}`)

	r := lpad.NewResource(nil, "", testServer.URL + "/myresource", nil)
	err := r.Get(nil)
	c.Assert(err, IsNil)

	r.SetField("a", 3)

	err = r.Patch()
	c.Assert(err, IsNil)
	c.Assert(r.Map(), Equals, map[string]interface{}{"new": "content"})

	req1 := testServer.WaitRequest()
	c.Assert(req1.Method, Equals, "GET")
	c.Assert(req1.URL.Path, Equals, "/myresource")

	req2 := testServer.WaitRequest()
	c.Assert(req2.Method, Equals, "PATCH")
	c.Assert(req2.URL.Path, Equals, "/myresource")
	c.Assert(req2.Header.Get("Accept"), Equals, "application/json")
	c.Assert(req2.Header.Get("Content-Type"), Equals, "application/json")

	var m M
	err = json.Unmarshal([]byte(body(req2)), &m)
	c.Assert(err, IsNil)
	c.Assert(m, Equals, M{"a": 3.0})
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
		r2, err := r1.Link("some_link")
		c.Assert(err, IsNil)
		c.Assert(r2.URL(), Equals, test.Result)
		c.Assert(r2.BaseURL(), Equals, test.BaseURL)
		c.Assert(r2.Session(), Equals, session)

		r3, err := r1.Link("bad_link")
		c.Assert(err, Matches, `Field "bad_link" not found in resource`)
		c.Assert(r3, Equals, nil)
	}
}

func (s *ResS) TestGetLocation(c *C) {
	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)
	r := lpad.NewResource(nil, "", "", nil)
	other, err := r.GetLocation(testServer.URL + "/link")
	c.Assert(err, IsNil)
	c.Assert(other.Map()["ok"], Equals, true)
}

func (s *ResS) TestGetLink(c *C) {
	m := M{
		"some_link": testServer.URL + "/link",
	}
	testServer.PrepareResponse(200, jsonType, `{"ok": true}`)
	r := lpad.NewResource(nil, "", "", m)
	other, err := r.GetLink("some_link")
	c.Assert(err, IsNil)
	c.Assert(other.Map()["ok"], Equals, true)
}

func (s *ResS) TestCollection(c *C) {
	data0 := `{
		"total_size": 5,
		"start": 1,
		"next_collection_link": "%s",
		"entries": [{"self_link": "http://self1"}, {"self_link": "http://self2"}]
	}`
	data1 := `{
		"total_size": 5,
		"start": 3,
		"entries": [{"self_link": "http://self3"}, {"self_link": "http://self4"}]
	}`
	testServer.PrepareResponse(200, jsonType, fmt.Sprintf(data0, testServer.URL + "/next?n=10"))
	testServer.PrepareResponse(200, jsonType, data1)

	r := lpad.NewResource(nil, "", testServer.URL + "/mycol", nil)

	err := r.Get(nil)
	c.Assert(err, IsNil)

	c.Assert(r.TotalSize(), Equals, 5)
	c.Assert(r.StartIndex(), Equals, 1)

	i := 1
	err = r.For(func(r lpad.Resource) os.Error {
		c.Assert(r.Map()["self_link"], Equals, "http://self" + strconv.Itoa(i))
		i++
		return nil
	})
	c.Assert(err, IsNil)
	c.Assert(i, Equals, 5)

	testServer.WaitRequest()
	req1 := testServer.WaitRequest()
	c.Assert(req1.Form["n"], Equals, []string{"10"})
}

func (s *ResS) TestCollectionGetError(c *C) {
	data := `{
		"total_size": 2,
		"start": 0,
		"next_collection_link": "%s",
		"entries": [{"self_link": "http://self1"}]
	}`
	testServer.PrepareResponse(200, jsonType, fmt.Sprintf(data, testServer.URL + "/next"))
	testServer.PrepareResponse(500, jsonType, "")

	r := lpad.NewResource(nil, "", testServer.URL + "/mycol", nil)

	err := r.Get(nil)
	c.Assert(err, IsNil)

	i := 0
	err = r.For(func(r lpad.Resource) os.Error {
		i++
		return nil
	})
	c.Assert(err, Matches, ".* returned 500 .*")
	c.Assert(i, Equals, 1)
}

func (s *ResS) TestCollectionNoEntries(c *C) {
	data := `{"total_size": 2, "start": 0}`
	testServer.PrepareResponse(200, jsonType, data)
	r := lpad.NewResource(nil, "", testServer.URL + "/mycol", nil)

	err := r.Get(nil)
	c.Assert(err, IsNil)

	i := 0
	err = r.For(func(r lpad.Resource) os.Error {
		i++
		return nil
	})
	c.Assert(err, Matches, "No entries found in resource")
	c.Assert(i, Equals, 0)
}

func (s *ResS) TestCollectionIterError(c *C) {
	data := `{
		"total_size": 2,
		"start": 0,
		"entries": [{"self_link": "http://self1"}, {"self_link": "http://self2"}]
	}`
	testServer.PrepareResponse(200, jsonType, data)
	r := lpad.NewResource(nil, "", testServer.URL + "/mycol", nil)

	err := r.Get(nil)
	c.Assert(err, IsNil)

	i := 0
	err = r.For(func(r lpad.Resource) os.Error {
		i++
		return os.ErrorString("Stop!")
	})
	c.Assert(err, Matches, "Stop!")
	c.Assert(i, Equals, 1)
}
