package main

import (
	"fmt"
	ewa "github.com/ewa-go/ewa"
	f "github.com/ewa-go/ewa-fiber"
	"github.com/ewa-go/ewa/example/fiber-api/controllers"
	"github.com/ewa-go/ewa/example/fiber-api/controllers/api/storage"
	"github.com/ewa-go/ewa/example/fiber-api/models"
	"github.com/ewa-go/ewa/security"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//BasicAuth
	basicAuthHandler := func(user string, pass string) bool {
		if user == "user" && pass == "Qq123456" {
			return true
		}
		return false
	}

	contextHandler := func(handler ewa.Handler) interface{} {
		return func(ctx *fiber.Ctx) error {
			return handler(ewa.NewContext(&f.Context{Ctx: ctx}))
		}
	}

	// Fiber
	app := fiber.New()
	// Cors
	app.Use(cors.New())
	server := &f.Server{App: app}
	// Конфиг
	cfg := ewa.Config{
		Port: 8070,
		Secure: &ewa.Secure{
			Path: "./cert",
			Key:  "key.pem",
			Cert: "cert.pem",
		},
		Authorization: security.Authorization{
			Basic: &security.Basic{
				Handler: basicAuthHandler,
			},
		},
		ContextHandler: contextHandler,
	}

	info := ewa.Info{
		Description: "Description",
		Version:     "1.0.0",
		Title:       "FiberApi",
		Contact: &ewa.Contact{
			Email: "user@mail.ru",
		},
		License: &ewa.License{
			Name: "License",
		},
	}

	hostname := ewa.Suffix{
		Index:       3,
		Value:       "{hostname}",
		Description: "Set hostname device",
	}

	//Инициализируем сервер
	ws := ewa.New(server, cfg)
	ws.Register(new(storage.User)).SetSuffix(hostname).SetDescription("Users")
	ws.Register(new(controllers.Home)).SetPath("/")
	// Swagger
	ws.Register(new(controllers.Api)).NotShow()

	// Описываем swagger
	ws.Swagger.SetInfo("10.28.0.73", &info, nil).SetBasePath("/api")
	ws.Swagger.SetModels(ewa.Models{
		models.ModelUser:     models.User{},
		models.ModelResponse: models.Response{},
	})

	// Канал для получения ошибки, если таковая будет
	errChan := make(chan error, 2)
	go func() {
		errChan <- ws.Start()
	}()

	// Ждем сигнал от ОС
	go getSignal(errChan)

	fmt.Println("Старт приложения")
	fmt.Printf("Остановка приложения. %s", <-errChan)
}

func getSignal(errChan chan error) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	errChan <- fmt.Errorf("%s", <-c)
}
