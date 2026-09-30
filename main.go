package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/sean9999/hermeti"
)

var _ hermeti.Runner = (*state)(nil)

type state struct {
	// Where the mirror should live
	rootDir      string
	// A list of selected organizations. If the list is not empty, only the selected organizations will be synced.
	selectedOrgs []string
	// A list of excluded organizations. Any organization in this list will be skipped.
	excludedOrgs []string
}

func (a *state) Run(env *hermeti.Env) {
	a.parseArgs(env.Args[1:])

	err := EnsureDir(env, a.rootDir)
	if err != nil {
		fmt.Fprintln(env.ErrStream, "You must pass in a valid directory")
		panic(err)
	}

	orgs := GetOrgs()

	for _, org := range orgs {
		if !a.shouldSyncOrg(org) {
			fmt.Fprintf(env.OutStream, "Skipping organization: %s\n", org.Name)
			continue
		}

		fmt.Fprintf(env.OutStream, "Syncing organization: %s\n", org.Name)

		repos, err := org.Repos(env)
		if err != nil {
			panic(err)
		}
		for _, repo := range repos {
			fmt.Fprintf(env.OutStream, "Syncing repository: %s/%s\n", org.Name, repo.Name)
			myDir := strings.Join([]string{a.rootDir, org.Name, repo.Name}, string(os.PathSeparator))
			err := EnsureSynced(env, repo, myDir)
			if err != nil {
				fmt.Fprintln(env.ErrStream, err)
			}
		}
	}
}

func (a *state) shouldSyncOrg(org Org) bool {
	if slices.Contains(a.excludedOrgs, org.Name) {
		return false
	}

	return len(a.selectedOrgs) == 0 || slices.Contains(a.selectedOrgs, org.Name)
}

func (a *state) parseArgs() {
	var selectedOrgs Strings
	flag.Var(&selectedOrgs, "with-org", "comma separated list of orgs to sync, can be specified multiple times")

	var excludedOrgs Strings
	flag.Var(&excludedOrgs, "without-org", "comma separated list of orgs to skip, can be specified multiple times")

	flag.Parse()

	rootDir := flag.Arg(0)
	if rootDir == "" {
		rootDir = "."
	}

	a.rootDir = rootDir
	a.selectedOrgs = selectedOrgs
	a.excludedOrgs = excludedOrgs
}

func main() {
	app := new(state)
	cli := hermeti.NewRealCli(app)
	cli.Run()
}
