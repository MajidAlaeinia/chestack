package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/MajidAlaeinia/chestack/app/http/resources"
	"github.com/MajidAlaeinia/chestack/utils"
	gintest "github.com/Valiben/gin_unit_test"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	UsersEndpoint = "/users"
)

type SuccessfulResponse struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Mobile          string    `json:"mobile"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
	Password        string    `json:"password"`
	RememberToken   string    `json:"remember_token"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func TestUserController_Show_SuccessfulResponse(t *testing.T) {
	db := utils.RunDatabase()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Println(fmt.Sprintf("Error: %v", err))
		}
	}(db)

	id := "1"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{
				Key:   "id",
				Value: id,
			},
		}
		new(UserController).Show(c)
	}))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error: %v", err.Error())
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var sr SuccessfulResponse
	err = json.Unmarshal(bodyBytes, &sr)
	if err != nil {
		t.Fatalf("Error: %v", err.Error())
	}
	assert.Equal(t, "alaeinia.majid@gmail.com", sr.Email)
}

func TestUserController_Create_SuccessfulResponse(t *testing.T) {
	db := utils.RunTxDb()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			t.Fatalf("Errro: %v", err)
		}
	}(db)

	router := gin.Default()
	router.POST(UsersEndpoint, new(UserController).Create)
	gintest.SetRouter(router)

	resp := resources.User{}

	name := "Name_4"
	email := "email_4@gmail.com"
	mobile := "09381201110"
	reqBody := map[string]interface{}{"name": name, "email": email, "mobile": mobile}
	err := gintest.TestHandlerUnMarshalResp("POST", UsersEndpoint, "json", reqBody, &resp)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	assert.Equal(t, name, resp.Name)
	assert.Equal(t, email, resp.Email)
	assert.Equal(t, mobile, resp.Mobile)
}
