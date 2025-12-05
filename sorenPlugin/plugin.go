package joernPlugin

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SorenHQ/joern-port/env"
	"github.com/SorenHQ/joern-port/etc"
	"github.com/SorenHQ/joern-port/models"
	gitServices "github.com/SorenHQ/joern-port/services/git"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	sdkv2 "github.com/sorenhq/go-plugin-sdk/gosdk"
	sdkv2Models "github.com/sorenhq/go-plugin-sdk/gosdk/models"
)

var PluginInstance *sdkv2.Plugin

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
		ReplyTo: "settings.config.submit",
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
				// Jsonui:     map[string]any{"type": "Control", "scope": "#/properties/project"},
				Jsonschema: map[string]any{"properties": map[string]any{"project": map[string]any{"enum": makeEnumsProject()}}},
			},
			RequestHandler: func(msg *nats.Msg) any {
				// for example in this step we register a job in local database or external system - mae a scan in Joern
				data := msg.Data
				fmt.Println(string(data))
				rawInput := sdkv2Models.ActionRequestContent{}
				err = sonic.Unmarshal(data, &rawInput)
				if err != nil {
					return msg.Respond([]byte(`{"details": {"error": "bad request"}}`))
				}
				selectedProject, ok := rawInput.Body["project"].(string)
				if !ok {
					return msg.Respond([]byte(`{"details": {"error": "invalid data"}}`))

				}
				inputDataFromSorenPlatform := models.GitRequest{
					Project: selectedProject,
					RepoURL: GetProjectUrl(selectedProject),
				}
				uuid, err := uuid.NewV6()
				if err != nil {
					return msg.Respond([]byte(`{"details": {"error": "service unavailable"}}`))
				}
				gitStatusHandler := make(chan models.GitResponse)
				git := gitServices.NewGitDetailsHandler(gitStatusHandler)
				fmt.Println(uuid.String())

				err = msg.Respond([]byte(fmt.Sprintf(`{"jobId":"%s"}`, uuid.String())))
				if err == nil {
					go workOnGitHandler(uuid.String(), gitStatusHandler)
					go git.ClonePull(inputDataFromSorenPlatform.Project, inputDataFromSorenPlatform.RepoURL, true)
				}
				return nil
			},
		},
		// Action #2 - scan code and create graph
		{
			Method: "scan.gen.graph",
			Title:  "Open Project in Joern and Generate Graph",
			Form: sdkv2Models.ActionFormBuilder{
				// Jsonui:     map[string]any{"type": "Control", "scope": "#/properties/project"},
				Jsonschema: map[string]any{"properties": map[string]any{"project": map[string]any{"enum": makeEnumsProject()}}},
			},
			RequestHandler: func(msg *nats.Msg) any {
				// for example in this step we register a job in local database or external system - mae a scan in Joern
				data := msg.Data
				fmt.Println(string(data))
				rawInput := sdkv2Models.ActionRequestContent{}
				err = sonic.Unmarshal(data, &rawInput)
				if err != nil {
					return msg.Respond([]byte(`{"details": {"error": "bad request"}}`))
				}
				selectedProject, ok := rawInput.Body["project"].(string)
				if !ok {
					return msg.Respond([]byte(`{"details": {"error": "invalid data"}}`))

				}

				uuid, err := uuid.NewV6()
				if err != nil {
					return msg.Respond([]byte(`{"details": {"error": "service unavailable"}}`))
				}
				fmt.Println(uuid.String())
				msg.Respond([]byte(fmt.Sprintf(`{"jobId":"%s"}`, uuid.String())))
				go openJoernProject(uuid.String(), selectedProject)
				
				return nil
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
			RequestHandler: func(msg *nats.Msg) any {
				// for example in this step we register a job in local database or external system - mae a scan in Joern
				data := msg.Data
				fmt.Println(string(data))
				rawInput := sdkv2Models.ActionRequestContent{}
				err = sonic.Unmarshal(data, &rawInput)
				if err != nil {
					return msg.Respond([]byte(`{"details": {"error": "bad request"}}`))
				}
				selectedProject, ok := rawInput.Body["project"].(string)
				if !ok {
					return msg.Respond([]byte(`{"details": {"error": "invalid data"}}`))

				}
				userQuery, ok := rawInput.Body["query"].(string)
				if !ok {
					return msg.Respond([]byte(`{"details": {"error": "invalid data"}}`))
				}
				uuid, err := uuid.NewV6()
				if err != nil {
					return msg.Respond([]byte(`{"details": {"error": "service unavailable"}}`))
				}
				// make a Joern Command Request
				
				msg.Respond([]byte(fmt.Sprintf(`{"jobId":"%s"}`, uuid.String())))
					time.Sleep(1*time.Second)
					fullQuery:=fmt.Sprintf(`workspace.project("%s").get.cpg.get.%s`, selectedProject, userQuery)
					url := fmt.Sprintf("%s/query", env.GetJoernUrl()) // async call
					req_uuid, err := etc.JoernAsyncCommand(PluginInstance.GetContext(), url, fullQuery)
					if err != nil {
						PluginInstance.Progress(uuid.String(), sdkv2Models.ProgressCommand, sdkv2Models.JobProgress{Progress: 100, Details: map[string]any{"error": err.Error()}})
						return nil
					}
					go workOnQueryGraph(uuid.String(),req_uuid)
				
				return nil

			},
		},
	})
	event := sdkv2.NewEventLogger(sdkInstance)
	event.Log("remote-mate-pc", sdkv2Models.LogLevelInfo, "Joern Plugin Started", nil)
	PluginInstance = plugin

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
		return map[string]any{"status": "error"}
	}
	err = os.WriteFile("my_database.json", msg.Data, 0644)
	if err != nil {
		fmt.Println("Error Writing Settings to File:", err)
		return map[string]any{"status": "not_accepted"}

	}

	return map[string]any{"status": "accepted"}
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

