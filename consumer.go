package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

func startConsumer(wg *sync.WaitGroup, opt *OptionSet, builder Builder) {
	c, err := NewConsumer(opt, builder)
	if err != nil {
		logmq.Warn("starting consumer", Fields{
			"error": err,
		})
		os.Exit(1)
	}

	if opt.declare {
		c.initDeclare()
	}

	err = c.Delivery()
	if err != nil {
		logmq.Warn("consumer diliver", Fields{
			"error": err,
		})
	}

	wg.Done()
}

type Consumer struct {
	ID         string
	connection *amqp.Connection
	channels   *amqp.Channel
	builder    Builder
	opt        *OptionSet
}

func NewConsumer(opt *OptionSet, builder Builder) (*Consumer, error) {
	c := &Consumer{
		builder: builder,
		opt:     opt,
	}

	logmq.Info("connecting to broker", Fields{
		"uri": opt.uri,
	})

	var err error
	c.connection, err = amqp.Dial(opt.uri)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	return c, nil
}

func (c *Consumer) initDeclare() error {
	c.opt.prepareMode = "all"
	return Prepare(c.opt)
}

// func (c *Consumer) Shutdown() error {
// 	// will close() the deliveries channel
// 	if err := c.channel.Cancel(c.opt.consumerTag, true); err != nil {
// 		return fmt.Errorf("Consumer cancel failed: %s", err)
// 	}

// 	if err := c.connection.Close(); err != nil {
// 		return fmt.Errorf("AMQP connection close error: %s", err)
// 	}

// 	defer log.Printf("AMQP shutdown OK")

// 	// wait for handle() to exit
// 	return <-c.done
// }

func (c *Consumer) Delivery() error {
	for ci := 0; ci < c.opt.ChannelCount; ci++ {
		channel, err := c.connection.Channel()
		if err != nil {
			return fmt.Errorf("hannel: %s", err)
		}
		for cc := 0; cc < c.opt.Channelconcurrency; cc++ {
			tag := c.opt.consumerTag + strconv.Itoa(cc)
			logmq.Info("queue bound to exchange", Fields{
				"queue": c.opt.queue,
				"tag":   tag,
			})
			deliveries, err := channel.Consume(
				c.opt.queue, // name
				tag,         // consumerTag,
				true,        // autoAck
				false,       // exclusive
				false,       // noLocal
				false,       // noWait
				nil,         // arguments
			)
			if err != nil {
				return fmt.Errorf("Queue Consume: %s", err)
			}

			go c.handle(deliveries, fmt.Sprintf("%v-%v", ci, cc))
		}
	}

	return nil
}

func (c *Consumer) handle(deliveries <-chan amqp.Delivery, id string) {
	start := time.Now()
	i := 0
	c.builder.Start()
	for d := range deliveries {
		logmq.Debug("get delivery", Fields{
			"len":  len(d.Body),
			"body": string(d.Body),
			"tag":  d.DeliveryTag,
		})
		i++
		if i%c.opt.consumerStatInvertal == 0 {
			now := time.Now()
			logmq.End("lllll", Fields{
				"delay": now.Sub(d.Timestamp).Milliseconds(),
				"now":   now.UnixNano(),
				"d":     d.Timestamp.UnixNano(),
			})
			c.builder.Parse(d.Body, id)
		}
	}

	v := time.Now().Sub(start).Seconds()
	logmq.End("consumer stat", Fields{
		"count": i,
		"time":  v,
		"rate":  float64(i) / v,
	})

	logmq.Warn("delivery channel closed", nil)
}

func (c *Consumer) shutdown() error {
	return nil
}
