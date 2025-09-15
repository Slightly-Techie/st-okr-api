package controllers

import (
	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
	"github.com/Slightly-Techie/st-okr-api/internal/response"
	"github.com/Slightly-Techie/st-okr-api/internal/services"
	"github.com/gin-gonic/gin"
)

type TeamController struct {
	teamService services.TeamService
}

func NewTeamController(teamService services.TeamService) *TeamController {
	return &TeamController{
		teamService: teamService,
	}
}

func (tctrl *TeamController) CreateTeam(c *gin.Context) {
	var teamDTO dto.CreateTeamRequest
	requestID := getRequestID(c)
	userID := getUserID(c)
	remoteIP := c.ClientIP()

	logger.Info("Team creation initiated",
		"request_id", requestID,
		"user_id", userID,
		"remote_ip", remoteIP,
	)

	if err := c.ShouldBindJSON(&teamDTO); err != nil {
		logger.Error("Team creation failed - invalid request",
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

	logger.Info("Creating team",
		"request_id", requestID,
		"user_id", userID,
		"team_name", teamDTO.Name,
		"company_id", teamDTO.CompanyID,
	)

	team, err := tctrl.teamService.CreateTeam(teamDTO)
	if err != nil {
		logger.Error("Team creation failed",
			"request_id", requestID,
			"user_id", userID,
			"team_name", teamDTO.Name,
			"company_id", teamDTO.CompanyID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to create team", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Team created successfully",
		"request_id", requestID,
		"user_id", userID,
		"team_id", team.ID,
		"team_name", team.Name,
		"company_id", team.CompanyID,
	)

	response.Created(c, team, "Team created successfully")
}

func (tctrl *TeamController) GetTeam(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Team retrieval requested",
		"request_id", requestID,
		"team_id", id,
		"user_id", userID,
	)

	team, err := tctrl.teamService.GetTeam("id", id)
	if err != nil {
		logger.Warn("Team not found",
			"request_id", requestID,
			"team_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.NotFound(c, "Team not found")
		return
	}

	logger.Info("Team retrieved successfully",
		"request_id", requestID,
		"team_id", team.ID,
		"team_name", team.Name,
		"user_id", userID,
	)

	response.OK(c, team, "Team retrieved successfully")
}

func (tctrl *TeamController) UpdateTeam(c *gin.Context) {
	var teamDTO dto.UpdateTeamRequest
	requestID := getRequestID(c)
	userID := getUserID(c)
	teamID := c.Param("id")

	logger.Info("Team update initiated",
		"request_id", requestID,
		"team_id", teamID,
		"user_id", userID,
	)

	if err := c.ShouldBindJSON(&teamDTO); err != nil {
		logger.Error("Team update failed - invalid request",
			"request_id", requestID,
			"team_id", teamID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	logger.Info("Updating team",
		"request_id", requestID,
		"team_id", teamID,
		"user_id", userID,
		"new_name", teamDTO.Name,
	)

	team, err := tctrl.teamService.UpdateTeam(teamDTO)
	if err != nil {
		logger.Error("Team update failed",
			"request_id", requestID,
			"team_id", teamID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to update team", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Team updated successfully",
		"request_id", requestID,
		"team_id", team.ID,
		"team_name", team.Name,
		"user_id", userID,
	)

	response.OK(c, team, "Team updated successfully")
}

func (tctrl *TeamController) DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Team deletion initiated",
		"request_id", requestID,
		"team_id", id,
		"user_id", userID,
	)

	err := tctrl.teamService.DeleteTeam(id)
	if err != nil {
		logger.Error("Team deletion failed",
			"request_id", requestID,
			"team_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to delete team", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Team deleted successfully",
		"request_id", requestID,
		"team_id", id,
		"user_id", userID,
	)

	response.OK(c, nil, "Team deleted successfully")
}

func (tctrl *TeamController) AddTeamMember(c *gin.Context) {
	var addMemberDTO dto.TeamMemberRequest
	requestID := getRequestID(c)
	userID := getUserID(c)
	remoteIP := c.ClientIP()

	logger.Info("Team member addition initiated",
		"request_id", requestID,
		"user_id", userID,
		"remote_ip", remoteIP,
	)

	if err := c.ShouldBindJSON(&addMemberDTO); err != nil {
		logger.Error("Team member addition failed - invalid request",
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

	logger.Info("Adding team member",
		"request_id", requestID,
		"user_id", userID,
		"team_id", addMemberDTO.TeamID,
		"member_user_id", addMemberDTO.UserID,
	)

	member, err := tctrl.teamService.AddMember(&addMemberDTO)
	if err != nil {
		if err.Error() == "user is already a member of the team" {
			logger.Warn("Team member addition failed - user already exists",
				"request_id", requestID,
				"user_id", userID,
				"team_id", addMemberDTO.TeamID,
				"member_user_id", addMemberDTO.UserID,
			)
			response.Conflict(c, "User is already a member of the team", nil)
			return
		}
		logger.Error("Team member addition failed",
			"request_id", requestID,
			"user_id", userID,
			"team_id", addMemberDTO.TeamID,
			"member_user_id", addMemberDTO.UserID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to add team member", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Team member added successfully",
		"request_id", requestID,
		"user_id", userID,
		"team_id", member.TeamID,
		"member_user_id", member.UserID,
	)

	response.Created(c, member, "Team member added successfully")
}

func (tctrl *TeamController) ListTeamMembers(c *gin.Context) {
	teamID := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Team members list requested",
		"request_id", requestID,
		"team_id", teamID,
		"user_id", userID,
	)

	members, err := tctrl.teamService.ListMembers("team_id", teamID)
	if err != nil {
		logger.Warn("Team members not found",
			"request_id", requestID,
			"team_id", teamID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.NotFound(c, "Team members not found")
		return
	}

	logger.Info("Team members retrieved successfully",
		"request_id", requestID,
		"team_id", teamID,
		"user_id", userID,
		"member_count", len(members),
	)

	response.OK(c, members, "Team members retrieved successfully")
}

func (tctrl *TeamController) RemoveMember(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Team member removal initiated",
		"request_id", requestID,
		"member_id", id,
		"user_id", userID,
	)

	err := tctrl.teamService.RemoveMember(id)
	if err != nil {
		logger.Error("Team member removal failed",
			"request_id", requestID,
			"member_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to remove team member", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Team member removed successfully",
		"request_id", requestID,
		"member_id", id,
		"user_id", userID,
	)

	response.OK(c, nil, "Team member removed successfully")
}
