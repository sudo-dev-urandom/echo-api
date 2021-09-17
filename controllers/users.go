package controllers

import (
	"fmt"
	"goecho/core"
	"goecho/helpers"
	"goecho/models"
	"io"
	"os"
	"path/filepath"

	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
)

type (
	CreateUserRequest struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
)

// UserList Get list of users
// @Summary Get list of users
// @Description Get list of users
// @Tags users
// @ID users-user-list
// @Accept mpfd
// @Produce plain
// @Param page query integer false "page number" default(1)
// @Param pageSize query integer false "number of bookings in single page" default(10)
// @Success 200 {object} []models.User
// @Router /api/v1/users [get]
func UserList(c echo.Context) error {
	defer c.Request().Body.Close()

	users := models.User{}
	rows, err := strconv.Atoi(c.QueryParam("pageSize"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := "created_at"
	sort := "DESC"

	var filter struct{}
	result, err := users.PagedFilterSearch(page, rows, orderby, sort, &filter)

	if err != nil {
		return helpers.Response(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, result.Data)
}

// UserCreate Create user
// @Summary Create user
// @Description Create user
// @Tags users
// @ID users-user-create
// @Accept mpfd
// @Produce plain
// @Param name formData string true "User name"
// @Param address formData string true "User address"
// @Success 200 {object} models.User{}
// @Router /api/v1/users/create [post]
func UserCreate(c echo.Context) error {
	defer c.Request().Body.Close()

	payloadRules := govalidator.MapData{
		"name":          []string{"required"},
		"username":      []string{"username"},
		"email":         []string{"email"},
		"password":      []string{"password"},
		"address":       []string{"address"},
		"phone_number":  []string{"phone_number"},
		"status_active": []string{"status_active"},
		"is_partner":    []string{"is_partner"},
	}

	validate := helpers.ValidateRequestFormData(c, payloadRules)
	if validate != nil {
		return helpers.Response(http.StatusUnprocessableEntity, validate, "Validation error")
	}

	user := models.User{
		Name:    c.FormValue("name"),
		Address: c.FormValue("address"),
	}

	tx := core.App.DB.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return helpers.Response(http.StatusUnprocessableEntity, err, "Error while saving the User data")
	}

	tx.Commit()

	return c.JSON(http.StatusCreated, user)
}

// UserUpdate Update user
// @Summary Update user
// @Description Update user
// @Tags users
// @ID users-user-update
// @Accept mpfd
// @Produce plain
// @Param id path int true "ID"
// @Param RequestBody body CreateUserRequest true "JSON Request Body"
// @Success 200 {object} models.User{}
// @Router /api/v1/users/update/{id} [patch]
func UserUpdate(c echo.Context) error {
	defer c.Request().Body.Close()

	payloadRules := govalidator.MapData{
		"name":    []string{"required"},
		"address": []string{"address"},
	}

	userID, _ := strconv.Atoi(c.Param("id"))

	user := models.User{
		Name:    c.FormValue("name"),
		Address: c.FormValue("address"),
	}

	if err := user.FindbyID(userID); err != nil {
		return helpers.Response(http.StatusNotFound, err, "User not found")
	}

	validate := helpers.ValidateRequestFormData(c, payloadRules)
	if validate != nil {
		return helpers.Response(http.StatusUnprocessableEntity, validate, "Validation error")
	}

	tx := core.App.DB.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return helpers.Response(http.StatusUnprocessableEntity, err, "Error while saving the user data")
	}

	if _, err := c.FormFile("profile_picture"); err != nil {
		tx.Rollback()
		return helpers.Response(http.StatusUnprocessableEntity, err, "Error while uploading profile picture")
	}

	file, err := c.FormFile("profile_picture")
	if err != nil {
		tx.Rollback()
		return helpers.Response(http.StatusUnprocessableEntity, err, "Error while uploading profile picture")
	}

	src, err := file.Open()
	if err != nil {
		tx.Rollback()
		return helpers.Response(http.StatusUnprocessableEntity, err, "Error while uploading profile picture")
	}
	defer src.Close()

	filename := fmt.Sprintf("%d%s", userID, filepath.Ext(file.Filename))
	fileLocation := filepath.Join("assets", filename)

	dst, err := os.Create(fileLocation)
	if err != nil {
		tx.Rollback()
		return helpers.Response(http.StatusUnprocessableEntity, err, "Error while uploading profile picture")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		tx.Rollback()
		return helpers.Response(http.StatusUnprocessableEntity, err, "Error while uploading profile picture")
	}

	if err := user.Save(); err != nil {
		return helpers.Response(http.StatusUnprocessableEntity, err, "Error while updating User data")
	}

	return c.JSON(http.StatusOK, user)
}

// UserDelete Delete user
// @Summary Delete user
// @Description Delete user
// @Tags users
// @ID users-user-delete
// @Accept mpfd
// @Produce plain
// @Param id path int true "ID"
// @Success 200 {object} models.User{}
// @Router /api/v1/users/delete/{id} [delete]
func UserDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	userID, _ := strconv.Atoi(c.Param("id"))
	user := models.User{}
	if err := user.FindbyID(userID); err != nil {
		return helpers.Response(http.StatusNotFound, err, "User not found")
	}

	if err := user.Delete(); err != nil {
		return helpers.Response(http.StatusUnprocessableEntity, err, "Delete user failed")
	}

	return c.JSON(http.StatusOK, user)
}
