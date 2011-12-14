package lpad

type SourcePackage struct {
	*Value
}

// Name returns the name of the package
func (s *SourcePackage) Name() string {
	return s.StringField("name")
}

// DisplayName returns the display name of the package
func (s *SourcePackage) DisplayName() string {
	return s.StringField("displayname")
}

// Component is the name of the component where the source was last published
func (s *SourcePackage) Component() string {
	return s.StringField("latest_published_component_name")
}

// WebPage is the webpage of this source package in Launchpad
func (s *SourcePackage) WebPage() string {
	return s.StringField("web_link")
}

// SelfLink returns the API URL of this source package entry
func (s *SourcePackage) SelfLink() string {
	return s.StringField("self_link")
}

// Distro returns the distribution of this souce package
func (s *SourcePackage) Distro() (*Distro, error) {
	d, err := s.Link("distribution_link").Get(nil)
	if err != nil {
		return nil, err
	}

	return &Distro{d}, nil
}

// DistroSeries returns the distribution of this souce package
func (s *SourcePackage) DistroSeries() (*DistroSeries, error) {
	d, err := s.Link("distroseries_link").Get(nil)
	if err != nil {
		return nil, err
	}

	return &DistroSeries{d}, nil
}

type DistroSourcePackage struct {
	*Value
}

// Name returns the name of the package
func (s *DistroSourcePackage) Name() string {
	return s.StringField("name")
}

// DisplayName returns the display name of the package
func (s *DistroSourcePackage) DisplayName() string {
	return s.StringField("display_name")
}

// Title is the name of the component where the source was last published
func (s *DistroSourcePackage) Title() string {
	return s.StringField("title")
}

// WebPage is the webpage of this source package in Launchpad
func (s *DistroSourcePackage) WebPage() string {
	return s.StringField("web_link")
}

// SelfLink returns the API URL of this source package entry
func (s *DistroSourcePackage) SelfLink() string {
	return s.StringField("self_link")
}

// Distro returns the distribution of this souce package
func (s *DistroSourcePackage) Distro() (*Distro, error) {
	d, err := s.Link("distribution_link").Get(nil)
	if err != nil {
		return nil, err
	}

	return &Distro{d}, nil
}
