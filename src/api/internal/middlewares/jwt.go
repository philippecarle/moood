package middlewares

import (
	"net/http"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/philippecarle/moood/api/internal/database/collections"
	"github.com/philippecarle/moood/api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// IdentityKey is where the user is stored in the context
const IdentityKey = "id"

// JWTMiddleWareFactory is a struct with the user collection providing a way to build a GinJWTMiddleware
type JWTMiddleWareFactory struct {
	UsersCollection collections.UsersCollection
}

// NewJWTMiddleWare returns a GinJWTMiddleware
func (f *JWTMiddleWareFactory) NewJWTMiddleWare() *jwt.GinJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "authenticated zone",
		Key:             []byte(os.Getenv("JWT_PRIVATE_KEY")),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour * 24 * 30,
		IdentityKey:     IdentityKey,
		PayloadFunc:     f.buildClaims,
		IdentityHandler: f.identify,
		Authenticator:   f.authenticate,
		TokenLookup:     "header: Authorization",
		LoginResponse:   f.loginResponse,
	})

	if err != nil {
		panic(err)
	}

	return authMiddleware
}

func (f *JWTMiddleWareFactory) buildClaims(data interface{}) jwt.MapClaims {
	if u, ok := data.(models.User); ok {
		return jwt.MapClaims{
			IdentityKey: u.Username,
			"mercure": map[string][]string{
				"subscribe": {"/users/" + u.Username},
			},
		}
	}
	return jwt.MapClaims{}
}

func (f *JWTMiddleWareFactory) identify(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	user, _ := f.UsersCollection.FindOneByUserName(claims[IdentityKey].(string))
	return user
}

func (f *JWTMiddleWareFactory) authenticate(c *gin.Context) (interface{}, error) {
	var creds models.Credentials
	if err := c.ShouldBind(&creds); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	user, err := f.UsersCollection.FindOneByUserName(creds.Username)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.EncodedPassword), []byte(creds.Password)); err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return user, nil
}

func (f *JWTMiddleWareFactory) loginResponse(c *gin.Context, code int, token string, t time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
