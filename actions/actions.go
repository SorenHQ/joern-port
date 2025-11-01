package actions

import (
	"context"
	"errors"
	"fmt"
	wsHandler "joern-output-parser/actions/ws"
	"joern-output-parser/etc"
	"joern-output-parser/models"
	"regexp"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func queryhandler(c *fiber.Ctx) error {
	input := models.CommandRequest{}
	if err := c.BodyParser(&input); err != nil {
		return c.JSON(models.Response{Data: nil, Error: fiber.ErrBadRequest})
	}
	url := fmt.Sprintf("%s/query-sync", input.Url)
	if input.Mode == "async" {
		url = fmt.Sprintf("%s/query", input.Url)
		err := wsHandler.ConnectToServer(input.Url, messageHandler)
		if err != nil {
			return c.JSON(models.Response{Data: nil, Error: err.Error()})
		}
	}
	res, err := etc.JoernCommand(c.Context(), url, input.Query)
	if err != nil {
		return c.JSON(models.Response{Data: res, Error: err.Error()})
	}
	return c.JSON(models.Response{Data: res, Error: nil})
}

func openProjectHandler(c *fiber.Ctx) error {
	input := models.OpenProjectRequest{}
	if err := c.BodyParser(&input); err != nil {
		return c.JSON(models.Response{Data: nil, Error: fiber.ErrBadRequest})
	}
	url := fmt.Sprintf("%s/query-sync", input.Url)
	body := map[string]any{"query": fmt.Sprintf(`open("%s").get.name`, input.Project)}
	bodyByte, _ := sonic.Marshal(body)
	resp, status, err := etc.CustomCall(c.Context(), "POST", bodyByte, url, nil)
	if err != nil {
		return c.JSON(models.Response{Data: nil, Error: fiber.ErrBadRequest})

	}
	if status != 200 {
		return c.JSON(models.Response{Data: nil, Error: errors.New("server is unavailable or bad request")})
	}
	respMap := map[string]any{}
	err = sonic.Unmarshal(resp, &respMap)
	if err != nil {
		return c.JSON(models.Response{Data: nil, Error: err})
	}
	if success, ok := respMap["success"].(bool); ok && success {
		re := regexp.MustCompile(`\s*"([^"]+)\"`)
		out := re.FindStringSubmatch(respMap["stdout"].(string))
		if len(out) > 1 {
			return c.JSON(models.Response{Data: out[1], Error: nil})
		}
	}
	return c.JSON(models.Response{Data: respMap, Error: err})
}

func messageHandler(url string, message []byte) {
	fmt.Println(string(message))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	url = fmt.Sprintf("%s/result/%s", url, string(message))
	response, _, _ := etc.CustomCall(ctx, "GET", nil, url, nil)
	fmt.Println(etc.ParseJoernStdout(response))
}
