package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

func startProducer(wg *sync.WaitGroup, opt *OptionSet, builder Builder) error {
	p, err := NewProducer(opt, builder)
	if err != nil {
		return err
	}
	if opt.declare {
		p.exchangeDeclare()
	}
	p.Publish(wg)
	wg.Done()
	return nil
}

type Producer struct {
	connection *amqp.Connection
	builder    Builder
	opt        *OptionSet
}

func NewProducer(opt *OptionSet, builder Builder) (*Producer, error) {
	p := &Producer{
		builder: builder,
		opt:     opt,
	}

	logmq.Printf(LevelInfo, "Connecting to %s", p.opt.uri)
	var err error
	p.connection, err = amqp.Dial(p.opt.uri)
	if err != nil {
		return nil, fmt.Errorf("Dial: %v", err)
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	// if reliable {
	// 	if err := p.channel.Confirm(false); err != nil {
	// 		return nil, fmt.Errorf("Channel could not be put into confirm mode: ", err)
	// 	}

	// 	ack, nack := p.channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	// 	// defer confirmOne(ack, nack)
	// }

	return p, nil
}

func (p *Producer) exchangeDeclare() error {
	channel, err := p.connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %v", err)
	}
	defer channel.Close()

	return DeclareExchange(channel, p.opt)
}

func (p *Producer) Publish(wg *sync.WaitGroup) error {
	for ci := 0; ci < p.opt.ChannelCount; ci++ {
		channel, err := p.connection.Channel()
		if err != nil {
			return fmt.Errorf("Channel: %v", err)
		}

		for cc := 0; cc < p.opt.Channelconcurrency; cc++ {
			wg.Add(1)
			go func() {
				start := time.Now()
				for i := 1; i <= p.opt.msgCount; i++ {
					err = p.publishOne(channel, p.opt.message)
					if err != nil {
						logmq.Warn("publish exchange", Fields{
							"error": err,
						})
						return
					}

					if p.opt.interval > 0 {
						time.Sleep(time.Duration(p.opt.interval) * time.Millisecond)
					}
				}

				v := time.Now().Sub(start).Seconds()
				logmq.End("producer stat", Fields{
					"count":    p.opt.msgCount,
					"time(s)":  v,
					"rate(/s)": int(float64(p.opt.msgCount) / v),
				})
				wg.Done()
			}()
		}
	}
	return nil
}

func (p *Producer) publishOne(channel *amqp.Channel, body string) error {
	logmq.Debug("publishing msg ", Fields{
		"exchange": p.opt.exchange,
		"key":      p.opt.routingKey,
		"body":     string(body),
		"len":      len(body),
	})

	pm := p.builder.Format(p.opt, body)
	err := channel.Publish(
		p.opt.exchange,   // publish to an exchange
		p.opt.routingKey, // routing to 0 or more queues
		false,            // mandatory
		false,            // immediate
		*pm,
	)

	if err != nil {
		return fmt.Errorf("Exchange Publish: %v", err)
	}

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(ack, nack chan uint64) {
	logmq.Printf(LevelDebug, "waiting for confirmation of one publishing")

	select {
	case tag := <-ack:
		log.Printf("confirmed delivery with delivery tag: %d", tag)
	case tag := <-nack:
		log.Printf("failed delivery of delivery tag: %d", tag)
	}
}
