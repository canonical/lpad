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
	me, err := root.GetLocation("/people/+me")
	return Person{me}, err
}

// The Person type encapsulates access to details about a person in Launchpad.
type Person struct {
	Resource
}

// DisplayName returns the person's name as it would be displayed
// throughout Launchpad.  Most people use their full name.
func (person Person) DisplayName() string {
	return person.StringField("display_name")
}

func (person Person) IRCNicks() (nicks []IRCNick, err os.Error) {
	list, err := person.GetLink("irc_nicknames_collection_link")
	if err != nil {
		return nil, err
	}
	list.ListIter(func(r Resource) os.Error {
		nicks = append(nicks, IRCNick{r})
		return nil
	})
	return
}

type IRCNick struct {
	Resource
}

func (nick IRCNick) Nick() string {
	return nick.StringField("nickname")
}

func (nick IRCNick) Network() string {
	return nick.StringField("network")
}
