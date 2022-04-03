package controller

import (
	"github.com/labstack/echo"
)

// interface 層から domain 層、usecase 層の呼び出しに関しては問題ない。
// しかし、infrastructure 層の呼び出しはできない。
// よって、直接呼び出すのではなく、interface を定義する。

// InitRouting Initialize Router
func InitRouting(e *echo.Echo, userController UserController) {
	e.GET("/users/:id", userController.Get())
	e.GET("/users", userController.GetAll())
	e.POST("/users", userController.Post())
	e.PUT("/users/:id", userController.Put())
	e.DELETE("/users/:id", userController.Delete())
}
