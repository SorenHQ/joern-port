package etc

import (
	"errors"
	"fmt"
	"regexp"


	"github.com/bytedance/sonic"
)
func ParseJoernResult(content string) ([]map[string]any, error) {

		ansiRegexp := regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)
		cleanText := ansiRegexp.ReplaceAllString(content, "")
	re := regexp.MustCompile(`(?s)"""(\[.*?\])"""`)
	out := re.FindStringSubmatch(cleanText)
	// fmt.Println(out)
	if len(out) > 1 {
		// Unescape JSON string
		// jsonStr, _ := strconv.Unquote(`"` + out[1] + `"`)
		// if err != nil {
		// 	fmt.Println("Error unquoting JSON:", err)
		// 	return nil, err
		// }
		// if jsonStr==""{
		// 	jsonStr = out[1]
		// }
		outinfo := []map[string]any{}
		err := sonic.Unmarshal([]byte(out[1]), &outinfo)

		return outinfo, err
	} else {
		return []map[string]any{{"result":content}}, nil
	}

}
func ParseJoernStdout(raw []byte) ([]map[string]any, error) {
	data := map[string]string{}
	sonic.Unmarshal(raw, &data)
		ansiRegexp := regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)
		cleanText := ansiRegexp.ReplaceAllString(data["stdout"], "")
	re := regexp.MustCompile(`(?s)"""(\[.*?\])"""`)
	out := re.FindStringSubmatch(cleanText)
	// fmt.Println(out)
	if len(out) > 1 {
		// Unescape JSON string
		// jsonStr, _ := strconv.Unquote(`"` + out[1] + `"`)
		// if err != nil {
		// 	fmt.Println("Error unquoting JSON:", err)
		// 	return nil, err
		// }
		// if jsonStr==""{
		// 	jsonStr = out[1]
		// }
		outinfo := []map[string]any{}
		err := sonic.Unmarshal([]byte(out[1]), &outinfo)

		return outinfo, err
	} else {
		return []map[string]any{{"result":data}}, nil
	}

}



func ParseJoernStdoutToString(raw []byte) (string, error) {
	data := map[string]string{}
	sonic.Unmarshal(raw, &data)
		// ansiRegexp := regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)
		// cleanText := ansiRegexp.ReplaceAllString(data["stdout"], "")
	re := regexp.MustCompile(`(?s)"""(\[.*?\])"""`)
	out := re.FindStringSubmatch(data["stdout"])
	// fmt.Println(out)
	if len(out) > 1 {
		// Unescape JSON string
		// jsonStr, _ := strconv.Unquote(`"` + out[1] + `"`)
		// if err != nil {
		// 	fmt.Println("Error unquoting JSON:", err)
		// 	return nil, err
		// }
		// if jsonStr==""{
		// 	jsonStr = out[1]
		// }
		// outinfo := []map[string]any{}
		// err := sonic.Unmarshal([]byte(out[1]), &outinfo)

		return out[1], nil
	} else {
		fmt.Println("No match found")
		return "", errors.New("No match found")
	}

}

// ParseJoernStdoutJsonObject parses stdout that contains a JSON object in triple quotes.
// It extracts JSON objects from patterns like: val res3: String = """{...}"""
// Returns the parsed JSON object as map[string]any
func ParseJoernStdoutJsonObject(raw []byte) (map[string]any, error) {
	// Parse the outer JSON structure (success, uuid, stdout)
	data := map[string]any{}
	err := sonic.Unmarshal(raw, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse outer JSON: %w", err)
	}

	// Extract stdout field
	stdout, ok := data["stdout"].(string)
	if !ok {
		return nil, errors.New("stdout field not found or not a string")
	}

	// Remove ANSI escape codes
	ansiRegexp := regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)
	cleanText := ansiRegexp.ReplaceAllString(stdout, "")

	// Find the pattern: val res\d+ : String = """{...}"""
	// Extract everything between the triple quotes
	re := regexp.MustCompile(`(?s)val\s+res\d+\s*:\s*String\s*=\s*"""(.*?)"""`)
	match := re.FindStringSubmatch(cleanText)
	if len(match) < 2 {
		return nil, errors.New("no JSON object pattern found in stdout")
	}

	// Parse the extracted JSON object
	jsonObj := map[string]any{}
	err = sonic.Unmarshal([]byte(match[1]), &jsonObj)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON object: %w", err)
	}
	return jsonObj, nil
}