package apictrs

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/zorbyte/agenda/utils"
)

// Task is a task for the program
type Task struct {
	Name    string `json:"name" validate:"required,min=3,max=32"`
	Content string `json:"content" validate:"required"`
}

// GetTask gets a task to the db
func GetTask(db *badger.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		taskName, err := utils.GetTaskName(c)
		if err != nil {
			return utils.SendRESTErr(c, *err)
		}

		return db.View(func(txn *badger.Txn) error {
			taskItem, err := txn.Get([]byte(taskName))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return utils.SendRESTErr(c, utils.RESTError{
						Code: 404,
						Msg:  fmt.Sprintf("The task with the name \"%v\" does not exist", taskName),
					})
				}

				return err
			}

			return taskItem.Value(func(val []byte) error {
				realVal := string(val)
				return c.SendString(realVal)
			})
		})
	}
}

// AddTask adds a task to the db
func AddTask(db *badger.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		t := new(Task)

		if err := c.BodyParser(t); err != nil {
			return utils.SendRESTErr(c, utils.RESTError{
				Code: 400,
				Msg:  err.Error(),
			})
		}

		return db.Update(func(txn *badger.Txn) error {
			_, err := txn.Get([]byte(t.Name))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					taskEntry := badger.NewEntry([]byte(t.Name), []byte(t.Content))
					err := txn.SetEntry(taskEntry)
					return err
				}

				return err
			}

			return utils.SendRESTErr(c, utils.RESTError{
				Code: 409,
				Msg:  "Task already exists",
			})
		})
	}
}
