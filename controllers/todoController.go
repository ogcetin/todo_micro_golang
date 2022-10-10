package controllers

import (
	"net/http"
	"strconv"
	"todomicro/core"
	"todomicro/models"

	"time"
	"todomicro/database"

	"github.com/gin-gonic/gin"
)

type Todo = models.Todo

func TodoAll(c *gin.Context) {

	rows, dbErr := database.DB.Query("select id, user_id, detail, creation_date, last_update, status from todo")
	if dbErr != nil {
		c.JSON(http.StatusOK, gin.H{"status": "", "message": "Todo Query Error", "data": dbErr})
		return
	}

	var todoList []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.UserID, &todo.Detail, &todo.CreationDate, &todo.LastUpdate, &todo.Status)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "", "message": "Todo Data Not Found", "data": err})
			return
		} else {
			todoList = append(todoList, todo)
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok", "data": todoList})

}
func TodoOne(c *gin.Context) {
	todo_id := c.Param("todo_id")
	TodoID, _ := strconv.Atoi(todo_id)

	row := database.DB.QueryRow("select id, detail, creation_date, last_update  from todo where id=?;", TodoID)

	var todo Todo
	err := row.Scan(&todo.ID, &todo.Detail, &todo.CreationDate, &todo.LastUpdate)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Todo Data Not Found", "data": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "done", "data": todo})
}

func TodoAdd(c *gin.Context) {
	type TodoAddForm struct {
		UserID string `json:"user_id"`
		Detail string `json:"detail" binding:"required,min=3,max=1000"`
		Status string `json:"status"`
	}
	var todo TodoAddForm
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}

	todo.Status = "active"
	var creationDate = time.Now().Format("2006-01-02 15:04:05")
	var user User = core.GetLoggedUser()

	result, err := database.DB.Exec("INSERT INTO todo (user_id,detail,creation_date,last_update,status) VALUES (?, ?, ?, ?, ?)", user.ID, todo.Detail, creationDate, creationDate, todo.Status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
	}

	todo_id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok", "data": todo_id})
}

func TodoUpdate(c *gin.Context) {
	todo_id := c.Param("todo_id")
	TodoID, _ := strconv.Atoi(todo_id)

	type todoUpdateForm struct {
		ID     int    `json:"id" binding:"required,numeric"`
		Detail string `json:"detail"`
		Status string `json: "status"`
	}

	var todo todoUpdateForm
	todo.ID = TodoID
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}

	if len(todo.Status) < 1 {
		_, updateErr := database.DB.Exec("update todo set detail=? where id=?", todo.Detail, TodoID)
		if updateErr != nil {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": updateErr.Error(), "data": ""})
			return
		}

	} else if len(todo.Detail) < 1 {
		_, updateErr := database.DB.Exec("update todo set status=? where id=?", todo.Status, TodoID)
		if updateErr != nil {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": updateErr.Error(), "data": ""})
			return
		}

	} else {
		_, updateErr := database.DB.Exec("update todo set detail=?, status=? where id=?", todo.Detail, todo.Status, TodoID)
		if updateErr != nil {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": updateErr.Error(), "data": ""})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok", "data": todo_id})
}

func TodoDelete(c *gin.Context) {
	todo_id := c.Param("todo_id")
	TodoID, _ := strconv.Atoi(todo_id)
	_, err := database.DB.Exec("delete from todo where id=?", TodoID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "data": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok", "data": todo_id})

}
