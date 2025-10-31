package actions
import(
				"github.com/gofiber/fiber/v2"

)
	func RegisterAll(api fiber.Router) {
		api.Post("/open",openProjectHandler)
		api.Post("/query",queryhandler)
	}
