package main

type ChatMap_t struct {
	user_topic map[int64]int
	topic_user map[int]int64
}

func NewChatMap() *ChatMap_t {
	return &ChatMap_t{
		user_topic: make(map[int64]int),
		topic_user: make(map[int]int64),
	}
}

func (e *ChatMap_t) Pair(chat_id int64, msg_id int) {
	e.user_topic[chat_id] = msg_id
	e.topic_user[msg_id] = chat_id
}

func (e *ChatMap_t) UnPair(chat_id int64, msg_id int) {
	delete(e.user_topic, chat_id)
	delete(e.topic_user, msg_id)
}

func (e *ChatMap_t) GetMsgTopic(chat_id int64) int {
	r, _ := e.user_topic[chat_id]
	return r
}

func (e *ChatMap_t) GetUserChat(msg_id int) int64 {
	r, _ := e.topic_user[msg_id]
	return r
}
