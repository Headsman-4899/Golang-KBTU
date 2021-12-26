package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	lru "github.com/hashicorp/golang-lru"
	"lectures/hw6/internal/message_broker"
	"lectures/hw6/internal/models"
	"log"
)

const (
	cacheTopic = "cache"
)

type (
	CacheBroker struct {
		syncProducer   sarama.SyncProducer
		consumerGroup  sarama.ConsumerGroup
		consumeHandler *cacheConsumeHandler
		clientID       string
	}

	cacheConsumeHandler struct {
		cache *lru.TwoQueueCache
		ready chan bool
	}
)

func NewCacheBroker(cache *lru.TwoQueueCache, clientID string) message_broker.CacheBroker {
	return &CacheBroker{
		clientID: clientID,
		consumeHandler: &cacheConsumeHandler{
			cache: cache,
			ready: make(chan bool),
		},
	}
}

func (c *CacheBroker) Connect(ctx context.Context, brokers []string) error {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 10
	producerConfig.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer(brokers, producerConfig)
	if err != nil {
		return err
	}
	c.syncProducer = syncProducer

	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	consumerGroup, err := sarama.NewConsumerGroup(brokers, c.clientID, consumerConfig)
	if err != nil {
		return err
	}
	c.consumerGroup = consumerGroup

	go func() {
		for {
			if err := c.consumerGroup.Consume(ctx, []string{cacheTopic}, c.consumeHandler); err != nil {
				log.Panicf("Error from consumer %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			c.consumeHandler.ready = make(chan bool)
		}
	}()
	<-c.consumeHandler.ready
	return nil
}

func (c *CacheBroker) Close() error {
	if err := c.syncProducer.Close(); err != nil {
		return err
	}
	if err := c.consumerGroup.Close(); err != nil {
		return err
	}
	return nil
}

func (c *CacheBroker) Remove(key interface{}) error {
	msg := &models.CacheMsg{
		Command: models.CacheCommandRemove,
		Key:     key,
	}

	msgRaw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, _, err = c.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: cacheTopic,
		Value: sarama.StringEncoder(msgRaw),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheBroker) Purge() error {
	msg := &models.CacheMsg{
		Command: models.CacheCommandPurge,
	}

	msgRaw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, _, err = c.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: cacheTopic,
		Value: sarama.StringEncoder(msgRaw),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *cacheConsumeHandler) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *cacheConsumeHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *cacheConsumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)

		cacheMsg := new(models.CacheMsg)
		if err := json.Unmarshal(msg.Value, cacheMsg); err != nil {
			return err
		}

		switch cacheMsg.Command {
		case models.CacheCommandRemove:
			c.cache.Remove(cacheMsg.Key)
		case models.CacheCommandPurge:
			c.cache.Purge()
		}

		session.MarkMessage(msg, "")
	}

	return nil
}
