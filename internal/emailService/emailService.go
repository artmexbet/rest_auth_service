package emailService

import "fmt"

// MocEmail ...
type MocEmail struct {
}

// New ...
func New() *MocEmail {
	return &MocEmail{}
}

// SendWarning ...
func (e *MocEmail) SendWarning(email string) error {
	fmt.Printf("Sanding message to email '%s'\n", email)
	return nil
}
