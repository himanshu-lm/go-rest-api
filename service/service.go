package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		ce.err = fmt.Errorf("Incorrect User ID Format. %v\n", err)
		return employee, ce
	}

	// return the employee after querying to the db
	row, err := d.db.Query("SELECT * FROM user WHERE user_id = (?)", user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "INCORRECT USER ID FORMAT")
		ce.msg = "Unable to query from the database"
		ce.err = err
		return employee, ce
	}

	defer row.Close()

	// using Scan to assign column data to struct fields in row variable.
	if err := row.Scan(&employee.User_id, &employee.Fname, &employee.Lname, &employee.Age); err != nil {
		if err == sql.ErrNoRows {
			// when the requested id is not present in the database //
			c.JSON(http.StatusNotFound, gin.H{"employee": user_id, "value": "no such employee exists with given user_id"})
			fmt.Errorf("GetUserWithID %d: no such employee exists", user_id)
		}
		c.JSON(http.StatusInternalServerError, "Something wrong with SQL server executing queries")
		fmt.Errorf("GetUserWithID %d: %w", user_id, err)
	}
	ce.msg = "All clear"
	return employee, ce
}

// ==========================================================================================

// ===================== SERVICE LAYER HANDLERS =============================================
func HandleGetAllEmployees(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		empDB := NewDatabase(db)
		employees, err := empDB.GetAllEmployees()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees"})
			return
		}
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
			errMsg := fmt.Errorf("an error occured while getting the requested Employee ID : %v", queryErr)
			c.JSON(http.StatusBadRequest, gin.H{"An error occured from the client side. Check the requested ID. ": errMsg})
		}
		c.JSON(http.StatusOK, gin.H{"Employee": user, "value": employee})
	}
}

// ==========================================================================================
