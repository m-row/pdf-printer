package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func (c *Controller) Meta(ctx echo.Context) error {
	commitCountInt, err := strconv.ParseInt(CommitCount, 10, 64)
	if err != nil {
		commitCountInt = 0
	}

	meta := map[string]any{
		"domain":      c.Settings.Domain,
		"api_version": Version,
		"commit": map[string]any{
			"describe": CommitDescribe,
			"count":    commitCountInt,
			"sha1":     c.Settings.CommitInfo.FullSHA1,
			"time":     c.Settings.CommitInfo.Time,
		},
		"project_code": c.Settings.AppCode,
		"project_name": c.Settings.AppName,
		"server_time": map[string]any{
			"go_now": TimeNow(),
			"utc_2":  TimeNow().Add(time.Hour * 2),
		},
	}

	return ctx.JSON(http.StatusOK, meta)
}
