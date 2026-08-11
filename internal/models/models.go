package models

import (
	"fmt"

	_ "gorm.io/gorm"
)

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
	ID           int    `gorm:"primaryKey"`
	Username     string `gorm:"unique"`
	PasswordHash string
	Email        string `grom:"unique"`
	Type         UserType
	Score        int
	ScoreCount   int
}

// TODO: Complete UpdateUserAverageScore method to update user's score!
func (user *User) UpdateUserAverageScore(new_score int) int {
	user.ScoreCount += 1
	return (user.Score + new_score) / user.ScoreCount
}

func (user *User) UpdateUserType(new_type UserType) bool {
	/*
		This method returns true if existing user's type has been changed.
		It also returns false if user had no type before!
	*/
	if user.Type != NoType {
		user.Type = new_type
		return false
	}
	user.Type = new_type
	return true

}

// NOTE: Defined a function to describe users
func (user *User) String() string {
	return fmt.Sprintf("user %v (%v)", user.Username, user.Type)
}

type Team struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	LeaderID int
	Leader   User   `gorm:"foreignKey:LeaderID"`
	Members  []User `gorm:"many2many:team_users;"`
}

// TODO: Complete AddMember method to add member to a team!
func (team *Team) AddMember(user_obj User) bool {
	if user_obj == team.Leader {
		return false
	}
	team.Members = append(team.Members, user_obj)
	return true
}

// TODO: Complete RemoveMember, so only leaders can remove a member!
func (team *Team) RemoveMember(user_obj User) bool {
	if user_obj == team.Leader {
		return false
	}
	for ind, val := range team.Members {
		if val == user_obj {
			team.Members = append(team.Members[:ind], team.Members[ind+1:]...)
		}
	}
	return true
}

// TODO: Complete ChangeLeader, so leaders be able to Change the leader of the team!
func (team *Team) ChangeLeader(user_obj User) bool {
	if user_obj == team.Leader {
		return false
	}
	team.LeaderID = user_obj.ID
	return true
}

// NOTE: Defined a function to describe teams
func (team *Team) String() string {
	return fmt.Sprintf("team %v (%v)", team.Name, team.Leader.Username)
}

type TaskStatus string

var (
	StatusPending   TaskStatus = "Pending"
	StatusCompleted TaskStatus = "Completed"
)

type Task struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	TeamID   int
	Team     Team `gorm:"foreignKey:TeamID"`
	Deadline string
	Status   TaskStatus
}

// TODO: Define a method named RenewDeadline so leader be able to renew the deadline!
func (task *Task) RenewDeadline(new_deadline string) bool {
	if new_deadline != "" {
		task.Deadline = new_deadline
		return true
	}
	return false
}

// TODO: Define a method named MarkAsCompleted to mark tasks as completed if they are not!
func (task *Task) MarkAsCompleted(user_id int, leader_id int) bool {
	if task.Status == StatusPending && user_id == leader_id {
		task.Status = StatusCompleted
		return true
	}
	return false
}

// NOTE: Defined a function to describe tasks
func (task *Task) String() string {
	return fmt.Sprintf("task %v (%v)", task.Name, task.Status)
}
