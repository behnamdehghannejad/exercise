package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Timestamped struct {
	CreatedAt time.Time
	SentAt    time.Time
}

type Retryable struct {
	Attempts    int
	MaxAttempts int
}

type Prioritized struct {
	Priority string
}

const (
	Low    = "LOW"
	Medium = "MEDIUM"
	High   = "HIGH"
)

type NotificationSender interface {
	Send(message string, recipient string) error
	GetType() string
}

type EmailNotification struct {
	Timestamped
	Retryable
	Prioritized
	EmailAddress string
}

func (e *EmailNotification) Send(message string, recipient string) error {
	return randomSuccess(e, message, recipient)
}

func (*EmailNotification) GetType() string {
	return "email"
}

type SMSNotification struct {
	Timestamped
	Retryable
	Prioritized
	PhoneNumber string
}

func (s *SMSNotification) Send(message string, recipient string) error {
	return randomSuccess(s, message, recipient)
}

func (*SMSNotification) GetType() string {
	return "sms"
}

type PushNotification struct {
	Timestamped
	Retryable
	Prioritized
	DeviceToken string
}

func (p *PushNotification) Send(message string, recipient string) error {
	return randomSuccess(p, message, recipient)
}

func (*PushNotification) GetType() string {
	return "push"
}

func randomSuccess(sender NotificationSender, message string, recipient string) error {
	for i := 0; i < 3; i++ {
		if n := rand.Intn(100); n <= 20 {
			return nil
		}
	}
	return fmt.Errorf("%s sent to %s: %s\n", sender.GetType(), recipient, message)
}

type NotificationService struct {
	sender NotificationSender
}

func NewNotificationService(sender NotificationSender) NotificationService {
	return NotificationService{
		sender: sender,
	}
}

type Order struct {
	ID        string
	Amount    float64
	recipient string
}

type OrderNotification struct {
	NotificationSender
	Order
}

func NotificationSystem() {
	email := &EmailNotification{
		EmailAddress: "behnam@gmail.com",
		Prioritized:  Prioritized{Priority: Medium},
		Retryable:    Retryable{MaxAttempts: 3},
		Timestamped:  Timestamped{CreatedAt: time.Now()},
	}

	push := &PushNotification{
		DeviceToken: "845484856698485",
		Prioritized: Prioritized{Priority: Low},
		Retryable:   Retryable{MaxAttempts: 2},
		Timestamped: Timestamped{CreatedAt: time.Now()},
	}

	sms := &SMSNotification{
		PhoneNumber: "091298765432",
		Prioritized: Prioritized{Priority: High},
		Retryable:   Retryable{MaxAttempts: 5},
		Timestamped: Timestamped{CreatedAt: time.Now()},
	}

	sendNotif(email)
	fmt.Println("")
	fmt.Println("=====================")
	sendNotif(push)
	fmt.Println("")
	fmt.Println("=====================")
	sendNotif(sms)
}

func sendNotif(sender NotificationSender) {
	service := NewNotificationService(sender)

	order := Order{ID: "12345", Amount: 100.0, recipient: "behnam"}
	orderNotification := OrderNotification{
		Order:              order,
		NotificationSender: sender,
	}

	err := service.Send("your order is ready", order.recipient)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = orderNotification.Send(
		fmt.Sprintf("Order %s confirmed, amount: %.2f to ", order.ID, order.Amount),
		sender.GetType(),
	)

	if err != nil {
		fmt.Println("order error:", err)
	}
}

func (s *NotificationService) Send(message, recipient string) error {
	switch v := s.sender.(type) {

	case *EmailNotification:
		fmt.Printf("priority of %s is %s\n:", s.sender.GetType(), v.Priority)
		for v.Attempts < v.MaxAttempts {
			err := s.sender.Send(message, recipient)
			v.Attempts++

			if err == nil {
				v.SentAt = time.Now()
				fmt.Println(message + " " + recipient + " at " + v.SentAt.Format(time.RFC3339))
				return nil
			}
		}
		err := fmt.Sprintf("failed after max retries by %s at %s | max retries = %d", s.sender.GetType(), v.SentAt.Format(time.RFC3339), v.MaxAttempts)
		return errors.New(err)

	case *SMSNotification:
		fmt.Printf("priority of %s is %s\n:", s.sender.GetType(), v.Priority)
		for v.Attempts < v.MaxAttempts {
			err := s.sender.Send(message, recipient)
			v.Attempts++

			if err == nil {
				v.SentAt = time.Now()
				fmt.Println(message + " " + recipient + " at " + v.SentAt.Format(time.RFC3339))
				return nil
			}
		}
		err := fmt.Sprintf("failed after max retries by %s at %s | max retries = %d", s.sender.GetType(), v.SentAt.Format(time.RFC3339), v.MaxAttempts)
		return errors.New(err)

	case *PushNotification:
		fmt.Printf("priority of %s is %s\n:", s.sender.GetType(), v.Priority)
		for v.Attempts < v.MaxAttempts {
			err := s.sender.Send(message, recipient)
			v.Attempts++

			if err == nil {
				v.SentAt = time.Now()
				fmt.Println(message + " " + recipient + " at " + v.SentAt.Format(time.RFC3339))
				return nil
			}
		}
		err := fmt.Sprintf("failed after max retries by %s at %s | max retries = %d", s.sender.GetType(), v.SentAt.Format(time.RFC3339), v.MaxAttempts)
		return errors.New(err)

	case *OrderNotification:
		fmt.Printf(message+" to : %s", recipient)
	}

	return nil
}
