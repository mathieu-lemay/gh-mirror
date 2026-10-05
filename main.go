package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/sean9999/hermeti"
	"golang.org/x/sync/errgroup"
)

var _ hermeti.Runner = (*state)(nil)

type state struct {
	// Where the mirror should live
	rootDir string
	// A list of selected organizations. If the list is not empty, only the selected organizations will be synced.
	selectedOrgs []string
	// A list of excluded organizations. Any organization in this list will be skipped.
	excludedOrgs []string
	// Number of threads to use for syncing
	threads int
}

func (a *state) Run(env *hermeti.Env) {
	ctx := context.Background()

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

		orgDir := strings.Join([]string{a.rootDir, org.Name}, string(os.PathSeparator))
		fmt.Fprintf(env.OutStream, "Syncing organization '%s' to '%s'\n", org.Name, orgDir)

		repos, err := org.Repos(env)
		if err != nil {
			fmt.Fprintf(env.ErrStream, "Failed to get repos for org %s: %v\n", org.Name, err)
			continue
		}

		eg, _ := errgroup.WithContext(ctx)
		eg.SetLimit(a.threads)

		for _, repo := range repos {
			eg.Go(func() error {
				fmt.Fprintf(env.OutStream, "Syncing repository: %s/%s\n", org.Name, repo.Name)
				myDir := strings.Join([]string{a.rootDir, org.Name, repo.Name}, string(os.PathSeparator))
				err := EnsureSynced(env, repo, myDir)
				if err != nil {
					fmt.Fprintln(env.ErrStream, err)
				}

				return nil
			})
		}

		err = eg.Wait()
		if err != nil {
			fmt.Fprintf(env.ErrStream, "Failed to sync org %s: %v\n", org.Name, err)
		}
	}
}

func (a *state) shouldSyncOrg(org Org) bool {
	if slices.Contains(a.excludedOrgs, org.Name) {
		return false
	}

	return len(a.selectedOrgs) == 0 || slices.Contains(a.selectedOrgs, org.Name)
}

func (a *state) parseArgs(args []string) {
	flagSet := flag.NewFlagSet("args", flag.PanicOnError)

	var selectedOrgs Strings
	flagSet.Var(&selectedOrgs, "with-org", "comma separated list of orgs to sync, can be specified multiple times")

	var excludedOrgs Strings
	flagSet.Var(&excludedOrgs, "without-org", "comma separated list of orgs to skip, can be specified multiple times")

	var threads int
	flagSet.IntVar(&threads, "threads", 4, "number of threads to use for syncing")

	flagSet.Parse(args)

	rootDir := flagSet.Arg(0)
	if rootDir == "" {
		rootDir = "."
	}

	a.rootDir = rootDir
	a.selectedOrgs = selectedOrgs
	a.excludedOrgs = excludedOrgs
	a.threads = threads
}

func main() {
	app := new(state)
	cli := hermeti.NewRealCli(app)
	cli.Run()
}
