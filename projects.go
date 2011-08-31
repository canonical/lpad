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

// ActiveMilestones returns the list of active milestones associated with
// the project, ordered by the target date.
func (p Project) ActiveMilestones() (milestones MilestoneList, err os.Error) {
	r, err := p.GetLink("active_milestones_collection_link")
	return MilestoneList{r}, err
}

// The Milestone type represents a milestone associated with a project.
type Milestone struct {
	Resource
}

// Name returns the milestone name, which consists of only
// letters, numbers, and simple punctuation.
func (ms Milestone) Name() string {
	return ms.StringField("name")
}

// CodeName returns the alternative name for the milestone, if any.
func (ms Milestone) CodeName() string {
	return ms.StringField("code_name")
}

// Title returns the milestone context title for pages.
func (ms Milestone) Title() string {
	return ms.StringField("title")
}

// Summary returns the summary of features and status of this milestone.
func (ms Milestone) Summary() string {
	return ms.StringField("summary")
}

// WebPage returns the web page link associated with this milestone.
func (ms Milestone) WebPage() string {
	return ms.StringField("web_link")
}

// Active returns true if the milestone is still active.
func (ms Milestone) Active() bool {
	return ms.BoolField("is_active")
}

// Date returns the target date for the milestone.
func (ms Milestone) Date() string {
	return ms.StringField("date_targeted")
}

// SetName changes the milestone name, which must consists of
// only letters, numbers, and simple punctuation.
func (ms Milestone) SetName(name string) {
	ms.SetField("name", name)
}

// SetCodeName sets the alternative name for the milestone.
func (ms Milestone) SetCodeName(name string) {
	ms.SetField("code_name", name)
}

// SetTitle changes the milestone's context title for pages.
func (ms Milestone) SetTitle(title string) {
	ms.SetField("title", title)
}

// SetSummary sets the summary of features and status of this milestone.
func (ms Milestone) SetSummary(summary string) {
	ms.SetField("summary", summary)
}

// SetWebPage sets the web page link associated with this milestone.
func (ms Milestone) SetWebPage(link string) {
	ms.SetField("web_link", link)
}

// SetActive sets whether the milestone is still active or not.
func (ms Milestone) SetActive(active bool) {
	ms.SetField("is_active", active)
}

// SetDate changes the target date for the milestone.
func (ms Milestone) SetDate(date string) {
	ms.SetField("date_targeted", date)
}

// The MilestoneList represents a list of project milestones.
type MilestoneList struct {
	Resource
}

// For iterates over the list of milestones and calls f for each one.
// If f returns a non-nil error, iteration will stop and the error will
// be returned as the result of For.
func (list MilestoneList) For(f func(t Milestone) os.Error) os.Error {
	return list.Resource.For(func(r Resource) os.Error {
		f(Milestone{r})
		return nil
	})
}
