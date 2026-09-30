package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/cli/go-gh/v2"
	"github.com/sean9999/hermeti"
)

type Org struct {
	Name string
}

// Repos fetches all repos in an Org
func (o Org) Repos(env *hermeti.Env) ([]Repo, error) {
	buf := bytes.NewBuffer(nil)
	err := runCli(buf, env.ErrStream, "repo", "list", o.Name, "--limit", "5000", "--json", "name,url,sshUrl,id")
	if err != nil {
		return nil, err
	}
	repos := make([]Repo, 0, 64)
	err = json.Unmarshal(buf.Bytes(), &repos)
	for i := range repos {
		repos[i].Org = o
	}
	return repos, err
}

func GetOrgs() []Org {
	args := []string{"org", "list"}
	stdOut, _, err := gh.Exec(args...)
	if err != nil {
		log.Fatal(err)
	}
	orgs := orgsFromString(stdOut.String())
	return orgs
}

func orgsFromString(s string) []Org {
	names := strings.Split(s, "\n")
	orgs := make([]Org, 0, len(names))
	for _, name := range names {
		if len(name) > 0 {
			orgs = append(orgs, Org{name})
		}
	}
	return orgs
}
