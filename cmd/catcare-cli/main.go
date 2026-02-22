package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	core "github.com/wastingnotime/zeroapps/core/catcare"
	projection "github.com/wastingnotime/zeroapps/projection/catcare"
	"github.com/wastingnotime/zeroapps/store"
	svc "github.com/wastingnotime/zeroapps/svc/catcare"
)

func main() {
	var (
		commandName = flag.String("cmd", "", "command name: register|log-weight|list-registered")
		dbPath      = flag.String("db", "catcare.db", "sqlite database path")
		aggregateID = flag.String("aggregate-id", "", "aggregate id (cat id)")
		commandID   = flag.String("command-id", "", "command id (required)")
		expected    = flag.Int("expected-version", -1, "expected stream version (optional)")
		name        = flag.String("name", "", "cat name (register)")
		birthDate   = flag.String("birth-date", "", "birth date (register)")
		at          = flag.String("at", "", "timestamp (log-weight)")
		grams       = flag.Int("grams", 0, "grams (log-weight)")
		notes       = flag.String("notes", "", "notes (log-weight)")
	)
	flag.Parse()

	if *commandName == "" {
		usageAndExit()
	}

	registeredCats := projection.NewRegisteredCats()
	eventStore, err := store.NewSQLiteStore(*dbPath)
	if err != nil {
		fail(err)
	}
	defer func() {
		if err := eventStore.Close(); err != nil {
			fail(err)
		}
	}()

	if err := eventStore.Replay(context.Background(), registeredCats); err != nil {
		fail(err)
	}

	service := svc.NewService(eventStore, registeredCats)

	if *commandName == "list-registered" {
		cats := registeredCats.ListRegisteredCats()
		fmt.Printf("registered_cats=%d\n", len(cats))
		for _, cat := range cats {
			fmt.Printf("- cat_id=%s name=%s birth_date=%s\n", cat.CatID, cat.Name, cat.BirthDate)
		}
		return
	}

	if *commandID == "" {
		usageAndExit()
	}

	command, err := buildCommand(*commandName, *commandID, *name, *birthDate, *at, *grams, *notes)
	if err != nil {
		fail(err)
	}

	if *aggregateID == "" {
		if *commandName == "register" {
			*aggregateID = "cat-" + *commandID
		} else {
			fail(fmt.Errorf("aggregate-id is required"))
		}
	}

	var expectedVersion *int
	if *expected >= 0 {
		expectedVersion = expected
	}

	result, err := service.HandleCommand(context.Background(), svc.CommandEnvelope{
		AggregateID:     *aggregateID,
		Command:         command,
		ExpectedVersion: expectedVersion,
	})
	if err != nil {
		fail(err)
	}

	if !result.Ok {
		fmt.Printf("rejected: %s\n", result.Rejection.Error())
		os.Exit(2)
	}

	fmt.Printf("ok: version=%d events=%d\n", result.NewVersion, len(result.Events))
	for _, event := range result.Events {
		fmt.Printf("- %s\n", eventSummary(event))
	}
}

func buildCommand(name, commandID, catName, birthDate, at string, grams int, notes string) (core.Command, error) {
	switch name {
	case "register":
		return core.RegisterCat{
			CommandID: commandID,
			Name:      catName,
			BirthDate: birthDate,
		}, nil
	case "log-weight":
		return core.LogWeight{
			CommandID: commandID,
			At:        at,
			Grams:     grams,
			Notes:     notes,
		}, nil
	default:
		return nil, fmt.Errorf("unknown cmd %q", name)
	}
}

func eventSummary(event core.Event) string {
	switch ev := event.(type) {
	case core.CatRegistered:
		return fmt.Sprintf("CatRegistered cat_id=%s name=%s birth_date=%s", ev.CatID, ev.Name, ev.BirthDate)
	case core.WeightLogged:
		return fmt.Sprintf("WeightLogged entry_id=%s at=%s grams=%d", ev.EntryID, ev.At, ev.Grams)
	default:
		return fmt.Sprintf("%T", event)
	}
}

func usageAndExit() {
	fmt.Println("Usage:")
	fmt.Println("  catcare-cli -db ./catcare.db -cmd register -command-id cmd-1 -name Miso -birth-date 2023-01-01")
	fmt.Println("  catcare-cli -db ./catcare.db -cmd log-weight -aggregate-id cat-cmd-1 -command-id cmd-2 -at 2026-02-14T10:00:00Z -grams 4200")
	fmt.Println("  catcare-cli -db ./catcare.db -cmd list-registered")
	os.Exit(1)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
