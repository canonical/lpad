package lpad

import "strconv"

// A BuildState holds the state a package build can be found in
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

// A Pocket represents the various distribution pockets where packages end up
type Pocket string

const (
	PocketAny       Pocket = ""
	PocketRelease   Pocket = "Release"
	PocketSecurity  Pocket = "Security"
	PocketUpdates   Pocket = "Updates"
	PocketProposed  Pocket = "Proposed"
	PocketBackports Pocket = "Backports"
)

// The Build type describes a package build
type Build struct {
	*Value
}

// The BuildList type represents a list of Build objects
type BuildList struct {
	*Value
}

// For calls the function f on each Build in the BuildList
func (bl *BuildList) For(f func(b *Build) error) error {
	return bl.Value.For(func(v *Value) error {
		return f(&Build{v})
	})
}

// Build returns the package build identified by distro, package name and id
func (root *Root) Build(distro string, source string, id int) (*Build, error) {
	v, err := root.Location("/" + distro + "/+source/" + source + "/" + strconv.Itoa(id)).Get(nil)
	if err != nil {
		return nil, err
	}
	return &Build{v}, nil
}

// Title is the title of this build
func (build *Build) Title() string {
	return build.StringField("title")
}

// ArchTag is the architecture of this build
func (build *Build) ArchTag() string {
	return build.StringField("arch_tag")
}

// Retry gives back a failed build to the builder farm
func (build *Build) Retry() error {
	_, err := build.Post(Params{"ws.op": "retry"})
	return err
}

// WebPage is the webpage of this build in Launchpad
func (build *Build) WebPage() string {
	return build.StringField("web_link")
}

// State returns the state of this build (successful, depwait, failed, etc.)
func (build *Build) State() string {
	return build.StringField("buildstate")
}

// BuildLogURL is the URL of the gzipped build log file of this build
func (build *Build) BuildLogURL() string {
	return build.StringField("build_log_url")
}

// UploadLogURL is the URL of upload log if there was an upload failure, None otherwise
func (build *Build) UploadLogURL() string {
	return build.StringField("upload_log_url")
}

// DateCreated is the date when this build was created
func (build *Build) DateCreated() string {
	return build.StringField("datecreated")
}

// DateBuilt is the date when the build finished
func (build *Build) DateBuilt() string {
	return build.StringField("datebuilt")
}

// SelfLink returns the API URL of this build entry
func (build *Build) SelfLink() string {
	return build.StringField("self_link")
}

// The SPPH type holds the source package publication history
type SPPH struct {
	*Value
}

// PackageName gets the source package name of a SPPH
func (s *SPPH) PackageName() string {
	return s.StringField("source_package_name")
}

// PackageVersion gets the  package version of a SPPH
func (s *SPPH) PackageVersion() string {
	return s.StringField("source_package_version")
}

// DistroSeries gets the distro series this is published in
func (s *SPPH) DistroSeries() (*DistroSeries, error) {
	v, err := s.Link("distro_series_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &DistroSeries{v}, nil
}

// Archive gets the archive this is published in
func (s *SPPH) Archive() (*Archive, error) {
	v, err := s.Link("archive_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &Archive{v}, nil
}

// Component gets the component this package was published in (i.e. main,universe)
func (s *SPPH) Component() string {
	return s.StringField("component_name")
}

// CurrentSourcePublicationLink gets the SPPH corresponding to a Build
func (build *Build) CurrentSourcePublicationLink() (*SPPH, error) {
	v, err := build.Link("current_source_publication_link").Get(nil)
	if err != nil {
		return nil, err
	}
	return &SPPH{v}, nil
}

// SPPHList is a list of SPPH objects used for iteration
type SPPHList struct {
	*Value
}

// For iterates over the list of SPPHs and calls f for each one.
func (list *SPPHList) For(f func(s *SPPH) error) error {
	return list.Value.For(func(v *Value) error {
		return f(&SPPH{v})
	})
}
