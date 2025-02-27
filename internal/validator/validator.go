package validator

import (
	"time"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/go-playground/validator/v10"
)

func KeyResultValidators(v *validator.Validate) {
	v.RegisterValidation("future", validateDueDate)
	v.RegisterValidation("metric_type", validateMetricType)
	v.RegisterValidation("assignee_type", validateAssigneeType)
}

func validateDueDate(f validator.FieldLevel) bool {
	date, ok := f.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return date.After(time.Now())
}

func validateMetricType(f validator.FieldLevel) bool {
	metricType := models.MetricType(f.Field().String())
	switch metricType {
	case models.MetricTypeNumeric,
		models.MetricTypePercentage,
		models.MetricTypeBinary,
		models.MetrictTypeCurrency:
		return true
	default:
		return false
	}
}

func validateAssigneeType(f validator.FieldLevel) bool {
	assigneeType := models.AssigneeType(f.Field().String())
	switch assigneeType {
	case models.AssigneeTypeIndividual,
		models.AssigneeTypeTeam:
		return true
	default:
		return false
	}
}
