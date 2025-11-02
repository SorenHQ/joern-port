package gitServices

// import (
// 	"fmt"
// 	"joern-output-parser/models"
// 	"testing"

// 	"github.com/go-git/go-git/v6"
// )

// func TestClone(t *testing.T) {
// 	gitRepo := "https://github.com/qxresearch/qxresearch-event-1"
// 	CloneRepo("py-mini-ev1",gitRepo,nil)

// }


// func TestPull(t *testing.T) {
// 	gitRepo := "py-mini-ev1"
// 	err:=PullRepo(gitRepo,nil)
// 	if err != nil && err != git.NoErrAlreadyUpToDate {
// 		return 
// 	}


// 	status := "success"
// 	msg := "Repository updated successfully"
// 	if err == git.NoErrAlreadyUpToDate {
// 		status = "up-to-date"
// 		msg = "Already up to date"
// 	}
	
// 	fmt.Println(models.GitResponse{Action: "pull", Status: status, Message: msg})

// }
