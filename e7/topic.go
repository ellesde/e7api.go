package e7

import "errors"

// Topic represents a hero's camping topics. The topics dictate the amount
// of morale gained or less depending on the party composition.
type Topic int

// Hero topic.
const (
	Criticism Topic = iota
	RealityCheck
	HeroicTale
	ComfortingCheer
	CuteCheer
	HeroicCheer
	SadMemory
	JoyfulMemory
	HappyMemory
	UniqueComment
	SelfIndulgent
	Occult
	Myth
	BizarreStory
	FoodStory
	HorrorStory
	Gossip
	Dream
	Advice
	Complain
	Belief
	InterestingStory
)

var topicStrings = map[Topic]string{
	Criticism:        "Criticism",
	RealityCheck:     "Reality Check",
	HeroicTale:       "Heroic Tale",
	ComfortingCheer:  "Comforting Cheer",
	CuteCheer:        "Cute Cheer",
	HeroicCheer:      "Heroic Cheer",
	SadMemory:        "Sad Memory",
	JoyfulMemory:     "Joyful Memory",
	HappyMemory:      "Happy Memory",
	UniqueComment:    "Unique Comment",
	SelfIndulgent:    "Self-Indulgent",
	Occult:           "Occult",
	Myth:             "Myth",
	BizarreStory:     "Bizarre Story",
	FoodStory:        "Food Story",
	HorrorStory:      "Horror Story",
	Gossip:           "Gossip",
	Dream:            "Dream",
	Advice:           "Advice",
	Complain:         "Complain",
	Belief:           "Belief",
	InterestingStory: "Interesting Story",
}

var topics = map[string]Topic{
	"Criticism":         Criticism,
	"Reality Check":     RealityCheck,
	"Heroic Tale":       HeroicTale,
	"Comforting Cheer":  ComfortingCheer,
	"Cute Cheer":        CuteCheer,
	"Heroic Cheer":      HeroicCheer,
	"Sad Memory":        SadMemory,
	"Joyful Memory":     JoyfulMemory,
	"Happy Memory":      HappyMemory,
	"Unique Comment":    UniqueComment,
	"Self-Indulgent":    SelfIndulgent,
	"Occult":            Occult,
	"Myth":              Myth,
	"Bizarre Story":     BizarreStory,
	"Food Story":        FoodStory,
	"Horror Story":      HorrorStory,
	"Gossip":            Gossip,
	"Dream":             Dream,
	"Advice":            Advice,
	"Complain":          Complain,
	"Belief":            Belief,
	"Interesting Story": InterestingStory,
}

func (t Topic) String() string {
	return topicStrings[t]
}

// MarshalJSON marshals t to a quoted JSON string.
func (t Topic) MarshalJSON() ([]byte, error) {
	buf := writeStringBuffer(topicStrings[t])
	return buf.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted JSON string to t.
func (t *Topic) UnmarshalJSON(b []byte) error {
	s, err := unmarshalJSON(b)
	if err != nil {
		return err
	}

	val, ok := topics[s]
	if !ok {
		return ErrUnknownTopic
	}
	*t = val
	return nil
}

// ErrUnknownTopic is returned when unmarshalling a quoted JSON string whose
// value is not in the list of defined topics.
var ErrUnknownTopic = errors.New("unknown topic")
