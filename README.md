This validates a GitLab CI ID Token JWT using the keys available at its jwks endpoint.

A GitLab CI ID Token JWT is a private string that can be used to authenticate a particular CI job in 3rd party services (like HashiCorp Vault).

Its available in a CI job as a [custom environment variable defined in the job `id_tokens` property](https://docs.gitlab.com/ee/ci/secrets/id_token_authentication.html), as, e.g.:

```yaml
example_job:
  id_tokens:
    EXAMPLE_ID_TOKEN:
      aud: https://example.com
  script:
    - echo $EXAMPLE_ID_TOKEN
```

A JWT is a structured string separated by dot characters; for example, a custom ID token JWT, something alike:

```
eyJraWQiOiJwTG9jeVlFWHBqX2FrQzNVcnRQNkNfMUpGMEpvU1Q3VFZwazdwQXZqdWJ3IiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYifQ.eyJuYW1lc3BhY2VfaWQiOiIxMCIsIm5hbWVzcGFjZV9wYXRoIjoiZXhhbXBsZSIsInByb2plY3RfaWQiOiIxIiwicHJvamVjdF9wYXRoIjoiZXhhbXBsZS9naXRsYWItY2ktdmFsaWRhdGUtand0IiwidXNlcl9pZCI6IjEiLCJ1c2VyX2xvZ2luIjoicm9vdCIsInVzZXJfZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsInBpcGVsaW5lX2lkIjoiMTgiLCJwaXBlbGluZV9zb3VyY2UiOiJwdXNoIiwiam9iX2lkIjoiMjkiLCJyZWYiOiJtYXN0ZXIiLCJyZWZfdHlwZSI6ImJyYW5jaCIsInJlZl9wYXRoIjoicmVmcy9oZWFkcy9tYXN0ZXIiLCJyZWZfcHJvdGVjdGVkIjoiZmFsc2UiLCJydW5uZXJfaWQiOjIsInJ1bm5lcl9lbnZpcm9ubWVudCI6InNlbGYtaG9zdGVkIiwic2hhIjoiZTNlMTM5NjBiNmMwMGNiMGIxZjI1NmI0OGQyOGE4OTY4MGQ4YjY2MCIsImp0aSI6ImMxMjYzYTJlLWIwMWItNDdmMy05MDdhLTE5ZTM5YmFkY2RmNiIsImlzcyI6Imh0dHBzOi8vZ2l0bGFiLmV4YW1wbGUuY29tIiwiaWF0IjoxNjg1MTI2Mzc0LCJuYmYiOjE2ODUxMjYzNjksImV4cCI6MTY4NTEyOTk3NCwic3ViIjoicHJvamVjdF9wYXRoOmV4YW1wbGUvZ2l0bGFiLWNpLXZhbGlkYXRlLWp3dDpyZWZfdHlwZTpicmFuY2g6cmVmOm1hc3RlciIsImF1ZCI6Imh0dHBzOi8vZXhhbXBsZS5jb20ifQ.MIJYPNidvTRLtl-jQDkgfJjeIJCr6gQHNcAFR0AA6ACBIbVRcZQ8xQIRT6JtDKvfKSZej4wx5PqDG73x80swgS7raAZpG4LdTWAhfkYtsi88TH050Zz7Ku3qdL-0KYp5ykdQoLPTm-JvkzTKYYmOs9VkUX-rcmpb-9Bqhp_o10cGYmflnq1wzxW7OXYrqmUA1lhZe8jFuhmULgoNwFWCa4Q7jEJmtqy6U5uo2PXJr5LeZB3umnk9062_KOft63JNkDvhpc7ZLXsk0gxaJVx4f5s04EC2T7yLZi_SP2n-Fle0XM1S1DIbO6d-agvn6KLyE95eLBjAMm-GR5zMrgWmhA
```

When split by dot and decoded it has a header, payload and signature.

In this case, the header is:

```json
{
    "typ": "JWT",
    "alg": "RS256",
    "kid": "pLocyYEXpj_akC3UrtP6C_1JF0JoST7TVpk7pAvjubw"
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
    "ref_path": "refs/heads/master",
    "ref_protected": "true",
    "runner_id": 2,
    "runner_environment": "self-hosted",
    "sha": "e3e13960b6c00cb0b1f256b48d28a89680d8b660",
    "jti": "c1263a2e-b01b-47f3-907a-19e39badcdf6",
    "iss": "https://gitlab.example.com",
    "iat": 1685126374,
    "nbf": 1685126369,
    "exp": 1685129974,
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
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v16.0.1/app/models/ci/build.rb
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v16.0.1/lib/gitlab/ci/jwt.rb
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v16.0.1/app/controllers/jwt_controller.rb
* JWKS (JSON Web Key Set) endpoint (e.g. https://gitlab.example.com/-/jwks) at https://gitlab.com/gitlab-org/gitlab-foss/blob/v16.0.1/config/routes.rb#L213-215
