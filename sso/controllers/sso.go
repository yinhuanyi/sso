/**
 * @Author: Robby
 * @File name: sso.go
 * @Create date: 2021-11-03
 * @Function:
 **/

package controllers

import (
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"sso/sso/model"
	"sso/sso/oauth2"
	"sso/sso/service"
	"strconv"

	"sso/sso/session"
	"sso/sso/utils"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// AuthorizeHandler Get接口，第一次调用
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

// ReAuthorizeHandler Get接口，第二次调用，数据从session中取出来
func ReAuthorizeHandler(c *gin.Context) {

	var err error
	var requestFormString string
	var requestForm url.Values

	if requestFormString, err = session.Get(c.Request, "RequestForm"); err != nil {
		zap.L().Error("[ReAuthorizeHandler]：session.Get", zap.Error(err))
		ResponseError(c, CodeServerInternalError)
		return
	}

	if requestForm, err = url.ParseQuery(requestFormString); err != nil {
		zap.L().Error("[ReAuthorizeHandler]：url.ParseQuery", zap.Error(err))
		ResponseError(c, CodeServerInternalError)
		return
	}

	// 给请求的form赋值
	c.Request.Form = requestForm

	if err = session.Delete(c.Writer, c.Request, "RequestForm"); err != nil {
		zap.L().Error("[ReAuthorizeHandler]：session.Delete", zap.Error(err))
		ResponseError(c, CodeServerInternalError)
		return
	}

	if err = oauth2.Srv.HandleAuthorizeRequest(c.Writer, c.Request); err != nil {
		zap.L().Error("[ReAuthorizeHandler]：oauth2.Srv.HandleAuthorizeRequest", zap.Error(err))
		ResponseError(c, CodeServerInternalError)
		return
	}

}

// 获取requestForm的数据, 为LoginHandler服务
func getRequestForm(c *gin.Context) (data *model.ClientScope, err error) {

	requestForm, err := session.Get(c.Request, "RequestForm")
	if err != nil {
		zap.L().Error("[LoginHandler]：session.Get", zap.Error(err))
		return nil, errors.New(CodeBadRequest.ToString())
	}

	if requestForm == "" {
		zap.L().Info("[LoginHandler]：requestForm == '' ")
		return nil, errors.New(CodeBadRequest.ToString())
	}

	decodeForm, err := url.ParseQuery(requestForm)
	if err != nil {
		return nil, errors.New(CodeServerInternalError.ToString())
	}

	// Get client_id and scope_name from user
	clientID := decodeForm.Get("client_id")
	scope := decodeForm.Get("scope")
	clientObj := utils.GetClientObj(clientID)
	scopeObj := utils.GetClientScope(clientID, scope)
	if scopeObj == nil {
		zap.L().Error("[LoginHandler]：bad scope")
		return nil, errors.New(CodeInvalidParam.ToString())
	}

	return &model.ClientScope{
		Client: clientObj,
		Scope:  scopeObj,
	}, nil

}

func LoginHandler(c *gin.Context) {

	switch c.Request.Method {

	case http.MethodGet:

		data, err := getRequestForm(c)

		if err != nil {
			code, err := strconv.Atoi(err.Error())
			if err != nil {
				zap.L().Error("[LoginHandler]：strconv.Atoi", zap.Error(err))
				ResponseError(c, CodeServerInternalError)
			}
			ResponseError(c, ResCode(code))
			return
		}

		tmpl, err := template.ParseFiles("sso/tpl/login.html")
		if err != nil {
			zap.L().Error("[LoginHandler]：template parse error")
			ResponseError(c, CodeServerInternalError)
			return
		}

		if err = tmpl.Execute(c.Writer, data); err != nil {
			zap.L().Error("[LoginHandler]：template execute error")
		}

	case http.MethodPost:

		// csrf token verify
		if c.PostForm("type") == "password" {

			// 如果传递的是空密码那么binding就会校验，依托的是github.com/gin-gonic/gin/binding库，自动校验参数返回
			userLoginParam := &model.UserLoginParam{
				Username: c.PostForm("username"),
				Password: c.PostForm("password"),
			}

			userID, err := service.GetUserIdByNamePwd(userLoginParam)

			if err != nil {
				zap.L().Error("[LoginHandler]：service.GetUserIdByNamePwd", zap.Error(err))
			}

			if userID == "" {

				tmpl, err := template.ParseFiles("sso/tpl/login.html")
				if err != nil {
					zap.L().Error("[LoginHandler]：template.ParseFiles", zap.Error(err))
					ResponseError(c, CodeServerInternalError)
					return
				}

				data, err := getRequestForm(c)

				if err != nil {
					code, err := strconv.Atoi(err.Error())
					if err != nil {
						zap.L().Error("[LoginHandler]：strconv.Atoi", zap.Error(err))
						ResponseError(c, CodeServerInternalError)
					}
					ResponseError(c, ResCode(code))
				}

				if data != nil {
					data.Error = "用户名或密码错误"
				}

				if err = tmpl.Execute(c.Writer, data); err != nil {
					zap.L().Error("[LoginHandler]：tmpl.Execute")
				}

			}

			if err = session.Set(c.Writer, c.Request, "LoggedInUserID", userID); err != nil {
				ResponseError(c, CodeServerInternalError)
				zap.L().Error("[LoginHandler]：session.Set", zap.Error(err))
				return
			}

			c.Redirect(http.StatusFound, "/api/v1/reauthorize")
			return

		}

		// csrf token error
		ResponseError(c, CodeBadRequest)

	}
}

func LogoutHandler(c *gin.Context) {

}

func TokenHandler(c *gin.Context) {

}

func VerifyHandler(c *gin.Context) {

}
