package main

import (
	"flag"
	"fmt"
	"github.com/cepwn/hotel-reservation/api"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "The listen address for the API server")
	flag.Parse()
	fmt.Print(reflect.TypeOf(*listenAddr))

	app := fiber.New()
	apiV1 := app.Group("/api/v1")

	apiV1.Get("/user", api.HandleGetUsers)
	apiV1.Get("/user/:id", api.HandleGetUser)

	err := app.Listen(*listenAddr)
	if err != nil {
		panic(err)
	}
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "working just fine"})
}

func handleUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"user": "James Foo"})
}
