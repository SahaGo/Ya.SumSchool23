package controllers

import (
	cerrors "Ya.SumSchool23/controllers/errors"
	"Ya.SumSchool23/rate_limiter"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PingController struct {
	rateLimiter *rate_limiter.RateLimiter
}

func NewPingController() *PingController {
	return &PingController{
		rateLimiter: rate_limiter.NewRateLimiter(),
	}
}

func (c *PingController) Ping(ctx echo.Context) error {
	if !c.rateLimiter.RegisterCall() {
		return cerrors.TooManyRequests.Newf("ping method overloaded")
	}
	return ctx.String(http.StatusOK, "pong")
}
