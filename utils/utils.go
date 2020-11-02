package utils

import "github.com/gofiber/fiber/v2"

// RESTMsg is an easy to use msg or error to use that automatically uses JSON when processed.
type RESTMsg struct {
	Code uint16 `json:"code"`
	Msg  string `json:"message"`
}

// SendRESTMsg sends a RESTMsg to the client with the appropriate headers.
func SendRESTMsg(c *fiber.Ctx, msg RESTMsg) error {
	return c.Status(int(msg.Code)).JSON(msg)
}

// GetTaskName gets the name of the task from the path segments.
func GetTaskName(c *fiber.Ctx) (string, *RESTMsg) {
	taskName := c.Params("name")
	if taskName == "" {
		return "", &RESTMsg{
			Code: 400,
			Msg:  "Field \"name\" was not supplied",
		}
	}

	return taskName, nil
}
