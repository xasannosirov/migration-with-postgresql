package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"migration/models"
	"migration/storage"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/user/create", CreateUser)
	router.PUT("/user/update", UpdateUser)
	router.DELETE("/user/delete", DeleteUser)
	router.GET("/user/get", GetUser)
	router.GET("/user/all", GetAllUsers)
	router.GET("/user/getrole", GetUserByRole)
	log.Println("Server is running...")
	if err := router.Run("localhost:8080"); err != nil {
		fmt.Println("Error while running server!")
	}
}

func CreateUser(c *gin.Context) {
	bodyByte, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("error while getting body", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var user *models.User
	if err = json.Unmarshal(bodyByte, &user); err != nil {
		log.Println("error while unmarshalling body", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	respUser, err := storage.CreateUser(user)
	if err != nil {
		log.Println("error while creating user", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusCreated, respUser)
}

func UpdateUser(c *gin.Context) {
	bodyByte, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("error while getting body", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var user *models.User
	if err = json.Unmarshal(bodyByte, &user); err != nil {
		log.Println("error while unmarshalling body", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	userId := c.Query("id")

	respUser, err := storage.UpdateUser(userId, user)
	if err != nil {
		log.Println("error while updating user", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, respUser)
}

func DeleteUser(c *gin.Context) {
	userId := c.Query("id")

	if err := storage.DeleteUser(userId); err != nil {
		log.Println("error while deleting user", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, "Deleted User")
}

func GetUser(c *gin.Context) {
	userId := c.Query("id")

	respUser, err := storage.GetUser(userId)
	if err != nil {
		log.Println("Error while getting user", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, respUser)
}

func GetAllUsers(c *gin.Context) {
	page := c.Request.URL.Query().Get("page")
	intPage, err := strconv.Atoi(page)
	if err != nil {
		log.Println("Error while converting page")
		c.AbortWithError(http.StatusBadRequest, err)
	}

	limit := c.Request.URL.Query().Get("limit")

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Println("Error while converting limit")
		c.AbortWithError(http.StatusBadRequest, err)
	}

	users, err := storage.GetAllUsers(intPage, intLimit)
	if err != nil {
		log.Println("Error while getting all users", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, users)
}

func GetUserByRole(c *gin.Context) {
	role := c.Query("role")
	role = strings.ToLower(role)

	page := c.Query("page")
	intPage, err := strconv.Atoi(page)
	if err != nil {
		log.Println("Error while converting page")
		c.AbortWithError(http.StatusBadRequest, err)
	}

	limit := c.Query("limit")
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Println("Error while converting limit")
		c.AbortWithError(http.StatusBadRequest, err)
	}

	respUsers, err := storage.GetUsersByRole(role, intPage, intLimit)
	if err != nil {
		log.Println("Error while getting user by role", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, respUsers)
}
