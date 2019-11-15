package main

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

// Prepare before bench test
func Prepare(opt *OptionSet) error {
	connection, err := amqp.Dial(opt.uri)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("hannel: %s", err)
	}
	defer channel.Close()

	

	switch opt.prepareMode {
	case "exchange":
		err = DeclareExchange(channel, opt)

	case "queue":
		err = DeclareQueue(channel, opt)

	case "binding":
		err = DeclareBind(channel, opt)

	case "all":
		err = DeclareExchange(channel, opt)
		if err != nil {
			return err
		}
		err = DeclareQueue(channel, opt)
		if err != nil {
			return err
		}
		err = DeclareBind(channel, opt)
		if err != nil {
			return err
		}

	default:
		err = errors.New("--help")
	}

	return err
}

// DeclareExchange declare exchange with option
func DeclareExchange(channel *amqp.Channel, opt *OptionSet) error {
	logmq.Info("declaring exchange", Fields{
		"name":       opt.exchange,
		"type":       opt.exchangeType,
		"durable":    opt.durable,
		"autoDelete": opt.autoDelete,
		"internal":   opt.internal,
		"noWait":     opt.noWait,
	})

	err := channel.ExchangeDeclare(
		opt.exchange,     // name of the exchange
		opt.exchangeType, // type
		opt.durable,      // durable
		opt.autoDelete,   // delete when complete
		opt.internal,     // internal
		opt.noWait,       // noWait
		nil,              // arguments
	)
	if err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}
	return nil
}

// DeclareQueue declare queue with option
func DeclareQueue(channel *amqp.Channel, opt *OptionSet) error {
	logmq.Info("declaring queue", Fields{
		"name":       opt.queue,
		"durable":    opt.durable,
		"autoDelete": opt.autoDelete,
		"internal":   opt.exclusive,
		"noWait":     opt.noWait,
	})

	state, err := channel.QueueDeclare(
		opt.queue,      // name of the queue
		opt.durable,    // durable
		opt.autoDelete, // delete when usused
		opt.exclusive,  // exclusive
		opt.noWait,     // noWait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	logmq.Info("queue state", Fields{
		"name":      state.Name,
		"message":   state.Messages,
		"consumers": state.Consumers,
	})
	return nil
}

// DeclareBind declare binding between exchange and queue
func DeclareBind(channel *amqp.Channel, opt *OptionSet) error {
	logmq.Info("declaring binding", Fields{
		"queue":      opt.queue,
		"routingKey": opt.routingKey,
		"exchange":   opt.exchange,
		"noWait":     opt.noWait,
	})

	err := channel.QueueBind(
		opt.queue,      // name of the queue
		opt.routingKey, // routingKey
		opt.exchange,   // sourceExchange
		opt.noWait,     // noWait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}
	return nil
}
