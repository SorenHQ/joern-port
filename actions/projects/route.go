package projectControllers

import (
	"github.com/SorenHQ/joern-port/models"

	validation "github.com/mehdi-shokohi/fiberValidation"

	"github.com/gofiber/fiber/v2"
)

func GitProjectsRoutes(api fiber.Router,) {
	api.Post("/git/repo", validation.ValidateBodyAs(models.GitRequest{}), cloneRepoHandler)

}
