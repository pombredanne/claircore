package claircore

type Package struct {
	// unique ID of this package. this will be created as discovered by the library
	// and used for persistence and hash map indexes
	ID int `json:"id"`
	// the name of the distribution
	Name string `json:"name"`
	// the version of the distribution
	Version string `json:"version"`
	// type of package. currently expectations are binary or source
	Kind string `json:"kind"`
	// if type is a binary package a source package maybe present which built this binary package.
	// must be a pointer to support recursive type:
	Source *Package `json:"source"`
	// the file system path or prefix where this package resides
	PackageDB string `json:"package_db"`
	// a hint on which repository this package was downloaded from
	RepositoryHint string `json:"repository_hint"`
}
