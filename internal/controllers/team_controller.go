package controllers

import (
	"github.com/Slightly-Techie/st-okr-api/internal/dto"
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

	if err := c.ShouldBindJSON(&teamDTO); err != nil {
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	team, err := tctrl.teamService.CreateTeam(teamDTO)
	if err != nil {
		response.BadRequest(c, "Failed to create team", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.Created(c, team, "Team created successfully")
}

func (tctrl *TeamController) GetTeam(c *gin.Context) {
	id := c.Param("id")

	team, err := tctrl.teamService.GetTeam("id", id)
	if err != nil {
		response.NotFound(c, "Team not found")
		return
	}

	response.OK(c, team, "Team retrieved successfully")
}

func (tctrl *TeamController) UpdateTeam(c *gin.Context) {
	var teamDTO dto.UpdateTeamRequest

	if err := c.ShouldBindJSON(&teamDTO); err != nil {
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	team, err := tctrl.teamService.UpdateTeam(teamDTO)
	if err != nil {
		response.BadRequest(c, "Failed to update team", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, team, "Team updated successfully")
}

func (tctrl *TeamController) DeleteTeam(c *gin.Context) {
	id := c.Param("id")

	err := tctrl.teamService.DeleteTeam(id)
	if err != nil {
		response.BadRequest(c, "Failed to delete team", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, nil, "Team deleted successfully")
}

func (tctrl *TeamController) AddTeamMember(c *gin.Context) {
	var addMemberDTO dto.TeamMemberRequest

	if err := c.ShouldBindJSON(&addMemberDTO); err != nil {
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	member, err := tctrl.teamService.AddMember(&addMemberDTO)
	if err != nil {
		if err.Error() == "user is already a member of the team" {
			response.Conflict(c, "User is already a member of the team", nil)
			return
		}
		response.BadRequest(c, "Failed to add team member", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.Created(c, member, "Team member added successfully")
}

func (tctrl *TeamController) ListTeamMembers(c *gin.Context) {
	teamID := c.Param("id")

	members, err := tctrl.teamService.ListMembers("team_id", teamID)
	if err != nil {
		response.NotFound(c, "Team members not found")
		return
	}

	response.OK(c, members, "Team members retrieved successfully")
}

func (tctrl *TeamController) RemoveMember(c *gin.Context) {
	id := c.Param("id")

	err := tctrl.teamService.RemoveMember(id)
	if err != nil {
		response.BadRequest(c, "Failed to remove team member", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, nil, "Team member removed successfully")
}
