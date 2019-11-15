package main

import (
	"os"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Option default option set
var Option = OptionSet{
	durable: true,
}

// OptionSet input set
type OptionSet struct {
	exchange   string
	routingKey string
	queue      string

	message      string
	uri          string
	exchangeType string
	consumerTag  string
	builder      string
	reliable     bool

	ChannelCount       int
	Channelconcurrency int
	interval           int
	msgCount           int

	consumerStatInvertal int
	consumerTimeout      time.Duration
	logLevel             string

	declare     bool
	prepareMode string
	durable     bool
	autoDelete  bool
	exclusive   bool
	noWait      bool
	internal    bool
}

func flagParse() string {
	app := kingpin.New("rabbitmq-bench", "for: test rabbitmq bench")
	app.HelpFlag.Short('h')

	app.Flag("log", "special log level - debug|info|warn|error|end").
		Default("end").EnumVar(&Option.logLevel, "debug", "info", "warn", "error", "end")
	app.Flag("uri", "speical rabbitmq broker uri").
		Default("amqp://guest:guest@127.0.0.1:5672").Short('u').StringVar(&Option.uri)
	app.Flag("type", "Exchange type - direct|fanout|topic").
		Default("direct").EnumVar(&Option.exchangeType, "direct", "fanout", "topic")
	app.Flag("channel-count", "channel count in only one tcp connection").
		Default("1").Short('c').IntVar(&Option.ChannelCount)
	app.Flag("concurrency", "Concurrency count in only one channel. if special greater 1, it mean that one channel will create multi goroutine").
		Default("1").Short('g').IntVar(&Option.Channelconcurrency)
	app.Flag("declare", "declare default exchange for producer. or declare default exchang queue direct binding for consumer").
		Default("false").BoolVar(&Option.declare)
	app.Flag("build", "build message with message prefix - time|order").Default("time").EnumVar(&Option.builder, "time", "order")

	// subcommand producer
	{
		producer := app.Command("producer", "producer mode")
		producer.Flag("interval", "send message freq(ms)").Default("0").Short('i').IntVar(&Option.interval)
		producer.Flag("count", "send message count").Default("10000").Short('n').IntVar(&Option.msgCount)
		producer.Arg("exchange", "special exchange name").Required().StringVar(&Option.exchange)
		producer.Arg("routingkey", "special routing-key name").Required().StringVar(&Option.routingKey)
		producer.Arg("message", "special message prefix").Required().StringVar(&Option.message)
	}

	// subcommand consumer
	{
		consumer := app.Command("consumer", "consumer mode")
		consumer.Flag("tag", "consumer tag prefix").Default("tag").StringVar(&Option.consumerTag)
		consumer.Flag("stat.interval", "send special count msg, stat the time and rate").Default("100").IntVar(&Option.consumerStatInvertal)
		consumer.Flag("stat.timeout", "special stat the time long").Default("10s").DurationVar(&Option.consumerTimeout)
		consumer.Arg("exchange", "special exchange name").Required().StringVar(&Option.exchange)
		consumer.Arg("queue", "special queue name").Required().StringVar(&Option.queue)
		consumer.Arg("routingkey", "special routing-key name. if not special, will use the same value as queue").
			StringVar(&Option.routingKey)
	}

	// subcommand prepare
	{
		prepare := app.Command("prepare", "prepare exchange or queue or binding")
		prepare.Flag("durable", "durable is or not").Default("false").BoolVar(&Option.durable)
		prepare.Flag("autoDelete", "autoDelete is or not").Default("false").BoolVar(&Option.autoDelete)
		prepare.Flag("internal", "only exchange").Default("false").BoolVar(&Option.internal)
		prepare.Flag("exclusive", "only queue").Default("false").BoolVar(&Option.exclusive)
		prepare.Flag("noWait", "noWait is or not").Default("false").BoolVar(&Option.noWait)

		// subcommand prepare all
		{
			all := prepare.Command("all", "declare exchange, routingkey, queue")
			all.Arg("exchange", "special exchange name").Required().StringVar(&Option.exchange)
			all.Arg("queue", "special queue name. ").Required().StringVar(&Option.queue)
			all.Arg("routingkey", "special routing-key name. if not special, will use the same value as queue").
				StringVar(&Option.routingKey)
		}

		// subcommand prepare exchange
		{
			exchange := prepare.Command("exchange", "only declare exchange")
			exchange.Arg("exchange", "special exchange name").Required().StringVar(&Option.exchange)
		}

		// subcommand prepare queue
		{
			queue := prepare.Command("queue", "only declare queue")
			queue.Arg("queue", "special exchange name").Required().StringVar(&Option.queue)
		}

		// subcommand prepare routingkey
		{
			routingkey := prepare.Command("routingkey", "only declare routingkey")
			routingkey.Arg("exchange", "special exchange name").Required().StringVar(&Option.exchange)
			routingkey.Arg("queue", "special queue name.").Required().StringVar(&Option.queue)
			routingkey.Arg("routingkey", "special routing-key name. if not special, will use the same value as queue").
				StringVar(&Option.routingKey)
		}
	}

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	flagFix(&Option)

	return cmd
}

func flagFix(opt *OptionSet) {
	if len(opt.queue) > 0 && len(opt.routingKey) == 0 {
		opt.routingKey = opt.queue
	}
}
