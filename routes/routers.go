package routes

import (
	a "nextlaundry_apis/controller/authController"
	member "nextlaundry_apis/controller/memberController"
	outlet "nextlaundry_apis/controller/outletController"
	product "nextlaundry_apis/controller/productController"
	trans "nextlaundry_apis/controller/transactionController"
	transdetail "nextlaundry_apis/controller/transactionDetailController"
	user "nextlaundry_apis/controller/usersController"
	h "nextlaundry_apis/helper"
	mw "nextlaundry_apis/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Content-Length", "Accept-Language", "Accept-Encoding", "Connection", "Access-Control-Allow-Origin", "Authorization", "Access-Control-Allow-Headers", "Headers"},
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
	}))
	public := router.Group("/api")
	{
		public.POST("/auth", a.AuthHandler)
		nextlaundry := public.Group("/nextlaundry")
		nextlaundry.GET("/validate", a.SecondValidate)
		{
			admin := nextlaundry.Group("/admin")
			admin.Use(mw.Admin())
			{
				admin.POST("/logout", a.LogoutHandler)

				//Dashboard GET DATA

				//Users Resource Data
				admin.GET("/users", user.Index)
				admin.GET("/users/:id", user.Show)
				admin.POST("/user", user.Create)
				admin.PUT("/user/:id", user.Update)
				admin.DELETE("/user", user.Delete)
				admin.GET("/users-export-excel", user.ExportToExcel)
				admin.GET("/users-export-pdf", user.ExportToPDF)

				//Outlet Resource Data
				admin.GET("/outlets", outlet.Index)
				admin.GET("/outlets/:id", outlet.Show)
				admin.POST("/outlet", outlet.Create)
				admin.PUT("/outlet/:id", outlet.Update)
				admin.DELETE("/outlet", outlet.Delete)
				admin.GET("/outlets-export-excel", outlet.ExportToExcel)
				admin.GET("/outlets-export-pdf", outlet.ExportToPDF)

				//Member Resource Data
				admin.GET("/members", member.Index)
				admin.GET("/members/:id", member.Show)
				admin.POST("/member", member.Create)
				admin.PUT("/member/:id", member.Update)
				admin.DELETE("/member", member.Delete)
				admin.GET("/members-export-excel", member.ExportToExcel)
				admin.GET("/members-export-pdf", member.ExportToPDF)

				//Package Resource Data
				admin.GET("/products", product.Index)
				admin.GET("/products/:id", product.Show)
				admin.POST("/product", product.Create)
				admin.PUT("/product/:id", product.Update)
				admin.DELETE("/product", product.Delete)
				admin.GET("/products-export-excel", product.ExportToExcel)
				admin.GET("/products-export-pdf", product.ExportToPDF)

				//Transaction Resource Data
				admin.GET("/transactions", trans.Index)
				admin.GET("/transactions/:id", trans.Show)
				admin.POST("/transaction", trans.Create)
				admin.PUT("transaction/:id", trans.Update)
				admin.DELETE("/transaction", trans.Delete)

				//Transaction Detail Resource Data
				admin.GET("/details", transdetail.Index)
				admin.GET("/details/:id", transdetail.Show)
				admin.POST("/detail", transdetail.Create)
				admin.PUT("detail/:id", transdetail.Update)
				admin.DELETE("/detail", transdetail.Delete)

				admin.GET("/transaction-report", h.ExportViewToExcel)

				//Log History READ DATA
				admin.GET("/logs", h.GetLogData)

			}
			cashier := nextlaundry.Group("/cashier")
			cashier.Use(mw.Cashier())
			{
				cashier.POST("/logout", a.LogoutHandler)

				//Dashboard GET DATA

				//Member Resource Data
				cashier.GET("/members", member.Index)
				cashier.GET("/members/:id", member.Show)
				cashier.POST("/member", member.Create)
				cashier.PUT("/member/:id", member.Update)
				cashier.DELETE("/member", member.Delete)
				cashier.GET("/members-export-excel", member.ExportToExcel)
				cashier.GET("/members-export-pdf", member.ExportToPDF)

				//Transaction Resource Data
				cashier.GET("/transactions", trans.Index)
				cashier.GET("/transactions/:id", trans.Show)
				cashier.POST("/transaction", trans.Create)
				cashier.PUT("transaction/:id", trans.Update)
				cashier.DELETE("/transaction", trans.Delete)

				//Transaction Detail Resource Data
				cashier.GET("/details", transdetail.Index)
				cashier.GET("/details/:id", transdetail.Show)
				cashier.POST("/detail", transdetail.Create)
				cashier.PUT("detail/:id", transdetail.Update)
				cashier.DELETE("/detail", transdetail.Delete)

				cashier.GET("/transaction-report", h.ExportViewToExcel)

			}
			owner := nextlaundry.Group("/owner")
			owner.Use(mw.Owner())
			{
				owner.POST("/logout", a.LogoutHandler)

				//Dashboard GET DATA

				//Transaction GET DATA (EXPORT TO EXCEL & PDF)
				owner.GET("/transaction-report", h.ExportViewToExcel)
			}
		}

	}
	return router
}
