package sample

import (
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func Main() {
	e := echo.New()
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}
	FileManagerInit(os.Getenv("PERSISTENT_FILE"))
	if err := buildOffset(); err != nil {
		log.Fatal("Unable to bootstrap server: " + err.Error())
		return
	}
	loadRestRoutes(e)
	e.Logger.Fatal(e.Start(":"+os.Getenv("APP_PORT")))
}
