// This simple example demonstrates how to get started using lpad to
// communicate with Launchpad in a console application:
//
//     root, err := lpad.Login(lpad.Production, &lpad.ConsoleOAuth{})
//     if err != nil {
//         panic(err)
//     }
//     me, err := root.Me()
//     if err != nil {
//         panic(err)
//     }
//     fmt.Println(me.DisplayName())
//
package lpad

import (
	"http"
	"os"
)

// The Auth interface is implemented by types which are able to Login and
// authenticate requests made against Launchpad.
type Auth interface {
	Login(baseURL string) (err os.Error)
	Sign(req *http.Request) (err os.Error)
}

const Production = "https://api.launchpad.net/1.0/"
const Staging = "https://api.staging.launchpad.net/1.0/"


// The Session type represents a session of communication with Launchpad,
// and carries the authenticator necessary to validate requests in the
// given session.  Creating sessions explicitly is generally not necessary.
// See the Login method for a convenient way to use lpad to access the
// Launchpad API.
type Session struct {
	auth Auth
}

// Create a new session using the auth authenticator.  Creating sessions
// explicitly is generally not necessary.  See the Login method for a
// convenient way to use lpad to access the Launchpad API.
func NewSession(auth Auth) *Session {
	return &Session{auth}
}

func (s *Session) Sign(req *http.Request) (err os.Error) {
	return s.auth.Sign(req)
}

// Login returns a Root object with a new session authenticated in Launchpad
// using the auth authenticator.  This is the primary method to start using
// the Launchpad API.
//
// This simple example demonstrates how to get a user's name in a console
// application:
//
//     root, err := lpad.Login(lpad.Production, &lpad.ConsoleOAuth{})
//     if err != nil {
//         panic(err)
//     }
//     me, err := root.Me()
//     if err != nil {
//         panic(err)
//     }
//     fmt.Println(me.DisplayName())
//
func Login(baseurl string, auth Auth) (root Root, err os.Error) {
	err = auth.Login(baseurl)
	if err != nil {
		return
	}
	return Root{&resource{session: NewSession(auth), baseurl: baseurl, url: baseurl}}, nil
}
