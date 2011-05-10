package lpad

import (
	"os"
)

// The Root type provides the entrance for the Launchpad API.
type Root struct {
	Resource
}

// Me returns the Person authenticated into Lauchpad in the current session.
func (root Root) Me() (p Person, err os.Error) {
	me := root.Location("/people/+me")
	return Person{me}, me.Get(nil)
}

// The Person type encapsulates access to details about a person in Launchpad.
type Person struct {
	Resource
}

// DisplayName returns the person's name as it would be displayed
// throughout Launchpad.  Most people use their full name.
func (person Person) DisplayName() string {
	if v, ok := person.Map()["display_name"].(string); ok {
		return v
	}
	return ""
}
