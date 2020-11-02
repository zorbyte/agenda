package utils

import "github.com/gofiber/fiber/v2"

// RESTError is an easy error to use that automatically uses JSON when processed.
type RESTError struct {
	Code uint16 `json:"code"`
	Msg  string `json:"message"`
}

// SendRESTErr sends a RESTError to the client with the appropriate headers.
func SendRESTErr(c *fiber.Ctx, err RESTError) error {
	return c.Status(int(err.Code)).JSON(err)
}

// GetTaskName gets the name of the task from the path segments.
func GetTaskName(c *fiber.Ctx) (string, *RESTError) {
	taskName := c.Params("name")
	if taskName == "" {
		return "", &RESTError{
			Code: 400,
			Msg:  "Field \"name\" was not supplied",
		}
	}

	return taskName, nil
}
