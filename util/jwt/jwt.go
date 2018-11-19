package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"net/http"
	"time"
)

var secret = []byte("oneshot")

type Claims struct {
	ID string
	jwt.StandardClaims
}

func BaseAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := ClaimsFromCookie(c)
		if claims == nil {
			c.JSON(404, gin.H{"errorMsg": "noSession"})
			c.Abort()
		}
		//c.Set("claims", *claims)
		c.Next()
	}
}

func ClaimsFromCookie(c *gin.Context) *Claims {
	cookie, e := c.Request.Cookie("gymtracker")

	//If e != nil no cookie found
	if e != nil {
		return nil
	}

	//Parse the cookie into token
	token, e := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("gymtracker"), nil
	})

	//Extract claims
	claims, ok := token.Claims.(*Claims)

	//Check for claims expired
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		//Recreate claims
		newClaims := NewClaims(claims.ID)

		//Recreate token
		newToken, _ := NewToken(newClaims)

		//Recreate cookie
		newCookie := NewCookie(newToken)

		//Set new cookie
		http.SetCookie(c.Writer, newCookie)

		//Return new claims
		return &newClaims
	}

	if ok && token.Valid {
		return claims
	}

	//Return false as fallback
	return nil

}

func NewCookie(token string) *http.Cookie {
	//Create the cookie
	return &http.Cookie{
		Name:     "gymtracker",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}
}

func NewToken(claims Claims) (string, error) {
	//Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Create signed string and return
	return token.SignedString(secret)
}

func NewClaims(ID string) Claims {
	return Claims{
		ID,
		jwt.StandardClaims{
			//	ExpiresAt: time.Now().Add(time.Second * 30).Unix(),
			Issuer: "gymtracker",
		},
	}
}
