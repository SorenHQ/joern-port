package projectControllers

import (
	"github.com/SorenHQ/joern-port/models"
	gitServices "github.com/SorenHQ/joern-port/services/git"

	"github.com/gofiber/fiber/v2"
)

func cloneRepoHandler(c *fiber.Ctx) error {
	input := models.GitRequest{}
	if err := c.BodyParser(&input); err != nil {
		return c.JSON(models.Response{Data: nil, Error: fiber.ErrBadRequest})
	}
	gitchan := make(chan models.GitResponse)

	git:=gitServices.NewGitDetailsHandler(gitchan)
	go git.ClonePull(input.Project, input.RepoURL, input.Pull)

	return c.JSON(models.Response{Data: "accepted", Error: nil})
}
