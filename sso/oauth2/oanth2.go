/**
 * @Author: Robby
 * @File name: oanth2.go
 * @Create date: 2021-11-03
 * @Function:
 **/

package oauth2

import (
	"log"
	"net/http"
	"sso/sso/model"
	"sso/sso/service"
	"sso/sso/session"
	"sso/sso/settings"
	"sso/sso/utils"

	"github.com/go-oauth2/oauth2/v4/errors"
	"go.uber.org/zap"

	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	oredis "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

var Manager *manage.Manager
var Srv *server.Server

func Init(cfg *settings.Oauth2Config) (err error) {

	Manager = manage.NewDefaultManager()
	Manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   1,
	}))

	Manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	Manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
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

	Manager.MapClientStorage(clientStore)

	Srv = server.NewDefaultServer(Manager)
	Srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler)
	Srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	Srv.SetAuthorizeScopeHandler(authorizeScopeHandler)
	Srv.SetInternalErrorHandler(internalErrorHandler)
	Srv.SetResponseErrorHandler(responseErrorHandler)

	return
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userId string, err error) {

	if userId, err = session.Get(r, "LoggedInUserID"); err != nil {
		return
	}

	if userId == "" {

		if err = r.ParseForm(); err != nil {
			return
		}

		if err = session.Set(w, r, "RequestForm", r.Form.Encode()); err != nil {
			return
		}

		http.Redirect(w, r, "/api/v1/login", http.StatusFound)

	}

	return
}

func passwordAuthorizationHandler(username, password string) (userId string, err error) {

	param := &model.UserLoginParam{
		Username: username,
		Password: password,
	}
	userId, err = service.GetUserIdByNamePwd(param)

	return
}

func authorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {

	if err = r.ParseForm(); err != nil {
		return
	}

	scopeObj := utils.GetClientScope(r.Form.Get("client_id"), r.Form.Get("scope"))
	if scopeObj == nil {
		http.Error(w, "Invalid Scope", http.StatusBadRequest)
		return
	}

	scope = utils.ScopeNameJoin(scopeObj)

	return
}

func internalErrorHandler(err error) (re *errors.Response) {

	zap.L().Error("Oauth2.0 Internal Error", zap.Error(err))

	return
}

func responseErrorHandler(re *errors.Response) {
	zap.L().Error("Oauth2.0 Response Error", zap.Error(re.Error))
}
