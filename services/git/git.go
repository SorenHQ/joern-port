package gitServices

import (
	"fmt"
	"log"

	"os"

	"github.com/SorenHQ/joern-port/env"
	"github.com/SorenHQ/joern-port/models"

	"github.com/go-git/go-git/v6"
)


type GitService struct {
	detailChan chan models.GitResponse
}

// impl io.Writer
func (gs *GitService) Write(p []byte) (n int, err error) {
	if gs.detailChan == nil {
		return len(p), nil // Return success but don't send if channel is nil
	}
	// Use recover to handle panics from sending to closed channels
	defer func() {
		if r := recover(); r != nil {
			// Channel was closed, ignore the error and continue
			log.Printf("Git progress channel closed, skipping message: %v", r)
		}
	}()
	// Use non-blocking send to avoid deadlock if channel is full
	select {
	case gs.detailChan <- models.GitResponse{Action: "log", Status: "success", Message: string(p)}:
		// Successfully sent
	default:
		// Channel is full, skip this message to avoid blocking
		// The receiver might not be reading fast enough
	}
	return len(p), nil
}

func (gs *GitService) ClonePull(project, url string, pull bool) {
	if gs.detailChan == nil {
		gs.detailChan = make(chan models.GitResponse) 
	}
	defer close(gs.detailChan)
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
			gs.detailChan <- models.GitResponse{Action: "clone", Status: "error", Error: err.Error()}
			return
		}

		ref, _ := repo.Head()
		gs.detailChan <- models.GitResponse{
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
		gs.detailChan <- models.GitResponse{Action: "open", Status: "error", Error: err.Error()}
		return
	}
	if pull {
		w, _ := repo.Worktree()
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			gs.detailChan <- models.GitResponse{Action: "pull", Status: "error", Error: err.Error()}
			return
		}

		ref, _ := repo.Head()
		status := "success"
		msg := "Repository updated successfully"
		if err == git.NoErrAlreadyUpToDate {
			status = "up-to-date"
			msg = "Already up to date"
		}

		gs.detailChan <- models.GitResponse{
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
	gs.detailChan <- models.GitResponse{
		Action:   "repo",
		Status:   status,
		Branch:   ref.Name().Short(),
		CommitID: ref.Hash().String(),
		Message:  msg,
	}
}

func NewGitDetailsHandler(detailChan chan models.GitResponse) *GitService {
    gitService :=&GitService{detailChan: detailChan}
	
	log.Default().Println("Git service initialized")
	return gitService
}

func NewGitLogHandler(detailChan chan models.GitResponse) *GitService {
	if detailChan == nil {
		log.Default().Println("git service not initialized")
		return  nil
	}
	log.Default().Println("Git Logger service initialized")
	return 	 &GitService{detailChan: detailChan}

}
