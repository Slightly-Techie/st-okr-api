package controllers

import (
	"net/http"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/services"
	"github.com/gin-gonic/gin"
)

type KeyResultController struct {
	keyResultService services.KeyResultService
}

func NewKeyResultController(keyResultService services.KeyResultService) *KeyResultController {
	return &KeyResultController{
		keyResultService: keyResultService,
	}
}

func (kctrl *KeyResultController) CreateKeyResult(c *gin.Context) {
	var req dto.CreateKeyResultRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kr, err := kctrl.keyResultService.CreateKeyResult(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Key result created successfully",
		"key_result": kr,
	})
}

func (kctrl *KeyResultController) GetKeyResult(c *gin.Context) {
	id := c.Param("id")

	kr, err := kctrl.keyResultService.GetData("id", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"key_result": kr,
	})
}

func (kctrl *KeyResultController) ListObjKeyResults(c *gin.Context) {
	odjID := c.Param("id")

	kr, err := kctrl.keyResultService.ListData("objective_id", odjID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"key_result": kr,
	})

}

func (kctrl *KeyResultController) ListAssigneeKeyResults(c *gin.Context) {
	assgnID := c.Param("id")

	kr, err := kctrl.keyResultService.ListData("assignee_id", assgnID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Assignee Key Results": kr,
	})
}

func (kctrl *KeyResultController) UpdateKeyResult(c *gin.Context) {
	var req dto.UpdateKeyResultRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kr, err := kctrl.keyResultService.UpdateKeyResult(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":            "Key result updated successfully",
		"updated key result": kr,
	})
}

func (kctrl *KeyResultController) DeleteKeyResult(c *gin.Context) {
	id := c.Param("id")

	err := kctrl.keyResultService.DeleteKeyResult(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Key result deleted successfully",
	})
}
