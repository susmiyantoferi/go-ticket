package routes

import (
	"ticket/handler"
	"ticket/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(UserHandler handler.UserHandler, EventHandler handler.EventHandler) *gin.Engine {
	router := gin.Default()

	public := router.Group("/api/v1/")
	{
		public.POST("users-register", UserHandler.Create)
		public.POST("login", UserHandler.Login)
		public.POST("refresh-token", UserHandler.TokenRefresh)
	}

	api := router.Group("/api/v1")
	api.Use(middleware.Authentication())
	{
		admin := api.Group("/")
		admin.Use(middleware.RoleAccessMiddleware("admin"))
		{
			admin.DELETE("users/:userId", UserHandler.Delete)
			admin.GET("users", UserHandler.FindAll)

			admin.POST("events", EventHandler.Create)
			admin.PUT("events/:id", EventHandler.Update)
			admin.DELETE("events/:id", EventHandler.Delete)
			admin.GET("events/:id", EventHandler.FindById)
		}

		cust := api.Group("/")
		cust.Use(middleware.RoleAccessMiddleware("customer", "admin"))
		{
			cust.PUT("users", UserHandler.Update)
			cust.GET("users/find", UserHandler.FindById)

			cust.GET("events", EventHandler.FindAll)
		}
	}

	return router

}
