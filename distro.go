package lpad

import (
	"os"
	"time"
	"url"
)

// Distro returns a distribution with the given name.
func (root Root) Distro(name string) (distribution Distro, err os.Error) {
	r, err := root.Location("/" + url.QueryEscape(name)).Get(nil)
	return Distro{r}, err
}

// The Distro type represents a distribution in Launchpad.
type Distro struct {
	*Value
}

// Name returns the distribution name, which is composed of at least one
// lowercase letter or number, followed by letters, numbers, dots,
// hyphens or pluses. This is a short name used in URLs.
func (d Distro) Name() string {
	return d.StringField("name")
}

// DisplayName returns the distribution name as it would be displayed
// in a paragraph. For example, a distribution's title might be
// "The Foo Distro" and its display name could be "Foo".
func (d Distro) DisplayName() string {
	return d.StringField("display_name")
}

// Title returns the distribution title as it might be used in isolation.
func (d Distro) Title() string {
	return d.StringField("title")
}

// Summary returns the distribution summary, which is a short paragraph
// to introduce the distribution's goals and highlights.
func (d Distro) Summary() string {
	return d.StringField("summary")
}

// Description returns the distribution description.
func (d Distro) Description() string {
	return d.StringField("description")
}

// WebPage returns the URL for accessing this distribution in a browser.
func (d Distro) WebPage() string {
	return d.StringField("web_link")
}

// SetName changes the distribution name, which must be composed of at
// least one lowercase letter or number, followed by letters, numbers,
// dots, hyphens or pluses. This is a short name used in URLs.
// Patch must be called to commit all changes.
func (d Distro) SetName(name string) {
	d.SetField("name", name)
}

// SetDisplayName changes the distribution name as it would be displayed
// in a paragraph. For example, a distribution's title might be
// "The Foo Distro" and its display name could be "Foo".
// Patch must be called to commit all changes.
func (d Distro) SetDisplayName(name string) {
	d.SetField("display_name", name)
}

// SetTitle changes the distribution title as it would be displayed
// in isolation. For example, the distribution title might be
// "The Foo Distro" and display name could be "Foo".
// Patch must be called to commit all changes.
func (d Distro) SetTitle(title string) {
	d.SetField("title", title)
}

// SetSummary changes the distribution summary, which is a short paragraph
// to introduce the distribution's goals and highlights.
// Patch must be called to commit all changes.
func (d Distro) SetSummary(title string) {
	d.SetField("summary", title)
}

// SetDescription changes the distributions's description.
// Patch must be called to commit all changes.
func (d Distro) SetDescription(description string) {
	d.SetField("description", description)
}

// ActiveMilestones returns the list of active milestones associated with
// the distribution, ordered by the target date.
func (d Distro) ActiveMilestones() (milestones MilestoneList, err os.Error) {
	r, err := d.Link("active_milestones_collection_link").Get(nil)
	return MilestoneList{r}, err
}

// AllSeries returns the list of series associated with the distribution.
func (d Distro) AllSeries() (series DistroSeriesList, err os.Error) {
	r, err := d.Link("series_collection_link").Get(nil)
	return DistroSeriesList{r}, err
}

// FocusDistroSeries returns the distribution series set as the current
// development focus.
func (d Distro) FocusSeries() (series DistroSeries, err os.Error) {
	r, err := d.Link("current_series_link").Get(nil)
	return DistroSeries{r}, err
}


// The DistroSeries type represents a series associated with a distribution.
type DistroSeries struct {
	*Value
}

// Name returns the series name, which is a unique name that identifies
// it and is used in URLs. It consists of only lowercase letters, digits,
// and simple punctuation.  For example, "2.0" or "trunk".
func (s DistroSeries) Name() string {
	return s.StringField("name")
}

// Title returns the series context title for pages.
func (s DistroSeries) Title() string {
	return s.StringField("title")
}

// Summary returns the summary for this distribution series.
func (s DistroSeries) Summary() string {
	return s.StringField("summary")
}

// WebPage returns the URL for accessing this distribution series in a browser.
func (s DistroSeries) WebPage() string {
	return s.StringField("web_link")
}

// Active returns true if this distribution series is still in active development.
func (s DistroSeries) Active() bool {
	return s.BoolField("is_active")
}

// SetName changes the series name, which must consists of only letters,
// numbers, and simple punctuation. For example: "2.0" or "trunk".
func (s DistroSeries) SetName(name string) {
	s.SetField("name", name)
}

// SetTitle changes the series title.
func (s DistroSeries) SetTitle(title string) {
	s.SetField("title", title)
}

// SetSummary changes the summary for this distribution series.
func (s DistroSeries) SetSummary(summary string) {
	s.SetField("summary", summary)
}

// SetActive sets whether the series is still in active development or not.
func (s DistroSeries) SetActive(active bool) {
	s.SetField("is_active", active)
}

// The DistroSeriesList represents a list of distribution series.
type DistroSeriesList struct {
	*Value
}

// For iterates over the list of series and calls f for each one.
// If f returns a non-nil error, iteration will stop and the error will
// be returned as the result of For.
func (list DistroSeriesList) For(f func(s DistroSeries) os.Error) os.Error {
	return list.Value.For(func(r *Value) os.Error {
		f(DistroSeries{r})
		return nil
	})
}
