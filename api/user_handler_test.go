package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cepwn/hotel-reservation/db"
	"github.com/cepwn/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http/httptest"
	"testing"
)

const testDbUri = "mongodb://localhost:27017"

type testDb struct {
	db.UserStore
}

func (tDb *testDb) teardown(t *testing.T) {
	if err := tDb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDbUri))
	if err != nil {
		log.Fatal(err)
	}
	return &testDb{
		UserStore: db.NewMongoUserStore(client, db.TestDbName),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@foo.com",
		FirstName: "James",
		LastName:  "Foo",
		Password:  "aslkdfjowdkjsasde",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		t.Error(err)
	}

	if len(user.ID.String()) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected the EncryptedPassword not to be included in JSON response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but for %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but for %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but for %s", params.Email, user.Email)
	}

	fmt.Println(resp.Status)
}
