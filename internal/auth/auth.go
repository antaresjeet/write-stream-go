package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"

	"write-stream-go/internal/config"
	"write-stream-go/internal/models"
)

type AuthService struct {
	config *config.Config
	db     *gorm.DB
}

func NewAuthService(cfg *config.Config, db *gorm.DB) *AuthService {
	return &AuthService{
		config: cfg,
		db:     db,
	}
}

func (s *AuthService) googleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.config.GoogleClientID,
		ClientSecret: s.config.GoogleClientSecret,
		RedirectURL:  s.config.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func (s *AuthService) GoogleAuthHandler(c *gin.Context) {
	url := s.googleOauthConfig().AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *AuthService) GoogleCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	token, err := s.googleOauthConfig().Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := s.googleOauthConfig().Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Picture  string `json:"picture"`
		SocialID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
		return
	}

	user := models.User{
		SocialID:  userInfo.SocialID,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		AvatarURL: userInfo.Picture,
	}

	if err := s.db.Where(models.User{SocialID: user.SocialID}).FirstOrCreate(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update user"})
		return
	}

	jwtToken, err := s.generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

func (s *AuthService) generateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})
}

func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
