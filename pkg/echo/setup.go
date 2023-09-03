package echo

import (
	"MedKick-backend/pkg/echo/dto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"

	_ "MedKick-backend/docs"
)

func Engine() *echo.Echo {
	e := echo.New()

	corsConfig := middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: true,
	}

	e.Use(middleware.CORSWithConfig(corsConfig))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
	}))

	return e
}

func Swagger(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}

// OnlineCheck godoc
// @Summary Check if API is online
// @Description Check if API is online
// @Tags General
// @Accept json
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Router / [get]
func OnlineCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "MedKick API is online",
	})
}
