package controllers

import (
	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
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
	requestID := getRequestID(c)
	userID := getUserID(c)
	remoteIP := c.ClientIP()

	logger.Info("Key result creation initiated",
		"request_id", requestID,
		"user_id", userID,
		"remote_ip", remoteIP,
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Key result creation failed - invalid request",
			"request_id", requestID,
			"user_id", userID,
			"remote_ip", remoteIP,
			"error", err.Error(),
		)
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	logger.Info("Creating key result",
		"request_id", requestID,
		"user_id", userID,
		"title", req.Title,
		"objective_id", req.ObjectiveID,
		"assignee_id", req.AssigneeID,
	)

	kr, err := kctrl.keyResultService.CreateKeyResult(req)
	if err != nil {
		logger.Error("Key result creation failed",
			"request_id", requestID,
			"user_id", userID,
			"title", req.Title,
			"objective_id", req.ObjectiveID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to create key result", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Key result created successfully",
		"request_id", requestID,
		"user_id", userID,
		"key_result_id", kr.ID,
		"title", kr.Title,
		"objective_id", kr.ObjectiveID,
	)

	response.Created(c, kr, "Key result created successfully")
}

func (kctrl *KeyResultController) GetKeyResult(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Key result retrieval requested",
		"request_id", requestID,
		"key_result_id", id,
		"user_id", userID,
	)

	kr, err := kctrl.keyResultService.GetData("id", id)
	if err != nil {
		logger.Warn("Key result not found",
			"request_id", requestID,
			"key_result_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.NotFound(c, "Key result not found")
		return
	}

	logger.Info("Key result retrieved successfully",
		"request_id", requestID,
		"key_result_id", kr.ID,
		"title", kr.Title,
		"user_id", userID,
	)

	response.OK(c, kr, "Key result retrieved successfully")
}

func (kctrl *KeyResultController) ListObjKeyResults(c *gin.Context) {
	objID := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Objective key results list requested",
		"request_id", requestID,
		"objective_id", objID,
		"user_id", userID,
	)

	kr, err := kctrl.keyResultService.ListData("objective_id", objID)
	if err != nil {
		logger.Warn("Key results not found for objective",
			"request_id", requestID,
			"objective_id", objID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.NotFound(c, "Key results not found for objective")
		return
	}

	logger.Info("Objective key results retrieved successfully",
		"request_id", requestID,
		"objective_id", objID,
		"user_id", userID,
		"key_result_count", len(kr),
	)

	response.OK(c, kr, "Objective key results retrieved successfully")
}

func (kctrl *KeyResultController) ListAssigneeKeyResults(c *gin.Context) {
	assigneeID := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Assignee key results list requested",
		"request_id", requestID,
		"assignee_id", assigneeID,
		"user_id", userID,
	)

	kr, err := kctrl.keyResultService.ListData("assignee_id", assigneeID)
	if err != nil {
		logger.Warn("Key results not found for assignee",
			"request_id", requestID,
			"assignee_id", assigneeID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.NotFound(c, "Key results not found for assignee")
		return
	}

	logger.Info("Assignee key results retrieved successfully",
		"request_id", requestID,
		"assignee_id", assigneeID,
		"user_id", userID,
		"key_result_count", len(kr),
	)

	response.OK(c, kr, "Assignee key results retrieved successfully")
}

func (kctrl *KeyResultController) UpdateKeyResult(c *gin.Context) {
	var req dto.UpdateKeyResultRequest
	requestID := getRequestID(c)
	userID := getUserID(c)
	keyResultID := c.Param("id")

	logger.Info("Key result update initiated",
		"request_id", requestID,
		"key_result_id", keyResultID,
		"user_id", userID,
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Key result update failed - invalid request",
			"request_id", requestID,
			"key_result_id", keyResultID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	logger.Info("Updating key result",
		"request_id", requestID,
		"key_result_id", keyResultID,
		"user_id", userID,
		"new_title", req.Title,
		"new_current_value", req.CurrentValue,
	)

	kr, err := kctrl.keyResultService.UpdateKeyResult(req)
	if err != nil {
		logger.Error("Key result update failed",
			"request_id", requestID,
			"key_result_id", keyResultID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to update key result", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Key result updated successfully",
		"request_id", requestID,
		"key_result_id", kr.ID,
		"title", kr.Title,
		"user_id", userID,
	)

	response.OK(c, kr, "Key result updated successfully")
}

func (kctrl *KeyResultController) DeleteKeyResult(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Key result deletion initiated",
		"request_id", requestID,
		"key_result_id", id,
		"user_id", userID,
	)

	err := kctrl.keyResultService.DeleteKeyResult(id)
	if err != nil {
		logger.Error("Key result deletion failed",
			"request_id", requestID,
			"key_result_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to delete key result", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Key result deleted successfully",
		"request_id", requestID,
		"key_result_id", id,
		"user_id", userID,
	)

	response.OK(c, nil, "Key result deleted successfully")
}
