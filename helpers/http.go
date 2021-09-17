package helpers

import (
	"github.com/labstack/echo/v4"
)

func Response(statuscode int, data interface{}, message string) error {
	response := map[string]interface{}{
		"code":    statuscode,
		"message": message,
		"data":    data,
	}

	return echo.NewHTTPError(statuscode, response)
}
