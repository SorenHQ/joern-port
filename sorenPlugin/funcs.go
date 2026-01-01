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
	sdkModel "github.com/sorenhq/go-plugin-sdk/gosdk/models"
)

func openJoernProject(jobId string, projectName string) {

	dir := fmt.Sprintf("%s/%s", env.GetProjectReposPath(), projectName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		PluginInstance.Done(jobId, map[string]any{"error": "Project not found"})
		return
	}
	url := fmt.Sprintf("%s/query-sync", env.GetJoernUrl())
	importresult, err := etc.ImportCode(context.Background(), dir, projectName, url)
	if err != nil {
		PluginInstance.Done(jobId, map[string]any{"info": importresult, "error": "timeout during joern operation"})
		return
	}
	body := map[string]any{"query": fmt.Sprintf(`open("%s").get.name`, projectName)}
	bodyByte, _ := sonic.Marshal(body)
	resp, status, err := etc.CustomCall(context.Background(), "POST", bodyByte, url, nil)
	if err != nil {
		PluginInstance.Done(jobId, map[string]any{"error": err.Error()})
		return
	}
	if status != 200 {
		PluginInstance.Done(jobId, map[string]any{"error": "server is unavailable or bad request"})
		return
	}
	respMap := map[string]any{}
	err = sonic.Unmarshal(resp, &respMap)
	if err != nil {
		PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Details: map[string]any{"error": "invalid response from joern server"}})
		return
	}
	if success, ok := respMap["success"].(bool); ok && success {
		re := regexp.MustCompile(`\s*"([^"]+)\"`)
		out := re.FindStringSubmatch(respMap["stdout"].(string))
		if len(out) > 1 {
			PluginInstance.Done(jobId, map[string]any{"opened_project":out[1],"result": respMap})
			return
		}
	}
	respMap["error"] = "unable to open project"
	PluginInstance.Done(jobId, respMap)
}
func workOnGitHandler(jobId string, git chan models.GitResponse) {
	progress := 0
	defer close(git)
	for {
		select {
		case status := <-git:
			fmt.Println(status.Status)
			if status.Status == "success" {
				fmt.Println("Done ", status.Status)
				PluginInstance.Done(jobId, map[string]any{"msg": status.Message, "action": status.Action, "branch": status.Branch})
				return
			}
			progress += 5
			PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: progress, Details: map[string]any{"msg": status.Message, "action": status.Action, "branch": status.Branch}})
			fmt.Println("GIT PROGRESS LOGG: ", status)
		case <-time.After(10 * time.Minute):
			PluginInstance.Done(jobId, map[string]any{"error": "timeout during git operation"})
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

			responsMap :=[] map[string]any{}
			err := sonic.Unmarshal([]byte(message), &responsMap)
			if err != nil {
				fmt.Println("Done With Err Response")
				PluginInstance.Done(jobId, map[string]any{"result": message})
				return
			}
			FinaleResponse := map[string]any{"results": responsMap}
			fmt.Println("Done With Success Response")

			PluginInstance.Done(jobId, FinaleResponse)
			return

		case <-time.After(10 * time.Minute):
			PluginInstance.Done(jobId, map[string]any{"error": "timeout waiting for Joern response"})
			return
		}
	}

}
