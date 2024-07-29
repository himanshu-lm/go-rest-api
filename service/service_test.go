package service

import (
	"database/sql"
	"fmt"

	// "fmt"
	"net/http"
	"testing"

	// "unittestexample/service/mocks"

	// "unittestexample/service

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var EmployeeTable *sql.DB

func TestGetAllEmployees(t *testing.T) {
	// Create a new instance of the mock
	// mockDb := service.NewMockDb
	ctrl := gomock.NewController(t)
	gin.SetMode(gin.TestMode)
	// defer ctrl.Finish() // no need to call as *testing .T is already included in the NewController Method.
	mockDb := NewMockDb(ctrl)
	// Create a new Gin router
	// r := gin.Default()

	// mycode
	var myerr error
	myerr = nil
	mockDb.EXPECT().GetAllEmployees().Return([]Employee{
		{
			User_id: 1,
			Fname:   "bob",
			Lname:   "marley",
			Age:     28,
		},
		{
			User_id: 1,
			Fname:   "bob",
			Lname:   "marley",
			Age:     28,
		},
	},
		myerr).AnyTimes()

	t.Run("Normal Test case for getting all", func(t *testing.T) {
		res, _ := mockDb.GetAllEmployees()
		assert.NotNil(t, res)
	})
	// mycode
	// Inject the mock into the handler
	// r.GET("/allemployees", func(c *gin.Context) {
	// 	employees, err := mockDb.GetAllEmployees()
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, employees)
	// })

	// // Create a mock response from the database
	// mockEmployees := []Employee{
	// 	{User_id: 1, Fname: "John", Lname: "Doe", Age: 30},
	// 	{User_id: 2, Fname: "Jane", Lname: "Smith", Age: 28},
	// }
	// mockDb.On("GetAllEmployees", mock.Anything).Return(mockEmployees, nil)
	// // Perform a GET request to the endpoint
	// req, err := http.NewRequest("GET", "/allemployees", nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// resp := httptest.NewRecorder()
	// r.ServeHTTP(resp, req)

	// // Assert the HTTP status code
	// assert.Equal(t, http.StatusOK, resp.Code)

	// // Assert the response body (assuming JSON response)
	// expectedBody := `[{"user_id":1,"fname":"John","lname":"Doe","age":30},{"user_id":2,"fname":"Jane","lname":"Smith","age":28}]`
	// assert.JSONEq(t, expectedBody, resp.Body.String())

	// // Assert the mock expectations
	// mockDb.AssertExpectations(t)
}

func TestGetOneWithId(t *testing.T) {
	// Create a new instance of the mock
	// mockDb := service.NewMockDb
	ctrl := gomock.NewController(t)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	// defer ctrl.Finish() // no need to call as *testing.T is already included in the NewController Method.
	mockDb := NewMockDb(ctrl)
	// Create a new Gin router
	// r := gin.Default()

	// mycode
	// mockDb is a mock implementation of the Db interface. It has all the methods that we had in the original Db interfaces.
	// mockDb is it's mock.

	// EXPECT() this expects one of the methods from the mockDb interface. Here it is GetOneWithId(). The Return part just overrides
	// methods definition and returns whatever we want to return

	// these below are the actual testcases. Whenever we call them, it checks with the EXPECTS method so that the test cases
	// are expeecting a particular method. The first string is just the name for the testcase.
	mockDb.EXPECT().GetOneWithId("6").Return(Employee{
		User_id: 6,
		Fname:   "bob",
		Lname:   "marley",
		Age:     28,
	}, CustomError{
		msg: "expected error",
		err: nil,
	}).AnyTimes()
	t.Run("Check if the response is not nil", func(t *testing.T) {
		res, _ := mockDb.GetOneWithId("6")
		assert.NotNil(t, res)
	})
	t.Run("Check if the length of response != 1", func(t *testing.T) {
		res, _ := mockDb.GetOneWithId("6")
		// assert.NotNil(t, res)
		// assert.Len(t, res, 1)
		assert.Equal(t, res, Employee{
			User_id: 6,
			Fname:   "bob",
			Lname:   "marley",
			Age:     28,
		})
	})

	t.Run("Input is either + or - only", func(t *testing.T) {
		mockDb.EXPECT().GetOneWithId("+").Return(Employee{
			User_id: -1,
			Fname:   "no",
			Lname:   "name",
			Age:     -1,
		}, CustomError{
			msg: "user id can not contain only + or -",
			err: fmt.Errorf("Incorrect User ID Format"),
		})

		_, err := mockDb.GetOneWithId("+")
		expectedErrorMsg := "Incorrect User ID Format"
		assert.EqualErrorf(t, err.err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err.msg)
	})
	t.Run("Input is a non-numeric string", func(t *testing.T) {
		mockDb.EXPECT().GetOneWithId("87ysf").Return(Employee{
			User_id: -1,
			Fname:   "no",
			Lname:   "name",
			Age:     -1,
		}, CustomError{
			msg: "Cannot use non-numeric as user ida",
			err: fmt.Errorf("Incorrect User ID Format"),
		})

		_, err := mockDb.GetOneWithId("87ysf")
		expectedErrorMsg := "Incorrect User ID Format"
		assert.EqualErrorf(t, err.err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err.msg)
	})
	// t.Run("")
	// Inject the mock into the handler
	r.GET("/allemployees", func(c *gin.Context) {
		employees, err := mockDb.GetAllEmployees()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees"})
			return
		}
		c.JSON(http.StatusOK, employees)
	})

}

func TestCreateUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	gin.SetMode(gin.TestMode)
	// r := gin.Default()
	mockDb := NewMockDb(ctrl)
	emp := []Employee{
		{
			User_id: 21,
			Fname:   "fname1",
			Lname:   "lname`",
			Age:     21,
		},
		{
			User_id: 22,
			Fname:   "fname2",
			Lname:   "lname2",
			Age:     32,
		},
	}
	mockDb.EXPECT().CreateUsers(emp, EmployeeTable).Return(true, CustomError{
		msg: "expected error",
		err: nil,
	}).AnyTimes()

	t.Run("testcase 1 : check if all are inserted\n", func(t *testing.T) {
		success, _ := mockDb.CreateUsers(emp, EmployeeTable)
		assert.True(t, success)
	})

	// t.Run("testcase 2 : check for errors", func(t *testing.T) {

	// })
}
