package handlers

import (
	"github.com/labstack/echo"
	. "echo-sample/models"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func GetEmployees(context echo.Context) error {
	employees, err := FindEmployees(context.Get("db").(*gorm.DB))

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Employees do not exist.")
	}
	return context.JSON(http.StatusOK, employees)
}

func GetEmployee(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, "Employee Id must be int")
	}

	employee := &Employee{Model: Model{ID: id}}
	if err := employee.Find(context.Get("db").(*gorm.DB)); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Employee does not exist.")
	}

	return context.JSON(http.StatusOK, employee)
}

func CreateEmployee(context echo.Context) error {
	employee := &Employee{}
	if err := context.Bind(employee); err != nil {
		return err
	}

	if err := employee.Create(context.Get("db").(*gorm.DB)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Employee can not Create.")
	}

	return context.JSON(http.StatusOK, employee)
}

func UpdateEmployee(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, "Employee Id must be int")
	}

	employee := &Employee{Model: Model{ID: id}}
	if err := context.Bind(employee); err != nil {
		return err
	}

	if err := employee.Update(context.Get("db").(*gorm.DB)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Employee can not update.")
	}

	return context.JSON(http.StatusOK, employee)
}

func DeleteEmployee(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, "Employee Id must be int")
	}

	employee := &Employee{Model: Model{ID: id}}

	if err := employee.Delete(context.Get("db").(*gorm.DB)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Employee can not delete.")
	}

	return context.JSON(http.StatusOK, employee)
}
