package git

import "github.com/go-git/go-git/v5"

func CloneGitRepo(repository, destination string) error {
	_, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL: repository,
	})
	if err != nil {
		return err
	}

	return nil
}
