package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/auth0/go-auth0/authentication"
	"github.com/auth0/go-auth0/authentication/database"
	"github.com/auth0/go-auth0/authentication/oauth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-memdb"
	"github.com/joho/godotenv"
)

// DRY on 'profile'
const (
	PROFILE = "profile"
	ID      = "id"
)

// model struct
type profile struct {
	Age                      int16  `json:"age,omitempty" validate:"omitempty,gte=0,lte=130"`
	Name                     string `json:"name,omitempty" validate:"omitempty,min=3,max=32"`
	Email                    string `json:"email,omitempty" validate:"omitempty,email,min=3,max=100"`
	FavouriteColor           string `json:"favcolor,omitempty" validate:"omitempty,min=3,max=50"`
	FavouriteOperatingSystem string `json:"favos,omitempty" validate:"omitempty,min=3,max=20"`
	Password                 string `json:"password,omitempty" validate:"omitempty,min=10,max=32"`
}

type request profile

// exportable
type App struct {
	auth      *authentication.Authentication
	db        *memdb.MemDB
	router    *gin.Engine
	validate  *validator.Validate
	demoEmail string
	demoPass  string
}

func (app *App) LogoutHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		scheme := "http"
		if ctx.Request.TLS != nil {
			scheme = "https"
		}

		returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		parameters := url.Values{}
		parameters.Add("returnTo", returnTo.String())
		parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
		logoutUrl.RawQuery = parameters.Encode()

		ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
	}
}

func (app *App) LoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// validate/sanitize request
		req := profile{}
		req.Email = ctx.PostForm("email")
		req.Password = ctx.PostForm("password")
		if err := app.validate.Struct(req); err != nil {
			log.Printf("Failed to validate request: %+v", err)
			ctx.String(http.StatusInternalServerError, "Failed to validate request")
			return
		}

		session := sessions.Default(ctx)

		tokens, err := app.auth.OAuth.LoginWithPassword(context.Background(), oauth.LoginWithPasswordRequest{
			Username: req.Email,
			Password: req.Password,
		}, oauth.IDTokenValidationOptions{})
		if err != nil {
			log.Printf("Failed to login to Auth0: %+v", err)
			ctx.String(http.StatusInternalServerError, "Failed to login")
			return
		}

		session.Set("access_token", tokens.AccessToken)

		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect to logged in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/profile")
	}
}

func (app *App) UpdateProfileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req profile
		if err := ctx.Bind(&req); err != nil {
			log.Printf("Failed to bind request: %+v", err)
			ctx.String(http.StatusInternalServerError, "Failed to bind request")
			return
		}
		if err := app.validate.Struct(req); err != nil {
			log.Printf("Failed to validate request: %+v", err)
			ctx.String(http.StatusInternalServerError, "Failed to validate request")
			return
		}

		session := sessions.Default(ctx)

		txn := app.db.Txn(true)
		defer txn.Abort()

		if err := txn.Insert(PROFILE, req); err != nil {
			log.Printf("Failed to update profile: %+v", err)
			ctx.String(http.StatusInternalServerError, "Failed to update profile")
			return
		}
		txn.Commit()

		session.Set(PROFILE, req)

		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.String(http.StatusOK, "Profile updated successfully")
	}
}

func (app *App) DeleteProfileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req request
		if err := ctx.BindJSON(&req); err != nil {
			log.Printf("Failed to bind json: %+v", err)
			ctx.String(http.StatusInternalServerError, "Failed to bind json")
			return
		}
		if err := app.validate.Struct(req); err != nil {
			log.Printf("Failed to validate request: %+v", err)
			ctx.String(http.StatusInternalServerError, "Failed to validate request")
			return
		}

		session := sessions.Default(ctx)

		txn := app.db.Txn(true)
		defer txn.Abort()

		if err := txn.Delete(PROFILE, req); err != nil {
			log.Printf("Failed to delete profile: %+v", err)
			ctx.String(http.StatusInternalServerError, "Failed to delete profile")
			return
		}
		txn.Commit()

		session.Delete(PROFILE)

		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.String(http.StatusOK, "Profile deleted successfully")
	}
}

// IsAuthenticated is a middleware that checks if
// the user has already been authenticated previously.
func IsAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if sessions.Default(ctx).Get(PROFILE) == nil {
			ctx.Redirect(http.StatusSeeOther, "/")
			return
		}
		ctx.Next()
	}
}

func NewDemoApp() *App {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load the env vars: %v", err)
	}

	app := &App{}

	domain := os.Getenv("AUTH0_DOMAIN")
	clientID := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	// assign validator
	app.validate = validator.New()

	// Initialize a new client using a domain, client ID and client secret.
	authAPI, err := authentication.New(
		context.TODO(), // Replace with a Context that better suits your usage
		domain,
		authentication.WithClientID(clientID),
		authentication.WithClientSecret(clientSecret), // Optional depending on the grants used
	)
	if err != nil {
		log.Fatalf("failed to initialize the auth0 authentication API client: %+v", err)
	}
	app.auth = authAPI

	// Now we can interact with the Auth0 Authentication API.
	// Sign up a user
	userData := database.SignupRequest{
		Connection: "Username-Password-Authentication",
		Password:   "faker.Password()",
		Email:      faker.Email(),
	}

	_, err = app.auth.Database.Signup(context.Background(), userData)
	if err != nil {
		log.Fatalf("failed to sign user up: %+v", err)
	}
	log.Printf("Created test user %v", userData.Email)

	// Create user profile from userData and persist
	userProfile := profile{
		Age:                      30,
		Name:                     faker.Name(),
		Email:                    userData.Email,
		FavouriteColor:           "black",
		FavouriteOperatingSystem: "linux",
	}

	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			PROFILE: &memdb.TableSchema{
				Name: PROFILE,
				Indexes: map[string]*memdb.IndexSchema{
					ID: &memdb.IndexSchema{
						Name:    ID,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"age": &memdb.IndexSchema{
						Name:    "age",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Age"},
					},
					"favcolor": &memdb.IndexSchema{
						Name:    "favcolor",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "FavouriteColor"},
					},
					"favos": &memdb.IndexSchema{
						Name:    "favos",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "FavouriteOperatingSystem"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		log.Fatalf("Failed to initiate profile database: %+v", err)
	}
	app.db = db

	// Create a write transaction
	txn := app.db.Txn(true)

	// Lookup by email
	err = txn.Insert(PROFILE, userProfile)
	if err != nil {
		log.Fatalf("failed to persist profile: %+v", err)
	}
	txn.Commit()

	app.router = gin.Default()
	app.demoEmail = userData.Email
	app.demoPass = userData.Password
	return app
}

func main() {
	app := NewDemoApp()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	app.router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	app.router.Use(gin.Recovery())

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(profile{})

	store := cookie.NewStore([]byte("secret"))
	app.router.Use(sessions.Sessions("auth-session", store))

	app.router.Static("/static", "./static")
	app.router.LoadHTMLGlob("template/*")

	app.router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})

	app.router.POST("/login", app.LoginHandler())

	// Profile Group
	profileGroup := app.router.Group("/profile")
	profileGroup.Use(IsAuthenticated())
	{
		profileGroup.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "profile.html", nil)
		})

		profileGroup.DELETE("/", app.DeleteProfileHandler())
		profileGroup.POST("/", app.UpdateProfileHandler())

		// logout
		profileGroup.POST("/logout", app.LogoutHandler())
	}

	log.Print("Server listening on http://localhost:3000/")
	log.Print("Please use the following email and password to login:")
	log.Print("=============================================")
	log.Printf("Email: %v", app.demoEmail)
	log.Printf("Password: %v", app.demoPass)
	log.Print("=============================================")

	if err := http.ListenAndServe("0.0.0.0:3000", app.router); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
