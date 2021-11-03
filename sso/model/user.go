/**
 * @Author: Robby
 * @File name: user.go
 * @Create date: 2021-11-03
 * @Function:
 **/

package model

type User struct {
	UserId   int64  `json:"user_id"  db:"user_id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type UserLoginParam struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
