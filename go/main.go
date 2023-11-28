package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nakaaaa/google-web-login/go/google"
)

var googleConfig google.Config

func init() {
	googleConfig.AppName = os.Getenv("GOOGLE_APP_NAME")
	googleConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	googleConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/hello", hello)
	e.GET("/auth", auth)
	e.GET("/token", token)
	e.GET("/verify", verify)

	e.Logger.Fatal(e.Start(":8080"))
}

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello Google"})
}

func auth(c echo.Context) error {
	url, err := googleConfig.Auth(c.Request().Context())
	if err != nil {
		return err
	}
	fmt.Println(url)
	return c.JSON(http.StatusFound, map[string]string{"url": url})
}

func token(c echo.Context) error {
	code := c.QueryParam("code")
	t, err := googleConfig.Token(c.Request().Context(), code)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func verify(e echo.Context) error {
	token := e.QueryParam("token")
	t, err := googleConfig.VerifyToken(e.Request().Context(), token)
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, t)
}
