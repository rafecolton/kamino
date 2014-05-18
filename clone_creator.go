package kamino

import (
	"fmt"
	"net/url"
	"os/exec"

	"github.com/modcloth/go-fileutils"
)

type cloneCreator struct {
	g       *genome
	workdir string
}

func (creator *cloneCreator) cachePath() string {
	return fmt.Sprintf("%s/%s/%s", creator.workdir, creator.g.Account, creator.g.Repo)
}

func (creator *cloneCreator) cloneCacheIfAvailable() (string, error) {
	if err := creator.updateToRef(creator.cachePath()); err != nil {
		return creator.cloneNoCache()
	}

	return creator.cachePath(), nil
}

func (creator *cloneCreator) cloneForceCache() (string, error) {
	if err := creator.updateToRef(creator.cachePath()); err != nil {
		return "", err
	}

	return creator.cachePath(), nil
}

func (creator *cloneCreator) cloneCreateCache() (string, error) {
	if err := creator.cloneRepo(creator.cachePath()); err != nil {
		return "", err
	}

	return creator.cachePath(), nil
}

func (creator *cloneCreator) cloneNoCache() (string, error) {
	uuid, err := nextUUID()
	if err != nil {
		return "", err
	}

	clonePath := fmt.Sprintf("%s/%s/%s", creator.workdir, creator.g.Account, uuid)

	if err = creator.cloneRepo(clonePath); err != nil {
		return "", err
	}

	return clonePath, nil
}

func (creator *cloneCreator) cloneRepo(dest string) error {
	repoURL := &url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   fmt.Sprintf("%s/%s", creator.g.Account, creator.g.Repo),
	}

	if creator.g.APIToken != "" {
		repoURL.User = url.User(creator.g.APIToken)
	}

	cloneCmd := exec.Command(
		"git", "clone",
		"--quiet",
		"--depth", creator.g.Depth,
		repoURL.String(),
		dest,
	)

	if err := cloneCmd.Run(); err != nil {
		return err
	}

	git, err := fileutils.Which("git")
	if err != nil {
		return err
	}

	checkoutCmd := &exec.Cmd{
		Path: git,
		Dir:  dest,
		Args: []string{"git", "checkout", "-qf", creator.g.Ref},
	}

	if err := checkoutCmd.Run(); err != nil {
		fmt.Printf("GOT HERE, err = %q, dest = %q\n", err, dest)
		return err
	}

	return nil
}

func (creator *cloneCreator) updateToRef(dest string) error {
	/*
		workflow as follows:
			git reset --hard
			git clean -df
			git fetch
			git pull --rebase || true # (in case it's a branch)
			git checkout -qf <ref>
	*/
	git, err := fileutils.Which("git")
	if err != nil {
		return err
	}

	cmds := []*exec.Cmd{
		&exec.Cmd{
			Path: git,
			Dir:  dest,
			Args: []string{"git", "reset", "--hard"},
		},
		&exec.Cmd{
			Path: git,
			Dir:  dest,
			Args: []string{"git", "clean", "-df"},
		},
		&exec.Cmd{
			Path: git,
			Dir:  dest,
			Args: []string{"git", "fetch"},
		},
	}

	for _, cmd := range cmds {
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	cmd := &exec.Cmd{
		Path: git,
		Dir:  dest,
		Args: []string{"git", "pull", "--rebase"},
	}

	cmd.Run() // ignore failure, since our ref may not be a branch

	checkout := &exec.Cmd{
		Path: git,
		Dir:  dest,
		Args: []string{"git", "checkout", "-qf", creator.g.Ref},
	}

	return checkout.Run()
}
