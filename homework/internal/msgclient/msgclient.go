package msgclient

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
	"net/url"
	"os"
	"strings"
	"time"
)

type MsgClient struct {
	conn *websocket.Conn
	send chan string
	recv chan string
	ctx  context.Context
}

func NewMsgClient() *MsgClient {
	return &MsgClient{
		send: make(chan string),
		recv: make(chan string),
	}
}

func (c *MsgClient) Run(ctx context.Context, host string) error {
	u := url.URL{
		Scheme: "ws",
		Host:   host,
		Path:   "/",
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	c.conn = conn
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	c.ctx = ctx

	eg.Go(c.sendFunc)
	eg.Go(c.recvFunc)
	eg.Go(c.readFunc)
	eg.Go(c.writeFunc)

	return eg.Wait()
}

func (c *MsgClient) sendFunc() error {
	for {
		select {
		case message := <-c.send:
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
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

func (c *MsgClient) readFunc() error {
	fmt.Println("!!! First enter your nickname, then enter your password, then enter messages, one on each line !!!")
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		c.send <- line
	}
}

func (c *MsgClient) recvFunc() error {
	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		default:
			mt, message, err := c.conn.ReadMessage()
			if err != nil {
				return err
			}
			if mt == websocket.CloseMessage {
				return errors.New("closed")
			}
			c.recv <- string(message)
		}
	}
}

func (c *MsgClient) writeFunc() error {
	for {
		select {
		case line := <-c.recv:
			fmt.Println(line)
		case <-c.ctx.Done():
			return c.ctx.Err()
		}
	}
}
