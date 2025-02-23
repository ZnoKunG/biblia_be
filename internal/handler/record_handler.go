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
//	@Description	return all record objects in the database
//	@Tags			record
//	@Accept			json
//	@Produce		json
//
// @Param userId query int false "user-owned record search by userId"
//
//	@Success		200	{array}	model.Record
//	@Router			/record [get]
func (handler *RecordHandler) GetRecords(c *gin.Context) {
	records := []model.Record{}
	userId, userIdExists := c.GetQuery("userId")
	user := model.User{}

	// Query all records with no condition
	if !userIdExists {
		handler.db.Find(&records)
		c.JSON(http.StatusOK, records)
		return
	}

	if err := handler.db.Find(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	if err := handler.db.Where("userId <> ?", user.ID).Find(&records).Error; err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, records)
}

// GetRecordById godoc
//
//	@Summary	Get record with its id
//	@Schemes
//	@Description	return one record with coresponding id
//	@Tags			record
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Record ID"
//	@Success		200	{object}	model.Record
//	@Router			/record/{id} [get]
func (handler *RecordHandler) GetRecord(c *gin.Context) {
	id := c.Param("id")
	record := model.Record{}

	if err := handler.db.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, record)
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
//	@Router			/record [post]
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
//	@Param			id	path		int	true	"Record ID"
//
// @Param record body model.UpdateRecord true "Updated Record data"
//
//	@Success		200	{object}	model.Record
//	@Router			/record/{id} [put]
func (handler *RecordHandler) UpdateRecord(c *gin.Context) {
	id := c.Param("id")
	record := model.Record{}
	updateRecord := model.UpdateRecord{}

	if err := handler.db.First(&record, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
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
//	@Param			id	path	int	true	"Record ID"
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Router			/record/{id} [delete]
func (handler *RecordHandler) DeleteRecord(c *gin.Context) {
	id := c.Param("id")
	record := model.Record{}

	if err := handler.db.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	if err := handler.db.Delete(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
