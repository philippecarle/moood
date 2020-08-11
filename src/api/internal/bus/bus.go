package bus

import (
	"github.com/streadway/amqp"
)

// Connection is a struct embedding the pointer to the Channel and a pointer method to publish
type Connection struct {
	channel *amqp.Channel
}

// Init opens the connection and the Channel
func (c *Connection) Init() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	c.channel = ch
}

// Publish will open the amqp channel and publish a message describing the new Entry
func (c *Connection) Write(j []byte) (int, error) {
	err := c.channel.Publish(
		"",
		"entries",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        j,
		})

	return len(j), err
}
