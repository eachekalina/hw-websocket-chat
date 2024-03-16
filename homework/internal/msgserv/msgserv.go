package msgserv

import (
	"context"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"homework/internal/model"
	"log"
	"net/http"
	"slices"
)

type MessageRepository interface {
	AddMessage(ctx context.Context, message model.Message) error
	GetLastMessages(ctx context.Context, n int) ([]model.Message, error)
}

type UserRepository interface {
	AddUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, nickname string) (model.User, error)
}

type MsgServ struct {
	clients       map[*client]struct{}
	in            chan model.Message
	newClients    chan *client
	closedClients chan *client
	ctx           context.Context
	msgRepo       MessageRepository
	userRepo      UserRepository
	upgrader      websocket.Upgrader
}

func NewMsgServ(msgRepo MessageRepository, userRepo UserRepository) *MsgServ {
	return &MsgServ{
		clients:       make(map[*client]struct{}),
		in:            make(chan model.Message, 64),
		newClients:    make(chan *client),
		closedClients: make(chan *client),
		msgRepo:       msgRepo,
		userRepo:      userRepo,
		upgrader:      websocket.Upgrader{},
	}
}

func (s *MsgServ) Serve(ctx context.Context, addr string) error {
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

func (s *MsgServ) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case c := <-s.newClients:
			s.clients[c] = struct{}{}
			go c.run(s.ctx)
			msgs, err := s.msgRepo.GetLastMessages(s.ctx, 10)
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
			err := s.msgRepo.AddMessage(context.Background(), msg)
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

func (s *MsgServ) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := newClient(s, conn)
	s.newClients <- c
}

func (s *MsgServ) authUser(ctx context.Context, nickname string, password []byte) error {
	dbUser, err := s.userRepo.GetUser(ctx, nickname)
	log.Println(err)
	if err != nil {
		passwordHash, err := bcrypt.GenerateFromPassword(password, 10)
		if err != nil {
			return err
		}
		return s.userRepo.AddUser(ctx, model.User{
			Nickname:     nickname,
			PasswordHash: passwordHash,
		})
	}

	return bcrypt.CompareHashAndPassword(dbUser.PasswordHash, password)
}
