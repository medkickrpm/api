package cron

import "github.com/labstack/echo/v4"

func Routes(r *echo.Group) {
	r.POST("/cron/clear-pwd-reset", clearPasswordResetTokens)
	r.POST("/cron/sync-devices", syncDevices)
	r.POST("/cron/trigger-cpt-worker", triggerCptWorker)
	r.POST("/cron/clear-test-billings", clearTestBillings)
}
