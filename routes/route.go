package routes

import (
	"ticket/handler"
	"ticket/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(UserHandler handler.UserHandler) *gin.Engine {
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
		}

		cust := api.Group("/")
		cust.Use(middleware.RoleAccessMiddleware("customer", "admin"))
		{
			cust.PUT("users", UserHandler.Update)
			cust.GET("users/find", UserHandler.FindById)
		}
	}

	return router

}
