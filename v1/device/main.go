package device

import "github.com/labstack/echo/v4"

func Routes(r echo.Group) {
	r.POST("/connect", IngestData)
}
