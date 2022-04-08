package handler

import (
	"chi-app/app/campaign"
	"chi-app/app/helper"
	"chi-app/app/key"
	"chi-app/app/user"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

func (h *campaignHandler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))

	campaigns, err := h.campaignService.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get campaigns", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	formatter := campaign.FormatCampaigns(campaigns)
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", formatter)
	helper.JSON(w, response, http.StatusOK)
}

func (h *campaignHandler) GetCampaignDetail(w http.ResponseWriter, r *http.Request) {
	campaignID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response := helper.APIResponse("Failed to get detail campaign", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	detailCampaign, err := h.campaignService.GetCampaignDetail(campaignID)
	if err != nil {
		response := helper.APIResponse("Failed to get campaigns", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	formatter := campaign.FormatCampaignDetail(detailCampaign)
	response := helper.APIResponse("Detail Campaign", http.StatusOK, "success", formatter)
	helper.JSON(w, response, http.StatusOK)
}

func (h *campaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		data := "Content Type must be application/json"
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", data)
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	v := validator.New()
	input := campaign.CreateCampaignInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		fmt.Println("Error when parsing")
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	err = v.Struct(input)
	if err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Error())
		}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errors)
		helper.JSON(w, response, http.StatusUnprocessableEntity)
		return
	}

	userCtx := r.Context().Value(key.CtxKeyAuth{}).(user.User)
	input.User = userCtx

	newCampaign, err := h.campaignService.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", err.Error())
		helper.JSON(w, response, http.StatusBadRequest)
		return
	}

	formatter := campaign.FormatCampaign(newCampaign)
	response := helper.APIResponse("Success to create campaign", http.StatusCreated, "success", formatter)
	helper.JSON(w, response, http.StatusCreated)
}
