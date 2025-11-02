package joernControllers

import (
	wsHandler "github.com/SorenHQ/joern-port/actions/joern/ws"

	"github.com/gofiber/fiber/v2"
)

var resultHandler *wsHandler.ResultHandlers

func JoernRouter(api fiber.Router, wsHandler *wsHandler.ResultHandlers) {
	resultHandler = wsHandler
	api.Post("/open", openProjectHandler)
	api.Post("/query", queryhandler)
}
