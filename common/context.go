package common

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type DawnCtx struct {
	FiberCtx *fiber.Ctx
}

func (ctx DawnCtx) INFO(message string) {
	INFO(ctx.FiberCtx, message)
}

func (ctx DawnCtx) DEBUG(message string) {
	DEBUG(ctx.FiberCtx, message)
}

func (ctx DawnCtx) TRACE(message string) {
	TRACE(ctx.FiberCtx, message)
}

func BuildCtx(c *fiber.Ctx) DawnCtx {
	return DawnCtx{
		FiberCtx: c,
	}
}

func (ctx DawnCtx) BodyParser(out interface{}) error {
	return ctx.FiberCtx.BodyParser(out)
}

var UNAUTHORIZED_TO_USER_ID = &DawnError{
	Name:        "UNAUTHORIZED_TO_USER_ID",
	Description: "user is not authorized to access this endpoint",
	Code:        403,
}

func (ctx DawnCtx) ValidateToUser(userId string) DawnCtx {
	if viper.GetBool("app.auth") {
		admin, _ := strconv.ParseBool(string(ctx.FiberCtx.Request().Header.Peek("admin")))
		if string(ctx.FiberCtx.Request().Header.Peek("user_id")) != userId && !admin {
			panic(UNAUTHORIZED_TO_USER_ID)
		}
	}
	return ctx
}

func (ctx DawnCtx) ValidateToAdmin(userId string) DawnCtx {
	if viper.GetBool("app.auth") {
		admin, _ := strconv.ParseBool(string(ctx.FiberCtx.Request().Header.Peek("admin")))
		if !admin {
			panic(UNAUTHORIZED_TO_USER_ID)
		}
	}
	return ctx
}
