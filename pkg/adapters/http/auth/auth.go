package auth

import (
	"crypto/ecdsa"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

var (
	ErrNoAuthHeader         = errors.New("Header Authorization not found")
	ErrInvalidSigningMethod = errors.New("JWT Signing method is different from server's")
)

type Auth struct {
	key *ecdsa.PrivateKey
}

func New(key *ecdsa.PrivateKey) *Auth {
	return &Auth{
		key: key,
	}
}

func (a *Auth) GetMiddleware() gin.HandlerFunc {
	return a.auth
}

func (a *Auth) auth(c *gin.Context) {
	data := c.Request.Header.Get("Authorization")

	splitted := strings.Split(data, " ")
	token := splitted[1]

	if token == "" {
		c.Error(ErrNoAuthHeader) //nolint: errcheck
		c.JSON(http.StatusUnauthorized, gin.H{"status": "authorization failed"})
		return
	}

	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("ES512") != t.Method {
			return nil, ErrInvalidSigningMethod
		}

		return a.key.PublicKey, nil
	})

	if err != nil {
		c.Error(err) //nolint: errcheck
		c.JSON(http.StatusUnauthorized, gin.H{"status": "authorization failed"})
		return
	}

	c.Next()
}

func (a *Auth) Login(c *gin.Context) {
	data := domain.LoginRequest{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.Error(err) //nolint: errcheck
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong body json"})
		return
	}

	if data.User != "Admin" || data.Password != "Admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong user credential"})
		return
	}

	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodES512,
		jwt.MapClaims{"exp": time.Now().Add(15 * time.Minute), "iat": time.Now(), "iss": "Pavel"},
	)

	access, err := accessToken.SignedString(a.key)
	if err != nil {
		log.Error().Err(err).Send()
	}

	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodES512,
		jwt.MapClaims{"exp": time.Now().Add(7 * 24 * time.Hour), "iat": time.Now(), "iss": "Pavel"},
	)

	refresh, err := refreshToken.SignedString(a.key)
	if err != nil {
		log.Error().Err(err).Send()
	}

	c.Header("X-Access-Token", access)
	c.Header("X-Refresh-Token", refresh)
	c.Status(http.StatusOK)
}
