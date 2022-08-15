package web

import (
	"fmt"
	ewa "github.com/ewa-go/ewa"
	"github.com/ewa-go/ewa/example/fiber/src/storage"
)

type Logout struct{}

func (Logout) Get(route *ewa.Route) {
	route.Session(ewa.Off)
	route.Handler = func(c *ewa.Context) error {
		if c.Session != nil {
			sessionId := c.Session.Value
			fmt.Println(sessionId)
			storage.DeleteStorage(sessionId)
		}
		return c.SendStatus(200)
	}
}
