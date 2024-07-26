package service

import (
	"database/sql"
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
	r := gin.Default()

	// mycode
	mockDb.EXPECT().GetAllEmployees().DoAndReturn(func(id int) Employee {
		return Employee{
			User_id: 1,
			Fname:   "bob",
			Lname:   "marley",
			Age:     28,
		}
	}).AnyTimes()

	t.Run("Normal Test case for getting all", func(t *testing.T) {
		res := HandleGetAllEmployees(EmployeeTable)
		assert.NotNil(t, res)
	})
	// mycode
	// Inject the mock into the handler
	r.GET("/allemployees", func(c *gin.Context) {
		employees, err := mockDb.GetAllEmployees()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees"})
			return
		}
		c.JSON(http.StatusOK, employees)
	})

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
