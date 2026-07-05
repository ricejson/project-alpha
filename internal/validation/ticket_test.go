package validation

import (
	"testing"

	"backend/internal/dto"
)

func TestParseListTicketsFilter(t *testing.T) {
	completed := false
	filter, err := ParseListTicketsFilter(dto.ListTicketsQuery{
		Title:     "  Search Me  ",
		TagIDsRaw: "1, 2,2,3",
		Completed: &completed,
		Page:      -1,
		PageSize:  500,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if filter.Title != "Search Me" {
		t.Fatalf("expected trimmed title, got %q", filter.Title)
	}
	if filter.Page != 1 {
		t.Fatalf("expected default page 1, got %d", filter.Page)
	}
	if filter.PageSize != MaxPageSize {
		t.Fatalf("expected page size capped to %d, got %d", MaxPageSize, filter.PageSize)
	}
	if filter.Completed == nil || *filter.Completed != false {
		t.Fatalf("expected completed=false filter")
	}
	if got, want := len(filter.TagIDs), 3; got != want {
		t.Fatalf("expected %d unique tag IDs, got %d", want, got)
	}
}

func TestParseListTicketsFilterRejectsInvalidTagIDs(t *testing.T) {
	_, err := ParseListTicketsFilter(dto.ListTicketsQuery{TagIDsRaw: "1,abc"})
	if err == nil {
		t.Fatal("expected invalid tagIds error")
	}
}

func TestValidateTicket(t *testing.T) {
	if err := ValidateTicket("title", "description", []uint{1, 2}); err != nil {
		t.Fatalf("expected valid ticket, got %v", err)
	}

	if err := ValidateTicket("", "", nil); err == nil {
		t.Fatal("expected missing title error")
	}

	if err := ValidateTicket("title", "", []uint{0}); err == nil {
		t.Fatal("expected invalid tag id error")
	}
}
