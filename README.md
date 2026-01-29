# workout_ledger_backend
Backend for Workout Ledger application.

## Environment variables

| Name | Description | Default |
| --- | --- | --- |
| `POSTGRES_DSN` | PostgreSQL connection string. | (required) |
| `HTTP_ADDR` | HTTP listen address. | `:8080` |
| `GOOGLE_OAUTH_CLIENT_ID` | Google OAuth client ID. | (required for OAuth) |
| `GOOGLE_OAUTH_CLIENT_SECRET` | Google OAuth client secret. | (required for OAuth) |
| `GOOGLE_OAUTH_REDIRECT_URL` | Google OAuth redirect URL. | (required for OAuth) |
| `GOOGLE_OAUTH_ALLOWED_DOMAIN` | Optional allowed email domain. | (optional) |
| `JWT_SECRET` | HMAC secret for signing JWTs. | (required) |
| `JWT_ACCESS_TTL_MINUTES` | Access token TTL in minutes. | `15` |
| `JWT_REFRESH_TTL_DAYS` | Refresh token TTL in days. | `30` |
