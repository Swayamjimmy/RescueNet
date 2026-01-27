package p2p

import (
	"context"
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

type ChatRoom struct {
	Messages chan *ChatMessage
	ctx      context.Context
	ps       *pubsub.PubSub
	topic    *pubsub.Topic
	sub      *pubsub.Subscription
	roomName string
	self     peer.ID
	nick     string
}

type ChatMessage struct {
	Message    string
	SenderID   string
	SenderNick string
}

const ChatRoomBufSize = 128

func topicName(roomName string) string {
	return "chat-room" + roomName
}

func JoinChatRoom(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID, nickname string, roomName string) (*ChatRoom, error) {
	topic, err := ps.Join(topicName(roomName))

	if err != nil {
		return nil, err
	}

	sub, err := topic.Subscribe()

	if err != nil {
		return nil, err
	}

	cr := &ChatRoom{
		ctx:      ctx,
		ps:       ps,
		topic:    topic,
		sub:      sub,
		self:     selfID,
		nick:     nickname,
		roomName: roomName,
		Messages: make(chan *ChatMessage, ChatRoomBufSize),
	}

	go cr.readLoop()

	return cr, nil
}

func (cr *ChatRoom) Publish(message string) error {

	m := ChatMessage{
		Message:    message,
		SenderID:   cr.self.String(),
		SenderNick: cr.nick,
	}

	msgBytes, err := json.Marshal(m)

	if err != nil {
		return err
	}

	return cr.topic.Publish(cr.ctx, msgBytes)
}

func (cr *ChatRoom) readLoop() {

	for {
		msg, err := cr.sub.Next(cr.ctx)

		if err != nil {
			close(cr.Messages)
			return
		}

		if msg.ReceivedFrom == cr.self {
			continue
		}

		cm := new(ChatMessage)

		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}
		cr.Messages <- cm
	}
}
