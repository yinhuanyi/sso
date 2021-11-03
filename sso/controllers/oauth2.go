/**
 * @Author: Robby
 * @File name: oauth2.go
 * @Create date: 2021-11-03
 * @Function:
 **/

package controllers

import (
	"net/http"
	"sso/sso/model"
	"sso/sso/service"
	"sso/sso/session"
	"sso/sso/utils"

	"go.uber.org/zap"

	"github.com/go-oauth2/oauth2/v4/errors"
)

func PasswordAuthorizationHandler(username, password string) (userId string, err error) {

	param := &model.UserLoginParam{
		Username: username,
		Password: password,
	}
	userId, err = service.GetUserIdByNamePwd(param)

	return
}

func UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userId string, err error) {

	if userId, err = session.Get(r, "LoggedInUserID"); err != nil {
		return
	}

	if userId == "" {

		if r.Form == nil {
			if err = r.ParseForm(); err != nil {
				return
			}
		}

		if err = session.Set(w, r, "RequestForm", r.Form.Encode()); err != nil {
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)

	}

	return
}

func AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {

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

func InternalErrorHandler(err error) (re *errors.Response) {

	zap.L().Error("Oauth2.0 Internal Error", zap.Error(err))

	return
}

func ResponseErrorHandler(re *errors.Response) {
	zap.L().Error("Oauth2.0 Response Error", zap.Error(re.Error))
}
