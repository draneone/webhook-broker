package storage

import (
	"database/sql"
	"time"

	"github.com/imyousuf/webhook-broker/storage/data"
)

// ProducerDBRepository is the producer repository implementation for RDBMS
type ProducerDBRepository struct {
	db *sql.DB
}

// Store either creates or updates the producer information
func (repo *ProducerDBRepository) Store(producer *data.Producer) (*data.Producer, error) {
	inProducer, err := repo.Get(producer.ProducerID)
	if err != nil {
		return repo.insertProducer(producer)
	}
	if producer.Name != inProducer.Name || producer.Token != inProducer.Token {
		if !producer.IsInValidState() {
			return nil, ErrInvalidStateToSave
		}
		return repo.updateProducer(inProducer, producer.Name, producer.Token)
	}
	return inProducer, err
}

func (repo *ProducerDBRepository) updateProducer(producer *data.Producer, name, token string) (*data.Producer, error) {
	err := transactionalExec(repo.db, func() {
		producer.Name = name
		producer.Token = token
		producer.UpdatedAt = time.Now()
	}, "UPDATE producer SET name = $1, token = $2, updatedAt = $3 WHERE producerId = $4", func() []interface{} {
		return []interface{}{producer.Name, producer.Token, producer.UpdatedAt, producer.ProducerID}
	})
	return producer, err
}

func (repo *ProducerDBRepository) insertProducer(producer *data.Producer) (*data.Producer, error) {
	producer.QuickFix()
	if !producer.IsInValidState() {
		return nil, ErrInvalidStateToSave
	}
	err := transactionalExec(repo.db, func() {}, "INSERT INTO producer (id, producerId, name, token, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6)", func() []interface{} {
		return []interface{}{producer.ID, producer.ProducerID, producer.Name, producer.Token, producer.CreatedAt, producer.UpdatedAt}
	})
	return producer, err
}

// Get retrieves the producer with matching producer id
func (repo *ProducerDBRepository) Get(producerID string) (*data.Producer, error) {
	producer := &data.Producer{}
	err := querySingleRow(repo.db, "SELECT ID, producerID, name, token, createdAt, updatedAt FROM producer WHERE producerID like $1", func() []interface{} {
		return []interface{}{producerID}
	}, func() []interface{} {
		return []interface{}{&producer.ID, &producer.ProducerID, &producer.Name, &producer.Token, &producer.CreatedAt, &producer.UpdatedAt}
	})
	return producer, err
}

// GetList retrieves the list of producer based on pagination params supplied. It will return a error if both after and before is present at the same time
func (repo *ProducerDBRepository) GetList(page *data.Pagination) ([]*data.Producer, *data.Pagination, error) {
	if page == nil || (page.Next != nil && page.Previous != nil) {
		return nil, nil, ErrPaginationDeadlock
	}
	producers := make([]*data.Producer, 0)
	pagination := &data.Pagination{}
	baseQuery := "SELECT ID, producerID, name, token, createdAt, updatedAt FROM producer" + getPaginationQueryFragment(page, false)
	scanArgs := func() []interface{} {
		producer := &data.Producer{}
		producers = append(producers, producer)
		return []interface{}{&producer.ID, &producer.ProducerID, &producer.Name, &producer.Token, &producer.CreatedAt, &producer.UpdatedAt}
	}
	err := queryRows(repo.db, baseQuery, nilArgs, scanArgs)
	if err == nil {
		producerCount := len(producers)
		if producerCount > 0 {
			pagination = data.NewPagination(producers[producerCount-1], producers[0])
		}
	}
	return producers, pagination, err
}

// NewProducerRepository returns a new producer repository
func NewProducerRepository(db *sql.DB) ProducerRepository {
	panicIfNoDBConnectionPool(db)
	return &ProducerDBRepository{db: db}
}
