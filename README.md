# oidc-cli

用于快速接入 OIDC 服务

## Usage

```go
var oidcConfig = oidc_cli.Config{
	ClientID:     "xxxx-xx-xxxxxxx-xxxx",
	ClientSecret: "xxxxxxxxx",
	RedirectURL:  "http://example-client.com/callback",
}

func OidcCallback(c echo.Context) error {
	code := c.QueryParam("code")
	callbackRes, err := oidc_cli.Callback(oidcConfig, code)
	if err != nil {
		...
	}
    
    log.Println(callbackRes.IDToken);
	log.Println(callbackRes.Claims.Subject);
    log.Println(callbackRes.Claims.Phone);
    
	... // use the id token & claims to further authorization
}

```

