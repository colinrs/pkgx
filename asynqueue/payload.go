package asynqueue

import "encoding/json"

type PayLoad[T any] struct {
	RequestId  string `json:"requestId"`
	MarketCode string `json:"marketCode"`
	UserId     int64  `json:"userId"`
	UserEmail  string `json:"userEmail"`
	Body       T      `json:"body"`
	Event      string `json:"eventName"`
}

func NewPayload[T any](event string, body T) *PayLoad[T] {
	return &PayLoad[T]{
		Event: event,
		Body:  body,
	}
}

func (t *PayLoad[T]) SetUserId(userId int64) {
	t.UserId = userId
}

func (t *PayLoad[T]) SetUserEmail(userEmail string) {
	t.UserEmail = userEmail
}

func (t *PayLoad[T]) SetRequestId(requestId string) {
	t.RequestId = requestId
}

func (t *PayLoad[T]) SetMarketCode(marketCode string) {
	t.MarketCode = marketCode
}

func (t *PayLoad[T]) Marshal() []byte {
	b, _ := json.Marshal(t)
	return b
}

func (t *PayLoad[T]) EventName() string {
	return t.Event
}
