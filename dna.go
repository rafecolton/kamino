package kamino

import (
	"errors"
	"fmt"
	"strconv"
)

type genome struct {
	APIToken string
	Account  string
	Depth    string
	Repo     string
	UseCache string
}

/*
NewGenome creates a new genome type, which is the options struct for a CloneFactory.

Valid fields for opts are as follows:

	* "account"
		- purpose: GitHub account
		- required: true
	* "cache"
		- purpose: whether or not to use the cached / previously cloned version of the repo
		- required: false
		- default: "no"
		- valid options:
			* "no" - do not use cache, create a uniquely named directory
			* "if_available" - use cache if already created, otherwise create a uniquely named directory
			* "create" - use cache if already created, create cache if not present
			* "force" - use cache if already created, fail if cache not present
	* "depth"
		- purpose: git clone `--depth` option
		- required: false
		- default: "50"
		- validation: must be empty string or parsable as a base 10 integer
	* "repo"
		- purpose: GitHub repo
		- required: true
	* "token"
		- purpose: GitHub API token for private repos
		- required: false (functionally required if your repo is private)
		- default: (not sent with request if empty)

*/
func NewGenome(opts map[string]string) (*genome, error) {
	g := &genome{}

	if depth, ok := opts["depth"]; ok && depth != "" {
		if _, err := strconv.Atoi(depth); err == nil {
			g.Depth = depth
		} else {
			return nil, fmt.Errorf("%q is not a valid clone depth", depth)
		}
	} else {
		g.Depth = "50"
	}

	if token, ok := opts["token"]; ok {
		g.APIToken = token
	}

	if account, ok := opts["account"]; ok && account != "" {
		g.Account = account
	} else {
		return nil, errors.New("account must be provided")
	}

	if repo, ok := opts["repo"]; ok && repo != "" {
		g.Repo = repo
	} else {
		return nil, errors.New("repo must be provided")
	}

	if cache, ok := opts["cache"]; ok && cache != "" {
		if validCacheOption(cache) {
			g.UseCache = cache
		} else {
			return nil, fmt.Errorf("%q is not a valid cache option", cache)
		}
	} else {
		g.UseCache = "no"
	}

	return g, nil
}

func validCacheOptions() []string {
	return []string{
		"create",
		"force",
		"if_available",
		"no",
	}
}

func validCacheOption(opt string) bool {
	for _, v := range validCacheOptions() {
		if v == opt {
			return true
		}
	}

	return false
}
