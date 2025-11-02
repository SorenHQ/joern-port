package gitServices

import (
	"errors"
	"fmt"
	"log"

	"joern-output-parser/env"
	"joern-output-parser/models"
	"os"

	"github.com/go-git/go-git/v6"
)

var gitService *GitService

type GitService struct{
	detailChan chan models.GitResponse
}
// impl io.Writer
func (gs *GitService) Write(p []byte) (n int, err error) {
	if gs.detailChan == nil {
		return 0, errors.New("detail channel not initialized")
	}
	gs.detailChan<- models.GitResponse{Action: "log", Status: "success", Message: string(p)}
	return len(p), nil
}


func (gs *GitService)ClonePull(project, url string, pull bool)  {
	if gs.detailChan == nil {
		gs.detailChan = make(chan models.GitResponse)
		defer close(gs.detailChan)
	}
	base := env.GetProjectReposPath()
	dir := fmt.Sprintf("%s/%s", base, project)
	// Check if repo already exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Clone new repo
		repo, err := git.PlainClone(dir, &git.CloneOptions{
			URL:      url,
			Progress: gs,
		})
		if err != nil {
			gs.detailChan<- models.GitResponse{Action: "clone", Status: "error", Error: err.Error()}
			return
		}

		ref, _ := repo.Head()
		gs.detailChan<- models.GitResponse{
			Action:   "clone",
			Status:   "success",
			Branch:   ref.Name().Short(),
			CommitID: ref.Hash().String(),
			Message:  "Repository cloned successfully",
		}
		return
	}

	// Open existing repo and pull latest
	repo, err := git.PlainOpen(dir)
	if err != nil {
		gs.detailChan<- models.GitResponse{Action: "open", Status: "error", Error: err.Error()}
		return
	}
	if pull {
		w, _ := repo.Worktree()
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			gs.detailChan<- models.GitResponse{Action: "pull", Status: "error", Error: err.Error()}
			return
		}

		ref, _ := repo.Head()
		status := "success"
		msg := "Repository updated successfully"
		if err == git.NoErrAlreadyUpToDate {
			status = "up-to-date"
			msg = "Already up to date"
		}

		gs.detailChan<- models.GitResponse{
			Action:   "pull",
			Status:   status,
			Branch:   ref.Name().Short(),
			CommitID: ref.Hash().String(),
			Message:  msg,
		}
	}
	ref, _ := repo.Head()
	status := "success"
	msg := "Repository updated successfully"
	gs.detailChan<- models.GitResponse{
		Action:   "repo",
		Status:   status,
		Branch:   ref.Name().Short(),
		CommitID: ref.Hash().String(),
		Message:  msg,
	}
}


func NewGitDetailsHandler(detailChan chan models.GitResponse) *GitService {
	if gitService==nil{
		gitService = &GitService{detailChan: detailChan}
	}
	log.Default().Println("Git service initialized")
	return gitService
}

func GitClonePull(dir, url string, pull bool) error{
	if gitService==nil{
		return errors.New("git service not initialized")
	}
	go gitService.ClonePull(dir, url, pull)
	return nil
}
