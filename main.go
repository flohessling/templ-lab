package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flohessling/templ-lab/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/joho/godotenv"
)

func main() {
	env := godotenv.Load()
	if env != nil {
		panic("cannot load environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":42069"
	} else {
		port = ":" + port
	}

	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})

	app.Use(compress.New())

	routes.SetRoutes(app)

	// start a server and listen for a shutdown
	go func() {
		// when testing use localhost to prevent macos firewall popup
		host := "localhost" + port
		if err := app.Listen(host); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // block the main thread until interrupt
	app.Shutdown()
	fmt.Println("shutting down server")
}
