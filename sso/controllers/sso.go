/**
 * @Author: Robby
 * @File name: sso.go
 * @Create date: 2021-11-03
 * @Function:
 **/

package controllers

import (
	"html/template"
	"net/http"
	"net/url"
	"sso/sso/model"
	"sso/sso/oauth2"

	"sso/sso/session"
	"sso/sso/utils"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// AuthorizeHandler Get接口
func AuthorizeHandler(c *gin.Context) {

	if err := session.Delete(c.Writer, c.Request, "RequestForm"); err != nil {
		zap.L().Error("[AuthorizeHandler]：session.Delete", zap.Error(err))
		ResponseError(c, CodeServerInternalError)
		return
	}

	if err := oauth2.Srv.HandleAuthorizeRequest(c.Writer, c.Request); err != nil {
		zap.L().Error("[AuthorizeHandler]：Srv.HandleAuthorizeRequest", zap.Error(err))
		ResponseError(c, CodeServerInternalError)
		return
	}

}

func ReAuthorizeHandler(c *gin.Context) {

}

func LoginHandler(c *gin.Context) {

	switch c.Request.Method {

	case http.MethodGet:
		requestForm, err := session.Get(c.Request, "RequestForm")
		if err != nil {
			zap.L().Error("[LoginHandler]：session.Get", zap.Error(err))
			ResponseError(c, CodeBadRequest)
			return
		}

		if requestForm == "" {
			zap.L().Error("[LoginHandler]：session.Get", zap.Error(err))
			ResponseError(c, CodeBadRequest)
			return
		}

		decodeForm, err := url.ParseQuery(requestForm)
		if err != nil {
			ResponseError(c, CodeServerInternalError)
			return
		}

		// Get client_id and scope_name from user
		clientID := decodeForm.Get("client_id")
		scope := decodeForm.Get("scope")
		clientObj := utils.GetClientObj(clientID)
		scopeObj := utils.GetClientScope(clientID, scope)
		if scopeObj == nil {
			zap.L().Error("[LoginHandler]：bad scope")
			ResponseError(c, CodeInvalidParam)
			return
		}
		data := model.ClientScope{
			Client: clientObj,
			Scope:  scopeObj,
		}

		tmpl, err := template.ParseFiles("sso/tpl/login.html")
		if err != nil {
			zap.L().Error("[LoginHandler]：template parse error")
			ResponseError(c, CodeServerInternalError)
			return
		}

		if err = tmpl.Execute(c.Writer, data); err != nil {
			zap.L().Error("[LoginHandler]：template execute error")
			//ResponseError(c, CodeServerInternalError)
			//return
		}
		//return

	case http.MethodPost:
		c.JSON(http.StatusOK, gin.H{"method": "POST"})
	}
}

func LogoutHandler(c *gin.Context) {

}

func TokenHandler(c *gin.Context) {

}

func VerifyHandler(c *gin.Context) {

}
