package gochannel

import (
	"context"
	"errors"
	"go-pubsub/message"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

type Config struct {
	ChannleBufferSize int64

	// need persistent events ?

}

var GoChanPubsub *GoChannelPubsub

func init() {

	cfg := Config{ChannleBufferSize: 100}

	GoChanPubsub = NewGoChannelPubsub(cfg)
	slog.Info("go chan pubsub started")

}

type GoChannelPubsub struct {
	config Config

	// logger slog add

	subscribersWg          sync.WaitGroup
	subscribers            map[string][]*Subscriber
	subscribersLock        sync.RWMutex
	subscribersByTopicLock sync.Map // map of *sync.Mutex

	closed     bool
	closedLock sync.Mutex
	closing    chan struct{}
}

func NewGoChannelPubsub(config Config) *GoChannelPubsub {

	return &GoChannelPubsub{
		config: config,

		subscribers:            make(map[string][]*Subscriber),
		subscribersByTopicLock: sync.Map{},

		closing: make(chan struct{}),
	}
}

func (g *GoChannelPubsub) Publish(topic string, messages ...*message.Message) error {

	if g.isClosed() {
		return errors.New("Pub/Sub closed")
	}

	messagesToPublish := make(message.Messages, len(messages))
	for i, msg := range messages {
		messagesToPublish[i] = msg.Copy()
	}

	g.subscribersLock.RLock()
	defer g.subscribersLock.RUnlock()

	subLock, _ := g.subscribersByTopicLock.LoadOrStore(topic, &sync.Mutex{})
	subLock.(*sync.Mutex).Lock()
	defer subLock.(*sync.Mutex).Unlock()

	for i := range messagesToPublish {
		msg := messagesToPublish[i]

		err := g.sendMessage(topic, msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GoChannelPubsub) sendMessage(topic string, message *message.Message) error {

	subscribers := g.topicSubscribers(topic)

	if len(subscribers) == 0 {
		slog.Info("no subscribers have")
		return nil
	}

	go func(subscribers []*Subscriber) {
		wg := &sync.WaitGroup{}

		for i := range subscribers {
			subscriber := subscribers[i]

			wg.Add(1)
			go func() {
				subscriber.sendMessageToSubscriber(message)
				wg.Done()
			}()
		}

		wg.Wait()

	}(subscribers)

	return nil
}

func (g *GoChannelPubsub) Subscribe(ctx context.Context, topic string) (chan *message.Message, error) {

	g.closedLock.Lock()

	if g.closed {
		g.closedLock.Unlock()
		return nil, errors.New("Pub/Sub closed")
	}

	g.subscribersWg.Add(1)
	g.closedLock.Unlock()

	g.subscribersLock.Lock()

	subLock, _ := g.subscribersByTopicLock.LoadOrStore(topic, &sync.Mutex{})
	subLock.(*sync.Mutex).Lock()

	s := &Subscriber{
		ctx:                  ctx,
		id:                   uuid.New(),
		CommunicationChannel: make(chan *message.Message, g.config.ChannleBufferSize),

		closing: make(chan struct{}),
	}

	defer g.subscribersLock.Unlock()
	defer subLock.(*sync.Mutex).Unlock()

	g.addSubscriber(topic, s)

	return s.CommunicationChannel, nil

}

func (g *GoChannelPubsub) addSubscriber(topic string, s *Subscriber) {
	if _, ok := g.subscribers[topic]; !ok {
		g.subscribers[topic] = make([]*Subscriber, 0)
	}
	g.subscribers[topic] = append(g.subscribers[topic], s)
}

func (g *GoChannelPubsub) topicSubscribers(topic string) []*Subscriber {
	subscribers, ok := g.subscribers[topic]
	if !ok {
		return nil
	}

	// let's do a copy to avoid race conditions and deadlocks due to lock
	subscribersCopy := make([]*Subscriber, len(subscribers))
	copy(subscribersCopy, subscribers)

	return subscribersCopy
}

func (g *GoChannelPubsub) isClosed() bool {
	g.closedLock.Lock()
	defer g.closedLock.Unlock()

	return g.closed
}

// Close closes the GoChannel Pub/Sub.
func (g *GoChannelPubsub) Close() error {
	g.closedLock.Lock()
	defer g.closedLock.Unlock()

	if g.closed {
		return nil
	}

	g.closed = true
	close(g.closing)

	// log
	g.subscribersWg.Wait()

	return nil
}

type Subscriber struct {
	id uuid.UUID

	sending              sync.Mutex
	CommunicationChannel chan *message.Message

	closed  bool
	closing chan struct{}

	ctx context.Context
}

func (s *Subscriber) Close() {
	if s.closed {
		return
	}
	close(s.closing)

	// ensuring that we are not sending to closed channel
	s.sending.Lock()
	defer s.sending.Unlock()

	s.closed = true

	close(s.CommunicationChannel)
}

func (s *Subscriber) sendMessageToSubscriber(msg *message.Message) {
	s.sending.Lock()
	defer s.sending.Unlock()

	ctx, cancelCtx := context.WithCancel(s.ctx)
	defer cancelCtx()

	// copy the message to prevent ack/nack propagation to other consumers
	// also allows to make retries on a fresh copy of the original message
	msgToSend := msg.Copy()
	msgToSend.SetContext(ctx)

	if s.closed {
		slog.Info("Pub/Sub closed, discarding msg ")
		return
	}

	select {
	case s.CommunicationChannel <- msgToSend:
		slog.Info("Sent message to subscriber ")
	case <-s.closing:
		slog.Info("Closing, message discarded ")
		return
	}

}
