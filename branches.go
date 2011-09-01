package lpad

import (
	"os"
	"strings"
	"url"
)

// Branch returns a branch for the provided id. The id should be
// prefixed with lp: as conventional for Launchpad.
func (root Root) Branch(id string) (project Project, err os.Error) {
	if !strings.HasPrefix(id, "lp:") {
		err = os.NewError("Invalid branch id provided: " + id)
		return
	}
	// Let's avoid a silly query injection issue here.
	parts := strings.Split(id, "/")
	for i, part := range parts {
		parts[i] = url.QueryEscape(part)
	}
	id = strings.Join(parts, "/")
	r, err := root.GetLocation("/" + id)
	return Project{r}, err
}

// The Branch type represents a project in Launchpad.
type Branch struct {
	Resource
}

// Id returns the shortest version for the branch name. If the branch
// is the development focus for a project, a lp:project form will be
// returned. If it's the development focus for a series, then a
// lp:project/series is returned. Otherwise, the unique name for the
// branch in the form lp:~user/project/branch-name is returned.
func (b Branch) Id() string {
	return b.StringField("bzr_identity")
}

// UniqueName returns the unique branch name, in the
// form lp:~user/project/branch-name.
func (b Branch) UniqueName() string {
	return b.StringField("unique_name")
}

type MergeStub struct {
	Description string
	CommitMessage string
	NeedsReview bool
	Target Branch
	PreReq Branch
}

// ProposeMerge proposes this branch for merging on another branch by
// creating the respective merge proposal.
func (b Branch) ProposeMerge(stub *MergeStub) (mp MergeProposal, err os.Error) {
	if stub.Target.Resource == nil {
		err = os.NewError("Missing target branch")
	}
	params := Params{
		"ws.op": "createMergeProposal",
		"target_branch": stub.Target.URL(),
	}
	if stub.Description != "" {
		params["initial_comment"] = stub.Description
	}
	if stub.CommitMessage != "" {
		params["commit_message"] = stub.CommitMessage
	}
	if stub.NeedsReview {
		params["needs_review"] = "true"
	}
	if stub.PreReq.Resource != nil {
		params["prerequisite_branch"] = stub.PreReq.URL()
	}
	r, err := b.Post(params)
	return MergeProposal{r}, err
}

type MergeProposal struct {
	Resource
}

// Description returns the merge proposal introductory comment.
func (mp MergeProposal) Description() string {
	return mp.StringField("description")
}

// Status returns the current status of the merge proposal.
// E.g. Needs review, Work In Progress, etc.
func (mp MergeProposal) Status() string {
	return mp.StringField("queue_status")
}

// CommitMessage returns the commit message to be used when merging
// the proposal.
func (mp MergeProposal) CommitMessage() string {
	return mp.StringField("commit_message")
}


// Email returns the unique email that may be used to add new comments
// to the merge proposal conversation.
func (mp MergeProposal) Email() string {
	return mp.StringField("address")
}

// Source returns the source branch that has additional code to land.
func (mp MergeProposal) Source() (branch Branch, err os.Error) {
	r, err := mp.GetLink("source_branch_link")
	return Branch{r}, err
}

// Target returns the branch where code will land on once merged.
func (mp MergeProposal) Target() (branch Branch, err os.Error) {
	r, err := mp.GetLink("target_branch_link")
	return Branch{r}, err
}

// PreReq returns the branch is the base (merged or not) for the code
// within the target branch.
func (mp MergeProposal) PreReq() (branch Branch, err os.Error) {
	r, err := mp.GetLink("prerequisite_branch_link")
	return Branch{r}, err
}

// WebPage returns the URL for accessing this merge proposal
// in a browser.
func (mp MergeProposal) WebPage() string {
	return mp.StringField("web_link")
}
