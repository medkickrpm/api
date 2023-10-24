package careplan

import (
	"MedKick-backend/pkg/echo/middleware"
	"github.com/labstack/echo/v4"
)

func Routes(r *echo.Group) {
	r.POST("/careplan", createCareplan, middleware.NotGuest, middleware.HasRole("nurse", "doctor", "admin"))
	r.GET("/careplan/:id", getCareplan, middleware.NotGuest)
	r.GET("/careplan/:id/file", downloadCareplan, middleware.NotGuest, middleware.HasRole("nurse", "doctor", "admin"))
	r.PUT("/careplan/:id", uploadCareplan, middleware.NotGuest, middleware.HasRole("nurse", "doctor", "admin"))
	r.DELETE("/careplan/:id", deleteCareplan, middleware.NotGuest, middleware.HasRole("nurse", "doctor", "admin"))
}
