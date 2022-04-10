package campaign

import "chi-app/app/user"

type GetCampaignDetailInput struct {
	ID int `uri:"id" validate:"required"`
}

type CreateCampaignInput struct {
	Name             string `json:"name" validate:"required"`
	ShortDescription string `json:"short_description" validate:"required"`
	Description      string `json:"description" validate:"required"`
	Perks            string `json:"perks" validate:"required"`
	GoalAmount       int    `json:"goal_amount" validate:"required"`
	User             user.User
}
