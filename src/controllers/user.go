package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"models"
	"net/http"
	"utility"
)

type RegisterRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type LoginData struct {
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
}

type LoginRes struct {
	Code    string    `json:"code"`
	Data    LoginData `json:"data"`
	Message string    `json:"message"`
}

type TokenInfo struct {
	Token string `json:"token"`
}

type LogoutRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type EditorRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type QueryRes struct {
	Code     string          `json:"code"`
	Message  string          `json:"message"`
	UserInfo models.UserInfo `json:"user_info"`
}

type PasswordModiInfo struct {
	Token       string `json:"token"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type PasswordModiRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type BrowseInfo struct {
	Token      string `json:"token"`
	MovieID    int    `json:"movie_id"`
	TimeOnSite string `json:"time_on_site"`
}

type BrowseRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var user models.User
	_ = json.Unmarshal(body, &user)

	res := models.Register(user)
	var info RegisterRes
	if res == models.SUCCESS {
		info = RegisterRes{Code: models.SUCCESS_CODE, Message: models.SUCCESS_MESS}
	} else if res == models.DB_ERROR {
		info = RegisterRes{Code: models.DB_ERROR_CODE, Message: models.DB_ERROR_MESS}
	} else if res == models.REQUIRE_FIELD_EMPTY {
		info = RegisterRes{Code: models.REQUIRE_FIELD_EMPTY_CODE, Message: models.REQUIRE_FIELD_EMPTY_MESS}
	} else if res == models.USER_EXIST {
		info = RegisterRes{Code: models.USER_EXIST_CODE, Message: models.USER_EXIST_MESS}
	}
	res_json, _ := json.Marshal(info)
	fmt.Fprint(w, string(res_json))
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var user_login_info models.UserLoginInfo
	_ = json.Unmarshal(body, &user_login_info)

	res := models.Login(user_login_info)
	var info LoginRes
	if res > 0 {
		token := utility.NewSession(res)
		info = LoginRes{Code: models.SUCCESS_CODE, Data: LoginData{Token: token, UserID: res}, Message: models.SUCCESS_MESS}
	} else if res == models.DB_ERROR {
		info = LoginRes{Code: models.DB_ERROR_CODE, Message: models.DB_ERROR_MESS}
	} else if res == models.REQUIRE_FIELD_EMPTY {
		info = LoginRes{Code: models.REQUIRE_FIELD_EMPTY_CODE, Message: models.REQUIRE_FIELD_EMPTY_MESS}
	} else if res == models.USER_NOT_EXIST {
		info = LoginRes{Code: models.USER_NOT_EXIST_CODE, Message: models.USER_NOT_EXIST_MESS}
	} else if res == models.PASSWORD_ERROR {
		info = LoginRes{Code: models.PASSWORD_ERROR_CODE, Message: models.PASSWORD_ERROR_MESS}
	}
	res_json, _ := json.Marshal(info)
	fmt.Fprint(w, string(res_json))
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var token_info TokenInfo
	_ = json.Unmarshal(body, &token_info)

	id := utility.DropSession(token_info.Token)
	var info LogoutRes
	if id < 0 {
		info = LogoutRes{Code: models.USER_NOT_EXIST_CODE, Message: models.USER_NOT_EXIST_MESS}
	} else {
		info = LogoutRes{Code: models.SUCCESS_CODE, Message: models.SUCCESS_MESS}
	}
	res_json, _ := json.Marshal(info)
	fmt.Fprint(w, string(res_json))
}

func UserEditor(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var user_editor_info models.UserEditorInfo
	_ = json.Unmarshal(body, &user_editor_info)

	id := utility.CheckSession(user_editor_info.Token)
	res := models.Editor(user_editor_info, id)
	var info EditorRes
	if res == models.SUCCESS {
		info = EditorRes{Code: models.SUCCESS_CODE, Message: models.SUCCESS_MESS}
	} else if res == models.DB_ERROR {
		info = EditorRes{Code: models.DB_ERROR_CODE, Message: models.DB_ERROR_MESS}
	} else if res == models.REQUIRE_FIELD_EMPTY {
		info = EditorRes{Code: models.REQUIRE_FIELD_EMPTY_CODE, Message: models.REQUIRE_FIELD_EMPTY_MESS}
	} else if res == models.USER_NOT_EXIST {
		info = EditorRes{Code: models.USER_NOT_EXIST_CODE, Message: models.USER_NOT_EXIST_MESS}
	}
	res_json, _ := json.Marshal(info)
	fmt.Fprint(w, string(res_json))
}

func UserQuery(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var token_info TokenInfo
	_ = json.Unmarshal(body, &token_info)

	id := utility.CheckSession(token_info.Token)
	res, user_info := models.Query(id)
	var info QueryRes
	if res == models.SUCCESS {
		info = QueryRes{Code: models.SUCCESS_CODE, Message: models.SUCCESS_MESS, UserInfo: user_info}
	} else if res == models.DB_ERROR {
		info = QueryRes{Code: models.DB_ERROR_CODE, Message: models.DB_ERROR_MESS}
	} else if res == models.USER_NOT_EXIST {
		info = QueryRes{Code: models.USER_NOT_EXIST_CODE, Message: models.USER_NOT_EXIST_MESS}
	}
	res_json, _ := json.Marshal(info)
	fmt.Fprint(w, string(res_json))
}

func UserPasswordModification(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var passwd_info PasswordModiInfo
	_ = json.Unmarshal(body, &passwd_info)

	id := utility.CheckSession(passwd_info.Token)
	res := models.PasswordModification(id, passwd_info.OldPassword, passwd_info.NewPassword)
	var info PasswordModiRes
	if res == models.SUCCESS {
		info = PasswordModiRes{Code: models.SUCCESS_CODE, Message: models.SUCCESS_MESS}
	} else if res == models.DB_ERROR {
		info = PasswordModiRes{Code: models.DB_ERROR_CODE, Message: models.DB_ERROR_MESS}
	} else if res == models.REQUIRE_FIELD_EMPTY {
		info = PasswordModiRes{Code: models.REQUIRE_FIELD_EMPTY_CODE, Message: models.REQUIRE_FIELD_EMPTY_MESS}
	} else if res == models.USER_NOT_EXIST {
		info = PasswordModiRes{Code: models.USER_NOT_EXIST_CODE, Message: models.USER_NOT_EXIST_MESS}
	} else if res == models.PASSWORD_ERROR {
		info = PasswordModiRes{Code: models.PASSWORD_ERROR_CODE, Message: models.PASSWORD_ERROR_MESS}
	}
	res_json, _ := json.Marshal(info)
	fmt.Fprint(w, string(res_json))
}

func UserBrowse(w http.ResponseWriter, r *http.Request) {
	utility.PreprocessXHR(&w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var browse_info BrowseInfo
	_ = json.Unmarshal(body, &browse_info)

	user_id := utility.CheckSession(browse_info.Token)
	footprint := models.Footprint{
		UserID:     user_id,
		MovieID:    browse_info.MovieID,
		TimeOnSite: browse_info.TimeOnSite,
	}
	res := models.InsertFootprint(footprint)

	var info BrowseRes
	if res == models.SUCCESS {
		info = BrowseRes{
			Code:    models.SUCCESS_CODE,
			Message: models.SUCCESS_MESS,
		}
	} else if res == models.DB_ERROR {
		info = BrowseRes{
			Code:    models.DB_ERROR_CODE,
			Message: models.DB_ERROR_MESS,
		}
	} else if res == models.REQUIRE_FIELD_EMPTY {
		info = BrowseRes{
			Code:    models.REQUIRE_FIELD_EMPTY_CODE,
			Message: models.REQUIRE_FIELD_EMPTY_MESS,
		}
	}
	res_json, _ := json.Marshal(info)
	fmt.Fprint(w, string(res_json))
}
