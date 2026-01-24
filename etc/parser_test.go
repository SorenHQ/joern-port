package etc

import (
	"fmt"
	"testing"

	"github.com/bytedance/sonic"
)

func TestParseJoernStdoutJsonObject(t *testing.T) {
	// Example input from user
	input := `{"success":true,"uuid":"75e5177c-00c1-4754-ac76-afa3cb452473","stdout":"\u001b[33mval\u001b[0m \u001b[36mres2\u001b[0m: \u001b[32mOption\u001b[0m[io.joern.console.workspacehandling.Project] = Some(\n  value = Project(\n    projectFile = ProjectFile(inputPath = \"/git_repo_location/projects/Digdag\", name = \"Digdag\"),\n    path = /home/ubuntu/workspace/Digdag,\n    cpg = Some(value = Cpg[Graph[259859 nodes]])\n  )\n)\ndef escapeForRegex(input: String): String\ndef writeJsonToFile(outputMap: Map[String, ujson.Obj]): Unit\ndef getFunctionFullNames(httpMethod: String): List[String]\ndef collectParams\n  (annotation: String, functionFullName: String):\n    Map[String, Map[String, String]]\ndef extractEndpointPath(functionFullName: String): String\ndef extractMethodComments(functionFullName: String): (String, String)\ndef extractConsumeAndProduceTypes(functionFullName: String): (String, String)\ndef collectResponses(httpMethod: String): Map[String, Map[String, String]]\ndef collectAuthorization(functionFullName: String): ujson.Arr\ndef buildJsonForFunction\n  (httpMethod: String, functionFullName: String): ujson.Obj\ndef createOutputMap(httpMethods: List[String]): Map[String, ujson.Obj]\ndef printOutputMapAsJson(outputMap: Map[String, ujson.Obj]): Unit\ndef main(): String\n\u001b[33mval\u001b[0m \u001b[36mres3\u001b[0m: \u001b[32mString\u001b[0m = \"\"\"{\n  \"io.digdag.server.rs.ScheduleResource.enableSchedule:io.digdag.client.api.RestScheduleSummary(int,io.digdag.client.api.RestScheduleEnableRequest)\": {\n    \"timestamp\": \"2026-01-21T15:30:57.687056926Z\",\n    \"http_method\": \"POST\",\n    \"func_name\": \"enableSchedule\",\n    \"file_path\": \"digdag-server/src/main/java/io/digdag/server/rs/ScheduleResource.java\",\n    \"file_name\": \"ScheduleResource.java\",\n    \"file_type\": \"java\",\n    \"api_linenum\": \"253\",\n    \"params\": {\n      \"path_params\": {},\n      \"query_params\": {},\n      \"header_params\": {},\n      \"body_params\": {},\n      \"cookie_params\": {},\n      \"matrix_params\": {},\n      \"form_params\": {}\n    },\n    \"endpoint_path\": \"/api/schedules/{id}/enable\",\n    \"summary\": \"No summary available.\",\n    \"description\": \"No description available.\",\n    \"consume_type\": \"application/json\",\n    \"produce_type\": \"application/json\",\n    \"authorization\": [],\n    \"responses\": {\n      \"201\": {\n        \"description\": \"Resource created\",\n        \"content_type\": \"application/json\"\n      },\n      \"400\": {\n        \"description\": \"Bad request\",\n        \"content_type\": \"application/json\"\n      }\n    }\n  },\n  \"io.digdag.server.rs.ScheduleResource.backfillSchedule:io.digdag.client.api.RestScheduleAttemptCollection(int,io.digdag.client.api.RestScheduleBackfillRequest)\": {\n    \"timestamp\": \"2026-01-21T15:30:57.684858024Z\",\n    \"http_method\": \"POST\",\n    \"func_name\": \"backfillSchedule\",\n    \"file_path\": \"digdag-server/src/main/java/io/digdag/server/rs/ScheduleResource.java\",\n    \"file_name\": \"ScheduleResource.java\",\n    \"file_type\": \"java\",\n    \"api_linenum\": \"188\",\n    \"params\": {\n      \"path_params\": {},\n      \"query_params\": {},\n      \"header_params\": {},\n      \"body_params\": {},\n      \"cookie_params\": {},\n      \"matrix_params\": {},\n      \"form_params\": {}\n    },\n    \"endpoint_path\": \"/api/schedules/{id}/backfill\",\n    \"summary\": \"No summary available.\",\n    \"description\": \"No description available.\",\n    \"consume_type\": \"application/json\",\n    \"produce_type\": \"application/json\",\n    \"authorization\": [],\n    \"responses\": {\n      \"201\": {\n        \"description\": \"Resource created\",\n        \"content_type\": \"application/json\"\n      },\n      \"400\": {\n        \"description\": \"Bad request\",\n        \"content_type\": \"application/json\"\n      }\n    }\n  },\n  \"io.digdag.server.rs.ScheduleResource.disableSchedule:io.digdag.client.api.RestScheduleSummary(int)\": {\n    \"timestamp\": \"2026-01-21T15:30:57.686051456Z\",\n    \"http_method\": \"POST\",\n    \"func_name\": \"disableSchedule\",\n    \"file_path\": \"digdag-server/src/main/java/io/digdag/server/rs/ScheduleResource.java\",\n    \"file_name\": \"ScheduleResource.java\",\n    \"file_type\": \"java\",\n    \"api_linenum\": \"221\",\n    \"params\": {\n      \"path_params\": {},\n      \"query_params\": {},\n      \"header_params\": {},\n      \"body_params\": {},\n      \"cookie_params\": {},\n      \"matrix_params\": {},\n      \"form_params\": {}\n    },\n    \"endpoint_path\": \"/api/schedules/{id}/disable\",\n    \"summary\": \"No summary available.\",\n    \"description\": \"No description available.\",\n    \"consume_type\": \"application/json\",\n    \"produce_type\": \"application/json\",\n    \"authorization\": [],\n    \"responses\": {\n      \"201\": {\n        \"description\": \"Resource created\",\n        \"content_type\": \"application/json\"\n      },\n      \"400\": {\n        \"description\": \"Bad request\",\n        \"content_type\": \"application/json\"\n      }\n    }\n  },\n  \"io.digdag.server.rs.ScheduleResource.skipSchedule:io.digdag.client.api.RestScheduleSummary(int,io.digdag.client.api.RestScheduleSkipRequest)\": {\n    \"timestamp\": \"2026-01-21T15:30:57.682916373Z\",\n    \"http_method\": \"POST\",\n    \"func_name\": \"skipSchedule\",\n    \"file_path\": \"digdag-server/src/main/java/io/digdag/server/rs/ScheduleResource.java\",\n    \"file_name\": \"ScheduleResource.java\",\n    \"file_type\": \"java\",\n    \"api_linenum\": \"135\",\n    \"params\": {\n      \"path_params\": {},\n      \"query_params\": {},\n      \"header_params\": {},\n      \"body_params\": {},\n      \"cookie_params\": {},\n      \"matrix_params\": {},\n      \"form_params\": {}\n    },\n    \"endpoint_path\": \"/api/schedules/{id}/skip\",\n    \"summary\": \"No summary available.\",\n    \"description\": \"No description available.\",\n    \"consume_type\": \"application/json\",\n    \"produce_type\": \"application/json\",\n    \"authorization\": [],\n    \"responses\": {\n      \"201\": {\n        \"description\": \"Resource created\",\n        \"content_type\": \"application/json\"\n      },\n      \"400\": {\n        \"description\": \"Bad request\",\n        \"content_type\": \"application/json\"\n      }\n    }\n  },\n  \"io.digdag.server.rs.AttemptResource.killAttempt:javax.ws.rs.core.Response(long)\": {\n    \"timestamp\": \"2026-01-21T15:30:57.616883919Z\",\n    \"http_method\": \"POST\",\n    \"func_name\": \"killAttempt\",\n    \"file_path\": \"digdag-server/src/main/java/io/digdag/server/rs/AttemptResource.java\",\n    \"file_name\": \"AttemptResource.java\",\n    \"file_type\": \"java\",\n    \"api_linenum\": \"355\",\n    \"params\": {\n      \"path_params\": {},\n      \"query_params\": {},\n      \"header_params\": {},\n      \"body_params\": {},\n      \"cookie_params\": {},\n      \"matrix_params\": {},\n      \"form_params\": {}\n    },\n    \"endpoint_path\": \"/api/attempts/{id}/kill\",\n    \"summary\": \"No summary available.\",\n    \"description\": \"No description available.\",\n    \"consume_type\": \"application/json\",\n    \"produce_type\": \"application/json\",\n    \"authorization\": [],\n    \"responses\": {\n      \"201\": {\n        \"description\": \"Resource created\",\n        \"content_type\": \"application/json\"\n      },\n      \"400\": {\n        \"description\": \"Bad request\",\n        \"content_type\": \"application/json\"\n      }\n    }\n  }\n}\"\"\""}`

	result, err := ParseJoernStdoutJsonObject([]byte(input))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Verify that we got a map with at least one key
	if len(result) == 0 {
		t.Fatal("Parsed result is empty")
	}

	// Verify that we can find one of the expected keys
	expectedKey := "io.digdag.server.rs.ScheduleResource.enableSchedule:io.digdag.client.api.RestScheduleSummary(int,io.digdag.client.api.RestScheduleEnableRequest)"
	if _, ok := result[expectedKey]; !ok {
		t.Fatalf("Expected key not found in result: %s", expectedKey)
	}

	// Verify the structure of one entry
	entry, ok := result[expectedKey].(map[string]any)
	if !ok {
		t.Fatal("Entry is not a map")
	}

	// Check some expected fields
	if entry["http_method"] != "POST" {
		t.Errorf("Expected http_method to be POST, got %v", entry["http_method"])
	}

	if entry["func_name"] != "enableSchedule" {
		t.Errorf("Expected func_name to be enableSchedule, got %v", entry["func_name"])
	}

	// Verify we can marshal it back to JSON
	jsonBytes, err := sonic.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}

	// Print the parsed payload
	fmt.Println("\n=== Parsed Payload ===")
	fmt.Println(string(jsonBytes))
	fmt.Println("=====================")
}
