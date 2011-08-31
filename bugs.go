package lpad

import (
	"os"
	"strconv"
	"strings"
)

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
func (root Root) Bug(id int) (bug Bug, err os.Error) {
	r, err := root.GetLocation("/bugs/" + strconv.Itoa(id))
	return Bug{r}, err
}

// CreateBug creates a new bug with an appropriate bug task and returns it.
func (root Root) CreateBug(stub *BugStub) (bug Bug, err os.Error) {
	params := Params{
		"ws.op": "createBug",
		"title": stub.Title,
		"description": stub.Description,
		"target": stub.Target.URL(),
	}
	if len(stub.Tags) > 0 {
		params["tags"] = strings.Join(stub.Tags, " ")
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

// LinkBranch associates a branch with this bug.
func (bug Bug) LinkBranch(branch Branch) os.Error {
	params := Params{
		"ws.op": "linkBranch",
		"branch": branch.URL(),
	}
	_, err := bug.Post(params)
	return err
}
