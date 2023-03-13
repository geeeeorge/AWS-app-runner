package handler

import (
	"github.com/labstack/echo/v4"
)

type Interface interface {
	GetApiHealthz(ec echo.Context) error
}
