package application

import (
	"context"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/event"
)

type EventSubscriptionError struct {
	Topic string
	Err   error
}

func (s *EventSubscriptionError) Error() string {
	return fmt.Sprintf("event subscription error: topic=%s, err=%s", s.Topic, s.Err.Error())
}

type EventHandlingError struct {
	Topic string
	Err   error
}

func (s *EventHandlingError) Error() string {
	return fmt.Sprintf("event handling error: topic=%s, err=%s", s.Topic, s.Err.Error())
}

type EventQueueManager interface {
	// Shutdown은 모든 구독을 취소하고, 모든 연결을 종료합니다.
	Shutdown(ctx context.Context) error
}

type EventPublisher interface {
	// Publish는 해당 topic에 event를 발행합니다.
	Publish(ctx context.Context, topic string, e *event.Event) error
}

type EventSubscribeCallback func(ctx context.Context, e *event.Event) error

type EventSubscriber interface {
	// Subscribe는 topic에 대한 event를 받아서 callback function을 호출하는 goroutine을 실행합니다.
	// 한 topic에 대해 하나의 callback만 등록할 수 있습니다.
	Subscribe(ctx context.Context, topic string, callback EventSubscribeCallback) error
	// SubscriptionError는 구독 중 발생하는 네트워킹 에러 등을 반환합니다.
	SubscriptionError() <-chan EventSubscriptionError
	// HandlingError는 이벤트 처리 중 발생하는 에러를 반환합니다.
	HandlingError() <-chan EventHandlingError
	// UnSubscribe는 topic에 대한 구독을 취소합니다.
	UnSubscribe(ctx context.Context, topic string) error
}

type EventQueue interface {
	EventQueueManager
	EventPublisher
	EventSubscriber
}
