package validation

import "testing"

func TestValidateTag(t *testing.T) {
	color := "#2563EB"
	if err := ValidateTag("frontend", &color); err != nil {
		t.Fatalf("expected valid tag, got %v", err)
	}

	if err := ValidateTag("", nil); err == nil {
		t.Fatal("expected missing name error")
	}

	invalidColor := "blue"
	if err := ValidateTag("frontend", &invalidColor); err == nil {
		t.Fatal("expected invalid color error")
	}
}

func TestNormalizeTagColor(t *testing.T) {
	empty := "   "
	if got := NormalizeTagColor(&empty); got != nil {
		t.Fatalf("expected empty color to normalize to nil, got %v", *got)
	}

	value := "  #16A34A  "
	got := NormalizeTagColor(&value)
	if got == nil || *got != "#16A34A" {
		t.Fatalf("expected trimmed color, got %v", got)
	}
}
