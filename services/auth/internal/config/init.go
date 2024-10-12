package config

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/supertokens/supertokens-golang/recipe/jwt"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	pb "marketplace/shared/protobuf"
	"os"
)

func InitSupertokens(r *chi.Mux, profileClient pb.ProfileServiceClient) {
	var (
		stUri        = os.Getenv("ST_URI")
		uri          = os.Getenv("AUTH_URI")
		frontUri     = os.Getenv("FRONT_URI")
		googleClient = os.Getenv("ST_GOOGLE_CLIENT")
		googleSecret = os.Getenv("ST_GOOGLE_SECRET")
		githubClient = os.Getenv("ST_GITHUB_CLIENT")
		githubSecret = os.Getenv("ST_GITHUB_SECRET")
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
				Override: &tpmodels.OverrideStruct{
					Functions: func(originalImplementation tpmodels.RecipeInterface) tpmodels.RecipeInterface {
						originalSignInUp := *originalImplementation.SignInUp

						*originalImplementation.SignInUp = func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens map[string]interface{}, rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider, tenantId string, userContext *map[string]interface{}) (tpmodels.SignInUpResponse, error) {

							response, err := originalSignInUp(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, tenantId, userContext)
							if err != nil {
								return tpmodels.SignInUpResponse{}, err
							}

							if response.OK != nil {
								user := response.OK.User

								if response.OK.CreatedNewUser {
									userInfo := rawUserInfoFromProvider.FromUserInfoAPI
									firstName := userInfo["given_name"].(string)
									lastName := userInfo["family_name"].(string)
									imageUrl := userInfo["picture"].(string)
									req := &pb.CreateProfileRequest{FirstName: firstName, LastName: lastName, ImageUrl: imageUrl}

									if _, err := profileClient.CreateProfile(context.Background(), req); err != nil {
										if err := supertokens.DeleteUser(user.ID); err != nil {
											return tpmodels.SignInUpResponse{}, err
										}
										return tpmodels.SignInUpResponse{}, err
									}
								}

							}
							return response, nil
						}

						return originalImplementation
					},
				},
			}),
			session.Init(nil),
			jwt.Init(nil),
		},
	})

	if err != nil {
		logrus.Fatalf("Error initializing Supertokens: %v", err)
	}

	r.Use(supertokens.Middleware)

	logrus.Info("Successfully initialized Supertokens")
}
