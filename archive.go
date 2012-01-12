package lpad

// Archive represents a package archive
type Archive struct {
	*Value
}

// Name returns the name of this archive
func (a *Archive) Name() string {
	return a.StringField("name")
}

// DisplayName returns an user friendly name of this archive
func (a *Archive) DisplayName() string {
	return a.StringField("displayname")
}

// Description returns a description string for this archive
func (a *Archive) Description() string {
	return a.StringField("description")
}

// Distro returns the distribution that uses this archive
func (a *Archive) Distro() (*Distro, error) {
	v, err := a.Link("distribution_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &Distro{v}, nil
}

// SelfLink returns the API URL of this Archive object
func (a *Archive) SelfLink() string {
	return a.StringField("self_link")
}

// WebPage returns the URL for accessing this archive in a browser.
func (a *Archive) WebPage() string {
	return a.StringField("web_link")
}

// GetPublishedSources returns a list of SPPH records of this archive
// that match the given criteria
func (a *Archive) GetPublishedSources(source string) (*SPPHList, error) {
	params := Params{
		"ws.op":       "getPublishedSources",
		"source_name": source,
		"exact_match": "true",
		"pocket":      "Release",
		"status":      "Published",
	}
	v, err := a.Location("").Get(params)
	if err != nil {
		return nil, err
	}
	return &SPPHList{v}, nil
}

// ArchiveList is a list of Archive objects used for iterating
type ArchiveList struct {
	*Value
}

//For iterates over a list of archives and calls a function on each
func (list *ArchiveList) For(f func(a *Archive) error) error {
	return list.Value.For(func(v *Value) error {
		return f(&Archive{v})
	})
}
