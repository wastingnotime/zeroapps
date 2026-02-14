package catcare

import "strings"

const (
	MinWeightGrams = 100
	MaxWeightGrams = 30000
)

const (
	CodeAlreadyRegistered = "already_registered"
	CodeNotRegistered     = "not_registered"
	CodeDuplicateCommand  = "duplicate_command"
	CodeInvalidCommand    = "invalid_command"
	CodeInvalidWeight     = "invalid_weight"
	CodeAbsurdWeight      = "absurd_weight"
	CodeInvalidName       = "invalid_name"
	CodeInvalidCommandID  = "invalid_command_id"
	CodeInvalidDate       = "invalid_date"
)

type Rejection struct {
	Code    string
	Message string
	Field   string
}

func (r Rejection) Error() string {
	if r.Field == "" {
		return r.Code + ": " + r.Message
	}
	return r.Code + ": " + r.Field + " " + r.Message
}

type Command interface {
	commandName() string
	commandID() string
}

type Event interface {
	eventName() string
	commandID() string
}

type RegisterCat struct {
	CommandID string
	Name      string
	BirthDate string
}

func (c RegisterCat) commandName() string { return "RegisterCat" }
func (c RegisterCat) commandID() string   { return c.CommandID }

type LogWeight struct {
	CommandID string
	At        string
	Grams     int
	Notes     string
}

func (c LogWeight) commandName() string { return "LogWeight" }
func (c LogWeight) commandID() string   { return c.CommandID }

type CatRegistered struct {
	CommandID string
	CatID     string
	Name      string
	BirthDate string
}

func (e CatRegistered) eventName() string { return "CatRegistered" }
func (e CatRegistered) commandID() string { return e.CommandID }

type WeightLogged struct {
	CommandID string
	EntryID   string
	At        string
	Grams     int
	Notes     string
}

func (e WeightLogged) eventName() string { return "WeightLogged" }
func (e WeightLogged) commandID() string { return e.CommandID }

type CatCare struct {
	CatID               string
	Name                string
	BirthDate           string
	Registered          bool
	WeightEntries       []WeightLogged
	processedCommandIDs map[string]struct{}
}

func New() *CatCare {
	return &CatCare{
		processedCommandIDs: map[string]struct{}{},
	}
}

func LoadFrom(events []Event) (*CatCare, error) {
	aggregate := New()
	for _, event := range events {
		if err := aggregate.Apply(event); err != nil {
			return nil, err
		}
	}
	return aggregate, nil
}

func (a *CatCare) Decide(command Command) ([]Event, error) {
	if command == nil {
		return nil, Rejection{Code: CodeInvalidCommand, Message: "command is required"}
	}
	if command.commandID() == "" {
		return nil, Rejection{Code: CodeInvalidCommandID, Message: "must not be empty", Field: "command_id"}
	}
	if _, exists := a.processedCommandIDs[command.commandID()]; exists {
		return nil, Rejection{Code: CodeDuplicateCommand, Message: "already applied", Field: "command_id"}
	}

	switch cmd := command.(type) {
	case RegisterCat:
		return a.decideRegisterCat(cmd)
	case LogWeight:
		return a.decideLogWeight(cmd)
	default:
		return nil, Rejection{Code: CodeInvalidCommand, Message: "unknown command"}
	}
}

func (a *CatCare) Apply(event Event) error {
	switch ev := event.(type) {
	case CatRegistered:
		a.CatID = ev.CatID
		a.Name = ev.Name
		a.BirthDate = ev.BirthDate
		a.Registered = true
		if ev.CommandID != "" {
			a.processedCommandIDs[ev.CommandID] = struct{}{}
		}
		return nil
	case WeightLogged:
		a.WeightEntries = append(a.WeightEntries, ev)
		if ev.CommandID != "" {
			a.processedCommandIDs[ev.CommandID] = struct{}{}
		}
		return nil
	default:
		return Rejection{Code: CodeInvalidCommand, Message: "unknown event"}
	}
}

func (a *CatCare) decideRegisterCat(cmd RegisterCat) ([]Event, error) {
	if a.Registered {
		return nil, Rejection{Code: CodeAlreadyRegistered, Message: "cat already registered"}
	}
	if strings.TrimSpace(cmd.Name) == "" {
		return nil, Rejection{Code: CodeInvalidName, Message: "must not be empty", Field: "name"}
	}

	catID := mintID("cat", cmd.CommandID)
	event := CatRegistered{
		CommandID: cmd.CommandID,
		CatID:     catID,
		Name:      strings.TrimSpace(cmd.Name),
		BirthDate: strings.TrimSpace(cmd.BirthDate),
	}
	return []Event{event}, nil
}

func (a *CatCare) decideLogWeight(cmd LogWeight) ([]Event, error) {
	if !a.Registered {
		return nil, Rejection{Code: CodeNotRegistered, Message: "cat must be registered first"}
	}
	if strings.TrimSpace(cmd.At) == "" {
		return nil, Rejection{Code: CodeInvalidDate, Message: "must not be empty", Field: "at"}
	}
	if cmd.Grams <= 0 {
		return nil, Rejection{Code: CodeInvalidWeight, Message: "must be positive", Field: "grams"}
	}
	if cmd.Grams < MinWeightGrams || cmd.Grams > MaxWeightGrams {
		return nil, Rejection{Code: CodeAbsurdWeight, Message: "outside allowed range", Field: "grams"}
	}

	event := WeightLogged{
		CommandID: cmd.CommandID,
		EntryID:   mintID("weight", cmd.CommandID),
		At:        strings.TrimSpace(cmd.At),
		Grams:     cmd.Grams,
		Notes:     strings.TrimSpace(cmd.Notes),
	}
	return []Event{event}, nil
}

func mintID(prefix string, commandID string) string {
	return prefix + "-" + commandID
}
