package message

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Slightly-Techie/st-okr-api/config"
	"github.com/Slightly-Techie/st-okr-api/internal/mailer"
	"github.com/streadway/amqp"
)


func TestRabbitMQConnection(cfg config.Config) error {

	// Load configuration

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RabbitUser, cfg.RabbitPassword, cfg.RabbitHost, cfg.RabbitPort)

	// Attempt to connect to RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Close the connection when done
	defer conn.Close()

	// If we reach here, the connection was successful
	return nil
}



func PublishMessage(eventType string, fields map[string]interface{}) error  {

	cfg := config.ENV

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RabbitUser, cfg.RabbitPassword, cfg.RabbitHost, cfg.RabbitPort)

	// Attempt to connect to RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}

	defer ch.Close()

	// map event type to queue
	queueName := getQueueName(eventType)

	if queueName == "" {
		return fmt.Errorf("invalid event type")
	}

	// Declare a queue
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	fields["event_type"] = eventType

	body, err := json.Marshal(fields)
	if err != nil {
		return fmt.Errorf("failed to marshal fields: %v", err)
	}

	log.Println("Publishing message to queue: ", queueName)

	// Publish a message
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	return nil
}


func ConsumeMessages()  {
	cfg := config.ENV

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RabbitUser, cfg.RabbitPassword, cfg.RabbitHost, cfg.RabbitPort)


	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	defer conn.Close()


	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}

	defer ch.Close()


	queues := []string{"sign_up"} // list of queues to consume messages from

	for _, queue := range queues {
		go ConsumeFromQueue(ch, queue)
	}

	select {}
	
}


func ConsumeFromQueue(ch *amqp.Channel,queueName string)  {
	// Declare the queue
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare queue %s: %v", queueName, err)
	}

	log.Println("Consuming from queues")

	// Consume messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Fatalf("Failed to register consumer for queue %s: %v", queueName, err)
	}


	// Process messages in a separate goroutine
	for msg := range msgs {
		handleMessage(msg, queueName)
	}
}



func handleMessage(msg amqp.Delivery, queueName string) {
	log.Printf("Received message")

	var fields map[string]interface{}
	if err := json.Unmarshal(msg.Body, &fields); err != nil {
		log.Println("Failed to unmarshal message from queue %s: %v", queueName, err)
		msg.Nack(false, false) // Nack the message and don't requeue
		return
	}

	log.Println("Message received from queue: ", queueName)

	eventType, ok := fields["event_type"].(string)
	if !ok {
		log.Println("Event type missing or invalid in message from queue %s", queueName)
		msg.Nack(false, false)
		return
	}

	switch eventType {
	case "sign_up":
		handleSignUpMailer(fields,msg)
	default:
		log.Println("Unknown event type: %s in queue %s", eventType, queueName)
		msg.Nack(false, false)
	}
	
}


func handleSignUpMailer(fields map[string]interface{},msg amqp.Delivery)  {
	userName,ok := fields["user_name"].(string)
	if !ok {
		log.Println("User name missing or invalid in message from queue sign_up")
		msg.Nack(false, false)
		return
	}

	userEmail,ok := fields["email"].(string)
	if !ok {
		log.Println("Email missing or invalid in message from queue sign_up")
		msg.Nack(false, false)
		return
	}


	if err := mailer.SendWelcomeEmail(userEmail, userName); err != nil {
		log.Println("Failed to send welcome email to user %s: %v", userName, err)
	}

	msg.Ack(false)
}




func getQueueName(eventType string) string {
	switch eventType {
	case "sign_up":
		return "sign_up"
	default:
		return ""
	}
}