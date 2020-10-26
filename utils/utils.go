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
