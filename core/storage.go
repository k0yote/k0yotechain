package core

type Storage interface {
	Put(*Block) error
}

type MemoryStore struct{}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Put(block *Block) error {
	return nil
}
