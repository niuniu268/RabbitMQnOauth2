package main

func main() {

}

//
//import (
//	"context"
//	"fmt"
//	"golang.org/x/oauth2"
//	"golang.org/x/oauth2/google"
//	"net/http"
//	"time"
//
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt/v4"
//)
//
//const (
//	clientID     = "your-client-id"
//	clientSecret = "your-client-secret"
//	redirectURL  = "http://localhost:8080/auth/callback" // Update with your redirect URL
//	jwtSecret    = "your-jwt-secret"                     // Replace with a secure secret for JWT signing
//)
//
//var (
//	oauthConfig *oauth2.Config
//	state       string
//)
//
//func init() {
//	oauthConfig = &oauth2.Config{
//		ClientID:     clientID,
//		ClientSecret: clientSecret,
//		RedirectURL:  redirectURL,
//		Scopes:       []string{"openid", "profile", "email"}, // Adjust scopes as needed
//		Endpoint:     google.Endpoint,                        // Replace with the OAuth2 provider's endpoint if not Google
//	}
//
//	state = "your-random-state" // Replace with a random string for better security
//}
//
//func main() {
//	r := gin.Default()
//
//	// Route to start the OAuth2 flow
//	r.GET("/login", func(c *gin.Context) {
//		url := oauthConfig.AuthCodeURL(state)
//		c.Redirect(http.StatusTemporaryRedirect, url)
//	})
//
//	// Callback route after successful OAuth2 authentication
//	r.GET("/auth/callback", func(c *gin.Context) {
//		code := c.Query("code")
//		if code == "" {
//			c.String(http.StatusBadRequest, "Code not provided")
//			return
//		}
//
//		// Verify state to prevent CSRF attacks
//		if st := c.Query("state"); st != state {
//			c.String(http.StatusBadRequest, "Invalid state")
//			return
//		}
//
//		// Exchange the authorization code for a token
//		token, err := oauthConfig.Exchange(context.Background(), code)
//		if err != nil {
//			c.String(http.StatusInternalServerError, fmt.Sprintf("Error exchanging code for token: %v", err))
//			return
//		}
//
//		// Create a JWT with user information
//		jwtToken, err := createJWT(token)
//		if err != nil {
//			c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating JWT: %v", err))
//			return
//		}
//
//		// Display the JWT
//		c.String(http.StatusOK, jwtToken)
//	})
//
//	// Run the web server on port 8080
//	err := r.Run(":8080")
//	if err != nil {
//		return
//	}
//}
//
//func createJWT(token *oauth2.Token) (string, error) {
//	claims := jwt.MapClaims{
//		"sub":           "user-subject",
//		"iss":           "your-issuer",                        // Replace with your JWT issuer
//		"aud":           "your-audience",                      // Replace with your JWT audience
//		"exp":           time.Now().Add(time.Hour * 1).Unix(), // Token expiration time
//		"iat":           time.Now().Unix(),                    // Token issuance time
//		"oauth2_token":  token.AccessToken,
//		"refresh_token": token.RefreshToken,
//	}
//
//	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	signedToken, err := jwtToken.SignedString([]byte(jwtSecret))
//	if err != nil {
//		return "", err
//	}
//
//	return signedToken, nil
//}
