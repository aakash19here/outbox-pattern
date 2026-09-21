package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aakash19here/outbox-pattern/internal/kafka/producer"
	repo "github.com/aakash19here/outbox-pattern/internal/repos"
	reposhared "github.com/aakash19here/outbox-pattern/internal/repos/repo-shared"
	"github.com/jmoiron/sqlx"
)

type OutboxService struct {
	eventRepo *repo.EventRepo
	producer  *producer.KafkaProducer
}

func NewOutbox(er *repo.EventRepo) *OutboxService {
	pr := producer.NewKafkaProducer()
	outbox := &OutboxService{
		eventRepo: er,
		producer:  pr,
	}

	go outbox.outboxLoop()
	return outbox
}

func (s *OutboxService) outboxLoop() {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	for range ticker.C {
		s.handlePending()
	}
}

func (s *OutboxService) handlePending() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	txrepo := s.eventRepo.GetRepo()
	_, err := reposhared.TxClosure(ctx, txrepo, func(ctx context.Context, tx *sqlx.Tx) (int, error) {
		events, err := s.eventRepo.GetAllPending(ctx, tx)
		if err != nil {
			return 0, err
		}

		toUpdateIds := []string{}
		for _, e := range events {
			// TODO -> add batch processing
			// TODO -> add binary protocol
			b, err := json.Marshal(e)
			if err != nil {
				fmt.Printf("error marshalling event %v\n", err)
				continue
			}
			err = s.producer.Produce(b)
			if err != nil {
				fmt.Printf("error producing event %v\n", err)
				// TODO -> add err status && save err
				continue
			}
			toUpdateIds = append(toUpdateIds, e.EventId)
		}

		rows, err := s.eventRepo.UpdateStatusByIds(ctx, tx, toUpdateIds, repo.EventStatus_Produced)
		if err != nil {
			return 0, err
		}
		if rows != len(toUpdateIds) {
			errMsg := "updated row count didn't match"
			fmt.Println(errMsg)
			return 0, fmt.Errorf("%v", errMsg)
		}
		return rows, nil
	})

	if err != nil {
		fmt.Printf("err outbox service %v\n", err)
	}
}
