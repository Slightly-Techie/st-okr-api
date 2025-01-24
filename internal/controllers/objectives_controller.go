package controllers

import (
	"fmt"
	"net/http"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/services"
	"github.com/gin-gonic/gin"
)

type ObjectivesController struct {
	objectiveService services.ObjectivesService
}

func NewObjectiveController(objectivesService services.ObjectivesService) *ObjectivesController {
	return &ObjectivesController{
		objectiveService: objectivesService,
	}
}

func (ctrl *ObjectivesController) CreateObjective(c *gin.Context) {
	type CCBody struct {
		Title    string `json:"title"`
		Deadline string `json:"deadline"`
	}
	var body CCBody

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqBody := dto.CreateObjectivesRequest{
		Title:     body.Title,
		CreatorID: userID.(string),
		Deadline:  body.Deadline,
	}

	fmt.Println(reqBody)

	data, err := ctrl.objectiveService.CreateObjectives(reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company created successfully",
		"data":    data,
	})
}

func (ctrl *ObjectivesController) GetObjectives(c *gin.Context) {

	id := c.Param("id")

	data, err := ctrl.objectiveService.GetObjectives("id", id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Objective found",
		"data":    data,
	})
}

func (ctrl *ObjectivesController) UpdateObjective(c *gin.Context) {
	var body dto.CreateObjectivesRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := ctrl.objectiveService.UpdateObjectives(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "OBjective created successfully",
		"data":    data,
	})
}

func (ctrl *ObjectivesController) DeleteObjective(c *gin.Context) {
	id := c.Param("id")
	err := ctrl.objectiveService.DeleteObjectives(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OBjective deleted successfully"})
}
