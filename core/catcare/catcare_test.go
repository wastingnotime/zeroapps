package catcare

import "testing"

func TestRegisterCatGivenEmptyStreamWhenRegisterCatThenEmitsCatRegistered(t *testing.T) {
	aggregate, err := LoadFrom(nil)
	if err != nil {
		t.Fatalf("load aggregate: %v", err)
	}

	events, err := aggregate.Decide(RegisterCat{
		CommandID: "cmd-1",
		Name:      "Miso",
		BirthDate: "2023-01-01",
	})
	if err != nil {
		t.Fatalf("decide register cat: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event, ok := events[0].(CatRegistered)
	if !ok {
		t.Fatalf("expected CatRegistered, got %T", events[0])
	}
	if event.CatID != "cat-cmd-1" {
		t.Fatalf("expected deterministic cat id, got %q", event.CatID)
	}
}

func TestRegisterCatGivenRegisteredCatWhenRegisterCatThenRejectsDuplicate(t *testing.T) {
	aggregate, err := LoadFrom([]Event{
		CatRegistered{
			CommandID: "cmd-registered",
			CatID:     "cat-cmd-registered",
			Name:      "Miso",
		},
	})
	if err != nil {
		t.Fatalf("load aggregate: %v", err)
	}

	_, err = aggregate.Decide(RegisterCat{
		CommandID: "cmd-2",
		Name:      "Taro",
	})
	if err == nil {
		t.Fatal("expected rejection")
	}
	rejection, ok := err.(Rejection)
	if !ok {
		t.Fatalf("expected Rejection, got %T", err)
	}
	if rejection.Code != CodeAlreadyRegistered {
		t.Fatalf("expected %q, got %q", CodeAlreadyRegistered, rejection.Code)
	}
}

func TestLogWeightGivenRegisteredCatWhenLogWeightWithInvalidGramsThenRejects(t *testing.T) {
	aggregate, err := LoadFrom([]Event{
		CatRegistered{
			CommandID: "cmd-register",
			CatID:     "cat-cmd-register",
			Name:      "Miso",
		},
	})
	if err != nil {
		t.Fatalf("load aggregate: %v", err)
	}

	cases := []struct {
		name string
		cmd  LogWeight
		code string
	}{
		{
			name: "non-positive grams",
			cmd: LogWeight{
				CommandID: "cmd-weight-0",
				At:        "2026-02-14T10:00:00Z",
				Grams:     0,
			},
			code: CodeInvalidWeight,
		},
		{
			name: "absurdly high grams",
			cmd: LogWeight{
				CommandID: "cmd-weight-high",
				At:        "2026-02-14T10:00:00Z",
				Grams:     MaxWeightGrams + 1,
			},
			code: CodeAbsurdWeight,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := aggregate.Decide(tc.cmd)
			if err == nil {
				t.Fatal("expected rejection")
			}
			rejection, ok := err.(Rejection)
			if !ok {
				t.Fatalf("expected Rejection, got %T", err)
			}
			if rejection.Code != tc.code {
				t.Fatalf("expected %q, got %q", tc.code, rejection.Code)
			}
		})
	}
}

func TestLogWeightGivenRegisteredCatWhenLogWeightValidThenEmitsWeightLogged(t *testing.T) {
	aggregate, err := LoadFrom([]Event{
		CatRegistered{
			CommandID: "cmd-register",
			CatID:     "cat-cmd-register",
			Name:      "Miso",
		},
	})
	if err != nil {
		t.Fatalf("load aggregate: %v", err)
	}

	events, err := aggregate.Decide(LogWeight{
		CommandID: "cmd-weight-1",
		At:        "2026-02-14T10:00:00Z",
		Grams:     4200,
		Notes:     "post breakfast",
	})
	if err != nil {
		t.Fatalf("decide log weight: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event, ok := events[0].(WeightLogged)
	if !ok {
		t.Fatalf("expected WeightLogged, got %T", events[0])
	}
	if event.EntryID != "weight-cmd-weight-1" {
		t.Fatalf("expected deterministic entry id, got %q", event.EntryID)
	}
}

func TestDuplicateCommandIDGivenAppliedCommandWhenDecideThenRejects(t *testing.T) {
	aggregate, err := LoadFrom([]Event{
		CatRegistered{
			CommandID: "cmd-register",
			CatID:     "cat-cmd-register",
			Name:      "Miso",
		},
	})
	if err != nil {
		t.Fatalf("load aggregate: %v", err)
	}

	_, err = aggregate.Decide(LogWeight{
		CommandID: "cmd-register",
		At:        "2026-02-14T10:00:00Z",
		Grams:     4000,
	})
	if err == nil {
		t.Fatal("expected rejection")
	}
	rejection, ok := err.(Rejection)
	if !ok {
		t.Fatalf("expected Rejection, got %T", err)
	}
	if rejection.Code != CodeDuplicateCommand {
		t.Fatalf("expected %q, got %q", CodeDuplicateCommand, rejection.Code)
	}
}
