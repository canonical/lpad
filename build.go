package lpad

import (
	"os"
	"strconv"
)

//The various states a package build can be found in
type BuildState string

const (
	BSNeedsBuilding            BuildState = "Needs building"
	BSSuccessfullyBuilt        BuildState = "Successfully built"
	BSFailedToBuild            BuildState = "Failed to build"
	BSDependencyWait           BuildState = "Dependency wait"
	BSChrootProblem            BuildState = "Chroot problem"
	BSBuildForSupersededSource BuildState = "Build for superseded source"
	BSCurrentlyBuilding        BuildState = "Currently building"
	BSFailedToUpload           BuildState = "Failed to upload"
	BSCurrentlyUploading       BuildState = "Currently uploading"
)

//The various distribution pockets where packages end up
type Pocket string

const (
	PocketAny       Pocket = ""
	PocketRelease   Pocket = "Release"
	PocketSecurity  Pocket = "Security"
	PocketUpdates   Pocket = "Updates"
	PocketProposed  Pocket = "Proposed"
	PocketBackports Pocket = "Backports"
)

//A Build describes a package build
type Build struct {
	*Value
}

//A BuildList is a list of Build objects
type BuildList struct {
	*Value
}

//Calls the function fn on each Build in the BuildList
func (bl BuildList) For(fn func(b Build) os.Error) os.Error {
	return bl.Value.For(func(v *Value) os.Error {
		fn(Build{v})
		return nil
	})
}

//Build returns the package build identified by distro, package name and id
func (root Root) Build(distro string, source string, id int) (build Build, err os.Error) {
	v, err := root.Location("/" + distro + "/+source/" + source + "/" + strconv.Itoa(id)).Get(nil)
	return Build{v}, err
}

//Title is the title of this build
func (build Build) Title() string {
	return build.StringField("title")
}

//ArchTag is the architecture of this build
func (build Build) ArchTag() string {
	return build.StringField("arch_tag")
}

//Retry gives back a failed build to the builder farm
func (build Build) Retry() os.Error {
	_, err := build.Post(Params{"ws_op": "retry"})
	return err
}

//WebPage is the webpage of this build in Launchpad
func (build Build) WebPage() string {
	return build.StringField("web_link")
}

//BuildLogURL is the URL of the gzipped build log file of this build
func (build Build) BuildLogURL() string {
	return build.StringField("build_log_url")
}

//DateCreated is the date when this build was created
func (build Build) DateCreated() string {
	return build.StringField("datecreated")
}

//DateBuilt is the date when the build finished
func (build Build) DateBuilt() string {
	return build.StringField("datebuilt")
}

//Source package publication history
type SPPH struct {
	*Value
}

//Gets the source package name of a SPPH
func (s SPPH) PackageName() string {
	return s.StringField("source_package_name")
}

//Gets the  package version of a SPPH
func (s SPPH) PackageVersion() string {
	return s.StringField("source_package_version")
}

//Gets the SPPH corresponding to a Build
func (build Build) CurrentSourcePublicationLink() (spph SPPH, err os.Error) {
	v, err := build.Link("current_source_publication_link").Get(nil)
	return SPPH{v}, err
}
