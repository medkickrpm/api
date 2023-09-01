package cron

import "github.com/labstack/echo/v4"

// TODO - Protect to only allow the system to call this
func Routes(r *echo.Group) {
	r.GET("/cron/clear-pwd-reset", clearPasswordResetTokens)
}
