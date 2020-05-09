package main

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hakula139/REALMS/internal/app/config"
	ctrl "github.com/hakula139/REALMS/internal/app/controllers"
	"github.com/hakula139/REALMS/internal/app/models"
)

const session = "mysession"

func main() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	dbcfg, err := config.LoadDbConfig("./configs/db_config.json")
	if err != nil {
		panic(err.Error())
	}
	db, err := models.DbSetup(dbcfg)
	if err != nil {
		panic(err.Error())
	}

	libcfg, err := config.LoadLibraryConfig("./configs/library_config.json")
	if err != nil {
		panic(err.Error())
	}

	// Provides variables to controllers
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("libcfg", libcfg)
		c.Next()
	})

	// Uses cookies to store sessions
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions(session, store))

	// Public
	r.POST("/login", ctrl.Login)
	r.GET("/logout", ctrl.Logout)

	r.GET("/books", ctrl.ShowBooks)
	r.GET("/books/:id", ctrl.ShowBook)
	r.POST("/books/find", ctrl.FindBooks)

	// User privilege required
	user := r.Group("/user")
	user.Use(ctrl.AuthRequired)
	{
		user.GET("/me", ctrl.Me)
		user.GET("/status", ctrl.Status)

		user.GET("/books", ctrl.ShowBookList)
		user.GET("/books/:id", ctrl.ShowBorrowed)
		user.GET("/overdue", ctrl.ShowOverdueList)
		user.GET("/history", ctrl.ShowHistory)
		user.POST("/books/:id", ctrl.BorrowBook)
		user.PATCH("/books/:id", ctrl.ExtendDeadline)
		user.DELETE("/books/:id", ctrl.ReturnBook)
	}

	// Admin privilege required
	admin := r.Group("/admin")
	admin.Use(ctrl.AdminRequired)
	{
		admin.POST("/books", ctrl.AddBook)
		admin.PATCH("/books/:id", ctrl.UpdateBook)
		admin.DELETE("/books/:id", ctrl.RemoveBook)

		admin.GET("/users", ctrl.ShowUsers)
		admin.GET("/users/:id", ctrl.ShowUser)
		admin.POST("/users", ctrl.AddUser)
		admin.PATCH("/users/:id", ctrl.UpdateUser)
		admin.DELETE("/users/:id", ctrl.RemoveUser)
	}

	if err := r.Run(":7274"); err != nil {
		fmt.Println("[error] realmsd: failed to start: " + err.Error())
	}
}
