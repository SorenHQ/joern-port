package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/SorenHQ/joern-port/env"
	"github.com/SorenHQ/joern-port/etc"
	"github.com/SorenHQ/joern-port/etc/db"
	"github.com/SorenHQ/joern-port/models"
	"github.com/bytedance/sonic"
	sdkModel "github.com/sorenhq/go-plugin-sdk/gosdk/models"
)

func openJoernProject(jobId string, projectName string) {

	dir := fmt.Sprintf("%s/%s", env.GetProjectReposPath(), projectName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: map[string]any{"error": "Project not found"}})
		return
	}
	url := fmt.Sprintf("%s/query-sync", env.GetJoernUrl())
	importresult, err := etc.ImportCode(context.Background(), dir, projectName, url)
	if err != nil {
			PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: map[string]any{"info":importresult,"error": "timeout during git operation"}})
		return
	}
	body := map[string]any{"query": fmt.Sprintf(`open("%s").get.name`, projectName)}
	bodyByte, _ := sonic.Marshal(body)
	resp, status, err := etc.CustomCall(context.Background(), "POST", bodyByte, url, nil)
	if err != nil {
		PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: map[string]any{"error": err.Error()}})
		return
	}
	if status != 200 {
		PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: map[string]any{"error": "server is unavailable or bad request"}})
		return
	}
	respMap := map[string]any{}
	err = sonic.Unmarshal(resp, &respMap)
	if err != nil {
		PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: map[string]any{"error": "invalid response from joern server"}})
		return
	}
	if success, ok := respMap["success"].(bool); ok && success {
		re := regexp.MustCompile(`\s*"([^"]+)\"`)
		out := re.FindStringSubmatch(respMap["stdout"].(string))
		if len(out) > 1 {
			PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: map[string]any{"result": out[1]}})
			return
		}
	}
	respMap["error"] = "unable to open project"
	PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: respMap})
}
func workOnGitHandler(jobId string, git chan models.GitResponse) {
	progress := 0
	defer close(git)
	for {
		select {
		case status := <-git:
			if status.Status == "success" {
				PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: map[string]any{"msg": status.Message, "action": status.Action, "branch": status.Branch}})
				continue
			}
			progress += 5
			PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: progress, Data: map[string]any{"msg": status.Message, "action": status.Action, "branch": status.Branch}})
			fmt.Println("GIT PROGRESS LOGG: ", status)
		case <-time.After(20 * time.Minute):
			PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: map[string]any{"error": "timeout during git operation"}})
			return
		}

	}

}
func workOnQueryGraph(jobId string) {
	// simulate long processing
	sub := db.GetRedisClient().Subscribe(context.Background(), JoernResultsTableInRedis)

	defer sub.Close()
	ch := sub.Channel()
	for {

		select {
		case msg := <-ch:
			payload := msg.Payload
			if !strings.HasPrefix(payload, jobId) {
				continue
			}
			splitPayload := strings.SplitN(payload, "||", 2)
			if len(splitPayload) != 2 {
				continue
			}
			message := splitPayload[1]
			res, err := etc.ParseJoernStdout([]byte(message))
			if err != nil {
				PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: map[string]any{"error": err.Error()}})
				return
			}

			FinaleResponse := map[string]any{"results": res}

			PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: FinaleResponse})
			return

		case <-time.After(45 * time.Minute):
			PluginInstance.Progress(jobId, sdkModel.ProgressCommand, sdkModel.JobProgress{Progress: 100, Data: map[string]any{"error": "timeout waiting for Joern response"}})
			return
		}
	}

}
