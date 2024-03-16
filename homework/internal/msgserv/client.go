package msgserv

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
	"homework/internal/model"
	"log"
	"time"
)

type client struct {
	s        *MsgServ
	conn     *websocket.Conn
	out      chan model.Message
	nickname string
	ctx      context.Context
	close    context.CancelFunc
}

func newClient(s *MsgServ, conn *websocket.Conn) *client {
	return &client{
		s:    s,
		conn: conn,
		out:  make(chan model.Message, 64),
	}
}

func (c *client) run(ctx context.Context) {
	ctx, cl := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	c.ctx = ctx
	c.close = cl

	eg.Go(c.inFunc)
	eg.Go(c.outFunc)

	err := eg.Wait()
	if err != nil {
		log.Println(err)
	}
	c.s.closedClients <- c
}

func (c *client) outFunc() error {
	defer close(c.out)
	for {
		select {
		case msg := <-c.out:
			bytes := []byte(fmt.Sprintf("%s: %s", msg.Nickname, msg.Message))
			err := c.conn.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				return err
			}
		case <-c.ctx.Done():
			closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
			err := c.conn.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second))
			if err != nil {
				return err
			}
			err = c.conn.Close()
			if err != nil {
				return err
			}
			return c.ctx.Err()
		}
	}
}

func (c *client) inFunc() error {
	var err error
	nickname, err := c.recvString(c.ctx)
	if err != nil {
		return err
	}
	log.Println(nickname)
	password, err := c.recvBytes(c.ctx)
	if err != nil {
		return err
	}
	log.Println("got password")
	err = c.s.authUser(c.ctx, nickname, password)
	if err != nil {
		return err
	}
	c.nickname = nickname
	for {
		msgStr, err := c.recvString(c.ctx)
		if err != nil {
			return err
		}
		msg := model.Message{
			SentTime: time.Now(),
			Nickname: c.nickname,
			Message:  msgStr,
		}
		select {
		case c.s.in <- msg:
		default:
			return errors.New("server failed")
		}
	}
}

func (c *client) recvString(ctx context.Context) (string, error) {
	bytes, err := c.recvBytes(ctx)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *client) recvBytes(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		mt, bytes, err := c.conn.ReadMessage()
		if err != nil {
			return nil, err
		}
		if mt == websocket.CloseMessage {
			return nil, errors.New("closed")
		}
		return bytes, nil
	}
}
