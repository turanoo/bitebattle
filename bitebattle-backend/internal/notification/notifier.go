package notification

import (
	"fmt"

	"github.com/google/uuid"
)

type Notifier struct {
	Service *Service
}

func NewNotifier(service *Service) *Notifier {
	return &Notifier{Service: service}
}

// Called when a user joins a group
func (n *Notifier) NotifyUserJoinedGroup(userID uuid.UUID, groupName string) {
	msg := fmt.Sprintf("You joined the group \"%s\"", groupName)
	_ = n.Service.Send(userID, msg)
}

// Called when a new poll is created
func (n *Notifier) NotifyPollCreated(userID uuid.UUID, groupName string) {
	msg := fmt.Sprintf("A new poll was created in group \"%s\"", groupName)
	_ = n.Service.Send(userID, msg)
}

// Called when someone votes in a poll
func (n *Notifier) NotifyVoted(userID uuid.UUID, restaurantName string) {
	msg := fmt.Sprintf("You voted for \"%s\"", restaurantName)
	_ = n.Service.Send(userID, msg)
}
