package models

import _ "gorm.io/gorm"

type UserType string

var (
	TypeLeader    UserType = "Leader"
	TypeSupporter UserType = "Supporter"
	TypeDoer      UserType = "Doer"
	TypeThinker   UserType = "Thinker"
	TypeConnector UserType = "Connector"
	NoType        UserType = "NoType"
)

type User struct {
	ID       int `gorm:"primaryKey"`
	Username string
	Email    string
	Type     UserType
	Score    int
}

type Team struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	LeaderID int
	Leader   User   `gorm:"foreignKey:LeaderID"`
	Members  []User `gorm:"many2many:team_users;"`
}

type Task struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	TeamID   int
	Team     Team `gorm:"many2many:task_teams;"`
	Deadline string
}
