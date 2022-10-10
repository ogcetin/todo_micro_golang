package routes

import (
	"todomicro/controllers"
	"todomicro/core"

	"github.com/gin-gonic/gin"
)

func Prepare(r *gin.Engine) {
	api := r.Group("/api")
	{
		//
		api.POST("/user/login", controllers.UserLogin)
		secured := api.Group("/").Use(core.JWTAuth())
		{
			secured.GET("/user", controllers.UserAll)
			secured.GET("/user/:user_id", controllers.UserOne)
			secured.POST("/user", controllers.UserAdd)
			secured.PUT("/user/:user_id", controllers.UserUpdate)
			secured.DELETE("/user/:user_id", controllers.UserDelete)

			secured.GET("/todo", controllers.TodoAll)
			secured.GET("/todo/:todo_id", controllers.TodoOne)
			secured.POST("/todo", controllers.TodoAdd)
			secured.PUT("/todo/:todo_id", controllers.TodoUpdate)
			secured.DELETE("/todo/:todo_id", controllers.TodoDelete)
		}

	}
}
