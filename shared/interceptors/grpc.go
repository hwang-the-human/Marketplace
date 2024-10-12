package interceptors

import (
	"context"
	"fmt"
	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	stJWT "github.com/supertokens/supertokens-golang/recipe/jwt"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
	"sync"
	"time"
)

var (
	jwksCache          *sessmodels.GetJWKSResult = nil
	mutex              sync.RWMutex
	JWKCacheMaxAgeInMs int64 = 60000
	coreUrl                  = os.Getenv("ST_URI") + "/.well-known/jwks.json"
)

var (
	jwtToken     string
	jwtExpiresAt time.Time
	mu           sync.Mutex
)

func getJWKSFromCacheIfPresent() *sessmodels.GetJWKSResult {
	mutex.RLock()
	defer mutex.RUnlock()
	if jwksCache != nil {
		// This means that we have valid JWKs for the given core path
		// We check if we need to refresh before returning
		currentTime := time.Now().UnixNano() / int64(time.Millisecond)

		// This means that the value in cache is not expired, in this case we return the cached value
		//
		// Note that this also means that the SDK will not try to query any other Core (if there are multiple)
		// if it has a valid cache entry from one of the core URLs. It will only attempt to fetch
		// from the cores again after the entry in the cache is expired
		if (currentTime - jwksCache.LastFetched) < JWKCacheMaxAgeInMs {
			return jwksCache
		}
	}

	return nil
}

func GetJWKS() (*keyfunc.JWKS, error) {
	resultFromCache := getJWKSFromCacheIfPresent()

	if resultFromCache != nil {
		return resultFromCache.JWKS, nil
	}

	mutex.Lock()
	defer mutex.Unlock()
	// RefreshUnknownKID - Fetch JWKS again if the kid in the header of the JWT does not match any in
	// the keyfunc library's cache
	jwks, jwksError := keyfunc.Get(coreUrl, keyfunc.Options{
		RefreshUnknownKID: true,
	})

	if jwksError == nil {
		jwksResult := sessmodels.GetJWKSResult{
			JWKS:        jwks,
			Error:       jwksError,
			LastFetched: time.Now().UnixNano() / int64(time.Millisecond),
		}

		// Dont add to cache if there is an error to keep the logic of checking cache simple
		//
		// This also has the added benefit where if initially the request failed because the core
		// was down and then it comes back up, the next time it will try to request that core again
		// after the cache has expired
		jwksCache = &jwksResult

		return jwksResult.JWKS, nil
	}

	// This means that fetching from the core failed
	return nil, jwksError
}

func JWTAuth(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata not provided")
	}

	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return nil, fmt.Errorf("authorization token not provided")
	}

	jwtString := authHeader[0][7:]
	jwks, err := GetJWKS()
	if err != nil {
		return nil, fmt.Errorf("could not get JWKS: %v", err)
	}

	parsedToken, parseError := jwt.Parse(jwtString, jwks.Keyfunc)
	if parseError != nil {
		return nil, fmt.Errorf("invalid token: %v", parseError)
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("token claims are not valid")
	}

	claimsMap := make(map[string]interface{})
	for key, value := range claims {
		claimsMap[key] = value
	}

	sourceClaim, ok := claimsMap["source"]
	if !ok || sourceClaim != "microservice" {
		return nil, fmt.Errorf("unauthorized: token not intended for microservice communication")
	}

	return handler(ctx, req)
}

func AttachJWT() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		jwtToken, err := getJWT()
		if err != nil {
			return err
		}
		md := metadata.Pairs("authorization", "Bearer "+jwtToken)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func getJWT() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	if jwtToken == "" || time.Now().After(jwtExpiresAt) {
		token, err := createJWT()
		if err != nil {
			return "", err
		}
		jwtToken = token
		jwtExpiresAt = time.Now().Add(1 * time.Hour)
	}

	return jwtToken, nil
}

func createJWT() (string, error) {
	jwtResponse, err := stJWT.CreateJWT(map[string]interface{}{
		"source": "microservice",
	}, nil, nil)
	if err != nil {
		return "", err
	}
	return jwtResponse.OK.Jwt, nil
}
