package echo

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
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

func OnlineCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "MedKick API is online",
	})
}
