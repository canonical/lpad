package lpad

import (
	"bytes"
	"http"
	"io/ioutil"
	"json"
	"os"
	"path"
	"strconv"
	"strings"
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

	// Location returns a new resource for a location which may be a
	// full URL, or an absolute path (based on the resource's BaseURL),
	// or a path relative to the resource itself (based on the
	// resource's URL).
	Location(url string) Resource

	// Link calls Location with a URL available in the given key
	// of the current resource's Map.  This is a convenient way to
	// navigate through *_link values.
	Link(key string) Resource

	// Get issues an HTTP GET to retrieve the content of this resource.
	// If params is not nil, it will provided as the query for the GET
	// request.
	Get(params Params) os.Error

	// Post issues an HTTP POST to perform a given action at the URL
	// specified by this resource.  If params is not nil, it will
	// provided as the parameters for the POST request.
	Post(params Params) (Resource, os.Error)
}

type resource struct {
	session *Session
	baseurl string
	url     string
	m       map[string]interface{}
}

// NewResource creates a new resource type.  Creating resources explicitly
// is generally not necessary.  If you're trying to access a location in
// the Launchpad API which is not covered by the static model yet, see the
// Link and Location methods on the Resource interface for more convenient
// ways to create resources.
func NewResource(session *Session, baseurl, url string, m map[string]interface{}) Resource {
	return &resource{session, baseurl, url, m}
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
	url, err := http.ParseURL(base)
	if err != nil {
		panic("Invalid URL: " + base)
	}
	url.Path = path.Join(url.Path, part)
	return url.String()
}

func (r *resource) Location(url string) Resource {
	return &resource{session: r.session, baseurl: r.baseurl, url: r.join(url)}
}

func (r *resource) Link(key string) Resource {
	location := r.m[key]
	if location == nil {
		panic(`Link("` + key + `") for ` + r.url + ` is nil.`)
	}
	return r.Location(location.(string))
}

func (r *resource) Get(params Params) os.Error {
	_, err := r.do("GET", params)
	return err
}

func (r *resource) Post(params Params) (other Resource, err os.Error) {
	return r.do("POST", params)
}

func (r *resource) do(method string, params map[string]string) (other Resource, err os.Error) {
	query := http.EncodeQuery(multimap(params))
	for redirect := 0; ; redirect++ {
		req, err := http.NewRequest(method, r.url, nil)
		req.Header["Accept"] = []string{"application/json"}
		if err != nil {
			return nil, err
		}

		if method == "POST" {
			req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(query)))
			req.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
			req.Header["Content-Length"] = []string{strconv.Itoa(len(query))}
			req.ContentLength = int64(len(query))

		} else {
			req.URL.RawQuery = query
		}

		if r.session != nil {
			err := r.session.Sign(req)
			if err != nil {
				return nil, err
			}
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		//dump, err := http.DumpResponse(resp, false)
		//println("Response:", string(dump))

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		location := resp.Header.Get("Location")

		if method == "POST" {
			// When posting do not follow redirects but instead
			// return a resource pointing to the new resource.
			other := resource{url: r.url, session: r.session}
			if location != "" {
				other.url = location
			}
			return &other, nil
		}

		if shouldRedirect(resp.StatusCode) && method == "GET" {
			if location == "" {
				msg := "Got redirection status " + strconv.Itoa(resp.StatusCode) + " without a Location"
				return nil, os.NewError(msg)
			}
			r.url = location
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return nil, &Error{resp.StatusCode, body}
		}

		ctype := resp.Header.Get("Content-Type")
		if ctype != "application/json" {
			return nil, os.NewError("Non-JSON content-type: " + ctype)
		}

		if err != nil {
			return nil, err
		}
		r.m = make(map[string]interface{})
		return nil, json.Unmarshal(body, &r.m)
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

func multimap(params map[string]string) map[string][]string {
	m := make(map[string][]string, len(params))
	for k, v := range params {
		m[k] = []string{v}
	}
	return m
}
