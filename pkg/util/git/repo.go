package git

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func Parse(rawurl string) (owner string, repo string, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", "", err
	}

	path := u.Path
	splits := strings.Split(path, "/")
	if len(splits) < 3 {
		return "", "", errors.New("invalid repository url format.")
	}

	owner = splits[1]
	repo = splits[2]

	repo = RemoveGitSuffix(repo)

	return owner, repo, nil
}

func RemoveGitSuffix(repo string) string {
	if strings.HasSuffix(repo, ".git") {
		repo = strings.Replace(repo, ".git", "", 1)
	}

	return repo
}
