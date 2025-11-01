package etc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bytedance/sonic"
)

func JoernCommand(ctx context.Context, url, command string) (any, error) {
	body := map[string]any{"query": command}

	bodyByte, _ := sonic.Marshal(body)
	resp, status, err := CustomCall(ctx, "POST", bodyByte, url, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, errors.New("server is unavailable or bad request")
	}
	response:=map[string]any{}
	sonic.Unmarshal(resp, &response)
	if _,ok:=response["stdout"];ok{
		return ParseJoernStdout(resp)
	}
	return response, nil
}
func CustomCall(ctx context.Context, method string, data []byte, url string, headers map[string]string) ([]byte, int, error) {
	bodyReader := bytes.NewReader(data)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancel()
	nreq, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), strings.TrimSpace(url), bodyReader)
	if err != nil {
		return nil, 0, err
	}

	// Add default headers only if not provided in custom headers
	if _, exists := headers["Accept"]; !exists {
		nreq.Header.Set("Accept", `application/json`)
	}
	if _, exists := headers["Content-Type"]; !exists {
		nreq.Header.Set("Content-Type", `application/json`)
	}

	// Add/override with custom headers (use Set to replace, not Add to append)
	for k, v := range headers {
		nreq.Header.Set(k, v)
	}

	client := &http.Client{}
	req, err := client.Do(nreq)
	if err != nil {
		fmt.Println("call api error", err)
		return nil, -1, err
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("response error : ", err)
		return nil, req.StatusCode, err
	}
	return body, req.StatusCode, err

}
