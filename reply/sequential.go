package reply

import (
	"fmt"
	"github.com/vitorsalgado/mocha/mock"
	"net/http"
)

type SequentialReply struct {
	replyOnNotFound mock.Reply
	replies         []mock.Reply
}

func Sequential() *SequentialReply {
	return &SequentialReply{replies: make([]mock.Reply, 0)}
}

func (mr *SequentialReply) ReplyOnSequenceEnded(reply mock.Reply) *SequentialReply {
	mr.replyOnNotFound = reply
	return mr
}

func (mr *SequentialReply) Add(reply ...mock.Reply) *SequentialReply {
	mr.replies = append(mr.replies, reply...)
	return mr
}

func (mr *SequentialReply) Then(reply mock.Reply) *SequentialReply {
	mr.Add(reply)
	return mr
}

func (mr *SequentialReply) Build(r *http.Request, m *mock.Mock) (*mock.Response, error) {
	size := len(mr.replies)
	if size == 0 {
		return nil,
			fmt.Errorf("you need to set at least one response when using multiple response builder")
	}

	var reply mock.Reply

	if m.Hits <= size {
		reply = mr.replies[m.Hits-1]
	}

	if reply == nil {
		if mr.replyOnNotFound != nil {
			return mr.replyOnNotFound.Build(r, m)
		}

		return nil,
			fmt.Errorf(
				"unable to obtain a response and no default response was set. request number: %d - sequence size: %d",
				m.Hits,
				size)
	}

	return reply.Build(r, m)
}