package main

import (
	"errors"

	"github.com/armed/mkdirp"
	"github.com/nodefortytwo/isgit"
	"github.com/sean9999/hermeti"
)

func sync(env *hermeti.Env, dir string) error {
	return runCli(env.OutStream, env.ErrStream, &dir, "repo", "sync")
}

func clone(env *hermeti.Env, repo Repo, dir string) error {
	return runCli(env.OutStream, env.ErrStream, nil, "repo", "clone", repo.SshUrl, dir)
}

// EnsureSynced ensures a folder is a git repo and is synced to upstream,
// cloning if necessary.
func EnsureSynced(env *hermeti.Env, repo Repo, dir string) error {
	err := EnsureDir(env, dir)
	if err != nil {
		return err
	}
	isRepo, err := isgit.IsGitRepo(dir)
	if err != nil {
		return err
	}
	if isRepo {
		return sync(env, dir)
	}
	return clone(env, repo, dir)
}

// EnsureDir ensures a directory exists by creating it or making sure it's already there.
func EnsureDir(env *hermeti.Env, dir string) error {
	info, err := env.Filesystem.Stat(dir)
	if err != nil {
		return mkdirp.Mk(dir, 0755)
	}
	if info.IsDir() == false {
		return errors.New("not a dir")
	}
	return nil
}
