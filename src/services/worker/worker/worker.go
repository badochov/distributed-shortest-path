package worker

import "gorm.io/gorm"

type Deps struct {
	Db *gorm.DB
}

// Worker All methods from link service and worker service should end up calling this interface.
type Worker interface {
	AssignSegment(id int32) error
	Add(a int32, b int32) (int32, error) // Example
}

type worker struct {
	db *gorm.DB
}

func (w *worker) Add(a int32, b int32) (int32, error) {
	return a + b, nil
}

func (w *worker) AssignSegment(id int32) error {
	//TODO implement me
	panic("implement me")
}

func New(deps Deps) Worker {
	return &worker{
		db: deps.Db,
	}
}
