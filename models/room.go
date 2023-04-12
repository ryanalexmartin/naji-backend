package models

import (
	"github.com/google/uuid"
	"sync"
)

type Room struct {
	ID    string
	Users sync.Map
}

func NewRoom() *Room {
	return &Room{
		ID: uuid.New().String(),
	}
}

func (r *Room) AddUser(user *User) {
	r.Users.Store(user.ID, user)
	user.Room = r.ID
}

func (r *Room) RemoveUser(userID string) {
	if user, ok := r.Users.Load(userID); ok {
		r.Users.Delete(userID) // user is still connected after being removed from the room
		user.(*User).Room = "" // therefore we need to update the user's room ID to empty
	}
}

func (r *Room) UserCount() int {
	userCount := 0
	r.Users.Range(func(_, _ interface{}) bool {
		userCount++
		return true
	})
	return userCount
}
