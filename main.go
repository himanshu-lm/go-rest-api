package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"unittestexample/service"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var EmployeeTable *sql.DB

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	// r.GET("/ping", service.HandlePing)

	// GET : get all users
	r.GET("/allemployees", service.HandleGetAllEmployees(EmployeeTable))

	// GET : get an employee with id
	r.GET("/employee/:User_id", service.GetUserWithID(EmployeeTable))

	// create an employee
	r.POST("/create", service.CreateEmployee(EmployeeTable))

	// DELETE :delete a requested user with User_id
	// r.DELETE("/delete/:User_id", service.DeleteEmployee)

	// PUT : Update existing user with User_id
	// r.PUT("update/:User_id", service.UpdateUser)

	// Authorized group (uses gin.BasicAuth() mUser_iddleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	// 	  "foo":  "bar",
	// 	  "manu": "123",
	// }))
	// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	// 	"foo":  "bar", // user:foo password:bar
	// 	"manu": "123", // user:manu password:123
	// }))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	// authorized.GET("bidenmissilebutton/", auth.HandleAuthFromBiden)
	// authorized.POST("admin", auth.HandleAdminPost)
	// middlewares
	// r.Use(sampleMiddleware())
	return r
}

func main() {
	/*SQL CONNNECTION*/
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "my_db_test",
	}
	// Get a database handle.
	var err error
	EmployeeTable, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := EmployeeTable.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to SQL!")
	// SQL CONNECTION DONE

	r := setupRouter()
	// Listen and Se`rver in 0.0.0.0:8080
	r.Run(":8080")
}
