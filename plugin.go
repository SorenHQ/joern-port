package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SorenHQ/joern-port/env"
	"github.com/SorenHQ/joern-port/etc"
	"github.com/SorenHQ/joern-port/models"
	gitServices "github.com/SorenHQ/joern-port/services/git"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
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
		Name:    "Code Analysis Plugin",
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
			"required": []string{"repository_name", "access_token","project"},
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
		RequestHandler: func(data []byte) any {
			// for example in this step we register a job in local database or external system - mae a scan in Joern
			inputDataFromSorenPlatform:=models.GitRequest{}
			err=sonic.Unmarshal(data,&inputDataFromSorenPlatform)
			if err!=nil{
				return map[string]any{"details": map[string]any{"error": "invalid project"}}
			}
			gitStatusHandler:=make(chan models.GitResponse)
			git:=gitServices.NewGitDetailsHandler(gitStatusHandler)
			git.ClonePull(inputDataFromSorenPlatform.Project, inputDataFromSorenPlatform.RepoURL, true)
			uuid,err:=uuid.NewV6()
			if err!=nil{
				return map[string]any{"jobId": uuid.String()}
			}
			go workOnGitHandler(uuid.String(), gitStatusHandler)
			return map[string]any{"jobId": uuid.String()}
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
		RequestHandler: func(data []byte) any {
			// for example in this step we register a job in local database or external system - mae a scan in Joern
			uuid,err:=uuid.NewV6()
			if err!=nil{
				return map[string]any{"jobId": uuid.String()}
			}
			project:=map[string]any{}
			err=sonic.Unmarshal(data,&project)
			if err!=nil{
				return map[string]any{"details": map[string]any{"error": "invalid project"}}
			}
			go openJoernProject(uuid.String(), project["project"].(string))
			return map[string]any{"details": map[string]any{"error": "service unavailable"}}
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
		RequestHandler: func(data []byte) any {
			// for example in this step we register a job in local database or external system - mae a scan in Joern

			query:=map[string]any{}
			err=sonic.Unmarshal(data,&query)
			if err!=nil{
				return map[string]any{"details": map[string]any{"error": "invalid query"}}
			}
				// make a Joern Command Request
			
					url := fmt.Sprintf("%s/query", env.GetJoernUrl())  // async call
					req_uuid, err := etc.JoernAsyncCommand(PluginInstance.GetContext(), url, query["query"].(string))
					if err != nil {
						return map[string]any{"details": map[string]any{"error": "service unavailable"}}
					}

				go workOnQueryGraph(req_uuid) 
				return map[string]any{"jobId": req_uuid}  // success reponse to workflow runtime - runtime will track job status based on jobId

		},
	},
	})
	event := sdkv2.NewEventLogger(sdkInstance)
	event.Log("remote-mate-pc", sdkv2Models.LogLevelInfo, "start plugin", nil)
	err=plugin.Start()
	if err!=nil{
		log.Fatalf("Failed to start plugin: %v", err)
	}
	PluginInstance = plugin
}

func settingsUpdateHandler(data []byte) any {
	fmt.Println("New Update As Settings : ", string(data))
	settings:=map[string]any{}
	err:=sonic.Unmarshal(data,&settings)
	if err!=nil{
		fmt.Println("Error Unmarshalling Settings:",err)
		return map[string]any{"status": "error"}
	}
	err=os.WriteFile("my_database.json",data,0644)
	if err!=nil{
		fmt.Println("Error Writing Settings to File:", err)
			return map[string]any{"status": "not_accepted"}

	}
	return map[string]any{"status": "accepted"}
}

func makeEnumsProject() []string {
	contentJson,err:=os.ReadFile("my_database.json")
	if err!=nil{
		return []string{}
	}
	savedSettings:=map[string]any{}
	err=sonic.Unmarshal(contentJson,&savedSettings)
	if err!=nil{
		return []string{}
	}
	if savedSettings["project"] == nil {
		return []string{}
	}
	return []string{savedSettings["project"].(string)}

}