package models

import (
	"database/sql"

	_ "gorm.io/gorm"
)

type UserType string

var (
	TypeLeader    UserType = "Leader"
	TypeSupporter UserType = "Supporter"
	TypeDoer      UserType = "Doer"
	TypeThinker   UserType = "Thinker"
	TypeConnector UserType = "Connector"
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
