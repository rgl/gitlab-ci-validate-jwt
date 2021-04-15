This validates a GitLab CI JWT using the keys available at its jwks endpoint.

A GitLab CI JWT is a private string that can be used to authenticate a particular CI job in 3rd party services (like HashiCorp Vault).

Its available in a CI job as the `CI_JOB_JWT` environment variable.

A JWT is a structured string separated by dot characters; something alike:

```
eyJhbGciOiJSUzI1NiIsImtpZCI6Il8ybnI0NTI1UzVBclAwS05YQ0xySDZwMG4zYXVDX0RZcVBJdU8zN2gzTkEiLCJ0eXAiOiJKV1QifQ.eyJuYW1lc3BhY2VfaWQiOiIxMCIsIm5hbWVzcGFjZV9wYXRoIjoiZXhhbXBsZSIsInByb2plY3RfaWQiOiI3IiwicHJvamVjdF9wYXRoIjoiZXhhbXBsZS9naXRsYWItY2ktdmFsaWRhdGUtand0IiwidXNlcl9pZCI6IjEiLCJ1c2VyX2xvZ2luIjoicm9vdCIsInVzZXJfZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsInBpcGVsaW5lX2lkIjoiMTIiLCJqb2JfaWQiOiIyMyIsInJlZiI6Im1hc3RlciIsInJlZl90eXBlIjoiYnJhbmNoIiwicmVmX3Byb3RlY3RlZCI6InRydWUiLCJqdGkiOiJlY2I3YjJhOS02ZTljLTQ4NmUtYmYxNC1mNjIyOTgyOTMwODAiLCJpc3MiOiJnaXRsYWIuZXhhbXBsZS5jb20iLCJpYXQiOjE2MDA1OTExMjgsIm5iZiI6MTYwMDU5MTEyMywiZXhwIjoxNjAwNTk0NzI4LCJzdWIiOiJqb2JfMjMifQ.O_5PjdarFNJQ1u8Xh17BoWdsrxHtmeKu8_GJHJVuFRG3PE66hDTC0cOrqCP4iGp5InygIp26DE-C-fJ1QzgAiCkROQY83vLCq3_aTDVozCpuKdvifg7rxM5kd9ZmccmLnRrSnMPFF3LZPxvwn8A50ajJJOEbdD1Cud_lJd5ViVYZRPaATy44gPTFC72yqBIFwsrl5cB5Tlir_iMQyY4iMNYj-OWHG--hMVovUVVr9lFmhU8CmcaWjEd7C9gngp7hQ-BqMTWqhnCUUcipy7hNeHEACTrYjARuJEKAUMQf_23p1WO_ELHBNGrKSrKDFWtY_VOuGi7nmNVXU-Af0HCPzeYcoDwX1ex6E8ucrH5cgwj0exOIknBrcROWrxd6OFGQLo7V0hwRJ5P6auZJr5lG_hc0n2Ijc-sr266LRBzgwrqcVD9pcgfr6hW1wuyt9fyuNDvnXSkNQFT4v_CjhByUHm13CNRm7WW2urVUSL_suKR5yjV1k1AAzHo3-x1SeH4e9J8RkWiAtRGkU3imPtaADR3FpHCSzkncp-DC4iRTtGIKVLLuaLNZqKQWtfbTT8bfP0PxV109sb404t7U_gXZ5cqgi8Jam0FoYUyO_qEuBwwQdyHsj1YvYFCBLIFz3Zcu7gfUgEjGHCcFyrr9SArlj5YUWMmnbns77B0mwvl0Y4M
```

When split by dot and decoded it has a header, payload and signature.

In this case, the header is:

```json
{
	"alg": "RS256",
	"kid": "_2nr4525S5ArP0KNXCLrH6p0n3auC_DYqPIuO37h3NA",
	"typ": "JWT"
}
```

The payload is:

```json
{
    "namespace_id": "10",
    "namespace_path": "example",
    "project_id": "7",
    "project_path": "example/gitlab-ci-validate-jwt",
    "user_id": "1",
    "user_login": "root",
    "user_email": "admin@example.com",
    "pipeline_id": "12",
    "job_id": "23",
    "ref": "master",
    "ref_type": "branch",
    "ref_protected": "true",
    "jti": "ecb7b2a9-6e9c-486e-bf14-f62298293080",
    "iss": "gitlab.example.com",
    "iat": 1600591128,
    "nbf": 1600591123,
    "exp": 1600594728,
    "sub": "job_23"
}
```

And the signature is the value from the 3rd part of the JWT string.

Before a JWT can be used it must be validated. In this particular example the JWT can be validated with:

```go
RSASHA256(
    base64UrlEncode(header) + "." + base64UrlEncode(payload),
    gitLabJwtKeySet.getPublicKey(header.kid))
```

The above public key should be retrieved from the GitLab jwks endpoint (e.g. https://gitlab.example.com/-/jwks).

To see how all of this can be done read the [main.go](main.go) file.

## Reference

* https://docs.gitlab.com/ce/ci/examples/authenticating-with-hashicorp-vault/
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v13.10.3/app/models/ci/build.rb
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v13.10.3/lib/gitlab/ci/jwt.rb
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v13.10.3/app/controllers/jwt_controller.rb
* JWKS (JSON Web Key Set) endpoint (e.g. https://gitlab.example.com/-/jwks) at https://gitlab.com/gitlab-org/gitlab-foss/blob/v13.10.3/config/routes.rb#L202-204
* https://www.vaultproject.io/docs/auth/jwt
