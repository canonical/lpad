package lpad

import (
	"os"
	"url"
)

// Project returns a project with the given name.
func (root Root) Project(name string) (project Project, err os.Error) {
	r, err := root.GetLocation("/" + url.QueryEscape(name))
	return Project{r}, err
}

// The Project type represents a project in Launchpad.
type Project struct {
	Resource
}

// Name returns the project name, which is composed of at least one
// lowercase letter or number, followed by letters, numbers, dots,
// hyphens or pluses. This is a short name used in URLs.
func (p Project) Name() string {
	return p.StringField("name")
}

// DisplayName returns the project name as it would be displayed
// in a paragraph.  For example, a project's title might be
// "The Foo Project" and its display name could be "Foo".
func (p Project) DisplayName() string {
	return p.StringField("display_name")
}

// Title returns the project title as it might be used in isolation.
// For example, a project's title might be "The Foo Project" and its
// display name could be "Foo".
func (p Project) Title() string {
	return p.StringField("title")
}

// Summary returns the project summary, which is a short paragraph
// to introduce the project's work.
func (p Project) Summary() string {
	return p.StringField("summary")
}

// Description returns the project description.
func (p Project) Description() string {
	return p.StringField("description")
}

// SetName changes the project name, which must be composed of at
// least one lowercase letter or number, followed by letters, numbers,
// dots, hyphens or pluses. This is a short name used in URLs.
// Patch must be called to commit all changes.
func (p Project) SetName(name string) {
	p.SetField("name", name)
}

// SetDisplayName changes the project name as it would be displayed
// in a paragraph. For example, a project's title might be
// "The Foo Project" and its display name could be "Foo".
// Patch must be called to commit all changes.
func (p Project) SetDisplayName(name string) {
	p.SetField("display_name", name)
}

// SetTitle changes the project title as it would be displayed
// in isolation. For example, the project title might be
// "The Foo Project" and display name could be "Foo".
// Patch must be called to commit all changes.
func (p Project) SetTitle(title string) {
	p.SetField("title", title)
}

// SetSummary changes the project summary, which is a short paragraph
// to introduce the project's work.
// Patch must be called to commit all changes.
func (p Project) SetSummary(title string) {
	p.SetField("summary", title)
}

// SetDescription changes the project's description.
// Patch must be called to commit all changes.
func (p Project) SetDescription(description string) {
	p.SetField("description", description)
}
