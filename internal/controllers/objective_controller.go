package controllers

import (
	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
	"github.com/Slightly-Techie/st-okr-api/internal/response"
	"github.com/Slightly-Techie/st-okr-api/internal/services"
	"github.com/gin-gonic/gin"
)

type ObjectiveController struct {
	objectiveService services.ObjectiveService
}

func NewObjectiveController(objectiveService services.ObjectiveService) *ObjectiveController {
	return &ObjectiveController{
		objectiveService: objectiveService,
	}
}

func (ctrl *ObjectiveController) CreateObjective(c *gin.Context) {
	var req dto.CreateObjectiveRequest
	requestID := getRequestID(c)
	userID := getUserID(c)
	remoteIP := c.ClientIP()

	logger.Info("Objective creation initiated",
		"request_id", requestID,
		"user_id", userID,
		"remote_ip", remoteIP,
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Objective creation failed - invalid request",
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

	logger.Info("Creating objective",
		"request_id", requestID,
		"user_id", userID,
		"title", req.Title,
		"company_id", req.CompanyID,
		"team_id", req.TeamID,
		"owner_id", req.OwnerID,
	)

	objective, err := ctrl.objectiveService.CreateObjective(req)
	if err != nil {
		logger.Error("Objective creation failed",
			"request_id", requestID,
			"user_id", userID,
			"title", req.Title,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to create objective", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Objective created successfully",
		"request_id", requestID,
		"user_id", userID,
		"objective_id", objective.ID,
		"title", objective.Title,
	)

	response.Created(c, objective, "Objective created successfully")
}

func (ctrl *ObjectiveController) GetObjective(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Objective retrieval requested",
		"request_id", requestID,
		"objective_id", id,
		"user_id", userID,
	)

	objective, err := ctrl.objectiveService.GetObjective("id", id)
	if err != nil {
		logger.Warn("Objective not found",
			"request_id", requestID,
			"objective_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.NotFound(c, "Objective not found")
		return
	}

	logger.Info("Objective retrieved successfully",
		"request_id", requestID,
		"objective_id", objective.ID,
		"title", objective.Title,
		"user_id", userID,
	)

	response.OK(c, objective, "Objective retrieved successfully")
}

func (ctrl *ObjectiveController) GetObjectiveWithKeyResults(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Objective with key results requested",
		"request_id", requestID,
		"objective_id", id,
		"user_id", userID,
	)

	objective, err := ctrl.objectiveService.GetObjectiveWithKeyResults(id)
	if err != nil {
		logger.Warn("Objective with key results not found",
			"request_id", requestID,
			"objective_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.NotFound(c, "Objective not found")
		return
	}

	logger.Info("Objective with key results retrieved successfully",
		"request_id", requestID,
		"objective_id", objective.ID,
		"title", objective.Title,
		"user_id", userID,
	)

	response.OK(c, objective, "Objective with key results retrieved successfully")
}

func (ctrl *ObjectiveController) UpdateObjective(c *gin.Context) {
	var req dto.UpdateObjectiveRequest
	requestID := getRequestID(c)
	userID := getUserID(c)
	id := c.Param("id")

	logger.Info("Objective update initiated",
		"request_id", requestID,
		"objective_id", id,
		"user_id", userID,
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Objective update failed - invalid request",
			"request_id", requestID,
			"objective_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	req.ID = id

	logger.Info("Updating objective",
		"request_id", requestID,
		"objective_id", id,
		"user_id", userID,
		"new_title", req.Title,
	)

	objective, err := ctrl.objectiveService.UpdateObjective(req)
	if err != nil {
		logger.Error("Objective update failed",
			"request_id", requestID,
			"objective_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to update objective", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Objective updated successfully",
		"request_id", requestID,
		"objective_id", objective.ID,
		"title", objective.Title,
		"user_id", userID,
	)

	response.OK(c, objective, "Objective updated successfully")
}

func (ctrl *ObjectiveController) DeleteObjective(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Objective deletion initiated",
		"request_id", requestID,
		"objective_id", id,
		"user_id", userID,
	)

	err := ctrl.objectiveService.DeleteObjective(id)
	if err != nil {
		logger.Error("Objective deletion failed",
			"request_id", requestID,
			"objective_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to delete objective", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Objective deleted successfully",
		"request_id", requestID,
		"objective_id", id,
		"user_id", userID,
	)

	response.OK(c, nil, "Objective deleted successfully")
}

func (ctrl *ObjectiveController) ListCompanyObjectives(c *gin.Context) {
	companyID := c.Param("company_id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Company objectives list requested",
		"request_id", requestID,
		"company_id", companyID,
		"user_id", userID,
	)

	objectives, err := ctrl.objectiveService.ListObjectivesByCompany(companyID)
	if err != nil {
		logger.Error("Failed to retrieve company objectives",
			"request_id", requestID,
			"company_id", companyID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to retrieve company objectives", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Company objectives retrieved successfully",
		"request_id", requestID,
		"company_id", companyID,
		"user_id", userID,
		"objective_count", len(objectives),
	)

	response.OK(c, objectives, "Company objectives retrieved successfully")
}

func (ctrl *ObjectiveController) ListTeamObjectives(c *gin.Context) {
	teamID := c.Param("team_id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Team objectives list requested",
		"request_id", requestID,
		"team_id", teamID,
		"user_id", userID,
	)

	objectives, err := ctrl.objectiveService.ListObjectivesByTeam(teamID)
	if err != nil {
		logger.Error("Failed to retrieve team objectives",
			"request_id", requestID,
			"team_id", teamID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to retrieve team objectives", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Team objectives retrieved successfully",
		"request_id", requestID,
		"team_id", teamID,
		"user_id", userID,
		"objective_count", len(objectives),
	)

	response.OK(c, objectives, "Team objectives retrieved successfully")
}

func (ctrl *ObjectiveController) ListOwnerObjectives(c *gin.Context) {
	ownerID := c.Param("owner_id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Owner objectives list requested",
		"request_id", requestID,
		"owner_id", ownerID,
		"user_id", userID,
	)

	objectives, err := ctrl.objectiveService.ListObjectivesByOwner(ownerID)
	if err != nil {
		logger.Error("Failed to retrieve owner objectives",
			"request_id", requestID,
			"owner_id", ownerID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to retrieve owner objectives", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Owner objectives retrieved successfully",
		"request_id", requestID,
		"owner_id", ownerID,
		"user_id", userID,
		"objective_count", len(objectives),
	)

	response.OK(c, objectives, "Owner objectives retrieved successfully")
}

func (ctrl *ObjectiveController) UpdateObjectiveProgress(c *gin.Context) {
	objectiveID := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Objective progress update initiated",
		"request_id", requestID,
		"objective_id", objectiveID,
		"user_id", userID,
	)

	err := ctrl.objectiveService.UpdateObjectiveProgress(objectiveID)
	if err != nil {
		logger.Error("Objective progress update failed",
			"request_id", requestID,
			"objective_id", objectiveID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to update objective progress", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Objective progress updated successfully",
		"request_id", requestID,
		"objective_id", objectiveID,
		"user_id", userID,
	)

	response.OK(c, nil, "Objective progress updated successfully")
}