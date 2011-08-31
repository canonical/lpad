package lpad

import (
	"bytes"
	"fmt"
	"http"
	"io/ioutil"
	"json"
	"os"
	"path"
	"strconv"
	"strings"
	"url"
)

// The Params type is a helper to pass parameter into the Resource request
// methods.  It may be used as:
//
//     resource.Get(lpad.Params{"name": "value"})
//
type Params map[string]string

type Error struct {
	StatusCode int    // HTTP status code (500, 403, ...)
	Body       []byte // Body of response
}

func (e *Error) String() string {
	if len(e.Body) == 0 {
		return "Server returned " + strconv.Itoa(e.StatusCode) + " and no body."
	}
	return "Server returned " + strconv.Itoa(e.StatusCode) + " and body: " + string(e.Body)
}

// The Resource interface is implemented by all model types and enables
// accessing the underlying representation of any value.  Besides being
// used internally to implement the static model itself, this interface
// enables accessing new features available in Launchpad which were not
// yet made available in lpad thorugh the more convenient static model.
type Resource interface {

	// Session returns the session for the interaction with Launchpad.
	// This session is used to sign any requests delivered to Launchpad.
	Session() *Session

	// BaseURL returns the URL the session is based on.  Absolute URLs
	// provided to Location and Link will be rooted at this place.
	BaseURL() string

	// URL returns the location of this resource.
	URL() string

	// Map returns the dynamic map with the content of this resource.
	Map() map[string]interface{}

	// StringField returns the named resource field if it exists and is
	// set to a string value, or the empty string otherwise.
	StringField(key string) string

	// IntField returns the named resource field if it exists and is
	// set to an int value, or zero otherwise.
	IntField(key string) int

	// FloatField returns the named resource field if it exists and is
	// set to a float64 value, or zero otherwise.
	FloatField(key string) float64

	// BoolField returns the named resource field if it exists and is
	// set to a bool value, or false otherwise.
	BoolField(key string) bool

	// SetField changes the named field with the provided value.
	SetField(key string, value interface{})

	// Location returns a new resource for a location which may be a
	// full URL, or an absolute path (based on the resource's BaseURL),
	// or a path relative to the resource itself (based on the
	// resource's URL).
	Location(url_ string) Resource

	// GetLocation builds a resource with Location and calls Get(nil)
	// on it.  It returns the loaded resource in case of success.
	GetLocation(url_ string) (r Resource, err os.Error)

	// Link calls Location with a URL available in the given key
	// of the current resource's Map.  It returns an error if the
	// requested key isn't found in the resource.  This is a convenient
	// way to navigate through *_link values.
	Link(key string) (r Resource, err os.Error)

	// GetLink builds a resource with Link and calls Get(nil) on it.
	// It returns the loaded resource in case of success.
	GetLink(key string) (r Resource, err os.Error)

	// Get issues an HTTP GET to retrieve the content of this resource.
	// If params is not nil, it will provided as the query for the GET
	// request.
	Get(params Params) os.Error

	// Post issues an HTTP POST to perform a given action at the URL
	// specified by this resource.  If params is not nil, it will
	// provided as the parameters for the POST request.
	Post(params Params) (Resource, os.Error)

	// Patch issues an HTTP PATCH request to modify the server resource
	// with the local changes.
	Patch() os.Error

	// TotalSize returns the total number of entries in a collection.
	TotalSize() int

	// StartIndex returns the offset of the first resource in a collection.
	StartIndex() int

	// For iterates over every element in a collection and calls the
	// provided function for each entry.  If the function returns a
	// non-nil err value, the iteration will stop.  Watch out for
	// very large collections!
	For(func(r Resource) os.Error) os.Error
}

type resource struct {
	session *Session
	baseurl string
	url     string
	m       map[string]interface{}
	patch   map[string]interface{}
}

// NewResource creates a new resource type.  Creating resources explicitly
// is generally not necessary.  If you're trying to access a location in
// the Launchpad API which is not covered by the static model yet, see the
// Link and Location methods on the Resource interface for more convenient
// ways to create resources.
func NewResource(session *Session, baseurl, url_ string, m map[string]interface{}) Resource {
	return &resource{session, baseurl, url_, m, nil}
}

func (r *resource) Session() *Session {
	return r.session
}

func (r *resource) BaseURL() string {
	return r.baseurl
}

func (r *resource) URL() string {
	return r.url
}

func (r *resource) Map() map[string]interface{} {
	if r.m == nil {
		r.m = make(map[string]interface{})
	}
	return r.m
}

func (r *resource) StringField(key string) string {
	if v, ok := r.Map()[key].(string); ok {
		return v
	}
	return ""
}

func (r *resource) IntField(key string) int {
	if v, ok := r.Map()[key].(float64); ok {
		return int(v)
	}
	return 0
}

func (r *resource) FloatField(key string) float64 {
	if v, ok := r.Map()[key].(float64); ok {
		return v
	}
	return 0
}

func (r *resource) BoolField(key string) bool {
	if v, ok := r.Map()[key].(bool); ok {
		return v
	}
	return false
}

func (r *resource) SetField(key string, value interface{}) {
	if r.patch == nil {
		r.patch = make(map[string]interface{})
	}
	p := r.patch
	m := r.Map()
	var newv interface{}
	switch v := value.(type) {
	case int:
		newv = float64(v)
	case string:
		newv = v
	case bool:
		newv = v
	default:
		panic(fmt.Sprintf("Unsupported value type for SetField: %#v", value))
	}
	p[key] = newv
	m[key] = newv
}

func (r *resource) join(part string) string {
	if part == "" {
		return r.url
	}
	if strings.HasPrefix(part, "http://") || strings.HasPrefix(part, "https://") {
		return part
	}
	base := r.baseurl
	if !strings.HasPrefix(part, "/") {
		// Relative to URL.
		base = r.url
	}
	url_, err := url.Parse(base)
	if err != nil {
		panic("Invalid URL: " + base)
	}
	url_.Path = path.Join(url_.Path, part)
	return url_.String()
}

func (r *resource) Location(url_ string) Resource {
	return &resource{session: r.session, baseurl: r.baseurl, url: r.join(url_)}
}

func (r *resource) Link(key string) (linkr Resource, err os.Error) {
	location, ok := r.m[key].(string)
	if !ok {
		return nil, os.NewError(fmt.Sprintf("Field %q not found in resource", key))
	}
	return r.Location(location), nil
}

func (r *resource) GetLocation(url_ string) (linkr Resource, err os.Error) {
	linkr = r.Location(url_)
	err = linkr.Get(nil)
	if err != nil {
		return nil, err
	}
	return linkr, nil
}

func (r *resource) GetLink(key string) (linkr Resource, err os.Error) {
	linkr, err = r.Link(key)
	if err != nil {
		return nil, err
	}
	err = linkr.Get(nil)
	if err != nil {
		return nil, err
	}
	return linkr, nil
}

func (r *resource) Get(params Params) os.Error {
	_, err := r.do("GET", params, nil)
	return err
}

func (r *resource) Post(params Params) (other Resource, err os.Error) {
	return r.do("POST", params, nil)
}

func (r *resource) Patch() os.Error {
	data, err := json.Marshal(r.patch)
	if err != nil {
		return err
	}
	_, err = r.do("PATCH", nil, data)
	return err
}

func (r *resource) TotalSize() int {
	return r.IntField("total_size")
}

func (r *resource) StartIndex() int {
	return r.IntField("start")
}

func (r *resource) For(f func(Resource) os.Error) os.Error {
	for {
		entries, ok := r.Map()["entries"].([]interface{})
		if !ok {
			return os.NewError("No entries found in resource")
		}
		for _, entry := range entries {
			m, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			url_, _ := m["self_link"].(string)
			err := f(&resource{session: r.session, baseurl: r.baseurl, url: url_, m: m})
			if err != nil {
				return err
			}
		}
		nextr, _ := r.Link("next_collection_link")
		if nextr == nil {
			break
		}
		err := nextr.Get(nil)
		if err != nil {
			return err
		}
		r = nextr.(*resource)
	}
	return nil
}

var stopRedir = os.NewError("Stop redirection marker.")

var httpClient = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) os.Error { return stopRedir },
}

func (r *resource) do(method string, params Params, body []byte) (res *resource, err os.Error) {
	res = r
	query := multimap(params).Encode()
	for redirect := 0; ; redirect++ {
		req, err := http.NewRequest(method, res.url, nil)
		req.Header["Accept"] = []string{"application/json"}
		if err != nil {
			return nil, err
		}

		ctype := "application/json"
		if method == "POST" {
			body = []byte(query)
			query = ""
			ctype = "application/x-www-form-urlencoded"
		} else {
			if req.URL.RawQuery != "" {
				req.URL.RawQuery += "&"
			}
			req.URL.RawQuery += query
		}

		if body != nil {
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			req.Header["Content-Type"] = []string{ctype}
			req.Header["Content-Length"] = []string{strconv.Itoa(len(body))}
			req.ContentLength = int64(len(body))
		}

		if r.session != nil {
			err := r.session.Sign(req)
			if err != nil {
				return nil, err
			}
		}

		resp, err := httpClient.Do(req)
		if urlerr, ok := err.(*url.Error); ok && urlerr.Error == stopRedir {
			// fine
		} else if err != nil {
			return nil, err
		}
		//dump, err := http.DumpResponse(resp, true)
		//println("Response:", string(dump))

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		location := resp.Header.Get("Location")

		if method == "POST" {
			res = &resource{url: r.url, baseurl: r.baseurl, session: r.session}
			if resp.StatusCode == 201 && location != "" {
				res.url = location
				return res.do("GET", nil, nil)
			}
		}

		if method == "GET" && shouldRedirect(resp.StatusCode) {
			if location == "" {
				msg := "Got redirection status " + strconv.Itoa(resp.StatusCode) + " without a Location"
				return nil, os.NewError(msg)
			}
			res.url = location
			continue
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != 209 {
			return nil, &Error{resp.StatusCode, body}
		}

		if method == "PATCH" && resp.StatusCode != 209 {
			return nil, nil
		}

		ctype = resp.Header.Get("Content-Type")
		if ctype != "application/json" {
			return nil, os.NewError("Non-JSON content-type: " + ctype)
		}

		if err != nil {
			return nil, err
		}
		res.m = make(map[string]interface{})
		return res, json.Unmarshal(body, &res.m)
	}

	panic("unreachable")
}

func shouldRedirect(statusCode int) bool {
	switch statusCode {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
		return true
	}
	return false
}

func multimap(params map[string]string) url.Values {
	m := make(url.Values, len(params))
	for k, v := range params {
		m[k] = []string{v}
	}
	return m
}
