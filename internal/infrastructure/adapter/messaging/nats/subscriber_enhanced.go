// Package nats
package nats

// type Subscriber struct {
// 	conn          *nats.Conn
// 	js            nats.JetStreamContext
// 	config        *config.NATSConfig
// 	registry      *event.Registry
// 	subscriptions map[string]*nats.Subscription
// 	mu            sync.RWMutex
// 	cancelFuncs   map[string]context.CancelFunc
// 	wg            sync.WaitGroup
// }
//
// func NewSubscriber(config *config.NATSConfig, registry *event.Registry) (*Subscriber, error) {
// 	opts := []nats.Option{
// 		nats.Timeout(config.ConnectionTimeout),
// 		nats.RetryOnFailedConnect(true),
// 		nats.MaxReconnects(-1),
// 		nats.ReconnectWait(2 * time.Second),
// 		nats.ClosedHandler(func(nc *nats.Conn) {
// 			log.Printf("[NATS] Connection closed: %s\n", nc.Opts.Url)
// 		}),
// 		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
// 			if err != nil {
// 				log.Printf("[NATS] Disconnected: %v\n", err)
// 			}
// 		}),
// 		nats.ReconnectHandler(func(nc *nats.Conn) {
// 			log.Printf("[NATS] Reconnected to %s\n", nc.ConnectedUrl())
// 		}),
// 	}
//
// 	conn, err := nats.Connect(config.URL, opts...)
// 	if err != nil {
// 		return nil, ErrConnectionFailed(err)
// 	}
//
// 	js, err := conn.JetStream(nats.MaxWait(config.RequestTimeout))
// 	if err != nil {
// 		conn.Close()
// 		return nil, ErrConnectionFailed(err)
// 	}
//
// 	subscriber := &Subscriber{
// 		conn:          conn,
// 		js:            js,
// 		config:        config,
// 		registry:      registry,
// 		subscriptions: make(map[string]*nats.Subscription),
// 		cancelFuncs:   make(map[string]context.CancelFunc),
// 	}
//
// 	return subscriber, nil
// }
//
// func (s *Subscriber) Subscribe(ctx context.Context, eventName string, handler repository.IHandler) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
//
// 	// Check if already subscribed
// 	if _, exists := s.subscriptions[eventName]; exists {
// 		return fmt.Errorf("already subscribed to event: %s", eventName)
// 	}
//
// 	subject := fmt.Sprintf("events.%s", eventName)
// 	consumerName := fmt.Sprintf("%s-%s", s.config.ConsumerName, strings.ReplaceAll(eventName, ".", "_"))
//
// 	// Create consumer config with proper configuration
// 	consumerConfig := &nats.ConsumerConfig{
// 		Durable:           consumerName,
// 		DeliverPolicy:     nats.DeliverAllPolicy,
// 		AckPolicy:         nats.AckExplicitPolicy,
// 		AckWait:           s.config.AckWait,
// 		MaxDeliver:        s.config.MaxDeliver,
// 		MaxAckPending:     s.config.MaxAckPending,
// 		ReplayPolicy:      nats.ReplayInstantPolicy,
// 		FilterSubject:     subject,
// 		InactiveThreshold: 10 * time.Minute, // Clean up inactive consumers
// 		MemoryStorage:     false,
// 	}
//
// 	// Create consumer
// 	_, err := s.js.AddConsumer(s.config.StreamName, consumerConfig)
// 	if err != nil {
// 		// If consumer exists, continue (it's ok)
// 		if apiErr, ok := err.(*nats.APIError); !ok || apiErr.ErrorCode != 10014 { // Consumer exists
// 			log.Printf("[NATS] Failed to create consumer: %s, error: %v", consumerName, err)
// 			return ErrConsumerCreationFailed(err)
// 		}
// 	}
//
// 	// Use async subscription instead of pull subscription for better performance
// 	sub, err := s.js.Subscribe(
// 		subject,
// 		s.createMsgHandler(handler),
// 		nats.Durable(consumerName),
// 		nats.ManualAck(),
// 		nats.MaxDeliver(s.config.MaxDeliver), // Let NATS handle max deliveries
// 		nats.AckWait(s.config.AckWait),
// 	)
// 	if err != nil {
// 		return ErrSubscribeFailed(err)
// 	}
//
// 	// Store the subscription
// 	s.subscriptions[eventName] = sub
//
// 	log.Printf("[NATS] Successfully subscribed to event: %s (consumer: %s)", eventName, consumerName)
// 	return nil
// }
//
// func (s *Subscriber) createMsgHandler(handler repository.IHandler) nats.MsgHandler {
// 	return func(msg *nats.Msg) {
// 		// Use context.Background() for message processing
// 		ctx := context.Background()
//
// 		meta, err := msg.Metadata()
// 		if err != nil {
// 			log.Printf("[NATS] Failed to get message metadata: %v", err)
// 			_ = msg.Nak() // Negative ack
// 			return
// 		}
//
// 		log.Printf("[NATS] Processing message (attempt %d/%d): %s",
// 			meta.NumDelivered, s.config.MaxDeliver, msg.Subject)
//
// 		// Deserialize event
// 		evt, err := s.registry.Unmarshal(msg.Data)
// 		if err != nil {
// 			log.Printf("[NATS] Failed to deserialize event: %v", err)
// 			// Mark message as bad and don't retry
// 			_ = msg.Term()
// 			return
// 		}
//
// 		// Check if handler can handle this event
// 		if !handler.CanHandle(evt) {
// 			log.Printf("[NATS] Handler cannot process event: %s", evt.EventName())
// 			_ = msg.Ack() // Ack to not reprocess
// 			return
// 		}
//
// 		// Process the event with the handler
// 		if err := handler.Handle(ctx, evt); err != nil {
// 			log.Printf("[NATS] Handler failed (attempt %d): %v", meta.NumDelivered, err)
//
// 			// Check if we've reached max delivery attempts
// 			if meta.NumDelivered >= uint64(s.config.MaxDeliver) {
// 				log.Printf("[NATS] Max delivery attempts reached for subject: %s, moving to DLQ", msg.Subject)
//
// 				// Move to DLQ or terminate
// 				dlqSubject := fmt.Sprintf("dlq.%s", msg.Subject)
// 				if pubErr := s.publishToDLQ(ctx, dlqSubject, msg.Data); pubErr != nil {
// 					log.Printf("[NATS] Failed to publish to DLQ: %v", pubErr)
// 				}
//
// 				// Ack the original message to stop reprocessing
// 				_ = msg.Ack()
// 				return
// 			}
//
// 			// Calculate backoff with jitter for better distribution
// 			backoff := s.calculateBackoff(int(meta.NumDelivered))
// 			log.Printf("[NATS] Retrying in %v...", backoff)
//
// 			// Nack with delay for retry
// 			_ = msg.NakWithDelay(backoff)
// 			return
// 		}
//
// 		// Successful processing
// 		if err := msg.Ack(); err != nil {
// 			log.Printf("[NATS] Failed to ACK message: %v", err)
// 		} else {
// 			log.Printf("[NATS] Message processed successfully: %s", msg.Subject)
// 		}
// 	}
// }
//
// func (s *Subscriber) publishToDLQ(ctx context.Context, subject string, data []byte) error {
// 	// Publish to DLQ stream
// 	_, err := s.js.Publish(subject, data)
// 	return err
// }
//
// func (s *Subscriber) calculateBackoff(attempt int) time.Duration {
// 	// Exponential backoff with jitter
// 	baseDelay := time.Second
// 	maxDelay := 30 * time.Second
//
// 	// Calculate exponential backoff
// 	delay := baseDelay * (1 << uint(attempt-1)) // 2^(attempt-1) * baseDelay
//
// 	// Cap the delay
// 	if delay > maxDelay {
// 		delay = maxDelay
// 	}
//
// 	// Add jitter (±25% of delay)
// 	jitter := time.Duration(float64(delay) * 0.25)
// 	delayWithJitter := delay + time.Duration(simpleRand())%jitter - jitter/2
//
// 	if delayWithJitter < baseDelay {
// 		delayWithJitter = baseDelay
// 	}
//
// 	return delayWithJitter
// }
//
// // Simple pseudo-random number generator for jitter
// func simpleRand() int64 {
// 	return time.Now().UnixNano() % 1000
// }
//
// func (s *Subscriber) Unsubscribe(eventName string) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
//
// 	sub, exists := s.subscriptions[eventName]
// 	if !exists {
// 		return fmt.Errorf("not subscribed to event: %s", eventName)
// 	}
//
// 	// Unsubscribe
// 	if err := sub.Unsubscribe(); err != nil {
// 		return fmt.Errorf("failed to unsubscribe from %s: %w", eventName, err)
// 	}
//
// 	// Remove from map
// 	delete(s.subscriptions, eventName)
//
// 	log.Printf("[NATS] Unsubscribed from event: %s", eventName)
// 	return nil
// }
//
// func (s *Subscriber) Close() error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
//
// 	// Close all subscriptions
// 	var errors []error
// 	for eventName, sub := range s.subscriptions {
// 		if err := sub.Unsubscribe(); err != nil {
// 			errors = append(errors, fmt.Errorf("failed to unsubscribe from %s: %w", eventName, err))
// 		}
// 	}
//
// 	// Clear maps
// 	s.subscriptions = make(map[string]*nats.Subscription)
// 	s.cancelFuncs = make(map[string]context.CancelFunc)
//
// 	// Close connection
// 	if s.conn != nil {
// 		s.conn.Close()
// 	}
//
// 	log.Println("[NATS] Subscriber closed")
//
// 	if len(errors) > 0 {
// 		return fmt.Errorf("errors occurred during close: %v", errors)
// 	}
// 	return nil
// }
//
// // HealthCheck returns the connection status
// func (s *Subscriber) HealthCheck(ctx context.Context) error {
// 	// Try a simple request to verify the connection is alive
// 	_, err := s.js.StreamInfo(s.config.StreamName, nats.Context(ctx))
// 	if err != nil {
// 		return fmt.Errorf("health check failed: %w", err)
// 	}
// 	return nil
// }
