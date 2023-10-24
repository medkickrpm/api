package interaction

import (
	"MedKick-backend/pkg/echo/middleware"
	"github.com/labstack/echo/v4"
)

func Routes(r *echo.Group) {
	r.POST("/interaction", createInteraction, middleware.NotGuest, middleware.HasRole("admin", "doctor", "nurse"))
	r.GET("/interaction/:id", getInteraction, middleware.NotGuest)
	r.PATCH("/interaction/:id", updateInteraction, middleware.NotGuest, middleware.HasRole("admin", "doctor", "nurse"))
	r.DELETE("/interaction/:id", deleteInteraction, middleware.NotGuest, middleware.HasRole("admin", "doctor", "nurse"))
}
