package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/bytedance/sonic"
)

func TestParser(t *testing.T) {

	dat, err := os.ReadFile("./out.json")
	if err != nil {
		t.Fatal(err)
	}
	data := map[string]string{}
	sonic.Unmarshal(dat, &data)
	ansiRegexp := regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)
	cleanText := ansiRegexp.ReplaceAllString(data["stdout"], "")
	re := regexp.MustCompile(`val\s+res\d{1,10}\s*:\s*String\s*=\s*"((?:\\.|[^\\"])*)"`)
	out := re.FindStringSubmatch(cleanText)
	// fmt.Println(out)
	if len(out) > 1 {
		// Unescape JSON string
		jsonStr, err := strconv.Unquote(`"` + out[1] + `"`)
		if err != nil {
			fmt.Println("Error unquoting JSON:", err)
			return
		}

		fmt.Println("Captured JSON:")
		// fmt.Println(jsonStr)
		outinfo := []map[string]any{}
		sonic.Unmarshal([]byte(jsonStr), &outinfo)
		b, _ := sonic.Marshal(outinfo)
		fmt.Println(string(b))
	} else {
		fmt.Println("No match found")
	}
}
