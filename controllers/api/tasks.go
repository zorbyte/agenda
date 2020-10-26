package apictrs

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/zorbyte/agenda/utils"
)

// AddTask adds a task to the db
func AddTask(db *badger.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		taskID := c.Params("id")
		if taskID == "" {
			return utils.SendRESTErr(c, utils.RESTError{
				Code: 400,
				Msg:  "Field \"id\" was not supplied",
			})
		}
	
		return db.View(func(txn *badger.Txn) error {
			testItem, err := txn.Get([]byte(taskID))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return utils.SendRESTErr(c, utils.RESTError{
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
	}

}
