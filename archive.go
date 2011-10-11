package lpad

import (
	"os"
)

type Archive struct {
	*Value
}

// Name of this archive
func (a Archive) Name() string {
	return a.StringField("name")
}

// User friendly name of this archive
func (a Archive) DisplayName() string {
	return a.StringField("displayname")
}

// Name of this archive
func (a Archive) Description() string {
	return a.StringField("description")
}

// The distribution that uses this archive
func (a Archive) Distribution() Distro {
	v, _ := a.Link("distribution_link").Get(nil)
	return Distro{v}
}

// URL of this archive
func (a Archive) SelfLink() string {
	return a.StringField("self_link")
}

// WebPage returns the URL for accessing this archive in a browser.
func (a Archive) WebPage() string {
	return a.StringField("web_link")
}

func (a Archive) GetPublishedSources(source string) (spph SPPHList, err os.Error) {
    p := Params{"ws.op" : "getPublishedSources",
                "source_name":source,
                "exact_match":"true",
                "pocket":"Release",
                "status":"Published"}
    v, err := a.Location(a.SelfLink()).Get(p)
    return SPPHList{v}, err
}

type ArchiveList struct {
	*Value
}

//For iterates over a list of archives and calls a function on each
func (list ArchiveList) For(fn func(a Archive) os.Error) os.Error {
	return list.Value.For(func(v *Value) os.Error {
		fn(Archive{v})
		return nil
	})
}
