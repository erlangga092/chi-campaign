package campaign

import (
	"strings"
)

type CampaignFormatter struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	ImageURL         string `json:"image_url"`
	CurrentAmount    int    `json:"current_amount"`
	GoalAmount       int    `json:"goal_amount"`
	UserID           int    `json:"user_id"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	formatter := CampaignFormatter{}
	formatter.ID = campaign.ID
	formatter.Name = campaign.Name
	formatter.ShortDescription = campaign.ShortDescription
	formatter.Description = campaign.Description
	formatter.ImageURL = ""
	formatter.CurrentAmount = campaign.CurrentAmount
	formatter.GoalAmount = campaign.GoalAmount
	formatter.UserID = campaign.UserID

	if len(campaign.CampaignImages) > 0 {
		formatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return formatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	formatters := []CampaignFormatter{}

	for _, campaign := range campaigns {
		formatter := FormatCampaign(campaign)
		formatters = append(formatters, formatter)
	}

	return formatters
}

type CampaignDetailFormatter struct {
	ID               int                   `json:"id"`
	UserID           int                   `json:"user_id"`
	Name             string                `json:"string"`
	ShortDescription string                `json:"short_description"`
	Description      string                `json:"description"`
	ImageURL         string                `json:"image_url"`
	CurrentAmount    int                   `json:"current_amount"`
	GoalAmount       int                   `json:"goal_amount"`
	Slug             string                `json:"slug"`
	Perks            []string              `json:"perks"`
	User             CampaignUserFormatter `json:"user"`
	Images           []CampaignImageFormatter
}

type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignImageFormatter struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	formatter := CampaignDetailFormatter{}

	formatter.ID = campaign.ID
	formatter.Name = campaign.Name
	formatter.ShortDescription = campaign.ShortDescription
	formatter.Description = campaign.Description
	formatter.ImageURL = ""
	formatter.CurrentAmount = campaign.CurrentAmount
	formatter.GoalAmount = campaign.GoalAmount
	formatter.UserID = campaign.UserID
	formatter.Slug = campaign.Slug
	formatter.Perks = []string{}

	imagesFormatter := []CampaignImageFormatter{}
	for _, campaignImage := range campaign.CampaignImages {
		if campaignImage.IsPrimary {
			formatter.ImageURL = campaignImage.FileName
		}

		imageFormatter := CampaignImageFormatter{}
		imageFormatter.ImageURL = campaignImage.FileName
		imageFormatter.IsPrimary = campaignImage.IsPrimary
		imagesFormatter = append(imagesFormatter, imageFormatter)
	}

	for _, perk := range strings.Split(campaign.Perks, ",") {
		formatter.Perks = append(formatter.Perks, strings.TrimSpace(perk))
	}

	userFormatter := CampaignUserFormatter{}
	userFormatter.Name = campaign.User.Name
	userFormatter.ImageURL = campaign.User.AvatarFileName

	formatter.Images = imagesFormatter
	formatter.User = userFormatter
	return formatter
}
