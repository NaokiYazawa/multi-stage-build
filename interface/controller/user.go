package controller

// router から呼び出され
// - リクエストで渡されたデータを usecase 層へ受け渡す。
// - 戻り値として受け取ったデータを、JSON 形式で返す。
// という役割を担っている。

import (
	"net/http"
	"strconv"

	"multi-stage-build/usecase"

	"github.com/labstack/echo"
)

// UserController Interface of an user controller
// userController は 下記の interface で定義されているメソッドが実装されている
type UserController interface {
	Post() echo.HandlerFunc
	Get() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	Put() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

// usecase.UserUsecase は interface
type userController struct {
	// インターフェースの usecase.UserUsecase 型を定義
	userUsecase usecase.UserUsecase
}

// NewUserController Constructor of an user controller
func NewUserController(userUsecase usecase.UserUsecase) UserController {
	// 構造体を返す
	return &userController{userUsecase: userUsecase}
}

// Input Data
type requestUser struct {
	Name string `json:"name"`
}

type responseUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Get Controller of getting user
func (userController *userController) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi((c.Param("id")))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		foundUser, err := userController.userUsecase.ReadByID(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		res := responseUser{
			ID:   foundUser.ID,
			Name: foundUser.Name,
		}

		return c.JSON(http.StatusOK, res)
	}
}

// GetAll Controller of getting users
func (userController *userController) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		foundUsers, err := userController.userUsecase.ReadAll()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		res := []responseUser{}
		for _, foundUser := range *foundUsers {
			res = append(res, responseUser{
				ID:   foundUser.ID,
				Name: foundUser.Name,
			})
		}

		return c.JSON(http.StatusOK, res)
	}
}

// Post Controller of posting user
// アプリケーションが要求するデータに、「入力を変換」 → Controller
func (userController *userController) Post() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Input Data???
		// 入力のデータを生成
		var req requestUser
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// Usecase Create が欲しがっているデータに変換して渡す
		// Input Boundary を Call
		// 実際に処理を行うのは Use Case Interacter
		createdUser, err := userController.userUsecase.Create(req.Name)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// View Model ??
		res := responseUser{
			ID:   createdUser.ID,
			Name: createdUser.Name,
		}

		// View
		return c.JSON(http.StatusCreated, res)
	}
}

// Put Controller of updating user
func (userController *userController) Put() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		var req requestUser
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		updatedUser, err := userController.userUsecase.Update(id, req.Name)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		res := responseUser{
			ID:   updatedUser.ID,
			Name: updatedUser.Name,
		}

		return c.JSON(http.StatusOK, res)
	}
}

// Delete Controller of deleting user
func (userController *userController) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err = userController.userUsecase.Delete(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
