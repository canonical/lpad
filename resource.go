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

// The Params type is a helper to pass parameter into the Value request
// methods.  It may be used as:
//
//     value.Get(lpad.Params{"name": "value"})
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

// The AnyValue interface is implemented by *Value and thus by all the
// more specific value types supported. See the Value type for the
// meaning of these methods.
type AnyValue interface {
	IsValid() bool
	Session() *Session
	BaseURL() string
	URL() string
	Map() map[string]interface{}
	StringField(key string) string
	IntField(key string) int
	FloatField(key string) float64
	BoolField(key string) bool
	SetField(key string, value interface{})
	Location(url_ string) *Value
	GetLocation(url_ string) (v *Value, err os.Error)
	Link(key string) (v *Value, err os.Error)
	GetLink(key string) (v *Value, err os.Error)
	Get(params Params) os.Error
	Post(params Params) (*Value, os.Error)
	Patch() os.Error
	TotalSize() int
	StartIndex() int
	For(func(v *Value) os.Error) os.Error
}

// The Value type is the underlying dynamic layer used as the foundation of
// all the more specific value types that support the Launchpad model.
// Besides being used internally to implement these types, the methods of
// this type also enable accessing new features available in Launchpad which
// were not yet made available in lpad thorugh more convenient methods.
type Value struct {
	session *Session
	baseurl string
	url     string
	m       map[string]interface{}
	patch   map[string]interface{}
}

// NewValue creates a new Value with the provided details. Creating values
// explicitly is generally not necessary.  If you're trying to access a
// location in the Launchpad API which is not covered by the supported
// types yet, see the Link and Location methods on the Value type for more
// convenient ways to create values.
func NewValue(session *Session, baseurl, url_ string, m map[string]interface{}) *Value {
	return &Value{session, baseurl, url_, m, nil}
}

// IsValid returns true if the value is initialized and thus not nil. This
// provided mainly as a convenience for all the types that embed a *Value.
func (v *Value) IsValid() bool {
	return v != nil
}

// Session returns the session for the interaction with Launchpad.
// This session is used to sign any requests delivered to Launchpad.
func (v *Value) Session() *Session {
	return v.session
}

// BaseURL returns the URL the session is based on.  Absolute URLs
// provided to Location and Link will be rooted at this place.
func (v *Value) BaseURL() string {
	return v.baseurl
}

// URL returns the location of this value.
func (v *Value) URL() string {
	return v.url
}

// Map returns the dynamic map with the content of this value.
func (v *Value) Map() map[string]interface{} {
	if v.m == nil {
		v.m = make(map[string]interface{})
	}
	return v.m
}

// StringField returns the named value field if it exists and is
// set to a string value, or the empty string otherwise.
func (v *Value) StringField(key string) string {
	if v, ok := v.Map()[key].(string); ok {
		return v
	}
	return ""
}

// IntField returns the named value field if it exists and is
// set to an int value, or zero otherwise.
func (v *Value) IntField(key string) int {
	if v, ok := v.Map()[key].(float64); ok {
		return int(v)
	}
	return 0
}

// FloatField returns the named value field if it exists and is
// set to a float64 value, or zero otherwise.
func (v *Value) FloatField(key string) float64 {
	if v, ok := v.Map()[key].(float64); ok {
		return v
	}
	return 0
}

// BoolField returns the named value field if it exists and is
// set to a bool value, or false otherwise.
func (v *Value) BoolField(key string) bool {
	if v, ok := v.Map()[key].(bool); ok {
		return v
	}
	return false
}

// SetField changes the named field with the provided value.
func (v *Value) SetField(key string, value interface{}) {
	if v.patch == nil {
		v.patch = make(map[string]interface{})
	}
	p := v.patch
	m := v.Map()
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

func (v *Value) join(part string) string {
	if part == "" {
		return v.url
	}
	if strings.HasPrefix(part, "http://") || strings.HasPrefix(part, "https://") {
		return part
	}
	base := v.baseurl
	if !strings.HasPrefix(part, "/") {
		// Relative to URL.
		base = v.url
	}
	url_, err := url.Parse(base)
	if err != nil {
		panic("Invalid URL: " + base)
	}
	url_.Path = path.Join(url_.Path, part)
	return url_.String()
}

// Location returns a new value for a location which may be a
// full URL, or an absolute path (based on the value's BaseURL),
// or a path relative to the value itself (based on the
// value's URL).
func (v *Value) Location(url_ string) *Value {
	return &Value{session: v.session, baseurl: v.baseurl, url: v.join(url_)}
}

// GetLocation builds a value with Location and calls Get(nil)
// on it.  It returns the loaded value in case of success.
func (v *Value) GetLocation(url_ string) (linkv *Value, err os.Error) {
	linkv = v.Location(url_)
	err = linkv.Get(nil)
	if err != nil {
		return nil, err
	}
	return linkv, nil
}

// Link calls Location with a URL available in the given key
// of the current value's Map.  It returns an error if the
// requested key isn't found in the value.  This is a convenient
// way to navigate through *_link values.
func (v *Value) Link(key string) (linkv *Value, err os.Error) {
	location, ok := v.m[key].(string)
	if !ok {
		return nil, os.NewError(fmt.Sprintf("Field %q not found in value", key))
	}
	return v.Location(location), nil
}

// GetLink builds a value with Link and calls Get(nil) on it.
// It returns the loaded value in case of success.
func (v *Value) GetLink(key string) (linkv *Value, err os.Error) {
	linkv, err = v.Link(key)
	if err != nil {
		return nil, err
	}
	err = linkv.Get(nil)
	if err != nil {
		return nil, err
	}
	return linkv, nil
}

// Get issues an HTTP GET to retrieve the content of this value.
// If params is not nil, it will provided as the query for the GET
// request.
func (v *Value) Get(params Params) os.Error {
	_, err := v.do("GET", params, nil)
	return err
}

// Post issues an HTTP POST to perform a given action at the URL
// specified by this value.  If params is not nil, it will
// provided as the parameters for the POST request.
func (v *Value) Post(params Params) (other *Value, err os.Error) {
	return v.do("POST", params, nil)
}

// Patch issues an HTTP PATCH request to modify the server value
// with the local changes.
func (v *Value) Patch() os.Error {
	data, err := json.Marshal(v.patch)
	if err != nil {
		return err
	}
	_, err = v.do("PATCH", nil, data)
	return err
}

// TotalSize returns the total number of entries in a collection.
func (v *Value) TotalSize() int {
	return v.IntField("total_size")
}

// StartIndex returns the offset of the first value in a collection.
func (v *Value) StartIndex() int {
	return v.IntField("start")
}

// For iterates over every element in a collection and calls the
// provided function for each entry.  If the function returns a
// non-nil err value, the iteration will stop.  Watch out for
// very large collections!
func (v *Value) For(f func(*Value) os.Error) os.Error {
	for {
		entries, ok := v.Map()["entries"].([]interface{})
		if !ok {
			return os.NewError("No entries found in value")
		}
		for _, entry := range entries {
			m, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			url_, _ := m["self_link"].(string)
			err := f(&Value{session: v.session, baseurl: v.baseurl, url: url_, m: m})
			if err != nil {
				return err
			}
		}
		nextv, _ := v.Link("next_collection_link")
		if nextv == nil {
			break
		}
		err := nextv.Get(nil)
		if err != nil {
			return err
		}
		v = nextv
	}
	return nil
}

var stopRedir = os.NewError("Stop redirection marker.")

var httpClient = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) os.Error { return stopRedir },
}

func (v *Value) do(method string, params Params, body []byte) (value *Value, err os.Error) {
	value = v
	query := multimap(params).Encode()
	for redirect := 0; ; redirect++ {
		req, err := http.NewRequest(method, value.url, nil)
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

		if v.session != nil {
			err := v.session.Sign(req)
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
			value = &Value{url: v.url, baseurl: v.baseurl, session: v.session}
			if resp.StatusCode == 201 && location != "" {
				value.url = location
				return value.do("GET", nil, nil)
			}
		}

		if method == "GET" && shouldRedirect(resp.StatusCode) {
			if location == "" {
				msg := "Got redirection status " + strconv.Itoa(resp.StatusCode) + " without a Location"
				return nil, os.NewError(msg)
			}
			value.url = location
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
		value.m = make(map[string]interface{})
		return value, json.Unmarshal(body, &value.m)
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
