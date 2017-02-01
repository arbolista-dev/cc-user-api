package models

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

type UserGoal struct {
	UserGoalID   uint      `json:"-" db:"user_goal_id,omitempty"`
  Key          string    `json:"key" db:"key"`
  UserID       uint      `json:"-" db:"user_id"`
  Status       string    `json:"status" db:"status"`
  TonsSaved    float64   `json:"tons_saved" db:"tons_saved"`
  DollarsSaved float64   `json:"dollars_saved" db:"dollars_saved"`
  UpfrontCost  float64   `json:"upfront_cost" db:"upfront_cost"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
}

type UserGoalUpdate struct {
  Key      string         `json:"key" db:"key"`
  Status   string         `json:"status" db:"status"`
  Details  types.JSONText `json:"details"`
}

type UserGoalsList struct {
	TotalCount		uint64	   `json:"total_count"`
	List					[]UserGoal `json:"list"`
}
