package main

import (
	"fmt"
	"log"

	joernControllers "github.com/SorenHQ/joern-port/actions/joern"
	wsHandler "github.com/SorenHQ/joern-port/actions/joern/ws"
	projectControllers "github.com/SorenHQ/joern-port/actions/projects"
	"github.com/SorenHQ/joern-port/env"
	"github.com/SorenHQ/joern-port/models"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// http server config
	app := fiber.New(fiber.Config{
		Prefork:      false,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		AppName:      "Joern-Proxy",
		ServerHeader: "JoernProxy",
	})
	app.Use(cors.New())

	app.Use(logger.New())
	resultHandler, err := wsHandler.NewJoernResultHandlers(env.GetJoernUrl(), &MessageHandler{})
	if err != nil {
		log.Fatalf("Could not create result handler: %v", err)
	}
	//Joern Api and websocket handler
	joernControllers.JoernRouter(app, resultHandler)

	// Git Api and Status Handler
	GitLogger := make(chan models.GitResponse)
	go gitStatus(GitLogger)
	projectControllers.GitProjectsRoutes(app, GitLogger)

	// Start Server
	if err := app.Listen(env.GetPort()); err != nil {
		log.Fatalf("Could not listen: %v", err)
	}
}

type MessageHandler struct{}

func (h *MessageHandler) Recv(message string) {
	// Implement message handling logic
	fmt.Println("Received Joern WebSocket message:", message)

}

func gitStatus(logger chan models.GitResponse) {
	for {
		status := <-logger
		fmt.Println("Received GIT status:", status)
	}
}
