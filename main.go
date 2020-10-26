package main

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	tasks_controllers "github.com/zorbyte/agenda/controllers/api/tasks"

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

	tasks_controllers.AddTask()

	restRoutes(app, db)
	app.Listen("127.0.0.1:8080")
}

func restRoutes(app *fiber.App, db *badger.DB) {
	log.Println("Registering API routes")
	app.Get("/api/tasks/:id", func(c *fiber.Ctx) error {
		taskID := c.Params("id")
		if taskID == "" {
			return sendError(c, httpError{
				Code: 400,
				Msg:  "Field \"id\" was not supplied",
			})
		}

		return db.View(func(txn *badger.Txn) error {
			testItem, err := txn.Get([]byte(taskID))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return sendError(c, httpError{
						Code: 404,
						Msg:  fmt.Sprintf("The requested task \"%v\" was not found", taskID),
					})
				}

				return err
			}

			return testItem.Value(func(val []byte) error {
				realVal := string(val)
				return c.SendString(realVal)
			})
		})
	})
}

