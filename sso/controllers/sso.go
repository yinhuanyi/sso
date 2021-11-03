/**
 * @Author: Robby
 * @File name: sso.go
 * @Create date: 2021-11-03
 * @Function:
 **/

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeHandler(c *gin.Context) {

}

func LoginHandler(c *gin.Context) {

	switch c.Request.Method {
	case http.MethodGet:
		c.JSON(http.StatusOK, gin.H{"method": "GET"})
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
