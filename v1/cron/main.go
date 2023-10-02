package cron

import "github.com/labstack/echo/v4"

func Routes(r *echo.Group) {
	r.GET("/cron/clear-pwd-reset", clearPasswordResetTokens)
	r.GET("/cron/sync-devices", syncDevices)
}
