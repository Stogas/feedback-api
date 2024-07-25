package feedbacktypes

import "gorm.io/gorm"

type Satisfaction struct {
	gorm.Model
	Satisfied *bool `json:"satisfied" binding:"required"`
}