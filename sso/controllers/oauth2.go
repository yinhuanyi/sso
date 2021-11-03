/**
 * @Author: Robby
 * @File name: oauth2.go
 * @Create date: 2021-11-03
 * @Function:
 **/

package controllers

import (
	"log"
	"net/http"
	"sso/sso/model"
	"sso/sso/service"
	"sso/sso/session"

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

func UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {

	v, err := session.Get(r, "LoggedInUserID")
	if err != nil {
		return
	}

	if v == nil {
		if r.Form == nil {
			r.ParseForm()
		}
		session.Set(w, r, "RequestForm", r.Form)

		http.Redirect(w, r, "/login", http.StatusFound)

		return
	}

	userID = v.(string)

	return
}

func AuthorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {

	if r.Form == nil {
		r.ParseForm()
	}

	s := config.ScopeFilter(r.Form.Get("client_id"), r.Form.Get("scope"))
	if s == nil {
		http.Error(w, "Invalid Scope", http.StatusBadRequest)
		return
	}
	scope = config.ScopeJoin(s)

	return
}

func InternalErrorHandler(err error) (re *errors.Response) {
	log.Println("Internal Error:", err.Error())
	return
}

func ResponseErrorHandler(re *errors.Response) {
	log.Println("Response Error:", re.Error.Error())
}
