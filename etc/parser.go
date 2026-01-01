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