package main

import (
	"errors"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	gitPath := getGitPath()
	repos, err := git.PlainOpen(gitPath)
	if err != nil {
		log.Fatal("Repos does not exist at ", gitPath)
	}

	headRef, err := repos.Head()
	if err != nil {
		log.Fatal("no commit yet on this repos")
	}
	log.Println("Head: ", headRef)

	commit, err := repos.CommitObject(headRef.Hash())

	release := getReleaseFrom(commit.Message)

	version, err := getLastVersion(repos)
	if err != nil {
		log.Fatal(err)
	}

	updateVersion(version, release)
	createNewVersion(version, repos)

}

func createNewVersion(version *semver.Version, repos *git.Repository) error {
	head, err := repos.Head()
	if err != nil {
		log.Printf("get HEAD error: %s", err)
		return err
	}

	_, err = repos.CreateTag(version.String(), head.Hash(), &git.CreateTagOptions{
		Message: "Automaticaly generated version",
		Tagger: &object.Signature{
			Name:  "Auto versionner",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})

	if err != nil {
		log.Printf("create tag error: %s", err)
		return err
	}
	log.Println("New version tagged is ", version)
	return nil
}

func updateVersion(version *semver.Version, release release) {
	switch release {
	case patch:
		version.IncrementPatch()
	case minor:
		version.IncrementMinor()
	case major:
		version.IncrementMajor()
	}

}

type release int

const (
	patch = iota
	minor
	major
)

func getReleaseFrom(message string) release {
	if isCommitMessageMajor(message) {
		return major
	}
	if isCommitMessageMinor(message) {
		return minor
	}
	return patch
}

func isCommitMessageMinor(message string) bool {
	re := regexp.MustCompile(`MINOR release`)
	return re.Match([]byte(message))
}

func isCommitMessageMajor(message string) bool {
	re := regexp.MustCompile(`MAJOR release`)
	return re.Match([]byte(message))
}

func getLastVersion(repos *git.Repository) (*semver.Version, error) {
	var result semver.Version
	var tag *object.Tag
	tagrefs, err := repos.Tags()
	if err != nil {
		log.Println("No tags found, creating version 0.0.1")
		result = semver.MustParse("0.0.1")
		return &result, nil
	}
	tagrefs.ForEach(func(t *plumbing.Reference) error {
		tagName := t.Name()
		version, err := semver.Parse(tagName.Short())
		if err != nil {
			return nil
		}
		if version.GT(result) {
			result = version
			tag, _ = repos.TagObject(t.Hash())
		}
		return nil
	})
	headref, _ := repos.Head()
	if tag.Target == headref.Hash() {
		return nil, errors.New("Head is allready tagged " + result.String())
	}
	return &result, nil
}

func getGitPath() string {
	args := os.Args[1:]
	switch len(args) {
	case 0:
		return "."
	case 1:
		return args[0]
	default:
		panic("args must be git root")

	}
}
