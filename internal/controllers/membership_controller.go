package controllers

import (
	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/Slightly-Techie/st-okr-api/internal/response"
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
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	data, err := ctrl.membershipService.CreateMembership(body)
	if err != nil {
		response.BadRequest(c, "Failed to create membership", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.Created(c, data, "Membership created successfully")
}

func (ctrl *MembershipController) GetMembership(c *gin.Context) {
	id := c.Param("id")

	data, err := ctrl.membershipService.GetMembership("id", id)
	if err != nil {
		response.NotFound(c, "Membership not found")
		return
	}

	response.OK(c, data, "Membership retrieved successfully")
}

func (ctrl *MembershipController) UpdateMembership(c *gin.Context) {
	var body dto.UpdateMembershipRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	data, err := ctrl.membershipService.UpdateMembership(body)
	if err != nil {
		response.BadRequest(c, "Failed to update membership", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, data, "Membership updated successfully")
}

func (ctrl *MembershipController) DeleteMembership(c *gin.Context) {
	id := c.Param("id")
	
	err := ctrl.membershipService.DeleteMembership(id)
	if err != nil {
		response.BadRequest(c, "Failed to delete membership", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, nil, "Membership deleted successfully")
}

func (ctrl *MembershipController) GetCompanyMembers(c *gin.Context) {
	companyID := c.Param("company_id")

	members, err := ctrl.membershipService.GetCompanyMembers(companyID)
	if err != nil {
		response.BadRequest(c, "Failed to retrieve company members", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, members, "Company members retrieved successfully")
}

func (ctrl *MembershipController) UpdateMembershipRole(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Role models.RoleType `json:"role" binding:"required,oneof=admin member viewer"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response.ValidationError(c, "Invalid role data", map[string]string{
			"role": "Must be one of: admin, member, viewer",
		})
		return
	}

	err := ctrl.membershipService.UpdateMembershipRole(id, body.Role)
	if err != nil {
		response.BadRequest(c, "Failed to update membership role", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, nil, "Membership role updated successfully")
}

func (ctrl *MembershipController) UpdateMembershipStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status models.StatusType `json:"status" binding:"required,oneof=active inactive suspended"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response.ValidationError(c, "Invalid status data", map[string]string{
			"status": "Must be one of: active, inactive, suspended",
		})
		return
	}

	err := ctrl.membershipService.UpdateMembershipStatus(id, body.Status)
	if err != nil {
		response.BadRequest(c, "Failed to update membership status", map[string]string{
			"service": err.Error(),
		})
		return
	}

	response.OK(c, nil, "Membership status updated successfully")
}