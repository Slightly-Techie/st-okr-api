package controllers

import (
	"github.com/Slightly-Techie/st-okr-api/internal/dto"
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

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	objective, err := ctrl.objectiveService.CreateObjective(req)
	if err != nil {
		response.BadRequest(c, "Failed to create objective", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.Created(c, objective, "Objective created successfully")
}

func (ctrl *ObjectiveController) GetObjective(c *gin.Context) {
	id := c.Param("id")

	objective, err := ctrl.objectiveService.GetObjective("id", id)
	if err != nil {
		response.NotFound(c, "Objective not found")
		return
	}

	response.OK(c, objective, "Objective retrieved successfully")
}

func (ctrl *ObjectiveController) GetObjectiveWithKeyResults(c *gin.Context) {
	id := c.Param("id")

	objective, err := ctrl.objectiveService.GetObjectiveWithKeyResults(id)
	if err != nil {
		response.NotFound(c, "Objective not found")
		return
	}

	response.OK(c, objective, "Objective with key results retrieved successfully")
}

func (ctrl *ObjectiveController) UpdateObjective(c *gin.Context) {
	var req dto.UpdateObjectiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	id := c.Param("id")
	req.ID = id

	objective, err := ctrl.objectiveService.UpdateObjective(req)
	if err != nil {
		response.BadRequest(c, "Failed to update objective", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, objective, "Objective updated successfully")
}

func (ctrl *ObjectiveController) DeleteObjective(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.objectiveService.DeleteObjective(id)
	if err != nil {
		response.BadRequest(c, "Failed to delete objective", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, nil, "Objective deleted successfully")
}

func (ctrl *ObjectiveController) ListCompanyObjectives(c *gin.Context) {
	companyID := c.Param("company_id")

	objectives, err := ctrl.objectiveService.ListObjectivesByCompany(companyID)
	if err != nil {
		response.BadRequest(c, "Failed to retrieve company objectives", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, objectives, "Company objectives retrieved successfully")
}

func (ctrl *ObjectiveController) ListTeamObjectives(c *gin.Context) {
	teamID := c.Param("team_id")

	objectives, err := ctrl.objectiveService.ListObjectivesByTeam(teamID)
	if err != nil {
		response.BadRequest(c, "Failed to retrieve team objectives", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, objectives, "Team objectives retrieved successfully")
}

func (ctrl *ObjectiveController) ListOwnerObjectives(c *gin.Context) {
	ownerID := c.Param("owner_id")

	objectives, err := ctrl.objectiveService.ListObjectivesByOwner(ownerID)
	if err != nil {
		response.BadRequest(c, "Failed to retrieve owner objectives", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, objectives, "Owner objectives retrieved successfully")
}

func (ctrl *ObjectiveController) UpdateObjectiveProgress(c *gin.Context) {
	objectiveID := c.Param("id")

	err := ctrl.objectiveService.UpdateObjectiveProgress(objectiveID)
	if err != nil {
		response.BadRequest(c, "Failed to update objective progress", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, nil, "Objective progress updated successfully")
}