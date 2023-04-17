package models

import "github.com/dgrijalva/jwt-go"

type GroupDetails struct {
	Id          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Visibility  string `json:"visibility" db:"visibility"`
	IsVerified  string `json:"isVerified" db:"is_verified"`
	CreatedBy   string `json:"createdBy" db:"created_by"`
	Url         string `json:"url" db:"url"`
	UploadType  string `json:"uploadType" db:"type"`
}

type Claims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}
