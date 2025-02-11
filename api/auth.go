package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// requireAuthentication checks incoming requests for tokens presented using the Authorization header
func (a *API) requireAuthentication(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	token, err := a.extractBearerToken(w, r)
	if err != nil {
		a.clearCookieToken(r.Context(), w)
		return nil, err
	}

	return a.parseJWTClaims(token, r, w)
}

type adminCheckParams struct {
	Aud string `json:"aud"`
}

func (a *API) requireAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	// Find the administrative user
	adminUser, err := getUserFromClaims(ctx, a.db)
	if err != nil {
		fmt.Println("adminUser")
		return nil, unauthorizedError("Invalid admin user").WithInternalError(err)
	}

	aud := a.requestAud(ctx, r)
	if r.Body != nil && r.Body != http.NoBody {
		c, err := addGetBody(w, r)
		if err != nil {
			return nil, internalServerError("Error getting body").WithInternalError(err)
		}
		ctx = c

		params := adminCheckParams{}
		bod, err := r.GetBody()
		if err != nil {
			return nil, internalServerError("Error getting body").WithInternalError(err)
		}
		err = json.NewDecoder(bod).Decode(&params)
		if err != nil {
			return nil, badRequestError("Could not decode admin user params: %v", err)
		}
		if params.Aud != "" {
			aud = params.Aud
		}
	}

	// Make sure user is admin
	if !a.isAdmin(ctx, adminUser, aud) {
		return nil, unauthorizedError("User not allowed")
	}
	return withAdminUser(ctx, adminUser), nil
}

func (a *API) extractBearerToken(w http.ResponseWriter, r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", unauthorizedError("This endpoint requires a Bearer token")
	}

	matches := bearerRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", unauthorizedError("This endpoint requires a Bearer token")
	}

	return matches[1], nil
}

func (a *API) parseJWTClaims(bearer string, r *http.Request, w http.ResponseWriter) (context.Context, error) {
	ctx := r.Context()
	config := a.getConfig(ctx)

	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}

	fmt.Println("")
	fmt.Println("try parsing claims....")
	fmt.Println("")

	token, err := p.ParseWithClaims(bearer, &GoTrueClaims{}, func(token *jwt.Token) (interface{}, error) {

		return []byte(config.JWT.Secret), nil
	})

	fmt.Println("")
	fmt.Println("tokenraw", token.Raw)
	fmt.Println("token claims: ", token.Claims)
	fmt.Println("token header: ", token.Header)
	fmt.Println("token valid: ", token.Valid)
	fmt.Println("token signedstring: ", token.SignedString)
	fmt.Println("token signingstring: ", token.SigningString)
	fmt.Println("token signature", token.Signature)
	fmt.Println("jwtsecret", config.JWT.Secret)
	fmt.Println("")

	if err != nil {
		fmt.Println(err)
		//a.clearCookieToken(ctx, w)
		return nil, unauthorizedError("Invalid token: %v", err)
	}

	return withToken(ctx, token), nil
}
