package controllers

import (
	"net/http"
	"strconv"
	"todomicro/core"
	"todomicro/models"

	"todomicro/database"

	"github.com/gin-gonic/gin"
)

type User = models.User

func UserLogin(c *gin.Context) {
	type LoginForm struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required,min=3,max=24"`
	}
	var user LoginForm
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}

	row := database.DB.QueryRow("select id from user where phone=? and password=?;", user.Phone, user.Password)

	var user_id int
	err := row.Scan(&user_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "", "message": "Something went wrong", "data": err})
		return
	}

	jwt_parsed := core.JWT{Phone: user.Phone, Password: user.Password}

	token, err := core.CreateToken(jwt_parsed)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "", "message": "Token not created", "data": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok", "data": token})
}

func UserAll(c *gin.Context) {
	rows, dbErr := database.DB.Query("select id, name, email, phone from user")
	if dbErr != nil {
		c.JSON(http.StatusOK, gin.H{"status": "", "message": "User Query Error", "data": dbErr})
		return
	}

	var userList []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "", "message": "User Data Not Found", "data": err})
			return
		} else {
			userList = append(userList, user)
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok", "data": userList})

}
func UserOne(c *gin.Context) {
	user_id := c.Param("user_id")
	UserID, _ := strconv.Atoi(user_id)

	row := database.DB.QueryRow("select id, name, email, phone from user where id=?;", UserID)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Phone)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Data Not Found", "data": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "done", "data": user})
}

func UserAdd(c *gin.Context) {
	type UserAddForm struct {
		Name   string `json:"name" binding:"required,min=3"`
		Status string `json:"status"`
		Email  string `json:"email" binding:"required,min=3"`
		Phone  string `json:"phone" binding:"required,min=3"`
	}
	var user UserAddForm
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}

	user.Status = "active"

	result, err := database.DB.Exec("INSERT INTO user (name, status, email, phone) VALUES (?, ?, ?, ?)", user.Name, user.Status, user.Email, user.Phone)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	user_id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok", "data": user_id})
}

func UserUpdate(c *gin.Context) {
	user_id := c.Param("user_id")
	UserID, _ := strconv.Atoi(user_id)

	type UserUpdateForm struct {
		ID     int    `json:"id" binding:"required,numeric"`
		Phone  string `json:"phone" binding:"required"`
		Name   string `json:"name" binding:"required"`
		Status string `json:"status" binding:"required"`
		Email  string `json:"email" binding:"required,email"`
	}

	var user UserUpdateForm
	user.ID = UserID
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}

	_, err := database.DB.Exec("update user set phone=?, name=?, status=?, email=? where id=?", user.Phone, user.Name, user.Status, user.Email, UserID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok", "data": user_id})
}

func UserDelete(c *gin.Context) {
	user_id := c.Param("user_id")
	UserID, _ := strconv.Atoi(user_id)
	_, err := database.DB.Exec("delete from user where id=?", UserID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok", "data": user_id})

}
