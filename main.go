package main

import (
	"chi-app/app/auth"
	"chi-app/app/campaign"
	"chi-app/app/handler"
	"chi-app/app/helper"
	"chi-app/app/key"
	"chi-app/app/user"
	"chi-app/database"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func init() {
	// https://github.com/joho/godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	db, err := database.GetConnection()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	fmt.Println("MySQL Connected!")

	// repository
	userRepository := user.NewUserRepository(db)
	campaignRepository := campaign.NewCampaignRepository(db)

	// service
	userService := user.NewUserService(userRepository)
	authService := auth.NewJwtService()
	campaignService := campaign.NewCampaignService(campaignRepository)

	// handler
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// route list
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello Chi!"))
		})

		// USERS
		r.Post("/users", userHandler.RegisterUser)
		r.Post("/sessions", userHandler.Login)
		r.Post("/email_checkers", userHandler.CheckEmailAvailable)
		r.With(func(h http.Handler) http.Handler { return authMiddleware(h, authService, userService) }).Post("/avatars", userHandler.UploadAvatar)

		// CAMPAIGNS
		r.Get("/campaigns/{id}", campaignHandler.GetCampaignDetail)
		r.Get("/campaigns", campaignHandler.GetCampaigns)
		r.With(func(h http.Handler) http.Handler { return authMiddleware(h, authService, userService) }).Post("/campaigns", campaignHandler.CreateCampaign)
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func authMiddleware(h http.Handler, authService auth.Service, userService user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			helper.JSON(w, response, http.StatusUnauthorized)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			helper.JSON(w, response, http.StatusUnauthorized)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			helper.JSON(w, response, http.StatusUnauthorized)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			helper.JSON(w, response, http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctxAuth := context.WithValue(ctx, key.CtxKeyAuth{}, user)

		// next to route
		h.ServeHTTP(w, r.WithContext(ctxAuth))
	})
}
