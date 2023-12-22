package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

func getEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("the %s environment variable must be set", name)
	}
	return v
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("you must pass the gitlab ci job id token environment variable name as the single command line argument")
	}
	ciJobJWT := getEnv(os.Args[1])
	ciServerURL := getEnv("CI_SERVER_URL")
	jwksURL := fmt.Sprintf("%s/oauth/discovery/keys", ciServerURL)
	boundIssuer := getEnv("CI_SERVER_HOST")

	// fetch the gitlab jwt key set.
	//
	// a key set is public object alike:
	//
	// 		{
	// 			"keys": [
	// 				{
	// 					"kty": "RSA",
	// 					"kid": "0uVtdPL846RW_BcI0YLO7Bjtipzyu5Z84vtCHybBRdY",
	// 					"e": "AQAB",
	// 					"n": "pydg59G3UOL2XX84st0hyDV7jPPG6zpikbFCP0Vwkj635LCVYx9tRTHE6BX0mu4QfPzDoj3qVp_2pjqUMuY3oaL0yIX1mI8stF6vljYU9kK1XjvM3Bxto-HqJx0RDFmg5bknLwwb6hhVz0Kh0-Hg1QLswv_Kop4a5jqNANnTQ1BhzNKIBztrYUHccfjpoWlD8P8Uu_7LU-Ka1OK5G5_wAA3vTgaTa-8R4aUhLISlinvAxPLfPoJ3pzYeplr_tPZwQNypN9IQQ0sJwRhYo2fXSldsPzsO1EmQXl7KihO9nc6Tq22ulb_YTtiUuC8BbiFKXyaVOTUvPiV4xdBvBNBHQQ",
	// 					"use": "sig",
	// 					"alg": "RS256"
	// 				}
	// 			]
	// 		}
	log.Printf("Getting the GitLab JWT public key set from the jwks endpoint at %s...", jwksURL)
	keySet, err := jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		log.Fatalf("failed to parse JWK from %s: %v", jwksURL, err)
	}
	if keySet.Len() < 1 {
		log.Fatalf("%s did not return any key", jwksURL)
	}

	// parse and validate the job id token jwt against the gitlab jwt key set.
	//
	// a job jwt is a private string alike:
	//
	// 		eyJhbGciOiJSUzI1NiIsImtpZCI6IjB1VnRkUEw4NDZSV19CY0kwWUxPN0JqdGlwenl1NVo4NHZ0Q0h5YkJSZFkiLCJ0eXAiOiJKV1QifQ.eyJuYW1lc3BhY2VfaWQiOiIxMiIsIm5hbWVzcGFjZV9wYXRoIjoiZXhhbXBsZSIsInByb2plY3RfaWQiOiIyIiwicHJvamVjdF9wYXRoIjoiZXhhbXBsZS9naXRsYWItY2ktdmFsaWRhdGUtand0IiwidXNlcl9pZCI6IjEiLCJ1c2VyX2xvZ2luIjoicm9vdCIsInVzZXJfZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsInBpcGVsaW5lX2lkIjoiMjIiLCJwaXBlbGluZV9zb3VyY2UiOiJwdXNoIiwiam9iX2lkIjoiMjkiLCJyZWYiOiJtYXN0ZXIiLCJyZWZfdHlwZSI6ImJyYW5jaCIsInJlZl9wcm90ZWN0ZWQiOiJ0cnVlIiwianRpIjoiNDZiMDRiOWItYThmMS00MzFmLTg4MjYtYzY5MTdjYTY3NzFmIiwiaXNzIjoiaHR0cHM6Ly9naXRsYWIuZXhhbXBsZS5jb20iLCJpYXQiOjE2NzcxODg1NDcsIm5iZiI6MTY3NzE4ODU0MiwiZXhwIjoxNjc3MTkyMTQ3LCJzdWIiOiJwcm9qZWN0X3BhdGg6ZXhhbXBsZS9naXRsYWItY2ktdmFsaWRhdGUtand0OnJlZl90eXBlOmJyYW5jaDpyZWY6bWFzdGVyIiwiYXVkIjoiaHR0cHM6Ly9leGFtcGxlLmNvbSJ9.M8RHBwvP5rVIJPEpztt4RKhgDLJJKiLP0O9XbBYvN5Fzt0zl-eNaoekTGbAAaZZK39jU3eZ-bkm8eKX8jdpYMz4di6WNGmVtJk_szDttXQpZ-HRYcNZz1EO83wuvPw0a9-ti9tfqgFy-xBc8jDVJhlS4bbRexrMbkwWiEWdHkApjopY9lnws61dl-OZ2-iPueIclzsQNEFSF1W9Us0lRi6OPp_dgSMPBVH6S3lqX-2p9V3FVlsW9aKqJr7_UuL1-F9dN-QCTztBra5A8GQeqIIVNalIhR8JDlon9vNyTFto34-elGmLfnk4jDkSbInQA7MKlN2to2vx18dkXWf3-DA
	//
	// and decoded as a private object is alike:
	//
	// 		header:
	//
	// 			{
	// 				"alg": "RS256",
	// 				"kid": "0uVtdPL846RW_BcI0YLO7Bjtipzyu5Z84vtCHybBRdY",
	// 				"typ": "JWT"
	// 			}
	//
	// 		payload:
	//
	// 			{
	// 				"namespace_id": "10",
	// 				"namespace_path": "example",
	// 				"project_id": "7",
	// 				"project_path": "example/gitlab-ci-validate-jwt",
	// 				"user_id": "1",
	// 				"user_login": "root",
	// 				"user_email": "admin@example.com",
	// 				"pipeline_id": "12",
	// 				"pipeline_source": "push",
	// 				"job_id": "23",
	// 				"ref": "master",
	// 				"ref_type": "branch",
	// 				"ref_protected": "true",
	// 				"jti": "46b04b9b-a8f1-431f-8826-c6917ca6771f",
	// 				"iss": "https://gitlab.example.com",
	// 				"iat": 1677188547,
	// 				"nbf": 1677188542,
	// 				"exp": 1677192147,
	// 				"sub": "project_path:example/gitlab-ci-validate-jwt:ref_type:branch:ref:master",
	// 				"aud": "https://example.com"
	//  		}
	//
	//		signature:
	//
	//			the value is the 3rd part of the jwt.
	//
	//			in this particular example the jwt can be validated with:
	//
	//				RSASHA256(
	//   				base64UrlEncode(header) + "." + base64UrlEncode(payload),
	//					gitLabJwtKeySet.getKey(header.kid))
	log.Println("Validating GitLab CI job JWT...")
	token, err := jwt.ParseString(ciJobJWT, jwt.WithIssuer(boundIssuer), jwt.WithKeySet(keySet))
	if err != nil {
		log.Fatalf("failed to validate the jwt: %v", err)
	}
	privateClaims := token.PrivateClaims()

	log.Printf("jwt is valid for project %s", privateClaims["project_path"])

	// dump the jwt claims (sorted by claim name).
	claims := make([]string, 0, len(privateClaims)+7)
	claims = append(claims,
		fmt.Sprintf("jti=%s", token.JwtID()),
		fmt.Sprintf("iss=%s", token.Issuer()),
		fmt.Sprintf("iat=%s", token.IssuedAt().Format("2006-01-02T15:04:05-0700")),
		fmt.Sprintf("nbf=%s", token.NotBefore().Format("2006-01-02T15:04:05-0700")),
		fmt.Sprintf("exp=%s", token.Expiration().Format("2006-01-02T15:04:05-0700")),
		fmt.Sprintf("sub=%s", token.Subject()),
		fmt.Sprintf("aud=%s", strings.Join(token.Audience(), ",")))
	for k := range privateClaims {
		claims = append(claims, fmt.Sprintf("%s=%v", k, privateClaims[k]))
	}
	sort.Strings(claims)
	for _, claim := range claims {
		log.Printf("jwt claim: %s", claim)
	}
}
