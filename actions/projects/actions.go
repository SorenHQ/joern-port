package projectControllers

import (
	"joern-output-parser/models"
	gitServices "joern-output-parser/services/git"

	"github.com/gofiber/fiber/v2"
)


func cloneRepoHandler(c *fiber.Ctx)error{
	input := models.GitRequest{}
	if err := c.BodyParser(&input); err != nil {
		return c.JSON(models.Response{Data: nil, Error: fiber.ErrBadRequest})
	}
	err:=gitServices.GitClonePull(input.Project,input.RepoURL,input.Pull)
	if err != nil {
		return c.JSON(models.Response{Data: nil, Error: err})
	}
	return c.JSON(models.Response{Data: "accepted", Error: nil})
}

