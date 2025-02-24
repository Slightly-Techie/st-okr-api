package validator

import (
	"fmt"
	"time"

	"github.com/Slightly-Techie/st-okr-api/internal/models"
	"github.com/go-playground/validator/v10"
)

func KeyResultValidators(v *validator.Validate) {
	v.RegisterValidation("due_date", validateDueDate)
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

func validateMetricValues(kr *models.KeyResult) error {
	switch kr.MetricType {

	case models.MetricTypeNumeric, models.MetrictTypeCurrency:
		if kr.CurrentValue < 0 {
			return fmt.Errorf("current value cannot be negative")
		}

	case models.MetricTypePercentage:
		if kr.TargetValue < 0 || kr.TargetValue > 100 {
			return fmt.Errorf("percentage target must be between 0 and 100")
		}
		if kr.CurrentValue < 0 || kr.CurrentValue > 100 {
			return fmt.Errorf("percentage current value must be between 0 and 100")
		}

	case models.MetricTypeBinary:
		if kr.TargetValue != 0 && kr.TargetValue != 1 {
			return fmt.Errorf("boolean target must be 0 or 1")
		}
		if kr.CurrentValue != 0 && kr.CurrentValue != 1 {
			return fmt.Errorf("boolean current value must be 0 or 1")
		}
	}
	return nil
}
