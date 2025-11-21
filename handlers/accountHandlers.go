package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SignUpCustomer godoc
// @Summary Sign up a new customer
// @Description Create a new customer account
// @Tags Account
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "Customer info"
// @Success 200 {object} map[string]string
// @Router /account/signup [post]
func (h *AccountHandler) SignUpCustomer(c *gin.Context) {
	_ = h.Service.SignUpCustomer(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetCustomer godoc
// @Summary Get a customer by ID
// @Description Retrieve customer info
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} map[string]interface{}
// @Router /account/customer/{id} [get]
func (h *AccountHandler) GetCustomer(c *gin.Context) {
	_ = h.Service.GetCustomer(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// UpdateCustomer godoc
// @Summary Update a customer
// @Description Update customer information
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param input body map[string]interface{} true "Updated customer info"
// @Success 200 {object} map[string]string
// @Router /account/customer/{id} [put]
func (h *AccountHandler) UpdateCustomer(c *gin.Context) {
	_ = h.Service.UpdateCustomer(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// DeleteCustomer godoc
// @Summary Delete a customer
// @Description Remove a customer by ID
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} map[string]string
// @Router /account/customer/{id} [delete]
func (h *AccountHandler) DeleteCustomer(c *gin.Context) {
	_ = h.Service.DeleteCustomer(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// ------------------- Employee Endpoints -------------------

// SignUpEmployee godoc
// @Summary Sign up a new employee
// @Description Create a new employee account
// @Tags Account
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "Employee info"
// @Success 200 {object} map[string]string
// @Router /account/employee/signup [post]
func (h *AccountHandler) SignUpEmployee(c *gin.Context) {
	_ = h.Service.SignUpEmployee(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetEmployee godoc
// @Summary Get an employee by ID
// @Description Retrieve employee info
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} map[string]interface{}
// @Router /account/employee/{id} [get]
func (h *AccountHandler) GetEmployee(c *gin.Context) {
	_ = h.Service.GetEmployee(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// UpdateEmployee godoc
// @Summary Update an employee
// @Description Update employee info
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param input body map[string]interface{} true "Updated employee info"
// @Success 200 {object} map[string]string
// @Router /account/employee/{id} [put]
func (h *AccountHandler) UpdateEmployee(c *gin.Context) {
	_ = h.Service.UpdateEmployee(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// UpdateEmployeePerformanceScore godoc
// @Summary Update employee performance score
// @Description Change performance score for employee
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param score query int true "New performance score"
// @Success 200 {object} map[string]string
// @Router /account/employee/{id}/performance [patch]
func (h *AccountHandler) UpdateEmployeePerformanceScore(c *gin.Context) {
	_ = h.Service.UpdateEmployeePerformanceScore(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// UpdateEmployeeStatus godoc
// @Summary Update employee status
// @Description Change status (active/inactive) for employee
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param status query string true "New status"
// @Success 200 {object} map[string]string
// @Router /account/employee/{id}/status [patch]
func (h *AccountHandler) UpdateEmployeeStatus(c *gin.Context) {
	_ = h.Service.UpdateEmployeeStatus(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// DeleteEmployee godoc
// @Summary Delete an employee
// @Description Remove an employee by ID
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} map[string]string
// @Router /account/employee/{id} [delete]
func (h *AccountHandler) DeleteEmployee(c *gin.Context) {
	_ = h.Service.DeleteEmployee(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}