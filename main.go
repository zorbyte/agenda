package main

import (
	// "fmt"
	"log"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	// "github.com/joho/godotenv"
	"github.com/markbates/pkger"
)

func main() {
	log.SetPrefix("agenda ")
	log.Println("Opening db @ ./test")

	db, err := badger.Open(badger.DefaultOptions("./db_data"))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	engine := html.NewFileSystem(pkger.Dir("/views"), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Listen("127.0.0.1:8080")
}
