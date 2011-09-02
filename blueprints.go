package lpad

import (
	"os"
)

// The Blueprint type represents a blueprint in Launchpad.
type Blueprint struct {
	*Value
}

// Name returns the blueprint name. May contain lower-case letters, numbers,
// and dashes. It will be used in the specification url.
// Examples: mozilla-type-ahead-find, postgres-smart-serial.
func (bp Blueprint) Name() string {
	return bp.StringField("name")
}

// SetName changes the blueprint name which must consist of lower-case
// letters, numbers, and dashes. It will be used in the specification url.
// Examples: mozilla-type-ahead-find, postgres-smart-serial.
// Patch must be called to commit all changes.
func (bp Blueprint) SetName(name string) {
	bp.SetField("name", name)
}

// Title returns the blueprint title that should describe the feature
// as clearly as possible, in up to 70 characters. This title is
// displayed in every feature list or report.
func (bp Blueprint) Title() string {
	return bp.StringField("title")
}

// SetTitle sets the blueprint title.  The title must describe the feature
// as clearly as possible, in up to 70 characters. This title is displayed
// in every feature list or report.
func (bp Blueprint) SetTitle(title string) {
	bp.SetField("title", title)
}

// Summary returns the blueprint summary which should consist of a single
// paragraph description of the feature.
func (bp Blueprint) Summary() string {
	return bp.StringField("summary")
}

// SetSummary changes the blueprint summary which must consist of a single
// paragraph description of the feature.
func (bp Blueprint) SetSummary(summary string) {
	bp.SetField("summary", summary)
}

// Whiteboard returns the blueprint whiteboard which contains any notes
// on the status of this specification.
func (bp Blueprint) Whiteboard() string {
	return bp.StringField("whiteboard")
}

// SetWhiteboard changes the blueprint whiteboard that may contain any
// notes on the status of this specification.
func (bp Blueprint) SetWhiteboard(whiteboard string) {
	bp.SetField("whiteboard", whiteboard)
}

// WebPage returns the URL for accessing this blueprint in a browser.
func (bp Blueprint) WebPage() string {
	return bp.StringField("web_link")
}

// LinkBranch associates a branch with this blueprint.
func (bp Blueprint) LinkBranch(branch Branch) os.Error {
	params := Params{
		"ws.op":  "linkBranch",
		"branch": branch.AbsLoc(),
	}
	_, err := bp.Post(params)
	return err
}
