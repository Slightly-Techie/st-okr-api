package controllers

import (
	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/response"
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
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	kr, err := kctrl.keyResultService.CreateKeyResult(req)
	if err != nil {
		response.BadRequest(c, "Failed to create key result", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.Created(c, kr, "Key result created successfully")
}

func (kctrl *KeyResultController) GetKeyResult(c *gin.Context) {
	id := c.Param("id")

	kr, err := kctrl.keyResultService.GetData("id", id)
	if err != nil {
		response.NotFound(c, "Key result not found")
		return
	}

	response.OK(c, kr, "Key result retrieved successfully")
}

func (kctrl *KeyResultController) ListObjKeyResults(c *gin.Context) {
	objID := c.Param("id")

	kr, err := kctrl.keyResultService.ListData("objective_id", objID)
	if err != nil {
		response.NotFound(c, "Key results not found for objective")
		return
	}

	response.OK(c, kr, "Objective key results retrieved successfully")
}

func (kctrl *KeyResultController) ListAssigneeKeyResults(c *gin.Context) {
	assigneeID := c.Param("id")

	kr, err := kctrl.keyResultService.ListData("assignee_id", assigneeID)
	if err != nil {
		response.NotFound(c, "Key results not found for assignee")
		return
	}

	response.OK(c, kr, "Assignee key results retrieved successfully")
}

func (kctrl *KeyResultController) UpdateKeyResult(c *gin.Context) {
	var req dto.UpdateKeyResultRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	kr, err := kctrl.keyResultService.UpdateKeyResult(req)
	if err != nil {
		response.BadRequest(c, "Failed to update key result", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, kr, "Key result updated successfully")
}

func (kctrl *KeyResultController) DeleteKeyResult(c *gin.Context) {
	id := c.Param("id")

	err := kctrl.keyResultService.DeleteKeyResult(id)
	if err != nil {
		response.BadRequest(c, "Failed to delete key result", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, nil, "Key result deleted successfully")
}
