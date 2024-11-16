package controllers

import (
	"net/http"

	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/Slightly-Techie/st-okr-api/internal/services"
	"github.com/gin-gonic/gin"
)

type MembershipController struct {
	membershipService services.MembershipService
}

func NewMembershipController(membershipService services.MembershipService) *MembershipController {
	return &MembershipController{
		membershipService: membershipService,
	}
}

func (ctrl *MembershipController) CreateMembership(c *gin.Context) {
	var body dto.CreateMembershipRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := ctrl.membershipService.CreateMembership(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Membership created successfully",
		"data":    data,
	})
}

func (ctrl *MembershipController) GetMembership(c *gin.Context) {
	id := c.Param("id")

	data, err := ctrl.membershipService.GetMembership("id", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Membership found",
		"data":    data,
	})
}

func (ctrl *MembershipController) UpdateMembership(c *gin.Context) {
	var body dto.UpdateMembershipRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := ctrl.membershipService.UpdateMembership(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Membership updated successfully",
		"data":    data,
	})
}

func (ctrl *MembershipController) DeleteMembership(c *gin.Context) {
	id := c.Param("id")
	
	err := ctrl.membershipService.DeleteMembership(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Membership deleted successfully"})
}

func (ctrl *MembershipController) GetCompanyMembers(c *gin.Context) {
	companyID := c.Param("company_id")

	members, err := ctrl.membershipService.GetCompanyMembers(companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company members retrieved successfully",
		"data":    members,
	})
}

func (ctrl *MembershipController) UpdateMembershipRole(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Role models.RoleType `json:"role" binding:"required,oneof=admin member viewer"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.membershipService.UpdateMembershipRole(id, body.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Membership role updated successfully"})
}

func (ctrl *MembershipController) UpdateMembershipStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status models.StatusType `json:"status" binding:"required,oneof=active inactive suspended"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.membershipService.UpdateMembershipStatus(id, body.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Membership status updated successfully"})
}