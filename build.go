package lpad

import (
	"fmt"
	"net/url"
)

// A BuildState holds the state a package build can be found in.
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

// A Pocket represents the various distribution pockets where packages end up.
type Pocket string

const (
	PocketAny       Pocket = ""
	PocketRelease   Pocket = "Release"
	PocketSecurity  Pocket = "Security"
	PocketUpdates   Pocket = "Updates"
	PocketProposed  Pocket = "Proposed"
	PocketBackports Pocket = "Backports"
)

// The Build type describes a package build.
type Build struct {
	*Value
}

// The BuildList type represents a list of Build objects.
type BuildList struct {
	*Value
}

// For calls the function f on each Build in the BuildList.
func (bl *BuildList) For(f func(b *Build) error) error {
	return bl.Value.For(func(v *Value) error {
		return f(&Build{v})
	})
}

// Build returns the identified package build.
func (root *Root) Build(distro string, source string, id int) (*Build, error) {
	distro = url.QueryEscape(distro)
	source = url.QueryEscape(source)
	path := fmt.Sprintf("/%s/+source/%s/%d/", distro, source, id)
	v, err := root.Location(path).Get(nil)
	if err != nil {
		return nil, err
	}
	return &Build{v}, nil
}

// Title is the title of build.
func (build *Build) Title() string {
	return build.StringField("title")
}

// Arch is the architecture of build.
func (build *Build) Arch() string {
	return build.StringField("arch_tag")
}

// Retry sends back a failed build to the builder farm.
func (build *Build) Retry() error {
	_, err := build.Post(Params{"ws.op": "retry"})
	return err
}

// WebPage is the webpage of build in Launchpad.
func (build *Build) WebPage() string {
	return build.StringField("web_link")
}

// State returns the state of build.
func (build *Build) State() BuildState {
	return BuildState(build.StringField("buildstate"))
}

// BuildLogURL is the URL of the gzipped build log file of build.
func (build *Build) BuildLogURL() string {
	return build.StringField("build_log_url")
}

// UploadLogURL is the URL of the upload log if there was an upload failure.
func (build *Build) UploadLogURL() string {
	return build.StringField("upload_log_url")
}

// DateCreated is the date when build was created.
func (build *Build) DateCreated() string {
	return build.StringField("datecreated")
}

// DateBuilt is the date when build finished.
func (build *Build) DateBuilt() string {
	return build.StringField("datebuilt")
}

// The PubHistory type holds the source package publication history.
type PubHistory struct {
	*Value
}

// PubHistory gets the the source publication history corresponding to build.
func (build *Build) PubHistory() (*PubHistory, error) {
	v, err := build.Link("current_source_publication_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &PubHistory{v}, nil
}

// PackageName gets the source package name of ph.
func (ph *PubHistory) PackageName() string {
	return ph.StringField("source_package_name")
}

// PackageVersion gets the  package version of ph.
func (ph *PubHistory) PackageVersion() string {
	return ph.StringField("source_package_version")
}

// DistroSeries returns the distro series packages are published into.
func (ph *PubHistory) DistroSeries() (*DistroSeries, error) {
	v, err := ph.Link("distro_series_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &DistroSeries{v}, nil
}

// Archive returns the archive packages are published into.
func (ph *PubHistory) Archive() (*Archive, error) {
	v, err := ph.Link("archive_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &Archive{v}, nil
}

// Component returns the component the packages are published into.
func (ph *PubHistory) Component() string {
	return ph.StringField("component_name")
}

// PubHistoryList represents a list of PubHistory objects.
type PubHistoryList struct {
	*Value
}

// For iterates over the list of PubHistory values and calls f for each one.
func (list *PubHistoryList) For(f func(s *PubHistory) error) error {
	return list.Value.For(func(v *Value) error {
		return f(&PubHistory{v})
	})
}
