package lpad

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
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

func (e *Error) Error() string {
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
	BaseLoc() string
	AbsLoc() string
	Map() map[string]interface{}
	StringField(key string) string
	IntField(key string) int
	FloatField(key string) float64
	BoolField(key string) bool
	SetField(key string, value interface{})
	Location(loc string) *Value
	Link(key string) *Value
	Get(params Params) (*Value, error)
	Post(params Params) (*Value, error)
	Patch() error
	TotalSize() int
	StartIndex() int
	For(func(v *Value) error) error
}

// The Value type is the underlying dynamic layer used as the foundation of
// all the more specific value types that support the Launchpad model.
// Besides being used internally to implement these types, the methods of
// this type also enable accessing new features available in Launchpad which
// were not yet made available in lpad thorugh more convenient methods.
type Value struct {
	session *Session
	baseloc string
	loc     string
	m       map[string]interface{}
	patch   map[string]interface{}
}

// NewValue creates a new Value with the provided details. Creating values
// explicitly is generally not necessary.  If you're trying to access a
// location in the Launchpad API which is not covered by the supported
// types yet, see the Link and Location methods on the Value type for more
// convenient ways to create values.
func NewValue(session *Session, baseloc, loc string, m map[string]interface{}) *Value {
	return &Value{session, baseloc, loc, m, nil}
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

// BaseLoc returns the API-oriented URL base for the session. Absolute
// paths provided to Location and Link will be rooted at this place.
func (v *Value) BaseLoc() string {
	return v.baseloc
}

// AbsLoc returns the API-oriented URL of this value.
func (v *Value) AbsLoc() string {
	if self := v.StringField("self_link"); self != "" {
		return self
	}
	return v.loc
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
		return v.loc
	}
	if strings.HasPrefix(part, "http://") || strings.HasPrefix(part, "https://") {
		return part
	}
	base := v.baseloc
	if !strings.HasPrefix(part, "/") {
		base = v.loc
	}
	u, err := url.Parse(base)
	if err != nil {
		panic("Invalid URL: " + base)
	}
	u.Path = path.Join(u.Path, part)
	return u.String()
}

// Location returns a new value for a location which may be a
// full URL, or an absolute path (based on the value's BaseLoc),
// or a path relative to the value itself (based on the
// value's URL).
func (v *Value) Location(loc string) *Value {
	return &Value{session: v.session, baseloc: v.baseloc, loc: v.join(loc)}
}

// Link calls Location with a URL available in the given key
// of the current value's Map.  It returns nil if the requested
// key isn't found in the value.  This is a convenient way to
// navigate through *_link fields in values.
func (v *Value) Link(key string) *Value {
	link, ok := v.m[key].(string)
	if !ok {
		return nil
	}
	return v.Location(link)
}

// Error returned when a Get, Post or Patch operation is done on
// a nil value.
var InvalidValue = errors.New("Invalid value")

// Get issues an HTTP GET to retrieve the content of this value,
// and returns itself and an error in case of problems. If params
// is not nil, it will provided as the query for the GET request.
//
// Since Get returns the value itself, it may be used as:
//
//     v, err := other.Link("some_link").Get(nil)
//
func (v *Value) Get(params Params) (same *Value, err error) {
	return v.do("GET", params, nil)
}

// Post issues an HTTP POST to perform a given action at the URL
// specified by this value.  If params is not nil, it will
// provided as the parameters for the POST request.
func (v *Value) Post(params Params) (other *Value, err error) {
	return v.do("POST", params, nil)
}

// Patch issues an HTTP PATCH request to modify the server value
// with the local changes.
func (v *Value) Patch() error {
	if v == nil {
		return InvalidValue
	}
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
func (v *Value) For(f func(*Value) error) (err error) {
	for {
		entries, ok := v.Map()["entries"].([]interface{})
		if !ok {
			return errors.New("No entries found in value")
		}
		for _, entry := range entries {
			m, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			link, _ := m["self_link"].(string)
			err := f(&Value{session: v.session, baseloc: v.baseloc, loc: link, m: m})
			if err != nil {
				return err
			}
		}
		nextv := v.Link("next_collection_link")
		if nextv == nil {
			break
		}
		v, err = nextv.Get(nil)
		if err != nil {
			return err
		}
	}
	return nil
}

var stopRedir = errors.New("Stop redirection marker.")

var httpClient = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error { return stopRedir },
}

func (v *Value) do(method string, params Params, body []byte) (value *Value, err error) {
	if v == nil {
		return nil, InvalidValue
	}
	value = v
	query := multimap(params).Encode()
	for redirect := 0; ; redirect++ {
		req, err := http.NewRequest(method, value.AbsLoc(), nil)
		if err != nil {
			return nil, err
		}
		req.Header["Accept"] = []string{"application/json"}

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

		if debugOn {
			printDump(httputil.DumpRequest(req, false))
		}

		resp, err := httpClient.Do(req)
		if urlerr, ok := err.(*url.Error); ok && urlerr.Err == stopRedir {
			// fine
		} else if err != nil {
			return nil, err
		}

		if debugOn {
			printDump(httputil.DumpResponse(resp, false))
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		location := resp.Header.Get("Location")

		if method == "POST" {
			value = &Value{baseloc: v.baseloc, session: v.session}
			if resp.StatusCode == 201 && location != "" {
				value.loc = location
				return value.do("GET", nil, nil)
			} else {
				value.loc = v.AbsLoc()
			}
		}

		if method == "GET" && shouldRedirect(resp.StatusCode) {
			if location == "" {
				msg := "Got redirection status " + strconv.Itoa(resp.StatusCode) + " without a Location"
				return nil, errors.New(msg)
			}
			value.loc = location
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
			return nil, errors.New("Non-JSON content-type: " + ctype)
		}

		if err != nil {
			return nil, err
		}
		value.m = make(map[string]interface{})
		if len(body) > 0 && body[0] == '[' {
			body = append([]byte(`{"value":`), body...)
			body = append(body, '}')
		}
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

var debugOn bool

// SetDebug enables debugging. With debug on requests and responses will all be
// dumped into the standard error output.
func SetDebug(debug bool) {
	debugOn = debug
}

func printDump(data []byte, err error) {
	s := string(data)
	if err != nil {
		s = err.Error()
	}
	fmt.Fprintf(os.Stderr, "===== DEBUG =====\n%s\n=================\n", s)
}
