package models

import (
	"database/sql"

	_ "gorm.io/gorm"
)

type UserType struct {
	NullString sql.NullString
}

func NullUserType(str string) UserType {
	return UserType{
		NullString: sql.NullString{
			String: str,
			Valid:  str != "",
		},
	}
}

var (
	TypeLeader    UserType = NullUserType("Leader")
	TypeSupporter UserType = NullUserType("Supporter")
	TypeDoer      UserType = NullUserType("Doer")
	TypeThinker   UserType = NullUserType("Thinker")
	TypeConnector UserType = NullUserType("Connector")
)

type User struct {
	ID       int `gorm:"primaryKey"`
	Username string
	Email    string
	Type     UserType
	Score    sql.NullInt64
}

type Team struct {
	ID      int `gorm:"primaryKey"`
	Name    string
	Leader  User
	Members []User
}

type Task struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	Team     Team
	Deadline string
}
