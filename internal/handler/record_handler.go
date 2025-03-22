package handler

import (
	"biblia-be/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecordHandler struct {
	db *gorm.DB
}

func (handler *RecordHandler) Initialize(db *gorm.DB) {
	handler.db = db
	db.AutoMigrate(&model.Record{})
}

// GetAllRecords godoc
//
//	@Summary	Get all records
//	@Schemes
//	@Description	return all record objects regarding userId and bookId
//	@Tags			record
//	@Accept			json
//	@Produce		json
//
// @Param userId query int false "the owner's user id of the record"
// @Param bookId query int false "book id of that user's record"
//
//	@Success		200	{array}	model.Record
//	@Router			/records [get]
func (handler *RecordHandler) GetRecords(c *gin.Context) {
	records := []model.Record{}
	userId, hasUserId := c.GetQuery("userId")
	bookId, hasBookId := c.GetQuery("bookId")

	// If neither userId nor bookId are provided, return all records
	if !hasUserId && !hasBookId {
		handler.db.Find(&records)
		c.JSON(http.StatusOK, records)
		return
	}

	// If only userId is provided, return all records for that user
	if hasUserId && !hasBookId {
		if err := handler.db.Where("userId = ?", userId).Find(&records).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User records not found"})
			return
		}
		c.JSON(http.StatusOK, records)
		return
	}

	// If both userId and bookId are provided, return the specific record
	if hasUserId && hasBookId {
		record := model.Record{}
		if err := handler.db.Where("userId = ? AND bookId = ?", userId, bookId).First(&record).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found for the given user and book"})
			return
		}
		c.JSON(http.StatusOK, record)
		return
	}

	// Default: return empty array if no conditions match
	c.JSON(http.StatusOK, []model.Record{})
}

// CreateRecord godoc
//
//	@Summary	Create a new record
//	@Schemes
//	@Description	Create a new record
//	@Tags			record
//	@Accept			json
//	@Produce		json
//
// @Param record body model.CreateRecord true "Record data"
//
//	@Success		201	{object}	model.Record
//	@Router			/records [post]
func (handler *RecordHandler) CreateRecord(c *gin.Context) {
	createRecord := model.CreateRecord{}

	if err := c.ShouldBindJSON(&createRecord); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	record := model.Record{
		BookID:        createRecord.BookID,
		UserID:        createRecord.UserID,
		Status:        createRecord.Status,
		Curr_progress: createRecord.Curr_chapter,
		Curr_chapter:  createRecord.Curr_chapter}

	if err := handler.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, record)
}

// UpdateRecordById godoc
//
//	@Summary	Update record with the id
//	@Schemes
//	@Description	update one record with coresponding id
//	@Tags			record
//	@Accept			json
//	@Produce		json
//
// @Param userId query int false "the owner's user id of the record"
// @Param bookId query int false "book id of that user's record"
//
// @Param record body model.UpdateRecord true "Updated Record data"
//
//	@Success		200	{object}	model.Record
//	@Router			/records/{id} [put]
func (handler *RecordHandler) UpdateRecord(c *gin.Context) {
	record := model.Record{}
	updateRecord := model.UpdateRecord{}
	userId, hasUserId := c.GetQuery("userId")
	bookId, hasBookId := c.GetQuery("bookId")

	if !hasUserId || !hasBookId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both userId and bookId need to be given."})
		return
	}

	if hasUserId && hasBookId {
		record := model.Record{}
		if err := handler.db.Where("userId = ? AND bookId = ?", userId, bookId).First(&record).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found for the given user and book"})
			return
		}
	}

	if err := c.ShouldBindJSON(&updateRecord); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	record.Status = updateRecord.Status
	record.Curr_chapter = updateRecord.Curr_chapter
	record.Curr_progress = updateRecord.Curr_progress
	if err := handler.db.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteRecordById godoc
//
//	@Summary	Delete record with its id
//	@Schemes
//	@Description	Delete one record with coresponding id
//	@Tags			record
//
// @Param userId query int false "the owner's user id of the record"
// @Param bookId query int false "book id of that user's record"
//
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Router			/records/{id} [delete]
func (handler *RecordHandler) DeleteRecord(c *gin.Context) {
	userId, hasUserId := c.GetQuery("userId")
	bookId, hasBookId := c.GetQuery("bookId")
	record := model.Record{}

	if !hasUserId || !hasBookId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both userId and bookId need to be given."})
		return
	}

	if err := handler.db.Where("userId = ? AND bookId = ?", userId, bookId).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found for the given user and book"})
		return
	}

	if err := handler.db.Delete(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
