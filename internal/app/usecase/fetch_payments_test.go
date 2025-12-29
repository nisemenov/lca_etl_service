package usecase

import (
	"context"
	"testing"

	"github.com/nisemenov/etl_service/internal/model"
	"github.com/stretchr/testify/require"
)

type mockPaymentFetcher struct{}

func newMockPaymentFetcher() *mockPaymentFetcher {
	return &mockPaymentFetcher{}
}

func (m *mockPaymentFetcher) Fetch(ctx context.Context) ([]model.Payment, error) {
	return []model.Payment{
		{ID: 1},
		{ID: 2},
	}, nil
}

type mockPaymentRepository struct {
	saved []model.Payment
}

func newMockPaymentRepo() *mockPaymentRepository {
	return &mockPaymentRepository{}
}

func (m *mockPaymentRepository) SaveBatch(ctx context.Context, payments []model.Payment) error {
	m.saved = append(m.saved, payments...)
	return nil
}

func TestFetchPaymentsUseCase(t *testing.T) {
	ctx := context.Background()

	fetcher := newMockPaymentFetcher()
	repo := newMockPaymentRepo()

	uc := NewFetchPaymentsUseCase(fetcher, repo)

	err := uc.Execute(ctx)
	require.NoError(t, err)

	require.Len(t, repo.saved, 2)
	require.Equal(t, model.StatusPending, repo.saved[0].Status)
}
