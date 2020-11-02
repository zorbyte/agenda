package apirts

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/gofiber/fiber/v2"
	apictrs "github.com/zorbyte/agenda/controllers/api"
)

// RegisterTasks no bloody errors!
func RegisterTasks(app *fiber.App, db *badger.DB) {
	app.Get("/tasks/:name", apictrs.GetTask(db))
	app.Delete("/tasks/:name", apictrs.GetTask(db))
	app.Post("/tasks/add", apictrs.AddTask(db))
}
