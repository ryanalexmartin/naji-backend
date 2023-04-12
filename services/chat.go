package services

import (
	"backend/models"
)

func MatchUsers(user1, user2 *models.User) *models.Room {
	room := models.NewRoom()

	room.AddUser(user1)
	room.AddUser(user2)

	return room
}

func BroadcastMessage(room *models.Room, sender *models.User, message string) {
	room.Users.Range(func(_, userValue interface{}) bool {
		user := userValue.(*models.User)
		err := user.Conn.WriteMessage(1, []byte(message))
		if err != nil {
			// Handle error if required
		}
		return true
	})
}

func LeaveChatSession(room *models.Room, user *models.User) {
	room.RemoveUser(user.ID)
	user.Room = ""

	if room.UserCount() == 0 {
		// Perform any room cleanup if necessary
	}
}
