package controllers

import (
	"github.com/Slightly-Techie/st-okr-api/internal/dto"
	"github.com/Slightly-Techie/st-okr-api/internal/logger"
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
	requestID := getRequestID(c)
	userID := getUserID(c)
	remoteIP := c.ClientIP()

	logger.Info("Membership creation initiated",
		"request_id", requestID,
		"user_id", userID,
		"remote_ip", remoteIP,
	)

	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Error("Membership creation failed - invalid request",
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

	logger.Info("Creating membership",
		"request_id", requestID,
		"user_id", userID,
		"member_user_id", body.UserID,
		"company_id", body.CompanyID,
		"role", body.Role,
	)

	data, err := ctrl.membershipService.CreateMembership(body)
	if err != nil {
		logger.Error("Membership creation failed",
			"request_id", requestID,
			"user_id", userID,
			"member_user_id", body.UserID,
			"company_id", body.CompanyID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to create membership", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Membership created successfully",
		"request_id", requestID,
		"user_id", userID,
		"membership_id", data.ID,
		"member_user_id", data.UserID,
		"company_id", data.CompanyID,
		"role", data.Role,
	)

	response.Created(c, data, "Membership created successfully")
}

func (ctrl *MembershipController) GetMembership(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Membership retrieval requested",
		"request_id", requestID,
		"membership_id", id,
		"user_id", userID,
	)

	data, err := ctrl.membershipService.GetMembership("id", id)
	if err != nil {
		logger.Warn("Membership not found",
			"request_id", requestID,
			"membership_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.NotFound(c, "Membership not found")
		return
	}

	logger.Info("Membership retrieved successfully",
		"request_id", requestID,
		"membership_id", data.ID,
		"member_user_id", data.UserID,
		"company_id", data.CompanyID,
		"user_id", userID,
	)

	response.OK(c, data, "Membership retrieved successfully")
}

func (ctrl *MembershipController) UpdateMembership(c *gin.Context) {
	var body dto.UpdateMembershipRequest
	requestID := getRequestID(c)
	userID := getUserID(c)
	membershipID := c.Param("id")

	logger.Info("Membership update initiated",
		"request_id", requestID,
		"membership_id", membershipID,
		"user_id", userID,
	)

	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Error("Membership update failed - invalid request",
			"request_id", requestID,
			"membership_id", membershipID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.ValidationError(c, "Invalid request data", map[string]string{
			"request": err.Error(),
		})
		return
	}

	logger.Info("Updating membership",
		"request_id", requestID,
		"membership_id", membershipID,
		"user_id", userID,
		"new_role", body.Role,
		"new_status", body.Status,
	)

	data, err := ctrl.membershipService.UpdateMembership(body)
	if err != nil {
		logger.Error("Membership update failed",
			"request_id", requestID,
			"membership_id", membershipID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to update membership", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Membership updated successfully",
		"request_id", requestID,
		"membership_id", data.ID,
		"role", data.Role,
		"status", data.Status,
		"user_id", userID,
	)

	response.OK(c, data, "Membership updated successfully")
}

func (ctrl *MembershipController) DeleteMembership(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Membership deletion initiated",
		"request_id", requestID,
		"membership_id", id,
		"user_id", userID,
	)
	
	err := ctrl.membershipService.DeleteMembership(id)
	if err != nil {
		logger.Error("Membership deletion failed",
			"request_id", requestID,
			"membership_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to delete membership", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Membership deleted successfully",
		"request_id", requestID,
		"membership_id", id,
		"user_id", userID,
	)

	response.OK(c, nil, "Membership deleted successfully")
}

func (ctrl *MembershipController) GetCompanyMembers(c *gin.Context) {
	companyID := c.Param("company_id")
	requestID := getRequestID(c)
	userID := getUserID(c)

	logger.Info("Company members list requested",
		"request_id", requestID,
		"company_id", companyID,
		"user_id", userID,
	)

	members, err := ctrl.membershipService.GetCompanyMembers(companyID)
	if err != nil {
		logger.Error("Failed to retrieve company members",
			"request_id", requestID,
			"company_id", companyID,
			"user_id", userID,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to retrieve company members", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Company members retrieved successfully",
		"request_id", requestID,
		"company_id", companyID,
		"user_id", userID,
		"member_count", len(members),
	)

	response.OK(c, members, "Company members retrieved successfully")
}

func (ctrl *MembershipController) UpdateMembershipRole(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)
	var body struct {
		Role models.RoleType `json:"role" binding:"required,oneof=admin member viewer"`
	}

	logger.Info("Membership role update initiated",
		"request_id", requestID,
		"membership_id", id,
		"user_id", userID,
	)

	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Error("Membership role update failed - invalid request",
			"request_id", requestID,
			"membership_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.ValidationError(c, "Invalid role data", map[string]string{
			"role": "Must be one of: admin, member, viewer",
		})
		return
	}

	logger.Info("Updating membership role",
		"request_id", requestID,
		"membership_id", id,
		"user_id", userID,
		"new_role", body.Role,
	)

	err := ctrl.membershipService.UpdateMembershipRole(id, body.Role)
	if err != nil {
		logger.Error("Membership role update failed",
			"request_id", requestID,
			"membership_id", id,
			"user_id", userID,
			"new_role", body.Role,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to update membership role", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Membership role updated successfully",
		"request_id", requestID,
		"membership_id", id,
		"user_id", userID,
		"new_role", body.Role,
	)

	response.OK(c, nil, "Membership role updated successfully")
}

func (ctrl *MembershipController) UpdateMembershipStatus(c *gin.Context) {
	id := c.Param("id")
	requestID := getRequestID(c)
	userID := getUserID(c)
	var body struct {
		Status models.StatusType `json:"status" binding:"required,oneof=active inactive suspended"`
	}

	logger.Info("Membership status update initiated",
		"request_id", requestID,
		"membership_id", id,
		"user_id", userID,
	)

	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Error("Membership status update failed - invalid request",
			"request_id", requestID,
			"membership_id", id,
			"user_id", userID,
			"error", err.Error(),
		)
		response.ValidationError(c, "Invalid status data", map[string]string{
			"status": "Must be one of: active, inactive, suspended",
		})
		return
	}

	logger.Info("Updating membership status",
		"request_id", requestID,
		"membership_id", id,
		"user_id", userID,
		"new_status", body.Status,
	)

	err := ctrl.membershipService.UpdateMembershipStatus(id, body.Status)
	if err != nil {
		logger.Error("Membership status update failed",
			"request_id", requestID,
			"membership_id", id,
			"user_id", userID,
			"new_status", body.Status,
			"error", err.Error(),
		)
		response.BadRequest(c, "Failed to update membership status", map[string]string{
			"service": err.Error(),
		})
		return
	}

	logger.Info("Membership status updated successfully",
		"request_id", requestID,
		"membership_id", id,
		"user_id", userID,
		"new_status", body.Status,
	)

	response.OK(c, nil, "Membership status updated successfully")
}