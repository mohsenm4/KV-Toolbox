package main

import "fmt"

// Mediator interface
type ChatMediator interface {
	SendMessage(msg string, user *User)
}

// Concrete Mediator
type ChatRoom struct {
	users []*User
}

func (c *ChatRoom) AddUser(user *User) {
	c.users = append(c.users, user)
}

func (c *ChatRoom) SendMessage(msg string, sender *User) {
	for _, user := range c.users {
		if user != sender {
			user.Receive(msg, sender)
		}
	}
}

// Colleague
type User struct {
	name     string
	mediator ChatMediator
}

func NewUser(name string, mediator ChatMediator) *User {
	return &User{name: name, mediator: mediator}
}

func (u *User) Send(msg string) {
	fmt.Printf("%s sends: %s\n", u.name, msg)
	u.mediator.SendMessage(msg, u)
}

func (u *User) Receive(msg string, sender *User) {
	fmt.Printf("%s receives from %s: %s\n", u.name, sender.name, msg)
}

// Main
func main() {
	chat := &ChatRoom{}

	alice := NewUser("Alice", chat)
	bob := NewUser("Bob", chat)
	charlie := NewUser("Charlie", chat)

	chat.AddUser(alice)
	chat.AddUser(bob)
	chat.AddUser(charlie)

	alice.Send("Hi everyone!")
	bob.Send("Hello Alice!")
}
