package config

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"log"
	"os"
)

func Init(r *chi.Mux) {
	var (
		stUri        = os.Getenv("AUTH_ST_URI")
		uri          = os.Getenv("AUTH_URI")
		frontUri     = os.Getenv("FRONT_URI")
		googleClient = os.Getenv("AUTH_ST_GOOGLE_CLIENT")
		googleSecret = os.Getenv("AUTH_ST_GOOGLE_SECRET")
		githubClient = os.Getenv("AUTH_ST_GITHUB_CLIENT")
		githubSecret = os.Getenv("AUTH_ST_GITHUB_SECRET")
	)

	apiBasePath := "/auth"
	websiteBasePath := "/auth"

	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: stUri,
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "Marketplace",
			APIDomain:       uri,
			WebsiteDomain:   frontUri,
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "google",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     googleClient,
										ClientSecret: googleSecret,
									},
								},
							},
						},
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "github",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     githubClient,
										ClientSecret: githubSecret,
									},
								},
							},
						},
					},
				},
			}),
		},
	})

	if err != nil {
		log.Fatalf("Error initializing Supertokens: %v", err)
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{frontUri},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: append([]string{"Content-Type"},
			supertokens.GetAllCORSHeaders()...),
		AllowCredentials: true,
	}))

	r.Use(supertokens.Middleware)

	log.Println("Successfully initialized Supertokens")
}
