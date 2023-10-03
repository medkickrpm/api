package device

import (
	"MedKick-backend/pkg/echo/middleware"
	"github.com/labstack/echo/v4"
)

func Routes(r *echo.Group) {
	r.POST("/connect", ingestData)
	r.GET("/device/:id", getDevice, middleware.NotGuest)
}
