package gitServices

import (
	"fmt"
	"testing"
	"time"

	"github.com/SorenHQ/joern-port/models"
)

func TestClone(t *testing.T) {
	gitRepo := "https://github.com/morozovsk/websocket"
	gitchan := make(chan models.GitResponse)
	git := NewGitDetailsHandler(gitchan)
	go git.ClonePull("ws", gitRepo, false)
done:
	for {
		select {
		case response, ok := <-gitchan:
			fmt.Println("GIT PROGRESS LOGG: ", response)
			if !ok {
				fmt.Println("Download Finished: ", response)

				break done
			}
		case <-time.After(time.Minute * 10):
			fmt.Println("finished test")
			break done
		}
	}
	fmt.Println("git clone done")
}
