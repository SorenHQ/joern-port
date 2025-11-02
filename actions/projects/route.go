package projectControllers

import (
	"joern-output-parser/models"
	gitServices "joern-output-parser/services/git"

	validation "github.com/mehdi-shokohi/fiberValidation"

	"github.com/gofiber/fiber/v2"
)
func GitProjectsRoutes(api fiber.Router, GitLogger chan models.GitResponse) {
	gitServices.NewGitDetailsHandler(GitLogger)
	api.Post("/git/repo",validation.ValidateBodyAs(models.GitRequest{}),cloneRepoHandler)


}

