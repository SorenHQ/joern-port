package main

import (
	"joern-output-parser/actions"
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
	actions.RegisterAll(app)
	if err := app.Listen(env.GetPort()); err != nil {
		log.Fatalf("Could not listen: %v", err)
	}
}
