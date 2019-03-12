package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

// Employee data
type Employee struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Salary string `json:"salary"`
	Age    string `json:"age"`
	Email  string `json:"email"`
}

// Employees object
type Employees struct {
	Employees []Employee `json:"employees"`
}

func init() {

	connStr := "host=localhost port=5432 user=postgres password=FR4uT dbname=employee_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("DB Connected...")
	}
}

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server connected!")
	})

	e.POST("/employee", addEmployee)
	e.PUT("/employee/:id", updateEmployee)
	e.DELETE("/employee/:id", deleteEmployee)
	e.GET("/employee/:id", getEmployee)
	e.GET("/employee", getAllEmployee)

	e.Logger.Fatal(e.Start(":1323"))

}

func addEmployee(c echo.Context) error {
	u := new(Employee)
	if err := c.Bind(u); err != nil {
		return err
	}
	sqlStatement := `INSERT INTO employees(name, salary, age, email) VALUES($1, $2, $3, $4)`
	res, err := db.Query(sqlStatement, u.Name, u.Salary, u.Age, u.Email)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
		return c.JSON(http.StatusCreated, u)
	}
	return c.String(http.StatusOK, "ok")
}

func getAllEmployee(c echo.Context) error {
	sqlStatement := `SELECT id, name, salary, age, email FROM employees ORDER BY id`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	result := Employees{}

	for rows.Next() {
		employee := Employee{}
		err2 := rows.Scan(&employee.ID, &employee.Name, &employee.Salary, &employee.Age, &employee.Email)
		if err2 != nil {
			return err2
		}
		result.Employees = append(result.Employees, employee)
	}
	return c.JSON(http.StatusCreated, result)
}

func updateEmployee(c echo.Context) error {
	u := new(Employee)
	if err := c.Bind(u); err != nil {
		return err
	}
	sqlStatement := `UPDATE employees SET name=$1, salary=$2, age=$3, email=$4 WHERE id=$5`
	res, err := db.Query(sqlStatement, u.Name, u.Salary, u.Age, u.Email, u.ID)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
		return c.JSON(http.StatusCreated, u)
	}
	return c.String(http.StatusOK, u.ID)
}

func getEmployee(c echo.Context) error {
	id := c.Param("id")
	sqlStatement := `SELECT FROM employees WHERE id=$1`
	res, err := db.Query(sqlStatement, id)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
		return c.JSON(http.StatusOK, "Selected")
	}
	return c.String(http.StatusOK, id+"Selected")
}

func deleteEmployee(c echo.Context) error {
	id := c.Param("id")
	sqlStatement := `DELETE FROM employees WHERE id=$1`
	res, err := db.Query(sqlStatement, id)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
		return c.JSON(http.StatusOK, "Deleted")
	}
	return c.String(http.StatusOK, id+"Deleted")
}
