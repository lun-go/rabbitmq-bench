package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

// Builder format parse for stat
type Builder interface {
	Format(*OptionSet, string) *amqp.Publishing
	Start()
	Parse(*OptionSet, *amqp.Delivery, string) interface{}
}

// BodyTimer implement Builder for cal delay
type BodyTimer struct {
	last time.Time
}

// Format format ampq publishing
func (b *BodyTimer) Format(opt *OptionSet, data string) *amqp.Publishing {
	// why not use Publishing.Timestamp?
	// because Publishing.Timestamp precision is second
	return &amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		Body:            []byte(fmt.Sprintf("%s||%v", data, time.Now().UnixNano())),
		DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
		Priority:        0,              // 0-9
	}
}

// Start time start
func (b *BodyTimer) Start() {
	b.last = time.Now()
}

// Parse parse deliver message
func (b *BodyTimer) Parse(opt *OptionSet, deliver *amqp.Delivery, id string) interface{} {
	now := time.Now()
	gap := now.Sub(b.last).Seconds()
	ss := strings.Split(string(deliver.Body), "||")[1]
	fromTime, _ := strconv.Atoi(ss)
	b.last = now
	logmq.End("consumer stat interval", Fields{
		"id":        id,
		"count":     opt.consumerStatInvertal,
		"time(s)":   gap,
		"rate(/s)":  int(float64(opt.consumerStatInvertal) / gap),
		"delay(ms)": (now.UnixNano() - int64(fromTime)) / int64(time.Millisecond),
	})

	return nil
}

// Body implement Builder
type Body struct {
	i    int
	last time.Time
}

// Format format ampq publishing
func (b *Body) Format(opt *OptionSet, data string) *amqp.Publishing {
	m := &amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		Body:            []byte(fmt.Sprintf("%s||%v", data, b.i)),
		DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
		Priority:        0,              // 0-9
	}
	b.i++
	return m
}

// Start time start
func (b *Body) Start() {
	b.last = time.Now()
}

// Parse parse deliver message
func (b *Body) Parse(opt *OptionSet, deliver *amqp.Delivery, id string) interface{} {
	now := time.Now()
	last := now.Sub(b.last).Seconds()
	i := Option.consumerStatInvertal
	b.last = now
	logmq.End("consumer stat interval", Fields{
		"count":    i,
		"time(s)":  last,
		"rate(/s)": int(float64(i) / last),
	})
	return nil
}
