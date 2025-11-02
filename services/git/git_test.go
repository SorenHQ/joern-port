package gitServices

import (
	"fmt"
	"testing"
	"time"

	"github.com/SorenHQ/joern-port/models"
)

func TestClone(t *testing.T) {
	gitRepo := "https://github.com/qxresearch/qxresearch-event-1"
	gitchan := make(chan models.GitResponse)
	NewGitDetailsHandler(gitchan)
	err := GitClonePull("python-ev-project", gitRepo, false)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	for {
		select {
		case response := <-gitchan:
			fmt.Println("GIT PROGRESS LOGG: ", response)
		case <-time.After(time.Second * 10):
			fmt.Println("finished test")
			return
		}
	}

}
