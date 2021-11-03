/**
 * @Author: Robby
 * @File name: oanth2.go
 * @Create date: 2021-11-03
 * @Function:
 **/

package oauth2

import (
	"log"
	"sso/sso/controllers"
	"sso/sso/settings"

	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

var manager *manage.Manager
var srv *server.Server

func Init(cfg *settings.Oauth2Config) (err error) {

	manager = manage.NewDefaultManager()
	manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   1,
	}))

	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	clientStore := store.NewClientStore()

	for _, v := range cfg.Client {
		if err = clientStore.Set(v.ClientId, &models.Client{
			ID:     v.ClientId,
			Secret: v.ClientSecret,
			Domain: v.ClientDomain,
		}); err != nil {
			log.Printf("客户端注册SSO失败：%s\n", err.Error())
			return
		}
	}

	manager.MapClientStorage(clientStore)

	srv = server.NewDefaultServer(manager)
	srv.SetPasswordAuthorizationHandler(controllers.PasswordAuthorizationHandler)
	srv.SetUserAuthorizationHandler(controllers.UserAuthorizeHandler)
	srv.SetAuthorizeScopeHandler(controllers.AuthorizeScopeHandler)
	srv.SetInternalErrorHandler(controllers.InternalErrorHandler)
	srv.SetResponseErrorHandler(controllers.ResponseErrorHandler)

	return
}
