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
					sdkv2.RejectWithBody(msg, map[string]any{"details": map[string]any{"error": "bad request"}})
					return
				}
				selectedProject, ok := rawInput.Body["project"].(string)
				if !ok {
					sdkv2.RejectWithBody(msg, map[string]any{"details": map[string]any{"error": "invalid data"}})
					return
				}
				inputDataFromSorenPlatform := models.GitRequest{
					Project: selectedProject,
					RepoURL: GetProjectUrl(selectedProject),
				}
				jobId := sdkv2.Accept(msg)
				gitStatusHandler := make(chan models.GitResponse, 100) // Buffered channel to prevent blocking
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
					sdkv2.RejectWithBody(msg, map[string]any{"details": map[string]any{"error": "bad request"}})
					return
				}
				selectedProject, ok := rawInput.Body["project"].(string)
				if !ok {
					sdkv2.RejectWithBody(msg, map[string]any{"details": map[string]any{"error": "invalid data"}})
					return
				}

				jobId := sdkv2.Accept(msg)
				go openJoernProject(jobId, selectedProject)

			},
		},

		// Action #3 - Query On Graph
		{
			Method: "graph.query",
			Title:  "Query On Graph",
			Form: sdkv2Models.ActionFormBuilder{
				// Make "query" a multi-line textarea in Soren Panel
				Jsonui: map[string]any{
					"type": "VerticalLayout",
					"elements": []map[string]any{
						{
							"type":  "Control",
							"scope": "#/properties/project",
						},
						{
							"type":  "Control",
							"scope": "#/properties/query",
							"options": map[string]any{
								"multi": true,
								"rows":  8,
							},
						},
					},
				},
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
							// Hint for Soren Panel schema renderer to use a textarea
							"multiline":   true,
							// Extra hint for other JSON-schema-based UIs
							"format":      "textarea",
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
					sdkv2.RejectWithBody(msg, map[string]any{"details": map[string]any{"error": "bad request"}})
					return
				}
				selectedProject, ok := rawInput.Body["project"].(string)
				if !ok {
					sdkv2.RejectWithBody(msg, map[string]any{"details": map[string]any{"error": "invalid data"}})
					return
				}
				userQuery, ok := rawInput.Body["query"].(string)
				if !ok {
					sdkv2.RejectWithBody(msg, map[string]any{"details": map[string]any{"error": "invalid data"}})
					return
				}
				jobId := sdkv2.Accept(msg)
				// make a Joern Command Request

				PluginInstance.Progress(jobId, sdkv2Models.ProgressCommand, sdkv2Models.JobProgress{Progress: 20, Details: map[string]any{"msg": "hello"}})

				fullQuery := fmt.Sprintf("open(\"%s\")\n%s", selectedProject, userQuery)
				url := fmt.Sprintf("%s/query", env.GetJoernUrl()) // async call
				req_uuid, err := etc.JoernAsyncCommand(PluginInstance.GetContext(), url, fullQuery)
				if err != nil {
					PluginInstance.Done(jobId, map[string]any{"error": err.Error()})
					return
				}
				go workOnQueryGraph(jobId, req_uuid)

			},
		},
	})
	event := sdkv2.NewEventLogger(sdkInstance)
	event.Log("remote-mate-pc", sdkv2Models.LogLevelInfo, fmt.Sprintf("%s Plugin Started , Version %s", plugin.Intro.Name, plugin.Intro.Version), nil)
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
		return msg.Respond([]byte(`{"status": "error , bad request"}`))
	}
	err = os.WriteFile("my_database.json", msg.Data, 0644)
	if err != nil {
		fmt.Println("Error Writing Settings to File:", err)
		return msg.Respond([]byte(`{"status": "not_accepted"}`))

	}

	// Keep in-memory settings in sync so @settings returns fresh data
	if PluginInstance != nil && PluginInstance.Settings != nil {
		PluginInstance.Settings.Data = settings
	}

	// Also refresh any enums that depend on settings (e.g. project list)
	updateProjectEnum(settings)

	return msg.Respond([]byte(`{"status": "accepted"}`))
}

func makeEnumsProject() []string {
	// Build the enum for the project field based on the latest saved settings
	savedSettings := getSavedData()
	project, ok := savedSettings["project"].(string)
	if !ok || project == "" {
		return []string{}
	}
	fmt.Println(savedSettings)
	return []string{project}

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

// updateProjectEnum updates any schema enums that depend on the "project" setting
// so that Soren Panel sees the latest value without needing a plugin restart.
func updateProjectEnum(settings map[string]any) {
	if PluginInstance == nil {
		return
	}

	project, ok := settings["project"].(string)
	if !ok || project == "" {
		return
	}
	enumValue := []string{project}

	// Update settings schema project enum if present
	if PluginInstance.Settings != nil && PluginInstance.Settings.Jsonschema != nil {
		if props, ok := PluginInstance.Settings.Jsonschema["properties"].(map[string]any); ok {
			if proj, ok := props["project"].(map[string]any); ok {
				proj["enum"] = enumValue
				props["project"] = proj
				PluginInstance.Settings.Jsonschema["properties"] = props
			}
		}
	}

	// Update all action forms that have a project field
	for i := range PluginInstance.Actions {
		formSchema := PluginInstance.Actions[i].Form.Jsonschema
		if formSchema == nil {
			continue
		}
		props, ok := formSchema["properties"].(map[string]any)
		if !ok {
			continue
		}
		proj, ok := props["project"].(map[string]any)
		if !ok {
			continue
		}
		proj["enum"] = enumValue
		props["project"] = proj
		formSchema["properties"] = props
		PluginInstance.Actions[i].Form.Jsonschema = formSchema
	}
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
