package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/BingyanStudio/oidc-cli/oidc"
)

var oidcCli = oidc.NewClient(
	&oidc.Config{
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
		RedirectURL:  "http://localhost:3000/callback", // Frontend Callback URI
	},
)

func UserTokenHandler(c echo.Context) error {
	code := c.QueryParam("code")
	tokens, err := oidcCli.RetrieveTokens(code)
	if err != nil {
		log.Printf("err: " + err.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	log.Println(tokens)

	// use the tokens & claims to further authorization

	return c.JSON(http.StatusOK, tokens)
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/user/token", UserTokenHandler)
	e.Start(":8000")
}
