package models

const SUCCESS int = 0
const DB_ERROR int = -1
const REQUIRE_FIELD_EMPTY int = -2
const USER_EXIST int = -3
const USER_NOT_EXIST int = -4
const PASSWORD_ERROR int = -5

const SUCCESS_MESS string = "success"
const DB_ERROR_MESS string = "database error"
const REQUIRE_FIELD_EMPTY_MESS string = "some required fields are empty"
const USER_EXIST_MESS string = "the user already exists"
const USER_NOT_EXIST_MESS string = "the user does not exist"
const PASSWORD_ERROR_MESS string = "password mistake"

const SUCCESS_CODE string = "200"
const DB_ERROR_CODE string = "401"
const REQUIRE_FIELD_EMPTY_CODE string = "402"
const USER_EXIST_CODE string = "403"
const USER_NOT_EXIST_CODE string = "404"
const PASSWORD_ERROR_CODE string = "405"
