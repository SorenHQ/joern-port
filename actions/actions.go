package actions

import (
	"errors"
	"fmt"
	"joern-output-parser/etc"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func queryhandler(c *fiber.Ctx) error {
	input := CommandRequest{}
	if err := c.BodyParser(&input); err != nil {
		return c.JSON(Response{Data: nil, Error: fiber.ErrBadRequest})
	}
	res, err := etc.JoernCommand(c.Context(), input.Url, input.Query)
	if err!=nil{
		return c.JSON(Response{Data: res, Error: err.Error()})
	}
	return c.JSON(Response{Data: res, Error: nil})
}

func openProjectHandler(c *fiber.Ctx) error {
	input := OpenProjectRequest{}
	if err := c.BodyParser(&input); err != nil {
		return c.JSON(Response{Data: nil, Error: fiber.ErrBadRequest})
	}
	body := map[string]any{"query": fmt.Sprintf(`open("%s")`, input.Project)}
	bodyByte, _ := sonic.Marshal(body)
	resp, status, err := etc.CustomCall(c.Context(), "POST", bodyByte, input.Url, nil)
	if err != nil {
		return c.JSON(Response{Data: nil, Error: fiber.ErrBadRequest})

	}
	if status != 200 {
		return c.JSON(Response{Data: nil, Error: errors.New("server is unavailable or bad request")})
	}
	return c.JSON(Response{Data: resp, Error: err})
}

type CommandRequest struct {
	Url   string `json:"url"`
	Query string `json:"query"`
}
type OpenProjectRequest struct {
	Url     string `json:"url"`
	Project string `json:"project"`
}
type Response struct {
	Data  any `json:"data"`
	Error any `json:"error"`
}
