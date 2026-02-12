package joernPlugin

import (
	"fmt"
	"log"
	"os"

	"github.com/SorenHQ/joern-port/env"
	"github.com/SorenHQ/joern-port/etc"
	"github.com/SorenHQ/joern-port/models"
	gitServices "github.com/SorenHQ/joern-port/services/git"
	"github.com/bytedance/sonic"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	sdkv2 "github.com/sorenhq/go-plugin-sdk/gosdk"
	sdkv2Models "github.com/sorenhq/go-plugin-sdk/gosdk/models"
)

// var sdkv2.GetPlugin() *sdkv2.Plugin

func LoadSorenPluginServer() {

	err := godotenv.Load(".env.soren")
	if err != nil {
		fmt.Println(err)
	}
	sdkInstance, err := sdkv2.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to create SDK: %v", err)
	}
	defer sdkInstance.Close()
	plugin := sdkv2.NewPlugin(sdkInstance)
	plugin.SetIntro(sdkv2Models.PluginIntro{
		Name:    "Joern Code Analyse Plugin",
		Version: "1.1.1",
		Author:  "Soren Team",
	}, nil)

	plugin.SetSettings(&sdkv2Models.Settings{
		Data: getSavedData(),
		Jsonui: map[string]any{
			"type": "VerticalLayout",
			"elements": []map[string]any{
				{
					"type":  "Control",
					"scope": "#/properties/project",
				},
				{
					"type":  "Control",
					"scope": "#/properties/repository_name",
				},
				{
					"type":  "Control",
					"scope": "#/properties/access_token",
				},
			},
		},
		Jsonschema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"project": map[string]any{
					"type":        "string",
					"title":       "Your Project Name",
					"description": "Project Name",
				},
				"repository_name": map[string]any{
					"type":        "string",
					"title":       "Your Repository Name",
					"description": "Github Respository name",
				},

				"access_token": map[string]any{
					"type":        "string",
					"title":       "Fine Grained Access Token",
					"description": "Github FineGrained Access Token",
				},
			},
			"required": []string{"repository_name", "access_token", "project"},
		},
	}, settingsUpdateHandler)
	plugin.AddActions([]sdkv2Models.Action{
		// Action #1 - prepare  (clone/pull repo)
		{
			Method: "prepare",
			Title:  "Clone/Pull Repo",
			Form: sdkv2Models.ActionFormBuilder{
				Jsonui:     map[string]any{"type": "Control", "scope": "#/properties/project"},
				Jsonschema: map[string]any{"type": "object","properties": map[string]any{"project": map[string]any{"enum": makeEnumsProject()}}},
			},
			RequestHandler: func(msg *nats.Msg) {
				// for example in this step we register a job in local database or external system - mae a scan in Joern
				data := msg.Data
				fmt.Println(string(data))
				rawInput := sdkv2Models.ActionRequestContent{}
				err = sonic.Unmarshal(data, &rawInput)
				if err != nil {
					sdkv2.RejectReq(msg, map[string]any{"details": map[string]any{"error": "bad request"}})
					return
				}
				selectedProject, ok := rawInput.Body["project"].(string)
				if !ok {
					sdkv2.RejectReq(msg, map[string]any{"details": map[string]any{"error": "invalid data"}})
					return
				}
				inputDataFromSorenPlatform := models.GitRequest{
					Project: selectedProject,
					RepoURL: GetProjectUrl(selectedProject),
				}
				jobId := sdkv2.AcceptReq(msg)
				gitStatusHandler := make(chan models.GitResponse)
				git := gitServices.NewGitDetailsHandler(gitStatusHandler)
				fmt.Println(jobId)

				go workOnGitHandler(jobId, gitStatusHandler)
				go git.ClonePull(inputDataFromSorenPlatform.Project, inputDataFromSorenPlatform.RepoURL, true)

			},
		},
		// Action #2 - scan code and create graph
		{
			Method: "scan.gen.graph",
			Title:  "Open Project in Joern and Generate Graph",
			Form: sdkv2Models.ActionFormBuilder{
				// Jsonui:     map[string]any{"type": "Control", "scope": "#/properties/project"},
				Jsonschema: map[string]any{"type": "object","properties": map[string]any{"project": map[string]any{"enum": makeEnumsProject()}}},
			},
			RequestHandler: func(msg *nats.Msg) {
				// for example in this step we register a job in local database or external system - mae a scan in Joern
				data := msg.Data
				fmt.Println(string(data))
				rawInput := sdkv2Models.ActionRequestContent{}
				err = sonic.Unmarshal(data, &rawInput)
				if err != nil {
					sdkv2.RejectReq(msg, map[string]any{"details": map[string]any{"error": "bad request"}})
					return
				}
				selectedProject, ok := rawInput.Body["project"].(string)
				if !ok {
					sdkv2.RejectReq(msg, map[string]any{"details": map[string]any{"error": "invalid data"}})
					return
				}

				jobId := sdkv2.AcceptReq(msg)
				go openJoernProject(jobId, selectedProject)

			},
		},

		// Action #3 - Query On Graph
		{
			Method: "graph.query",
			Title:  "Query On Graph",
			Form: sdkv2Models.ActionFormBuilder{
				// Jsonui:     map[string]any{"type": "Control", "scope": "#/properties/query"},
				Jsonschema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"project": map[string]any{
							"type":        "enum",
							"title":       "Your Project Name",
							"description": "Project Name",
							"enum":        makeEnumsProject(),
						},
						"query": map[string]any{
							"type":        "string",
							"title":       "Query to Execute",
							"description": "Query to Execute on Joern Graph",
						},
					},
					"required": []string{"project", "query"},
				},
			},
			RequestHandler: func(msg *nats.Msg) {
				// for example in this step we register a job in local database or external system - mae a scan in Joern
				data := msg.Data
				fmt.Println(string(data))
				rawInput := sdkv2Models.ActionRequestContent{}
				err = sonic.Unmarshal(data, &rawInput)
				if err != nil {
					sdkv2.RejectReq(msg, map[string]any{"details": map[string]any{"error": "bad request"}})
					return
				}
				selectedProject, ok := rawInput.Body["project"].(string)
				if !ok {
					sdkv2.RejectReq(msg, map[string]any{"details": map[string]any{"error": "invalid data"}})
					return
				}
				userQuery, ok := rawInput.Body["query"].(string)
				if !ok {
					sdkv2.RejectReq(msg, map[string]any{"details": map[string]any{"error": "invalid data"}})
					return
				}
				jobId := sdkv2.AcceptReq(msg)
				// make a Joern Command Request

				sdkv2.GetPlugin().Progress(jobId,  sdkv2Models.CommandPayload{Progress: 20, Details: map[string]any{"msg": "hello"}})

				fullQuery := fmt.Sprintf(`workspace.project("%s").get.cpg.get.%s`, selectedProject, userQuery)
				url := fmt.Sprintf("%s/query", env.GetJoernUrl()) // async call
				req_uuid, err := etc.JoernAsyncCommand(sdkv2.GetPlugin().GetContext(), url, fullQuery)
				if err != nil {
					sdkv2.GetPlugin().Done(jobId, map[string]any{"error": err.Error()})
					return
				}
				go workOnQueryGraph(jobId, req_uuid)

			},
		},
	})
	event := sdkv2.NewEventLogger(sdkInstance)
	event.Log("remote-mate-pc", sdkv2Models.LogLevelInfo, fmt.Sprintf("%s Plugin Started , Version %s", plugin.Intro.Name, plugin.Intro.Version), nil)

	err = plugin.Start()
	if err != nil {
		log.Fatalf("Failed to start plugin: %v", err)
	}
}

func settingsUpdateHandler(msg *nats.Msg) any {
	fmt.Println("New Update As Settings : ", string(msg.Data))
	settings := map[string]any{}
	err := sonic.Unmarshal(msg.Data, &settings)
	if err != nil {
		fmt.Println("Error Unmarshalling Settings:", err)
		return msg.Respond([]byte(`{"status": "error , bad request"}`))
	}
	err = os.WriteFile("my_database.json", msg.Data, 0644)
	if err != nil {
		fmt.Println("Error Writing Settings to File:", err)
		return msg.Respond([]byte(`{"status": "not_accepted"}`))

	}
	return msg.Respond([]byte(`{"status": "accepted"}`))
}

func makeEnumsProject() []string {
	contentJson, err := os.ReadFile("my_database.json")
	if err != nil {
		return []string{}
	}

	savedSettings := map[string]any{}
	err = sonic.Unmarshal(contentJson, &savedSettings)
	if err != nil {
		return []string{}
	}
	if savedSettings["project"] == nil {
		return []string{}
	}
	fmt.Println(savedSettings)
	return []string{savedSettings["project"].(string)}

}
func getSavedData() map[string]any {
	contentJson, err := os.ReadFile("my_database.json")
	if err != nil {
		return map[string]any{}
	}

	savedSettings := map[string]any{}
	err = sonic.Unmarshal(contentJson, &savedSettings)
	if err != nil {
		return map[string]any{}
	}
	return savedSettings

}
func GetProjectUrl(project string) string {
	contentJson, err := os.ReadFile("my_database.json")
	if err != nil {
		return ""
	}

	savedSettings := map[string]any{}
	err = sonic.Unmarshal(contentJson, &savedSettings)
	if err != nil {
		return ""
	}
	if savedSettings["repository_name"] == nil {
		return ""
	}
	return savedSettings["repository_name"].(string)
}
