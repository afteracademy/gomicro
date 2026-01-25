package dto

import (
	"time"
)

type HealthCheck struct {
	Timestamp time.Time `json:"timestamp" binding:"required"`
	Status    string    `json:"status" binding:"required"`
}
