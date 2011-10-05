package lpad

import (
	"os"
)

//https://launchpad.net/builders
//Not all info presented on that page is available via the LP API though.

//Stands for an individual machine that builds packages
type Builder struct {
	*Value
}

//Builder returns a builder by its name
func (root Root) Builder(buildername string) (builder Builder, err os.Error) {
	v, err := root.Location("/builders").Get(Params{"ws.op": "getByName", "name": buildername})
	return Builder{v}, err
}

//Name gets the name of a builder
func (builder Builder) Name() string {
	return builder.StringField("name")
}

//Title gets the title of the machine
func (builder Builder) Title() string {
	return builder.StringField("title")
}

//Active says whether the builder is active or not
func (builder Builder) Active() bool {
	return builder.BoolField("active")
}

//BuilderOK says whether the builder is fine
func (builder Builder) BuilderOK() bool {
	return builder.BoolField("builderok")
}

//Virtualized says if it is running under Xen
func (builder Builder) Virtualized() bool {
	return builder.BoolField("virtualized")
}

//VMHost gets the name of the VM host machine
func (builder Builder) VMHost() string {
	return builder.StringField("vm_host")
}

//WebPage is the webpage of this builder in Launchpad
func (build Builder) WebPage() string {
	return build.StringField("web_link")
}

//A list of builders
type BuilderList struct {
	*Value
}

//Builders gets all the builders
func (root Root) Builders() (list BuilderList, err os.Error) {
	v, err := root.Location("/builders").Get(nil)
	return BuilderList{v}, err
}

//For calls a given function on each builder
func (list BuilderList) For(f func(b Builder) os.Error) os.Error {
	return list.Value.For(func(v *Value) os.Error {
		f(Builder{v})
		return nil
	})
}
