package handler

import (
	"biblia-be/internal/model"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RecordHandler manages book record operations
type RecordHandler struct {
	db *gorm.DB
}

// Initialize sets up the handler with a database connection and performs migrations
func (handler *RecordHandler) Initialize(db *gorm.DB) {
	handler.db = db
	db.AutoMigrate(&model.Record{})
}

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// GetRecords godoc
//
//	@Summary	Get user reading records
//	@Schemes
//	@Description	Returns reading records filtered by userId and/or ISBN
//	@Tags			records
//	@Accept			json
//	@Produce		json
//
// @Param userId query int false "User ID to filter records (optional)"
// @Param isbn query int false "ISBN to filter records (optional)"
//
//	@Success		200	{object} Response{data=[]model.Record} "Successfully retrieved records"
//	@Failure		400	{object} Response "Invalid request parameters"
//	@Failure		500	{object} Response "Internal server error"
//	@Router			/records [get]
func (handler *RecordHandler) GetRecords(c *gin.Context) {
	var records []model.Record
	userIdParam, hasUserId := c.GetQuery("userId")
	isbnParam, hasIsbn := c.GetQuery("isbn")

	query := handler.db

	// Apply filters if they exist
	if hasUserId {
		userId, err := strconv.ParseUint(userIdParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Error:   "Invalid userId format",
			})
			return
		}
		query = query.Where("user_id = ?", userId)
	}

	if hasIsbn {
		isbn, err := strconv.ParseUint(isbnParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Error:   "Invalid ISBN format",
			})
			return
		}
		query = query.Where("isbn = ?", isbn)
	}

	// Execute query
	result := query.Find(&records)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve records",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    records,
		Message: "Records retrieved successfully",
	})
}

// GetRecordByUserAndISBN godoc
//
//	@Summary	Get a specific reading record
//	@Schemes
//	@Description	Returns a single reading record for a specific user and ISBN
//	@Tags			records
//	@Accept			json
//	@Produce		json
//
// @Param userId query int true "User ID of the record owner"
// @Param isbn query int true "ISBN of the book"
//
//	@Success		200	{object} Response{data=model.Record} "Successfully retrieved record"
//	@Failure		400	{object} Response "Invalid request parameters"
//	@Failure		404	{object} Response "Record not found"
//	@Failure		500	{object} Response "Internal server error"
//	@Router			/records/detail [get]
func (handler *RecordHandler) GetRecordByUserAndISBN(c *gin.Context) {
	userIdParam, hasUserId := c.GetQuery("userId")
	isbnParam, hasIsbn := c.GetQuery("isbn")

	// Validate required parameters
	if !hasUserId || !hasIsbn {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Both userId and isbn parameters are required",
		})
		return
	}

	// Parse and validate userId
	userId, err := strconv.ParseUint(userIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid userId format",
		})
		return
	}

	// Parse and validate isbn
	isbn, err := strconv.ParseUint(isbnParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ISBN format",
		})
		return
	}

	// Query the record
	var record model.Record
	result := handler.db.Where("user_id = ? AND isbn = ?", userId, isbn).First(&record)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, Response{
				Success: false,
				Error:   "Record not found for the specified user and ISBN",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve record",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    record,
		Message: "Record retrieved successfully",
	})
}

// CreateRecord godoc
//
//	@Summary	Create a new reading record
//	@Schemes
//	@Description	Creates a new reading record for a user
//	@Tags			records
//	@Accept			json
//	@Produce		json
//
// @Param record body model.CreateRecord true "Reading record data"
//
//	@Success		201	{object} Response{data=model.Record} "Record created successfully"
//	@Failure		400	{object} Response "Invalid request body"
//	@Failure		409	{object} Response "Record already exists"
//	@Failure		500	{object} Response "Internal server error"
//	@Router			/records [post]
func (handler *RecordHandler) CreateRecord(c *gin.Context) {
	var createRecord model.CreateRecord

	// Parse request body
	if err := c.ShouldBindJSON(&createRecord); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if createRecord.UserID == 0 || createRecord.ISBN == "" || createRecord.Title == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "UserID, ISBN, and Title are required fields",
		})
		return
	}

	// Check if record already exists
	var existingRecord model.Record
	result := handler.db.Where("user_id = ? AND isbn = ?", createRecord.UserID, createRecord.ISBN).First(&existingRecord)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, Response{
			Success: false,
			Error:   "A record for this user and ISBN already exists",
		})
		return
	}

	// Create new record
	record := model.Record{
		ISBN:        createRecord.ISBN,
		UserID:      createRecord.UserID,
		Title:       createRecord.Title,
		Author:      createRecord.Author,
		Cover:       createRecord.Cover,
		Genre:       createRecord.Genre,
		Status:      createRecord.Status,
		CurrentPage: createRecord.CurrentPage,
		TotalPages:  createRecord.TotalPages,
		DateAdded:   time.Now(),
	}

	if err := handler.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to create record: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    record,
		Message: "Record created successfully",
	})
}

// UpdateRecord godoc
//
//	@Summary	Update a reading record
//	@Schemes
//	@Description	Updates reading progress for a specific user's book
//	@Tags			records
//	@Accept			json
//	@Produce		json
//
// @Param userId query int true "User ID of the record owner"
// @Param isbn query int true "ISBN of the book"
// @Param record body model.UpdateRecord true "Updated reading status"
//
//	@Success		200	{object} Response{data=model.Record} "Record updated successfully"
//	@Failure		400	{object} Response "Invalid request parameters"
//	@Failure		404	{object} Response "Record not found"
//	@Failure		500	{object} Response "Internal server error"
//	@Router			/records [put]
func (handler *RecordHandler) UpdateRecord(c *gin.Context) {
	var updateRecord model.UpdateRecord

	userIdParam, hasUserId := c.GetQuery("userId")
	isbnParam, hasIsbn := c.GetQuery("isbn")

	// Validate required parameters
	if !hasUserId || !hasIsbn {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Both userId and isbn parameters are required",
		})
		return
	}

	// Parse and validate userId
	userId, err := strconv.ParseUint(userIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid userId format",
		})
		return
	}

	// Parse and validate isbn
	isbn, err := strconv.ParseUint(isbnParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ISBN format",
		})
		return
	}

	// Parse request body
	if err := c.ShouldBindJSON(&updateRecord); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Find the record to update
	var record model.Record
	result := handler.db.Where("user_id = ? AND isbn = ?", userId, isbn).First(&record)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, Response{
				Success: false,
				Error:   "Record not found for the specified user and ISBN",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve record for update",
		})
		return
	}

	// Validate update data
	if updateRecord.CurrentPage > record.TotalPages {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Current page cannot exceed total pages",
		})
		return
	}

	// Update record fields
	record.Status = updateRecord.Status
	record.CurrentPage = updateRecord.CurrentPage

	// Save updated record
	if err := handler.db.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to update record: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    record,
		Message: "Record updated successfully",
	})
}

// DeleteRecord godoc
//
//	@Summary	Delete a reading record
//	@Schemes
//	@Description	Removes a reading record for a specific user and book
//	@Tags			records
//	@Accept			json
//	@Produce		json
//
// @Param userId query int true "User ID of the record owner"
// @Param isbn query int true "ISBN of the book"
//
//	@Success		200	{object} Response "Record deleted successfully"
//	@Failure		400	{object} Response "Invalid request parameters"
//	@Failure		404	{object} Response "Record not found"
//	@Failure		500	{object} Response "Internal server error"
//	@Router			/records [delete]
func (handler *RecordHandler) DeleteRecord(c *gin.Context) {
	userIdParam, hasUserId := c.GetQuery("userId")
	isbnParam, hasIsbn := c.GetQuery("isbn")

	// Validate required parameters
	if !hasUserId || !hasIsbn {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Both userId and isbn parameters are required",
		})
		return
	}

	// Parse and validate userId
	userId, err := strconv.ParseUint(userIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid userId format",
		})
		return
	}

	// Parse and validate isbn
	isbn, err := strconv.ParseUint(isbnParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ISBN format",
		})
		return
	}

	// Find the record to delete
	var record model.Record
	result := handler.db.Where("user_id = ? AND isbn = ?", userId, isbn).First(&record)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, Response{
				Success: false,
				Error:   "Record not found for the specified user and ISBN",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve record for deletion",
		})
		return
	}

	// Delete the record
	if err := handler.db.Delete(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to delete record: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Record deleted successfully",
	})
}
