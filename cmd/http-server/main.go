package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	dbpostgres "github.com/aakash19here/outbox-pattern/internal/db/postgres"
	repo "github.com/aakash19here/outbox-pattern/internal/repos"
	"github.com/aakash19here/outbox-pattern/internal/service"
	pkgtypes "github.com/aakash19here/outbox-pattern/pkg/types"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	addr           string
	msgCH          chan *pkgtypes.Message[*repo.Event]
	paymentService *service.PaymentService
	outboxService  *service.OutboxService
	exitCH         chan struct{}
}

func NewServer(addr string, db *sqlx.DB) *Server {
	eventRepo := repo.NewEventRepo(db)
	pr := repo.NewPaymentRepo(db)
	ps := service.NewPaymentService(pr, eventRepo)
	os := service.NewOutbox(eventRepo)
	return &Server{
		msgCH:  make(chan *pkgtypes.Message[*repo.Event], 64),
		exitCH: make(chan struct{}),
		addr:   addr,

		paymentService: ps,
		outboxService:  os,
	}
}

func (s *Server) handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	_, err := s.createPayment()
	if err != nil {
		fmt.Printf("err creating PMNT %v\n", err)
	}
}

func (s *Server) createPayment() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	amount := rand.Intn(100_000)
	orderN := strconv.Itoa(amount)
	paym := repo.NewPayment(orderN, amount, "created")
	return s.paymentService.Save(ctx, paym)
}

func main() {
	dbOpts := new(dbpostgres.DBPostgresOptions)
	db, err := dbpostgres.NewDBConn(dbOpts)
	if err != nil {
		panic(fmt.Sprintf("unable to conn to db, err = %v\n", err))
	}
	defer db.Close()

	s := NewServer(":7576", db)
	mux := http.NewServeMux()
	mux.HandleFunc("/payment", s.handleCreatePayment)

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			s.createPayment()
		}
	}()

	log.Fatal(http.ListenAndServe(s.addr, nil))
}
