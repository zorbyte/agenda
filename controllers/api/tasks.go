package apictrs

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/zorbyte/agenda/utils"
)

// Task is a task for the program.
type Task struct {
	Name    string `json:"name" validate:"required,min=3,max=32"`
	Content string `json:"content" validate:"required"`
}

// GetTask gets a task.
func GetTask(db *badger.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		taskName, err := utils.GetTaskName(c)
		if err != nil {
			return utils.SendRESTMsg(c, *err)
		}

		return db.View(func(txn *badger.Txn) error {
			taskItem, err := txn.Get([]byte(taskName))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return utils.SendRESTMsg(c, utils.RESTMsg{
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

// AddTask adds a task.
func AddTask(db *badger.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		t := new(Task)
		if err := c.BodyParser(t); err != nil {
			return utils.SendRESTMsg(c, utils.RESTMsg{
				Code: 400,
				Msg:  err.Error(),
			})
		}

		return db.Update(func(txn *badger.Txn) error {
			_, err := txn.Get([]byte(t.Name))
			if err != nil {
				// Ensure that it doesn't already exist.
				if err == badger.ErrKeyNotFound {
					taskEntry := badger.NewEntry([]byte(t.Name), []byte(t.Content))
					err := txn.SetEntry(taskEntry)
					return err
				}

				return err
			}

			return utils.SendRESTMsg(c, utils.RESTMsg{
				Code: 409,
				Msg:  "Task already exists",
			})
		})
	}
}

// RemTask removes a task.
func RemTask(db *badger.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		taskName, err := utils.GetTaskName(c)
		if err != nil {
			return utils.SendRESTMsg(c, *err)
		}

		return db.Update(func(txn *badger.Txn) error {
			encodedTaskName := []byte(taskName)
			_, err := txn.Get(encodedTaskName)
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return utils.SendRESTMsg(c, utils.RESTMsg{
						Code: 404,
						Msg:  fmt.Sprintf("The task with the name \"%v\" does not exist", taskName),
					})
				}

				return err
			}

			txn.Delete(encodedTaskName)

			return utils.SendRESTMsg(c, utils.RESTMsg{
				Code: 200,
				Msg:  "Successfully deleted task",
			})
		})
	}
}
