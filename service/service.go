package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// "github.com/eapache/go-resiliency/retrier"
	"github.com/gin-gonic/gin"
)

// ==============================EMPLOYEE STRUCT AND NEWEMPLOYEE METHOD========================================
type Employee struct {
	User_id int    `json:"User_id"`
	Fname   string `json:"Fname"`
	Lname   string `json:"Lname"`
	Age     int    `json:"Age"`
}

func NewEmployee(id int, fname, lname string, age int) Employee {
	emp := Employee{
		User_id: id,
		Fname:   fname,
		Lname:   lname,
		Age:     age,
	}

	return emp
}

// ============================================================================================================

type CustomError struct {
	msg string
	err error
}

type Database struct {
	db *sql.DB
}

type Db interface {
	GetAllEmployees() ([]Employee, error)
	GetOneWithId(id string) (Employee, CustomError)
	CreateUsers(emp []Employee, db *sql.DB) (bool, CustomError)
}

func NewDatabase(db *sql.DB) Db {
	return &Database{db: db}
}

// ====================IMPLEMENTING THE Db INTERFACE METHODS==============================

func (d *Database) GetAllEmployees() ([]Employee, error) {
	// Example SQL query to fetch all employees
	rows, err := d.db.Query("SELECT user_id, fname, lname, age FROM user")
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	var employees []Employee

	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.User_id, &emp.Fname, &emp.Lname, &emp.Age); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return employees, nil
}

func (d *Database) GetOneWithId(user string) (Employee, CustomError) {
	var c *gin.Context
	var employee Employee
	var ce CustomError
	// convert string User_id to int User_id
	user_id, err := strconv.Atoi(user)
	if err != nil {
		ce.msg = "Incorrect User ID Format"
		ce.err = fmt.Errorf("incorrect User ID Format. %v", err)
		return employee, ce
	}
	// return the employee after querying to the db
	row, err := d.db.Query("SELECT * FROM user WHERE user_id = (?)", user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "unable to query from the database")
		ce.msg = "unable to query from the database"
		ce.err = err
		return employee, ce
	}
	fmt.Print("\n\nReached here\n\n")

	defer row.Close()

	// using Scan to assign column data to struct fields in row variable.
	if err := row.Scan(&employee.User_id, &employee.Fname, &employee.Lname, &employee.Age); err != nil {
		if err == sql.ErrNoRows {
			// when the requested id is not present in the database //
			c.JSON(http.StatusNotFound, gin.H{"employee": user_id, "value": "no such employee exists with given user_id"})
			// fmt.Errorf("GetUserWithID %d: no such employee exists", user_id)
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, "Something wrong with SQL server executing queries")
		}
		// fmt.Errorf("GetUserWithID %d: %w", user_id, err)
	}
	ce.msg = "All clear"
	return employee, ce
}

func (d *Database) CreateUsers(newEmployee []Employee, db *sql.DB) (bool, CustomError) {
	var ce CustomError
	for i := 0; i < len(newEmployee); i++ {
		insertResult, err := d.db.ExecContext(context.Background(), "INSERT INTO `user` VALUES(?, ?, ?, ?)", newEmployee[i].User_id, newEmployee[i].Fname, newEmployee[i].Lname, newEmployee[i].Age)
		if err != nil {
			log.Fatalf("Unable to insert values(%v, %v, %v, %v)", newEmployee[i].User_id, newEmployee[i].Fname, newEmployee[i].Lname, newEmployee[i].Age)
			ce.msg = "unsuccessful insert operation"
			ce.err = err
			return false, ce
		}
		id, err := insertResult.LastInsertId()
		if err != nil {
			// c.JSON(http.)  // cannot find the right http status code
			log.Fatalf("impossible to retrieve last inserted id: %s", err)
			ce.msg = "impossible to retrieve last inserted id\n"
			ce.err = err
			return false, ce
		}
		log.Printf("inserted id: %d", id)
	}
	return true, ce
}

// ==========================================================================================

// ===================== SERVICE LAYER HANDLERS =============================================
func HandleGetAllEmployees(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		empDB := NewDatabase(db)
		employees, err := empDB.GetAllEmployees()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees"})
			log.Fatalf("failed to retrieve message : %v", err)
			return
		}
		log.Println("successfully got all the employees")
		c.JSON(http.StatusOK, employees)
	}
}

func GetUserWithID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// send the converted json data

		user := c.Params.ByName("User_id")

		empDB := NewDatabase(db)

		employee, queryErr := empDB.GetOneWithId(user)
		if queryErr.err != nil {
			errMsg := fmt.Errorf("an error occured while getting the requested Employee ID : %v\n", queryErr)
			log.Fatalf("%v\n", errMsg)
			c.JSON(http.StatusBadRequest, gin.H{"An error occured from the client side. Check the requested ID. ": errMsg})
		}
		log.Println("successfully got the employee")
		c.JSON(http.StatusOK, gin.H{"Employee": user, "value": employee})
	}
}

func CreateEmployee(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newEmployee []Employee
		empDB := NewDatabase(db)
		// Call BindJSON to bind the received JSON to the employee struct
		if err := c.BindJSON(&newEmployee); err != nil {
			return
		}

		empDB.CreateUsers(newEmployee, db)
		// employeeData = append(employeeData, newEmployee...)
		c.IndentedJSON(http.StatusCreated, newEmployee)
	}
}

// ==========================================================================================
