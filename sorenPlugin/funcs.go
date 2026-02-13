package joernPlugin

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	joernMessHandler "github.com/SorenHQ/joern-port/actions/joern/ws"

	"github.com/SorenHQ/joern-port/env"
	"github.com/SorenHQ/joern-port/etc"
	"github.com/SorenHQ/joern-port/models"
	"github.com/bytedance/sonic"
	sdkv2 "github.com/sorenhq/go-plugin-sdk/gosdk"
	sdkModel "github.com/sorenhq/go-plugin-sdk/gosdk/models"
)

func openJoernProject(jobId string, projectName string) {

	dir := fmt.Sprintf("%s/%s", env.GetProjectReposPath(), projectName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		sdkv2.GetPlugin().Done(jobId, map[string]any{"error": "Project not found"})
		return
	}
	url := fmt.Sprintf("%s/query-sync", env.GetJoernUrl())
	importresult, err := etc.ImportCode(context.Background(), dir, projectName, url)
	if err != nil {
		sdkv2.GetPlugin().Done(jobId, map[string]any{"info": importresult, "error": "timeout during joern operation"})
		return
	}
	body := map[string]any{"query": fmt.Sprintf(`open("%s").get.name`, projectName)}
	bodyByte, _ := sonic.Marshal(body)
	resp, status, err := etc.CustomCall(context.Background(), "POST", bodyByte, url, nil)
	if err != nil {
		sdkv2.GetPlugin().Done(jobId, map[string]any{"error": err.Error()})
		return
	}
	if status != 200 {
		sdkv2.GetPlugin().Done(jobId, map[string]any{"error": "server is unavailable or bad request"})
		return
	}
	respMap := map[string]any{}
	err = sonic.Unmarshal(resp, &respMap)
	if err != nil {
		sdkv2.GetPlugin().Progress(jobId,  sdkModel.CommandPayload{Progress: 100, Details: map[string]any{"error": "invalid response from joern server"}})
		return
	}
	if success, ok := respMap["success"].(bool); ok && success {
		re := regexp.MustCompile(`\s*"([^"]+)\"`)
		out := re.FindStringSubmatch(respMap["stdout"].(string))
		if len(out) > 1 {
			sdkv2.GetPlugin().Done(jobId, map[string]any{"opened_project":out[1],"result": respMap})
			return
		}
	}
	respMap["error"] = "unable to open project"
	sdkv2.GetPlugin().Done(jobId, respMap)
}
func workOnGitHandler(jobId string, git chan models.GitResponse) {
	progress := 0
	// Don't close the channel here - let ClonePull close it when done
	for {
		select {
		case status, ok := <-git:
			if !ok {
				// Channel closed by ClonePull, exit
				return
			}
			fmt.Println(status.Status)
			if status.Status == "success" && (status.Action == "clone" || status.Action == "pull" || status.Action == "repo") {
				fmt.Println("Done ", status.Status)
				sdkv2.GetPlugin().Done(jobId, map[string]any{"msg": status.Message, "action": status.Action, "branch": status.Branch})
				return
			}
			if status.Status == "error" {
				sdkv2.GetPlugin().Done(jobId, map[string]any{"error": status.Error, "action": status.Action})
				return
			}
			progress += 5
			if progress >99{
				progress = 99
			}
			sdkv2.GetPlugin().Progress(jobId, sdkModel.CommandPayload{Progress: progress, Details: map[string]any{"msg": status.Message, "action": status.Action, "branch": status.Branch}})
			fmt.Println("GIT PROGRESS LOGG: ", status)
		case <-time.After(10 * time.Minute):
			sdkv2.GetPlugin().Done(jobId, map[string]any{"error": "timeout during git operation"})
			return
		}

	}

}
func workOnQueryGraph(jobId, joern_uuid string) {

	joern:=joernMessHandler.GetJoernResultHandler()
	if joern==nil{
		fmt.Println("joern websocket handler not found.")
		return
	}
	ch := joern.GetResultChannel()
	for {

		select {
		case payload := <-ch:
			// fmt.Println("Received payload:", payload)
			if !strings.HasPrefix(payload, joern_uuid) {
				continue
			}
			splitPayload := strings.SplitN(payload, "||", 2)
			if len(splitPayload) != 2 {
				continue
			}
			message := splitPayload[1]

			responsMap := map[string]any{}
			err := sonic.Unmarshal([]byte(message), &responsMap)
			if err != nil {
				fmt.Println("Done With Err Response")
				sdkv2.GetPlugin().Done(jobId, map[string]any{"result": message})
				return
			}
			FinaleResponse := map[string]any{"results": responsMap}
			fmt.Println("Done With Success Response")

			sdkv2.GetPlugin().Done(jobId, FinaleResponse)
			return

		case <-time.After(10 * time.Minute):
			sdkv2.GetPlugin().Done(jobId, map[string]any{"error": "timeout waiting for Joern response"})
			return
		}
	}

}
