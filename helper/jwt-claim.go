package helper

import (
	"container/heap"
	"errors"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JwtKey = []byte("SECRET_KEY")

type JWTClaim struct {
	UID      int    `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type BlacklistedToken struct {
	Token          string
	ExpirationTime time.Time
}

type BlacklistedTokenHeap []BlacklistedToken

func (h BlacklistedTokenHeap) Len() int { return len(h) }
func (h BlacklistedTokenHeap) Less(i, j int) bool {
	return h[i].ExpirationTime.Before(h[j].ExpirationTime)
}
func (h BlacklistedTokenHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *BlacklistedTokenHeap) Push(x interface{}) {
	*h = append(*h, x.(BlacklistedToken))
}

func (h *BlacklistedTokenHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

var BlacklistedTokens = make(map[string]*time.Time)
var TokenHeap BlacklistedTokenHeap

func GenerateJWT(uid int, username string, role string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		UID:      uid,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func ValidateToken(c *gin.Context) (tokenString string, claim *JWTClaim, err error) {
	signedToken := ExtractToken(c)
	if signedToken == "" {
		c.JSON(401, gin.H{"error": "request does not contain an access token"})
		c.Abort()
		return
	}

	if _, ok := BlacklistedTokens[signedToken]; ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has been revoked"})
		return
	}

	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		return
	}

	claim, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claim.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token has expired")
		return
	}
	return signedToken, claim, nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func CleanTokenHeap() {
	for {
		now := time.Now()

		for TokenHeap.Len() > 0 {
			token := TokenHeap[0]
			if now.Before(token.ExpirationTime) {
				break
			}
			heap.Pop(&TokenHeap)
			delete(BlacklistedTokens, token.Token)
		}

		time.Sleep(1 * time.Minute)
	}
}
