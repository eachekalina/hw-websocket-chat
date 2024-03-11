package msgserv

import (
	"context"
	"github.com/gorilla/websocket"
	"homework/internal/model"
	"homework/internal/repository"
	"log"
	"net/http"
	"slices"
)

type MsgServ interface {
	Serve(ctx context.Context, addr string) error
}

type serv struct {
	clients       map[*client]bool
	in            chan model.Message
	newClients    chan *client
	closedClients chan *client
	ctx           context.Context
	repo          repository.MessageRepository
}

func NewMsgServ(repo repository.MessageRepository) MsgServ {
	return &serv{
		clients:       make(map[*client]bool),
		in:            make(chan model.Message, 64),
		newClients:    make(chan *client),
		closedClients: make(chan *client),
		repo:          repo,
	}
}

func (s *serv) Serve(ctx context.Context, addr string) error {
	s.ctx = ctx
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)
	httpServ := &http.Server{
		Handler: mux,
		Addr:    addr,
	}
	go func() {
		err := httpServ.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()
	s.run()
	return httpServ.Shutdown(context.Background())
}

func (s *serv) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case c := <-s.newClients:
			s.clients[c] = true
			go c.run(s.ctx)
			msgs, err := s.repo.GetLastMessages(s.ctx, 10)
			if err != nil {
				log.Println(err)
				return
			}
			slices.Reverse(msgs)
			for _, msg := range msgs {
				select {
				case c.out <- msg:
				default:
					c.close()
				}
			}
		case c := <-s.closedClients:
			delete(s.clients, c)
		case msg := <-s.in:
			err := s.repo.AddMessage(context.Background(), msg)
			if err != nil {
				log.Println(err)
				return
			}
			for c := range s.clients {
				select {
				case c.out <- msg:
				default:
					c.close()
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{}

func (s *serv) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := newClient(s, conn)
	s.newClients <- c
}
