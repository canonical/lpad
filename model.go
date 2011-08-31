package lpad

import (
	"os"
	"strings"
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

// A BugStub holds details necessary for creating a new bug in Launchpad.
type BugStub struct {
	Title           string   // Required
	Description     string   // Required
	Target          Resource // Project, source package, or distribution
	Private         bool
	SecurityRelated bool
	Tags            []string
}

// CreateBug creates a new bug with an appropriate bug task and returns it.
func (root Root) CreateBug(stub *BugStub) (bug Bug, err os.Error) {
	params := Params{
		"ws.op": "createBug",
		"title": stub.Title,
		"description": stub.Description,
		"target": stub.Target.URL(),
		"tags": strings.Join(stub.Tags, " "),
	}
	if stub.Private {
		params["private"] = "true"
	}
	if stub.SecurityRelated {
		params["security_related"] = "true"
	}
	r, err := root.Location("/bugs").Post(params)
	return Bug{r}, err
}

// FindPeople returns a PersonList containing all Person accounts whose
// Name, DisplayName or email address match text.
func (root Root) FindPeople(text string) (list PersonList, err os.Error) {
	list = PersonList{root.Location("/people")}
	err = list.Get(Params{"ws.op": "findPerson", "text": text})
	return
}

// FindTeams returns a TeamList containing all Team accounts whose
// Name, DisplayName or email address match text.
func (root Root) FindTeams(text string) (list TeamList, err os.Error) {
	list = TeamList{root.Location("/people")}
	err = list.Get(Params{"ws.op": "findTeam", "text": text})
	return
}

// FindMembers returns a MemberList containing all Person or Team accounts
// whose Name, DisplayName or email address match text.
func (root Root) FindMembers(text string) (list MemberList, err os.Error) {
	list = MemberList{root.Location("/people")}
	err = list.Get(Params{"ws.op": "find", "text": text})
	return
}

// The MemberList type encapsulates a mixed list containing Person and Team
// elements for iteration.
type MemberList struct {
	Resource
}

// For iterates over the list of people and teams and calls f for each one.
// If f returns a non-nil error, iteration will stop and the error will be
// returned as the result of For.
func (list MemberList) For(f func(r Resource) os.Error) os.Error {
	return list.Resource.For(func(r Resource) os.Error {
		if r.BoolField("is_team") {
			f(Team{r})
		} else {
			f(Person{r})
		}
		return nil
	})
}

// The PersonList type encapsulates a list of Person elements for iteration.
type PersonList struct {
	Resource
}

// For iterates over the list of people and calls f for each one.  If f
// returns a non-nil error, iteration will stop and the error will be
// returned as the result of For.
func (list PersonList) For(f func(p Person) os.Error) os.Error {
	return list.Resource.For(func(r Resource) os.Error {
		f(Person{r})
		return nil
	})
}

// The TeamList type encapsulates a list of Team elements for iteration.
type TeamList struct {
	Resource
}

// For iterates over the list of teams and calls f for each one.  If f
// returns a non-nil error, iteration will stop and the error will be
// returned as the result of For.
func (list TeamList) For(f func(t Team) os.Error) os.Error {
	return list.Resource.For(func(r Resource) os.Error {
		f(Team{r})
		return nil
	})
}

// The Person type represents a person in Launchpad.
type Person struct {
	Resource
}

// DisplayName returns the person's name as it would be displayed
// throughout Launchpad.  Most people use their full name.
func (person Person) DisplayName() string {
	return person.StringField("display_name")
}

// SetDisplayName changes the person's name as it would be displayed
// throughout Launchpad.  Most people use their full name.
// Patch must be called to commit all changes.
func (person Person) SetDisplayName(name string) {
	person.SetField("display_name", name)
}

// IRCNicks returns a list of all IRC nicks for the person.
func (person Person) IRCNicks() (nicks []IRCNick, err os.Error) {
	list, err := person.GetLink("irc_nicknames_collection_link")
	if err != nil {
		return nil, err
	}
	list.For(func(r Resource) os.Error {
		nicks = append(nicks, IRCNick{r})
		return nil
	})
	return
}

type IRCNick struct {
	Resource
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
	Resource
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

// The Bug type represents a bug in Launchpad.
type Bug struct {
	Resource
}

// Id returns the bug numeric identifier (the bug # itself).
func (bug Bug) Id() int {
	return bug.IntField("id")
}

// Title returns the short bug summary.
func (bug Bug) Title() string {
	return bug.StringField("title")
}

// Description returns the main bug description.
func (bug Bug) Description() string {
	return bug.StringField("description")
}

// Tags returns the set of tags associated with the bug.
func (bug Bug) Tags() []string {
	return strings.Split(bug.StringField("tags"), " ")
}

// Private returns true if the bug is flagged as private.
func (bug Bug) Private() bool {
	return bug.BoolField("private")
}

// SecurityRelated returns true if the bug describes sensitive
// information about a security vulnerability.
func (bug Bug) SecurityRelated() bool {
	return bug.BoolField("security_related")
}

// SetTitle changes the bug title.
// Patch must be called to commit all changes.
func (bug Bug) SetTitle(title string) {
	bug.SetField("title", title)
}

// SetDescription changes the bug description.
// Patch must be called to commit all changes.
func (bug Bug) SetDescription(description string) {
	bug.SetField("description", description)
}

// SetTags changes the bug tags.
// Patch must be called to commit all changes.
func (bug Bug) SetTags(tags []string) {
	bug.SetField("tags", strings.Join(tags, " "))
}

// SetPrivate changes the bug private flag.
// Patch must be called to commit all changes.
func (bug Bug) SetPrivate(private bool) {
	bug.SetField("private", private)
}

// SetSecurityRelated sets to related the flag that tells if
// a bug is security sensitive or not.
// Patch must be called to commit all changes.
func (bug Bug) SetSecurityRelated(related bool) {
	bug.SetField("security_related", related)
}
