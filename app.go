package main

import (
	"fmt"
	"joern-output-parser/actions"
	wsHandler "joern-output-parser/actions/ws"
	"joern-output-parser/env"
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// http server config
	app := fiber.New(fiber.Config{
		Prefork:               false,
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
		AppName:               "Joern-Proxy",
		ServerHeader:          "JoernProxy",
	})
	app.Use(cors.New())

	app.Use(logger.New())
	resultHandler,err:=wsHandler.NewResultHandlers("http://localhost:8081",&MessageHandler{})
	if err!=nil{
		log.Fatalf("Could not create result handler: %v", err)
	}
	actions.RegisterAll(app,resultHandler)
	if err := app.Listen(env.GetPort()); err != nil {
		log.Fatalf("Could not listen: %v", err)
	}
}



type MessageHandler struct {}

func (h *MessageHandler) Recv(message string)  {
	// Implement message handling logic
	fmt.Println("Received message:", message)

}
