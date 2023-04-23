package google

import (
	"context"
	"net/http"
	"os"
	"project_altabe4_1/lib/databases"
	"project_altabe4_1/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	userInfo "google.golang.org/api/oauth2/v2"
)

type Google struct {
}

func GoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GCP_CLIENT_ID"),
		ClientSecret: os.Getenv("GCP_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_AUTH_CALLBACK"),
		Scopes: []string{
			os.Getenv("GCP_SCOPE_EMAIL"),
			os.Getenv("GCP_SCOPE_PROFILE"),
		},
		Endpoint: google.Endpoint,
	}
}

func (g *Google) HandleGoogleLogin(c echo.Context) (err error) {
	url := GoogleOauthConfig().AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func getUserInfo(client *http.Client) (*userInfo.Userinfo, error) {
	service, err := userInfo.New(client)
	if err != nil {
		return nil, err
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (g *Google) HandleGoogleCallback(c echo.Context) (err error) {
	token := &oauth2.Token{}
	ctx := context.Background()
	code := c.QueryParam("code")
	token, err = GoogleOauthConfig().Exchange(ctx, code)
	if err != nil {
		return
	}

	client := GoogleOauthConfig().Client(context.Background(), token)
	userInfo, err := getUserInfo(client)
	if err != nil {
		return err
	}

	user := &models.Users{
		Email: userInfo.Email,
		Nama:  userInfo.Name,
		Token: token.AccessToken,
	}

	_, err = databases.CreateUser(user)
	if err != nil {
		return err
	}

	return
}

func NewGoogleAuth() *Google {
	return &Google{}
}
