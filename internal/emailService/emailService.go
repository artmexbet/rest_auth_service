package emailService

import "fmt"

// Email ...
type Email struct {
}

// New ...
func New() *Email {
	return &Email{}
}

// SendWarning ...
func (e *Email) SendWarning(email string) error {
	fmt.Printf("Sanding message to email '%s'\n", email)
	return nil
}
