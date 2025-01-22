package controllers

import (
	"net/http"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := tctrl.teamService.CreateTeam(teamDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Team created successfully",
		"team":    team,
	})
}

func (tctrl *TeamController) GetTeam(c *gin.Context) {
	id := c.Param("id")

	team, err := tctrl.teamService.GetTeam("id", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Team found",
		"team":    team,
	})
}

func (tctrl *TeamController) UpdateTeam(c *gin.Context) {
	var teamDTO dto.UpdateTeamRequest

	if err := c.ShouldBindJSON(&teamDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := tctrl.teamService.UpdateTeam(teamDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Team updated successfully",
		"Updated team": team,
	})
}

func (tctrl *TeamController) DeleteTeam(c *gin.Context) {
	id := c.Param("id")

	err := tctrl.teamService.DeleteTeam(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Team deleted successfully",
	})
}

func (tctrl *TeamController) AddTeamMember(c *gin.Context) {
	var addMemberDTO dto.TeamMemberRequest

	if err := c.ShouldBindJSON(&addMemberDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := tctrl.teamService.AddMember(&addMemberDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Team member added successfully",
		"member":  member,
	})
}

func (tctrl *TeamController) ListTeamMembers(c *gin.Context) {
	teamID := c.Param("teamId")

	members, err := tctrl.teamService.ListMembers(teamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"members": members,
	})
}

func (tctrl *TeamController) RemoveMember(c *gin.Context) {
	id := c.Param("id")

	err := tctrl.teamService.RemoveMember(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Team member removed successfully",
	})
}
