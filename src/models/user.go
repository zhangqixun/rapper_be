package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"utility"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Motto    string `json:"motto"`
}

type UserEditorInfo struct {
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Motto    string `json:"motto"`
	Token    string `json:"token"`
}

type UserLoginInfo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserInfo struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Motto    string `json:"motto"`
}

func Register(user User) int {
	db, err := gorm.Open("mysql", utility.DBAddr)
	if err != nil {
		fmt.Println(err)
		return DB_ERROR
	}
	defer db.Close()

	if len(user.Name) == 0 || len(user.Password) == 0 || len(user.Nickname) == 0 || len(user.Phone) == 0 {
		return REQUIRE_FIELD_EMPTY
	}

	var existingUsers []User
	var count int
	db.Where("name=?", user.Name).Find(&existingUsers).Count(&count)
	if count > 0 {
		return USER_EXIST
	}

	db.Create(&user)
	return SUCCESS
}

func Editor(userdi UserEditorInfo, id int) int {
	db, err := gorm.Open("mysql", utility.DBAddr)
	if err != nil {
		fmt.Println(err)
		return DB_ERROR
	}
	defer db.Close()

	if len(userdi.Nickname) == 0 || len(userdi.Phone) == 0 {
		return REQUIRE_FIELD_EMPTY
	}

	var existingUsers []User
	var count int
	db.Where("id=?", id).Find(&existingUsers).Count(&count)
	if count != 1 {
		return USER_NOT_EXIST
	}

	existingUsers[0].Nickname = userdi.Nickname
	existingUsers[0].Phone = userdi.Phone
	existingUsers[0].Email = userdi.Email
	existingUsers[0].Motto = userdi.Motto
	db.Save(&existingUsers[0])
	return SUCCESS
}

func Login(user_login_info UserLoginInfo) int {
	db, err := gorm.Open("mysql", utility.DBAddr)
	if err != nil {
		fmt.Println(err)
		return DB_ERROR
	}
	defer db.Close()

	if len(user_login_info.Name) == 0 || len(user_login_info.Password) == 0 {
		return REQUIRE_FIELD_EMPTY
	}

	var existingUsers []User
	var count int
	db.Where("name=?", user_login_info.Name).Find(&existingUsers).Count(&count)
	if count != 1 {
		return USER_NOT_EXIST
	}

	if existingUsers[0].Password != user_login_info.Password {
		return PASSWORD_ERROR
	} else {
		return existingUsers[0].ID
	}
}

func Query(id int) (res int, info UserInfo) {
	db, err := gorm.Open("mysql", utility.DBAddr)
	if err != nil {
		fmt.Println(err)
		res = DB_ERROR
		return
	}
	defer db.Close()

	var user User
	err_q := db.Where("id=?", id).First(&user).Error
	if err_q != nil {
		res = USER_NOT_EXIST
		return
	} else {
		info = UserInfo{Name: user.Name, Nickname: user.Nickname, Phone: user.Phone, Email: user.Email, Motto: user.Motto}
		res = SUCCESS
		return
	}
}

func PasswordModification(id int, old_passwd string, new_passwd string) int {
	db, err := gorm.Open("mysql", utility.DBAddr)
	if err != nil {
		fmt.Println(err)
		return DB_ERROR
	}
	defer db.Close()

	if len(new_passwd) == 0 {
		return REQUIRE_FIELD_EMPTY
	}

	var existingUsers []User
	var count int
	db.Where("id=?", id).Find(&existingUsers).Count(&count)
	if count != 1 {
		return USER_NOT_EXIST
	}
	if old_passwd != existingUsers[0].Password {
		return PASSWORD_ERROR
	}
	db.Model(&existingUsers[0]).Update("password", new_passwd)
	return SUCCESS
}
