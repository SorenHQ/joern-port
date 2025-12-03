package main

import (
	"context"
	"fmt"
	"log"

	joernControllers "github.com/SorenHQ/joern-port/actions/joern"
	wsHandler "github.com/SorenHQ/joern-port/actions/joern/ws"
	projectControllers "github.com/SorenHQ/joern-port/actions/projects"
	"github.com/SorenHQ/joern-port/env"
	"github.com/SorenHQ/joern-port/etc/db"
	"github.com/SorenHQ/joern-port/models"
	joernPlugin "github.com/SorenHQ/joern-port/sorenPlugin"

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
	go joernPlugin.LoadSorenPluginServer()
	// Start Server
	if err := app.Listen(env.GetPort()); err != nil {
		log.Fatalf("Could not listen: %v", err)
	}
}

type MessageHandler struct{}

func (h *MessageHandler) Recv(req_uuid,message string) {
	// Implement message handling logic
	fmt.Println("Received Joern WebSocket message:", message)
	db.GetRedisClient().Publish(context.Background(),joernPlugin.JoernResultsTableInRedis, req_uuid+"||"+message)

}

func gitStatus(logger chan models.GitResponse) {
	for {
		status := <-logger
		fmt.Println("Received GIT status:", status)
		//call python api 
	}
}
