This validates a GitLab CI ID Token JWT using the keys available at its jwks endpoint.

A GitLab CI ID Token JWT is a secret string that can be used to authenticate a particular CI job in 3rd party services (like HashiCorp Vault).

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
eyJraWQiOiJjYVJnUUFlSGl1dEgzNXlVcXJMbHpST3RBcGZfYzhaWjZYSHN3RkJ5MERNIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYifQ.eyJuYW1lc3BhY2VfaWQiOiIxMCIsIm5hbWVzcGFjZV9wYXRoIjoiZXhhbXBsZSIsInByb2plY3RfaWQiOiIxIiwicHJvamVjdF9wYXRoIjoiZXhhbXBsZS9naXRsYWItY2ktdmFsaWRhdGUtand0IiwidXNlcl9pZCI6IjEiLCJ1c2VyX2xvZ2luIjoicm9vdCIsInVzZXJfZW1haWwiOiJnaXRsYWJfYWRtaW5fMDkzYjgzQGV4YW1wbGUuY29tIiwidXNlcl9hY2Nlc3NfbGV2ZWwiOiJvd25lciIsInBpcGVsaW5lX2lkIjoiMTgiLCJwaXBlbGluZV9zb3VyY2UiOiJwdXNoIiwiam9iX2lkIjoiMzQiLCJyZWYiOiJtYXN0ZXIiLCJyZWZfdHlwZSI6ImJyYW5jaCIsInJlZl9wYXRoIjoicmVmcy9oZWFkcy9tYXN0ZXIiLCJyZWZfcHJvdGVjdGVkIjoidHJ1ZSIsInJ1bm5lcl9pZCI6MiwicnVubmVyX2Vudmlyb25tZW50Ijoic2VsZi1ob3N0ZWQiLCJzaGEiOiI5NWQxOGQ2NmFmZDJjMDYwOWY2YzQxYmQ1MzdhODI3YmViNjk4ZTY0IiwicHJvamVjdF92aXNpYmlsaXR5IjoicHVibGljIiwiY2lfY29uZmlnX3JlZl91cmkiOiJnaXRsYWIuZXhhbXBsZS5jb20vZXhhbXBsZS9naXRsYWItY2ktdmFsaWRhdGUtand0Ly8uZ2l0bGFiLWNpLnltbEByZWZzL2hlYWRzL21hc3RlciIsImNpX2NvbmZpZ19zaGEiOiI5NWQxOGQ2NmFmZDJjMDYwOWY2YzQxYmQ1MzdhODI3YmViNjk4ZTY0IiwianRpIjoiNGJkODc2N2UtM2Q2Ni00OTU4LThiODMtNzA5N2RhZWJjMWE3IiwiaWF0IjoxNzQwODIxMTQyLCJuYmYiOjE3NDA4MjExMzcsImV4cCI6MTc0MDgyNDc0MiwiaXNzIjoiaHR0cHM6Ly9naXRsYWIuZXhhbXBsZS5jb20iLCJzdWIiOiJwcm9qZWN0X3BhdGg6ZXhhbXBsZS9naXRsYWItY2ktdmFsaWRhdGUtand0OnJlZl90eXBlOmJyYW5jaDpyZWY6bWFzdGVyIiwiYXVkIjoiaHR0cHM6Ly9leGFtcGxlLmNvbSJ9.C4lwXJEppq_h0JBgJZGOWyuf19RQL4nE8Jckv90Exy92XTk53y0IXhepYtKgrCTwWmLByzlikfbz6nmfpxDeggBIyAbn6H1NDzg7SqFCMXTx8caQkDyyX4uizKEnB9_z0V1hxwJdhYAiyVB9itHyXQWylwWUsihYORFSgjhK0JsrKp2VTKd304dH4XZ4jc3MSrrUK0251GV6IF0CIUZCSZlkLd1Z6gND4EY7PjSlBpYh9wd9IA9_bXt7_0F34Gvb1P_Ne9y3PFfWRmBOm-9qnHSMxn-rdYeb0LTUUsFc3w4ooO3q61u6XrW-xBgafoeahSqmlLjxPqQPwhM9NKEdiA
```

When split by dot and decoded it has a header, payload and signature.

In this case, the header is:

```json
{
  "kid": "caRgQAeHiutH35yUqrLlzROtApf_c8ZZ6XHswFBy0DM",
  "typ": "JWT",
  "alg": "RS256"
}
```

The payload is:

```json
{
  "namespace_id": "10",
  "namespace_path": "example",
  "project_id": "1",
  "project_path": "example/gitlab-ci-validate-jwt",
  "user_id": "1",
  "user_login": "root",
  "user_email": "gitlab_admin_093b83@example.com",
  "user_access_level": "owner",
  "pipeline_id": "18",
  "pipeline_source": "push",
  "job_id": "34",
  "ref": "master",
  "ref_type": "branch",
  "ref_path": "refs/heads/master",
  "ref_protected": "true",
  "runner_id": 2,
  "runner_environment": "self-hosted",
  "sha": "95d18d66afd2c0609f6c41bd537a827beb698e64",
  "project_visibility": "public",
  "ci_config_ref_uri": "gitlab.example.com/example/gitlab-ci-validate-jwt//.gitlab-ci.yml@refs/heads/master",
  "ci_config_sha": "95d18d66afd2c0609f6c41bd537a827beb698e64",
  "jti": "4bd8767e-3d66-4958-8b83-7097daebc1a7",
  "iat": 1740821142,
  "nbf": 1740821137,
  "exp": 1740824742,
  "iss": "https://gitlab.example.com",
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

The above public key should be retrieved from the GitLab jwks endpoint (e.g. https://gitlab.example.com/oauth/discovery/keys).

To see how all of this can be done read the [main.go](main.go) file.

## Reference

* https://docs.gitlab.com/ce/ci/examples/authenticating-with-hashicorp-vault/
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v18.2.1/app/models/ci/build.rb
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v18.2.1/lib/gitlab/ci/jwt.rb
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v18.2.1/lib/gitlab/ci/jwt_v2.rb
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v18.2.1/app/controllers/jwt_controller.rb
* https://gitlab.com/gitlab-org/gitlab-foss/blob/v18.2.1/app/controllers/jwks_controller.rb
* JWKS (JSON Web Key Set) endpoint (e.g. https://gitlab.example.com/oauth/discovery/keys) at https://gitlab.com/gitlab-org/gitlab-foss/blob/v18.2.1/config/routes.rb#L47
