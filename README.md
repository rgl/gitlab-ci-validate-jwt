This validates a GitLab CI JWT using the keys available at its jwks endpoint.

A GitLab CI JWT is a private string that can be used to authenticate a particular CI job in 3rd party services (like HashiCorp Vault).

Its available in a CI job as the `CI_JOB_JWT` environment variable.

It can also be available [as a custom ID token environment variable, with a custom `aud` claim](https://docs.gitlab.com/ee/ci/secrets/id_token_authentication.html).

A JWT is a structured string separated by dot characters; for example, a custom ID token JWT, something alike:

```
eyJhbGciOiJSUzI1NiIsImtpZCI6IjB1VnRkUEw4NDZSV19CY0kwWUxPN0JqdGlwenl1NVo4NHZ0Q0h5YkJSZFkiLCJ0eXAiOiJKV1QifQ.eyJuYW1lc3BhY2VfaWQiOiIxMiIsIm5hbWVzcGFjZV9wYXRoIjoiZXhhbXBsZSIsInByb2plY3RfaWQiOiIyIiwicHJvamVjdF9wYXRoIjoiZXhhbXBsZS9naXRsYWItY2ktdmFsaWRhdGUtand0IiwidXNlcl9pZCI6IjEiLCJ1c2VyX2xvZ2luIjoicm9vdCIsInVzZXJfZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsInBpcGVsaW5lX2lkIjoiMjIiLCJwaXBlbGluZV9zb3VyY2UiOiJwdXNoIiwiam9iX2lkIjoiMjkiLCJyZWYiOiJtYXN0ZXIiLCJyZWZfdHlwZSI6ImJyYW5jaCIsInJlZl9wcm90ZWN0ZWQiOiJ0cnVlIiwianRpIjoiNDZiMDRiOWItYThmMS00MzFmLTg4MjYtYzY5MTdjYTY3NzFmIiwiaXNzIjoiaHR0cHM6Ly9naXRsYWIuZXhhbXBsZS5jb20iLCJpYXQiOjE2NzcxODg1NDcsIm5iZiI6MTY3NzE4ODU0MiwiZXhwIjoxNjc3MTkyMTQ3LCJzdWIiOiJwcm9qZWN0X3BhdGg6ZXhhbXBsZS9naXRsYWItY2ktdmFsaWRhdGUtand0OnJlZl90eXBlOmJyYW5jaDpyZWY6bWFzdGVyIiwiYXVkIjoiaHR0cHM6Ly9leGFtcGxlLmNvbSJ9.M8RHBwvP5rVIJPEpztt4RKhgDLJJKiLP0O9XbBYvN5Fzt0zl-eNaoekTGbAAaZZK39jU3eZ-bkm8eKX8jdpYMz4di6WNGmVtJk_szDttXQpZ-HRYcNZz1EO83wuvPw0a9-ti9tfqgFy-xBc8jDVJhlS4bbRexrMbkwWiEWdHkApjopY9lnws61dl-OZ2-iPueIclzsQNEFSF1W9Us0lRi6OPp_dgSMPBVH6S3lqX-2p9V3FVlsW9aKqJr7_UuL1-F9dN-QCTztBra5A8GQeqIIVNalIhR8JDlon9vNyTFto34-elGmLfnk4jDkSbInQA7MKlN2to2vx18dkXWf3-DA
```

When split by dot and decoded it has a header, payload and signature.

In this case, the header is:

```json
{
    "alg": "RS256",
    "kid": "0uVtdPL846RW_BcI0YLO7Bjtipzyu5Z84vtCHybBRdY",
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
    "pipeline_source": "push",
    "job_id": "23",
    "ref": "master",
    "ref_type": "branch",
    "ref_protected": "true",
    "jti": "46b04b9b-a8f1-431f-8826-c6917ca6771f",
    "iss": "https://gitlab.example.com",
    "iat": 1677188547,
    "nbf": 1677188542,
    "exp": 1677192147,
    "sub": "project_path:example/gitlab-ci-validate-jwt:ref_type:branch:ref:master",
    "aud": "https://example.com"
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
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v15.9.0/app/models/ci/build.rb
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v15.9.0/lib/gitlab/ci/jwt.rb
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v15.9.0/app/controllers/jwt_controller.rb
* JWKS (JSON Web Key Set) endpoint (e.g. https://gitlab.example.com/-/jwks) at https://gitlab.com/gitlab-org/gitlab-foss/blob/v15.9.0/config/routes.rb#L234-236
* https://www.vaultproject.io/docs/auth/jwt
