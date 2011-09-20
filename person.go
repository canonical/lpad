package lpad

import (
	"os"
	"url"
)

// The Root type provides the entrance for the Launchpad API.
type Root struct {
	*Value
}

// Me returns the Person authenticated into Lauchpad in the current session.
func (root Root) Me() (p Person, err os.Error) {
	me, err := root.Location("/people/+me").Get(nil)
	return Person{me}, err
}

// Person returns the Person with the provided username.
func (root Root) Person(username string) (p Person, err os.Error) {
	v, err := root.Location("/~" + url.QueryEscape(username)).Get(nil)
	if err == nil && v.BoolField("is_team") {
		err = os.NewError(username + " is a team, not a person")
	}
	return Person{v}, err
}

// Team returns the Team with the provided name.
func (root Root) Team(name string) (p Team, err os.Error) {
	v, err := root.Location("/~" + url.QueryEscape(name)).Get(nil)
	if err == nil && !v.BoolField("is_team") {
		err = os.NewError(name + " is not a team")
	}
	return Team{v}, err
}

// Member returns the Team or Person with the provided name or username.
func (root Root) Member(name string) (member AnyValue, err os.Error) {
	v, err := root.Location("/~" + url.QueryEscape(name)).Get(nil)
	if err != nil {
		return nil, err
	}
	if v.BoolField("is_team") {
		return Team{v}, nil
	}
	return Person{v}, nil
}

// FindPeople returns a PersonList containing all Person accounts whose
// Name, DisplayName or email address match text.
func (root Root) FindPeople(text string) (list PersonList, err os.Error) {
	v, err := root.Location("/people").Get(Params{"ws.op": "findPerson", "text": text})
	return PersonList{v}, err
}

// FindTeams returns a TeamList containing all Team accounts whose
// Name, DisplayName or email address match text.
func (root Root) FindTeams(text string) (list TeamList, err os.Error) {
	v, err := root.Location("/people").Get(Params{"ws.op": "findTeam", "text": text})
	return TeamList{v}, err
}

// FindMembers returns a MemberList containing all Person or Team accounts
// whose Name, DisplayName or email address match text.
func (root Root) FindMembers(text string) (list MemberList, err os.Error) {
	v, err := root.Location("/people").Get(Params{"ws.op": "find", "text": text})
	return MemberList{v}, err
}

// The MemberList type encapsulates a mixed list containing Person and Team
// elements for iteration.
type MemberList struct {
	*Value
}

// For iterates over the list of people and teams and calls f for each one.
// If f returns a non-nil error, iteration will stop and the error will be
// returned as the result of For.
func (list MemberList) For(f func(v AnyValue) os.Error) os.Error {
	return list.Value.For(func(v *Value) os.Error {
		if v.BoolField("is_team") {
			f(Team{v})
		} else {
			f(Person{v})
		}
		return nil
	})
}

// The PersonList type encapsulates a list of Person elements for iteration.
type PersonList struct {
	*Value
}

// For iterates over the list of people and calls f for each one.  If f
// returns a non-nil error, iteration will stop and the error will be
// returned as the result of For.
func (list PersonList) For(f func(p Person) os.Error) os.Error {
	return list.Value.For(func(v *Value) os.Error {
		f(Person{v})
		return nil
	})
}

// The TeamList type encapsulates a list of Team elements for iteration.
type TeamList struct {
	*Value
}

// For iterates over the list of teams and calls f for each one.  If f
// returns a non-nil error, iteration will stop and the error will be
// returned as the result of For.
func (list TeamList) For(f func(t Team) os.Error) os.Error {
	return list.Value.For(func(v *Value) os.Error {
		f(Team{v})
		return nil
	})
}

// The Person type represents a person in Launchpad.
type Person struct {
	*Value
}

// DisplayName returns the person's name as it would be displayed
// throughout Launchpad.  Most people use their full name.
func (person Person) DisplayName() string {
	return person.StringField("display_name")
}

// WebPage returns the URL for accessing this person's page in a browser.
func (person Person) WebPage() string {
	return person.StringField("web_link")
}

// SetDisplayName changes the person's name as it would be displayed
// throughout Launchpad.  Most people use their full name.
// Patch must be called to commit all changes.
func (person Person) SetDisplayName(name string) {
	person.SetField("display_name", name)
}

// IRCNicks returns a list of all IRC nicks for the person.
func (person Person) IRCNicks() (nicks []IRCNick, err os.Error) {
	list, err := person.Link("irc_nicknames_collection_link").Get(nil)
	if err != nil {
		return nil, err
	}
	list.For(func(v *Value) os.Error {
		nicks = append(nicks, IRCNick{v})
		return nil
	})
	return
}

type IRCNick struct {
	*Value
}

// Nick returns the person's nick on an IRC network.
func (nick IRCNick) Nick() string {
	return nick.StringField("nickname")
}

// SetNick changes the person's nick on an IRC network.
// Patch must be called to commit all changes.
func (nick IRCNick) SetNick(n string) {
	nick.SetField("nickname", n)
}

// Network returns the IRC network this nick is associated to.
func (nick IRCNick) Network() string {
	return nick.StringField("network")
}

// SetNetwork changes the IRC network this nick is associated to.
// Patch must be called to commit all changes.
func (nick IRCNick) SetNetwork(n string) {
	nick.SetField("network", n)
}

// The Team type encapsulates access to details about a team in Launchpad.
type Team struct {
	*Value
}

// Name returns the team's name.  This is a short unique name, beginning with a
// lower-case letter or number, and containing only letters, numbers, dots,
// hyphens, or plus signs.
func (team Team) Name() string {
	return team.StringField("name")
}

// SetName changes the team's name.  This is a short unique name, beginning
// with a lower-case letter or number, and containing only letters, numbers,
// dots, hyphens, or plus signs.  Patch must be called to commit all changes.
func (team Team) SetName(name string) {
	team.SetField("name", name)
}

// DisplayName returns the team's name as it would be displayed
// throughout Launchpad.
func (team Team) DisplayName() string {
	return team.StringField("display_name")
}

// SetDisplayName changes the team's name as it would be displayed
// throughout Launchpad.  Patch must be called to commit all changes.
func (team Team) SetDisplayName(name string) {
	team.SetField("display_name", name)
}

// WebPage returns the URL for accessing this team's page in a browser.
func (team Team) WebPage() string {
	return team.StringField("web_link")
}
