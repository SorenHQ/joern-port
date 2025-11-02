package joernControllers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/SorenHQ/joern-port/env"
	"github.com/SorenHQ/joern-port/etc"
	"github.com/SorenHQ/joern-port/models"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func queryhandler(c *fiber.Ctx) error {
	input := models.CommandRequest{}
	if err := c.BodyParser(&input); err != nil {
		return c.JSON(models.Response{Data: nil, Error: fiber.ErrBadRequest})
	}
	if input.Url == "" {
		input.Url = env.GetJoernUrl()
	}
	url := fmt.Sprintf("%s/query-sync", input.Url)
	if input.Mode == "async" {
		url = fmt.Sprintf("%s/query", input.Url)
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
	if input.Url == "" {
		input.Url = env.GetJoernUrl()
	}
	dir := fmt.Sprintf("%s/%s", env.GetProjectReposPath(), input.Project)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return c.JSON(models.Response{Data: nil, Error: "Project not found"})
	}
	url := fmt.Sprintf("%s/query-sync", input.Url)
	importresult, err := importCode(c.Context(), dir, input.Project, url)
	if err != nil {
		return c.JSON(models.Response{Data: importresult, Error: err.Error()})
	}
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

func importCode(ctx context.Context, dir, project, url string) (map[string]any, error) {

	body := map[string]any{"query": fmt.Sprintf(`importCode("%s", "%s")`, dir, project)}
	bodyByte, _ := sonic.Marshal(body)
	resp, status, err := etc.CustomCall(ctx, "POST", bodyByte, url, nil)
	if err != nil {
		return nil, err

	}
	if status != 200 {
		return nil, errors.New("server is unavailable or bad request")
	}
	respMap := map[string]any{}
	err = sonic.Unmarshal(resp, &respMap)
	if err != nil {
		return nil, err
	}

	return respMap, nil
}
