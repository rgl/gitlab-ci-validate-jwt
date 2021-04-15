package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

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
	ciJobJWT := getEnv("CI_JOB_JWT")
	ciServerURL := getEnv("CI_SERVER_URL")
	jwksURL := fmt.Sprintf("%s/-/jwks", ciServerURL)
	boundIssuer := getEnv("CI_SERVER_HOST")

	// fetch the gitlab jwt key set.
	//
	// a key set is public object alike:
	//
	// 		{
	// 			"keys": [
	// 				{
	// 					"kty": "RSA",
	// 					"kid": "_2nr4525S5ArP0KNXCLrH6p0n3auC_DYqPIuO37h3NA",
	// 					"e": "AQAB",
	// 					"n": "rYyQl7zEPEQPpxbhZpVkxFD-rEZHtyXnsr203AbgY1ks-8E3FuYCmFIZ7GrExlQjbhQkvEF7jdBHgALkb2Gu_p1DiBuWlTMSCXIbXzUjlh6_ULM05MUCXBYmrSaxRsgM6JH5T9awEHp9C48Sow9ZD5f9DcSIWptcjC8sPyxalje3xcPCrY4wzijUHisDNjjVqC7xSdhksdw_EByDcRWZ35vycm4FQnMXya745kwOFRKlpOXYvQJqXB1edONWZ3LbfuweAauNZfInWUjQiXSAUWEV3OqdANY8MnGAYn9ajT9O8k86DX14Qk4ZBR5m4A0y_cd-AKSHeC41NU45Tg70H2zvdx-L9R0Nc4HPpbBk-ELiKKDMpYkxRhnN9jsrv2bF66obMhtPd_jmRevAPJlWgNN60lhFZEK_Zck_xx1UYPMEyL5lh3hg2-UOuhG02ryfHYk1wQMAfJTQXH6ZTo5Kf-6ky9XG3XNRG2-VztBD4Xuo-KDMxPApT2J2QQ60fHlfRH5AVLE94NIl1vCu_rm6Vui3Ag94uWsFj27PQmf3UK5SJ-GwE8IvGgZtxNloGHkSVm5UjHW0dl4kJlnFZ-itxdNkDey_cUeT10JoswYwwQNKLpjKlplCgMC5VvQ_ddtSob4FeKpO71EZZZpwqgVIacs8M-5IiJ8JXcS9VeKE0Wk",
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

	// parse and validate the job jwt against the gitlab jwt key set.
	//
	// a job jwt is a private string alike:
	//
	// 		eyJhbGciOiJSUzI1NiIsImtpZCI6Il8ybnI0NTI1UzVBclAwS05YQ0xySDZwMG4zYXVDX0RZcVBJdU8zN2gzTkEiLCJ0eXAiOiJKV1QifQ.eyJuYW1lc3BhY2VfaWQiOiIxMCIsIm5hbWVzcGFjZV9wYXRoIjoiZXhhbXBsZSIsInByb2plY3RfaWQiOiI3IiwicHJvamVjdF9wYXRoIjoiZXhhbXBsZS9naXRsYWItY2ktdmFsaWRhdGUtand0IiwidXNlcl9pZCI6IjEiLCJ1c2VyX2xvZ2luIjoicm9vdCIsInVzZXJfZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsInBpcGVsaW5lX2lkIjoiMTIiLCJqb2JfaWQiOiIyMyIsInJlZiI6Im1hc3RlciIsInJlZl90eXBlIjoiYnJhbmNoIiwicmVmX3Byb3RlY3RlZCI6InRydWUiLCJqdGkiOiJlY2I3YjJhOS02ZTljLTQ4NmUtYmYxNC1mNjIyOTgyOTMwODAiLCJpc3MiOiJnaXRsYWIuZXhhbXBsZS5jb20iLCJpYXQiOjE2MDA1OTExMjgsIm5iZiI6MTYwMDU5MTEyMywiZXhwIjoxNjAwNTk0NzI4LCJzdWIiOiJqb2JfMjMifQ.O_5PjdarFNJQ1u8Xh17BoWdsrxHtmeKu8_GJHJVuFRG3PE66hDTC0cOrqCP4iGp5InygIp26DE-C-fJ1QzgAiCkROQY83vLCq3_aTDVozCpuKdvifg7rxM5kd9ZmccmLnRrSnMPFF3LZPxvwn8A50ajJJOEbdD1Cud_lJd5ViVYZRPaATy44gPTFC72yqBIFwsrl5cB5Tlir_iMQyY4iMNYj-OWHG--hMVovUVVr9lFmhU8CmcaWjEd7C9gngp7hQ-BqMTWqhnCUUcipy7hNeHEACTrYjARuJEKAUMQf_23p1WO_ELHBNGrKSrKDFWtY_VOuGi7nmNVXU-Af0HCPzeYcoDwX1ex6E8ucrH5cgwj0exOIknBrcROWrxd6OFGQLo7V0hwRJ5P6auZJr5lG_hc0n2Ijc-sr266LRBzgwrqcVD9pcgfr6hW1wuyt9fyuNDvnXSkNQFT4v_CjhByUHm13CNRm7WW2urVUSL_suKR5yjV1k1AAzHo3-x1SeH4e9J8RkWiAtRGkU3imPtaADR3FpHCSzkncp-DC4iRTtGIKVLLuaLNZqKQWtfbTT8bfP0PxV109sb404t7U_gXZ5cqgi8Jam0FoYUyO_qEuBwwQdyHsj1YvYFCBLIFz3Zcu7gfUgEjGHCcFyrr9SArlj5YUWMmnbns77B0mwvl0Y4M
	//
	// and decoded as a private object is alike:
	//
	// 		header:
	//
	// 			{
	// 				"alg": "RS256",
	// 				"kid": "_2nr4525S5ArP0KNXCLrH6p0n3auC_DYqPIuO37h3NA",
	// 				"typ": "JWT"
	// 			}
	//
	//		payload:
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
	// 				"job_id": "23",
	// 				"ref": "master",
	// 				"ref_type": "branch",
	// 				"ref_protected": "true",
	// 				"jti": "ecb7b2a9-6e9c-486e-bf14-f62298293080",
	// 				"iss": "gitlab.example.com",
	// 				"iat": 1600591128,
	// 				"nbf": 1600591123,
	// 				"exp": 1600594728,
	// 				"sub": "job_23"
	// 			}
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
	claims := token.PrivateClaims()

	log.Printf("jwt is valid for project %s", claims["project_path"])

	// dump the jwt claims (sorted by claim name).
	keys := make([]string, 0, len(claims))
	for k := range claims {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := claims[k]
		log.Printf("jwt claim: %s=%v", k, v)
	}
}
