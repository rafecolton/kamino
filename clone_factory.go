package kamino

import "os"

/*
A CloneFactory generates your clones for you.  Create a clone factory with NewCloneFactory().
*/
type CloneFactory struct {
	workdir string
}

/*
NewCloneFactory creates a new CloneFactory, ready to do some cloning for you
(into its specified workdir).
*/
func NewCloneFactory(workdir string) (*CloneFactory, error) {

	if err := os.MkdirAll(workdir, 0755); err != nil {
		return nil, err
	}

	return &CloneFactory{
		workdir: workdir,
	}, nil
}

/*
Clone clones the repo as specified by the genome parameters.
*/
func (factory *CloneFactory) Clone(g *genome) (path string, err error) {
	creator := &cloneCreator{
		g:       g,
		workdir: factory.workdir,
	}

	switch g.UseCache {
	case "no":
		return creator.cloneNoCache()
	case "create":
		return creator.cloneCreateCache()
	case "force":
		return creator.cloneForceCache()
	case "if_available":
		return creator.cloneCacheIfAvailable()
	default:
		return creator.cloneNoCache()
	}
}
