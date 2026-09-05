package lookup

import (
	"testing"

	"github.com/abhigyan-chatterjee/aurora/internal/resolve"
)

func TestResolveOne_Official(t *testing.T) {
	res, err := ResolveOne("bash")
	if err != nil {
		t.Fatalf("unexpected error resolving bash: %v", err)
	}
	if res == nil {
		t.Fatalf("expected non-nil result for bash")
	}
	if res.ChosenSource != resolve.SourceOfficial {
		t.Errorf("expected SourceOfficial, got %v", res.ChosenSource)
	}
	if res.PacmanResult == nil || res.PacmanResult.Name != "bash" {
		t.Errorf("expected PacmanResult with name bash, got %+v", res.PacmanResult)
	}
}

func TestResolveOne_NotFound(t *testing.T) {
	const nonexistent = "definitely_nonexistent_pkg_xyz987654321"
	res, err := ResolveOne(nonexistent)
	if err != nil {
		t.Fatalf("expected nil error for nonexistent package, got: %v", err)
	}
	if res == nil {
		t.Fatalf("expected non-nil ResolvedPackage struct for nonexistent package")
	}
	if res.ChosenSource != resolve.SourceUnknown {
		t.Errorf("expected SourceUnknown, got %v", res.ChosenSource)
	}
	if res.PacmanResult != nil {
		t.Errorf("expected nil PacmanResult, got %+v", res.PacmanResult)
	}
	if res.AURResult != nil {
		t.Errorf("expected nil AURResult, got %+v", res.AURResult)
	}
	if res.Query != nonexistent {
		t.Errorf("expected Query %q, got %q", nonexistent, res.Query)
	}
}

func TestResolveBatch_OrderPreserved(t *testing.T) {
	queries := []string{
		"bash",
		"definitely_nonexistent_pkg_1",
		"coreutils",
	}

	results := ResolveBatch(queries)
	if len(results) != len(queries) {
		t.Fatalf("expected %d results, got %d", len(queries), len(results))
	}

	for i, q := range queries {
		if results[i] == nil {
			t.Fatalf("result at index %d is nil", i)
		}
		if results[i].Query != q {
			t.Errorf("expected query at index %d to be %q, got %q", i, q, results[i].Query)
		}
	}

	if results[0].ChosenSource != resolve.SourceOfficial {
		t.Errorf("expected bash to be SourceOfficial, got %v", results[0].ChosenSource)
	}
	if results[1].ChosenSource != resolve.SourceUnknown {
		t.Errorf("expected nonexistent to be SourceUnknown, got %v", results[1].ChosenSource)
	}
	if results[2].ChosenSource != resolve.SourceOfficial {
		t.Errorf("expected coreutils to be SourceOfficial, got %v", results[2].ChosenSource)
	}
}
