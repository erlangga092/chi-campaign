package campaign

import (
	"chi-app/app/user"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Repository interface {
	Save(campaign Campaign) (Campaign, error)
	GetCampaignByID(ID int) (Campaign, error)
	GetCampaigns() ([]Campaign, error)
	GetCampaignsByUserID(userID int) ([]Campaign, error)
	FindCampaignImagesByCampaignID(campaignID int) ([]CampaignImage, error)
	Update(campaign Campaign) (Campaign, error)
}

type repository struct {
	DB *sql.DB
}

const (
	layoutDateTime string = "2006-01-02 15:04:05"
)

func NewCampaignRepository(DB *sql.DB) Repository {
	return &repository{DB}
}

func (r *repository) Save(campaign Campaign) (Campaign, error) {
	sqlQuery := sq.Insert("campaigns").Columns(
		"user_id",
		"name",
		"short_description",
		"description",
		"perks",
		"backer_count",
		"goal_amount",
		"current_amount",
		"slug",
		"created_at",
		"updated_at").
		Values(
			campaign.UserID,
			campaign.Name,
			campaign.ShortDescription,
			campaign.Description,
			campaign.Perks,
			campaign.BackerCount,
			campaign.GoalAmount,
			campaign.CurrentAmount,
			campaign.Slug,
			time.Now().Format(layoutDateTime),
			time.Now().Format(layoutDateTime)).
		RunWith(r.DB)

	result, err := sqlQuery.Exec()
	if err != nil {
		return campaign, err
	}

	campaignID, err := result.LastInsertId()
	if err != nil {
		return campaign, err
	}

	newCampaign, err := r.GetCampaignByID(int(campaignID))
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (r *repository) GetCampaignByID(ID int) (Campaign, error) {
	campaign := Campaign{}
	user := user.User{}

	sqlQuery := sq.Select(
		"campaigns.id",
		"campaigns.user_id",
		"campaigns.name",
		"campaigns.short_description",
		"campaigns.description",
		"campaigns.perks",
		"campaigns.backer_count",
		"campaigns.goal_amount",
		"campaigns.current_amount",
		"campaigns.slug",
		"campaigns.created_at",
		"campaigns.updated_at",
		"users.name",
		"users.avatar_file_name").
		From("campaigns").
		Join("users ON users.id = campaigns.user_id").
		Where(sq.Eq{"campaigns.id": ID})

	rows, err := sqlQuery.RunWith(r.DB).Query()
	if err != nil {
		return campaign, err
	}

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(
			&campaign.ID,
			&campaign.UserID,
			&campaign.Name,
			&campaign.ShortDescription,
			&campaign.Description,
			&campaign.Perks,
			&campaign.BackerCount,
			&campaign.GoalAmount,
			&campaign.CurrentAmount,
			&campaign.Slug,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
			&user.Name,
			&user.AvatarFileName,
		)

		if err != nil {
			return campaign, err
		}

		campaignImages, err := r.FindCampaignImagesByCampaignID(campaign.ID)
		if err != nil {
			return campaign, err
		}

		campaign.CampaignImages = campaignImages
	}

	campaign.User = user
	return campaign, nil
}

func (r *repository) GetCampaigns() ([]Campaign, error) {
	campaigns := []Campaign{}

	sqlQuery := sq.Select(
		"id",
		"user_id",
		"name",
		"short_description",
		"description",
		"perks",
		"backer_count",
		"goal_amount",
		"current_amount",
		"slug",
		"created_at",
		"updated_at").
		From("campaigns")

	rows, err := sqlQuery.RunWith(r.DB).Query()
	if err != nil {
		return campaigns, err
	}

	defer rows.Close()

	for rows.Next() {
		campaign := Campaign{}

		err := rows.Scan(
			&campaign.ID,
			&campaign.UserID,
			&campaign.Name,
			&campaign.ShortDescription,
			&campaign.Description,
			&campaign.Perks,
			&campaign.BackerCount,
			&campaign.GoalAmount,
			&campaign.CurrentAmount,
			&campaign.Slug,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		)

		if err != nil {
			return campaigns, err
		}

		campaignImages, err := r.FindCampaignImagesByCampaignID(campaign.ID)
		if err != nil {
			return campaigns, err
		}

		campaign.CampaignImages = campaignImages
		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

func (r *repository) GetCampaignsByUserID(userID int) ([]Campaign, error) {
	campaigns := []Campaign{}

	sqlQuery := sq.Select(
		"id",
		"user_id",
		"name",
		"short_description",
		"description",
		"perks",
		"backer_count",
		"goal_amount",
		"current_amount",
		"slug",
		"created_at",
		"updated_at").
		From("campaigns").
		Where(sq.Eq{"user_id": userID})

	rows, err := sqlQuery.RunWith(r.DB).Query()
	if err != nil {
		return campaigns, err
	}

	defer rows.Close()

	for rows.Next() {
		campaign := Campaign{}

		err := rows.Scan(
			&campaign.ID,
			&campaign.UserID,
			&campaign.Name,
			&campaign.ShortDescription,
			&campaign.Description,
			&campaign.Perks,
			&campaign.BackerCount,
			&campaign.GoalAmount,
			&campaign.CurrentAmount,
			&campaign.Slug,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		)

		if err != nil {
			return campaigns, err
		}

		campaignImages, err := r.FindCampaignImagesByCampaignID(campaign.ID)
		if err != nil {
			return campaigns, err
		}

		campaign.CampaignImages = campaignImages
		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

func (r *repository) FindCampaignImagesByCampaignID(campaignID int) ([]CampaignImage, error) {
	campaignImages := []CampaignImage{}

	sqlQuery := sq.Select(
		"id",
		"campaign_id",
		"file_name",
		"is_primary",
		"created_at",
		"updated_at").
		From("campaign_images").
		Where(sq.Eq{"campaign_id": campaignID})

	rows, err := sqlQuery.RunWith(r.DB).Query()
	if err != nil {
		return campaignImages, err
	}

	defer rows.Close()

	for rows.Next() {
		var isPrimaryNum int
		campaignImage := CampaignImage{}

		rows.Scan(
			&campaignImage.ID,
			&campaignImage.CampaignID,
			&campaignImage.FileName,
			&isPrimaryNum,
			&campaignImage.CreatedAt,
			&campaignImage.UpdatedAt,
		)

		isPrimary := false
		if isPrimaryNum == 1 {
			isPrimary = true
		}

		campaignImage.IsPrimary = isPrimary
		campaignImages = append(campaignImages, campaignImage)
	}

	return campaignImages, nil
}

func (r *repository) Update(campaign Campaign) (Campaign, error) {
	sqlQuery := sq.Update("campaigns").
		Set("name", campaign.Name).
		Set("short_description", campaign.ShortDescription).
		Set("perks", campaign.Perks).
		Set("backer_count", campaign.BackerCount).
		Set("goal_amount", campaign.GoalAmount).
		Set("current_amount", campaign.CurrentAmount).
		Set("slug", campaign.Slug).
		Set("updated_at", time.Now().Format(layoutDateTime)).
		Where(sq.Eq{"id": campaign.ID}).RunWith(r.DB)

	result, err := sqlQuery.Exec()
	if err != nil {
		return campaign, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return campaign, err
	}

	if int(affected) > 0 {
		updatedCampaign, err := r.GetCampaignByID(campaign.ID)
		if err != nil {
			return updatedCampaign, err
		}

		return updatedCampaign, nil
	} else {
		return Campaign{}, errors.New("failed to update campaign")
	}
}
