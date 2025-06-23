package notification

import (
	"fmt"
)

type Notifier struct {
	Service *Service
}

func NewNotifier(service *Service) *Notifier {
	return &Notifier{Service: service}
}

func (n *Notifier) NotifyUserJoinedGroup(userID string, groupName string) {
	msg := fmt.Sprintf("You joined the group \"%s\"", groupName)
	_ = n.Service.Send(userID, msg)
}

func (n *Notifier) NotifyPollCreated(userID string, groupName string) {
	msg := fmt.Sprintf("A new poll was created in group \"%s\"", groupName)
	_ = n.Service.Send(userID, msg)
}

func (n *Notifier) NotifyVoted(userID string, restaurantName string) {
	msg := fmt.Sprintf("You voted for \"%s\"", restaurantName)
	_ = n.Service.Send(userID, msg)
}
