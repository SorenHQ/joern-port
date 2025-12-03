package projectControllers

import (
	"github.com/SorenHQ/joern-port/models"
	gitServices "github.com/SorenHQ/joern-port/services/git"

	validation "github.com/mehdi-shokohi/fiberValidation"

	"github.com/gofiber/fiber/v2"
)

func GitProjectsRoutes(api fiber.Router, GitLogger chan models.GitResponse) {
	gitServices.NewGitLogHandler(GitLogger)
	api.Post("/git/repo", validation.ValidateBodyAs(models.GitRequest{}), cloneRepoHandler)

}
