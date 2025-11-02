package actions

import (
	wsHandler "joern-output-parser/actions/ws"

	"github.com/gofiber/fiber/v2"
)
var resultHandler *wsHandler.ResultHandlers
	func RegisterAll(api fiber.Router,wsHandler *wsHandler.ResultHandlers) {
		resultHandler = wsHandler
		api.Post("/open",openProjectHandler)
		api.Post("/query",queryhandler)
	}
