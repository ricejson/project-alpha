package repositories

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backend/internal/dto"
	"backend/internal/models"
)

func TestTicketRepositoryApplyListFiltersBuildsSearchAndStatusQuery(t *testing.T) {
	db := newDryRunDB(t)
	repo := NewTicketRepository(db)
	completed := true

	var tickets []models.Ticket
	tx := repo.applyListFilters(db.Model(&models.Ticket{}), dto.ListTicketsFilter{
		Title:     "list",
		Completed: &completed,
	}).Find(&tickets)

	sql := tx.Statement.SQL.String()
	assertSQLContains(t, sql, "tickets.title ILIKE")
	assertSQLContains(t, sql, "tickets.completed =")
}

func TestTicketRepositoryApplyListFiltersBuildsMultiTagAndQuery(t *testing.T) {
	db := newDryRunDB(t)
	repo := NewTicketRepository(db)

	var tickets []models.Ticket
	tx := repo.applyListFilters(db.Model(&models.Ticket{}), dto.ListTicketsFilter{
		TagIDs: []uint{1, 2},
	}).Find(&tickets)

	sql := tx.Statement.SQL.String()
	assertSQLContains(t, sql, "JOIN ticket_tags ON ticket_tags.ticket_id = tickets.id")
	assertSQLContains(t, sql, "ticket_tags.tag_id IN")
	assertSQLContains(t, sql, "GROUP BY")
	assertSQLContains(t, sql, "tickets.id")
	assertSQLContains(t, sql, "COUNT(DISTINCT ticket_tags.tag_id) =")
}

func TestNewPaginationResponseComputesTotalPages(t *testing.T) {
	pagination := dto.NewPaginationResponse(2, 20, 41)
	if pagination.TotalPages != 3 {
		t.Fatalf("expected 3 total pages, got %d", pagination.TotalPages)
	}
}

func newDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: fakeConnPool{}}), &gorm.Config{
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("open dry run db: %v", err)
	}
	return db
}

func assertSQLContains(t *testing.T, sql string, want string) {
	t.Helper()

	if !strings.Contains(sql, want) {
		t.Fatalf("expected SQL to contain %q, got %s", want, sql)
	}
}

type fakeConnPool struct{}

func (fakeConnPool) PrepareContext(context.Context, string) (*sql.Stmt, error) {
	return nil, nil
}

func (fakeConnPool) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (fakeConnPool) QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (fakeConnPool) QueryRowContext(context.Context, string, ...interface{}) *sql.Row {
	return nil
}
