package campaign

import (
	"errors"
	"strings"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignDetail(ID GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	Update(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
}

type service struct {
	campaignRepository Repository
}

func NewCampaignService(campaignRepository Repository) Service {
	return &service{campaignRepository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	if userID != 0 {
		campaigns, err := s.campaignRepository.GetCampaignsByUserID(userID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	campaigns, err := s.campaignRepository.GetCampaigns()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignDetail(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.campaignRepository.GetCampaignByID(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID

	slug := strings.ToLower(strings.Join(strings.Split(campaign.Name, " "), "-"))
	campaign.Slug = slug

	newCampaign, err := s.campaignRepository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *service) Update(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	campaign, err := s.campaignRepository.GetCampaignByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("not an owner of the campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	slug := strings.ToLower(strings.Join(strings.Split(campaign.Name, " "), "-"))
	campaign.Slug = slug

	updatedCampaign, err := s.campaignRepository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}
