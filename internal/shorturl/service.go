package shorturl

import "context"

// Repository defines the persistence operations used by the short URL service.
type Repository interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string) error
	Delete(context.Context, string) error
	GetAllKeyValues(context.Context) (map[string]string, error)
}

// Service composes short URL behavior with a persistence implementation.
type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Get(ctx context.Context, key string) (string, error) {
	return s.repository.Get(ctx, key)
}

func (s *Service) Set(ctx context.Context, key, value string) error {
	return s.repository.Set(ctx, key, value)
}

func (s *Service) Delete(ctx context.Context, key string) error {
	return s.repository.Delete(ctx, key)
}

func (s *Service) GetAllKeyValues(ctx context.Context) (map[string]string, error) {
	return s.repository.GetAllKeyValues(ctx)
}
