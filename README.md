# oidc-cli

用于快速接入 OIDC 服务

## Usage

在某个神秘的地方注册好 Client，配置 redirect url 并获取 id 和 secret

```go
import(
    "github.com/BingyanStudio/oidc-cli/oidc"
)

var oidcConfig = oidc.Config{
    ClientID:     "xxxx-xx-xxxxxxx-xxxx",
    ClientSecret: "xxxxxxxxx",
    RedirectURL:  "http://example-client.com/callback",
}

func OidcCallbackHandler(c echo.Context) error {
    code := c.QueryParam("code")
    callbackRes, err := oidc.Callback(oidcConfig, code)
    if err != nil {
        ...
    }

    log.Println(callbackRes.IDToken);
    log.Println(callbackRes.Claims.Subject);
    log.Println(callbackRes.Claims.Phone);
    
    ... // use the token & claims to further authorization
}

```

