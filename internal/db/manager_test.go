package db

import (
	"context"
	"testing"
	"time"

	"github.com/hydrocode-de/datailama/internal/sql"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

type mockManager struct {
	*sql.Queries
	mock pgxmock.PgxPoolIface
}

func (m *mockManager) Close() {
	m.mock.Close()
}

func TestSearchPaperByTitle(t *testing.T) {
	// Create mock database
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	// Create test data
	testTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	testData := []sql.SearchPaperByTitleRow{
		{
			ID:      1,
			Title:   "Paper A",
			Doi:     "10.1234/A",
			Journal: "Journal X",
			Published: pgtype.Date{
				Time:  testTime,
				Valid: true,
			},
		},
		{
			ID:      2,
			Title:   "Paper B",
			Doi:     "10.1234/B",
			Journal: "Journal Y",
			Published: pgtype.Date{
				Time:  testTime.AddDate(0, 1, 0),
				Valid: true,
			},
		},
	}

	t.Run("ascending order", func(t *testing.T) {
		// Set up expectations for ascending order
		rows := mock.NewRows([]string{"match", "cosine_distance", "id", "title", "doi", "url", "journal", "author", "published", "citations", "citations_year"}).
			AddRow("", int32(0), testData[0].ID, testData[0].Title, testData[0].Doi, testData[0].Url, testData[0].Journal,
				testData[0].Author, testData[0].Published, testData[0].Citations, testData[0].CitationsYear).
			AddRow("", int32(0), testData[1].ID, testData[1].Title, testData[1].Doi, testData[1].Url, testData[1].Journal,
				testData[1].Author, testData[1].Published, testData[1].Citations, testData[1].CitationsYear)

		mock.ExpectQuery(`-- name: SearchPaperByTitle :many SELECT`).
			WithArgs(int32(10), "test", "", "published", "asc").
			WillReturnRows(rows)

		// Create manager with mock db
		manager := &mockManager{
			Queries: sql.New(mock),
			//mock:    mock,
		}

		// Test ascending order
		results, err := manager.SearchPaperByTitle(context.Background(), sql.SearchPaperByTitleParams{
			Limit:     10,
			Title:     "test",
			OrderBy:   "published",
			Direction: "asc",
		})

		assert.NoError(t, err)
		if assert.Len(t, results, 2) {
			assert.Equal(t, testData[0].Title, results[0].Title)
			assert.Equal(t, testData[1].Title, results[1].Title)
		}
	})

	t.Run("descending order", func(t *testing.T) {
		// Set up expectations for descending order
		rows := mock.NewRows([]string{"match", "cosine_distance", "id", "title", "doi", "url", "journal", "author", "published", "citations", "citations_year"}).
			AddRow("", int32(0), testData[1].ID, testData[1].Title, testData[1].Doi, testData[1].Url, testData[1].Journal,
				testData[1].Author, testData[1].Published, testData[1].Citations, testData[1].CitationsYear).
			AddRow("", int32(0), testData[0].ID, testData[0].Title, testData[0].Doi, testData[0].Url, testData[0].Journal,
				testData[0].Author, testData[0].Published, testData[0].Citations, testData[0].CitationsYear)

		mock.ExpectQuery(`-- name: SearchPaperByTitle :many SELECT`).
			WithArgs(int32(10), "test", "", "published", "desc").
			WillReturnRows(rows)

		// Create manager with mock db
		manager := &mockManager{
			Queries: sql.New(mock),
			// mock:    mock,
		}

		// Test descending order
		results, err := manager.SearchPaperByTitle(context.Background(), sql.SearchPaperByTitleParams{
			Limit:     10,
			Title:     "test",
			OrderBy:   "published",
			Direction: "desc",
		})

		assert.NoError(t, err)
		if assert.Len(t, results, 2) {
			assert.Equal(t, testData[1].Title, results[0].Title)
			assert.Equal(t, testData[0].Title, results[1].Title)
		}
	})
}
