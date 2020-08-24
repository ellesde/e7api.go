package e7

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTopic_String(t *testing.T) {
	tests := []struct {
		in   Topic
		want string
	}{
		{in: Criticism, want: "Criticism"},
		{in: RealityCheck, want: "Reality Check"},
		{in: HeroicTale, want: "Heroic Tale"},
		{in: ComfortingCheer, want: "Comforting Cheer"},
		{in: CuteCheer, want: "Cute Cheer"},
		{in: HeroicCheer, want: "Heroic Cheer"},
		{in: SadMemory, want: "Sad Memory"},
		{in: JoyfulMemory, want: "Joyful Memory"},
		{in: HappyMemory, want: "Happy Memory"},
		{in: UniqueComment, want: "Unique Comment"},
		{in: SelfIndulgent, want: "Self-Indulgent"},
		{in: Occult, want: "Occult"},
		{in: Myth, want: "Myth"},
		{in: BizarreStory, want: "Bizarre Story"},
		{in: FoodStory, want: "Food Story"},
		{in: HorrorStory, want: "Horror Story"},
		{in: Gossip, want: "Gossip"},
		{in: Dream, want: "Dream"},
		{in: Advice, want: "Advice"},
		{in: Complain, want: "Complain"},
		{in: Belief, want: "Belief"},
		{in: InterestingStory, want: "Interesting Story"},
	}

	for _, tt := range tests {
		got := tt.in.String()
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Topic.String mismatch (-want, +got):\n%s", diff)
		}
	}
}

func TestTopic_MarshalJSON(t *testing.T) {
	tests := []struct {
		in   Topic
		want []byte
	}{
		{in: Criticism, want: []byte(`"Criticism"`)},
		{in: RealityCheck, want: []byte(`"Reality Check"`)},
		{in: HeroicTale, want: []byte(`"Heroic Tale"`)},
		{in: ComfortingCheer, want: []byte(`"Comforting Cheer"`)},
		{in: CuteCheer, want: []byte(`"Cute Cheer"`)},
		{in: HeroicCheer, want: []byte(`"Heroic Cheer"`)},
		{in: SadMemory, want: []byte(`"Sad Memory"`)},
		{in: JoyfulMemory, want: []byte(`"Joyful Memory"`)},
		{in: HappyMemory, want: []byte(`"Happy Memory"`)},
		{in: UniqueComment, want: []byte(`"Unique Comment"`)},
		{in: SelfIndulgent, want: []byte(`"Self-Indulgent"`)},
		{in: Occult, want: []byte(`"Occult"`)},
		{in: Myth, want: []byte(`"Myth"`)},
		{in: BizarreStory, want: []byte(`"Bizarre Story"`)},
		{in: FoodStory, want: []byte(`"Food Story"`)},
		{in: HorrorStory, want: []byte(`"Horror Story"`)},
		{in: Gossip, want: []byte(`"Gossip"`)},
		{in: Dream, want: []byte(`"Dream"`)},
		{in: Advice, want: []byte(`"Advice"`)},
		{in: Complain, want: []byte(`"Complain"`)},
		{in: Belief, want: []byte(`"Belief"`)},
		{in: InterestingStory, want: []byte(`"Interesting Story"`)},
	}

	for _, tt := range tests {
		got, err := tt.in.MarshalJSON()
		if err != nil {
			t.Errorf("Topic.MarshalJSON returned error: %v", err)
		}
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Topic.MarshalJSON mismatch (-want, +got):\n%s", diff)
		}
	}
}

func TestTopic_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		in   []byte
		want Topic
	}{
		{in: []byte(`"Criticism"`), want: Criticism},
		{in: []byte(`"Reality Check"`), want: RealityCheck},
		{in: []byte(`"Heroic Tale"`), want: HeroicTale},
		{in: []byte(`"Comforting Cheer"`), want: ComfortingCheer},
		{in: []byte(`"Cute Cheer"`), want: CuteCheer},
		{in: []byte(`"Heroic Cheer"`), want: HeroicCheer},
		{in: []byte(`"Sad Memory"`), want: SadMemory},
		{in: []byte(`"Joyful Memory"`), want: JoyfulMemory},
		{in: []byte(`"Happy Memory"`), want: HappyMemory},
		{in: []byte(`"Unique Comment"`), want: UniqueComment},
		{in: []byte(`"Self-Indulgent"`), want: SelfIndulgent},
		{in: []byte(`"Occult"`), want: Occult},
		{in: []byte(`"Myth"`), want: Myth},
		{in: []byte(`"Bizarre Story"`), want: BizarreStory},
		{in: []byte(`"Food Story"`), want: FoodStory},
		{in: []byte(`"Horror Story"`), want: HorrorStory},
		{in: []byte(`"Gossip"`), want: Gossip},
		{in: []byte(`"Dream"`), want: Dream},
		{in: []byte(`"Advice"`), want: Advice},
		{in: []byte(`"Complain"`), want: Complain},
		{in: []byte(`"Belief"`), want: Belief},
		{in: []byte(`"Interesting Story"`), want: InterestingStory},
	}

	for _, tt := range tests {
		tp := new(Topic)
		err := tp.UnmarshalJSON(tt.in)
		if err != nil {
			t.Errorf("Topic.UnmarshalJSON returned error: %v", err)
		}
		if diff := cmp.Diff(tt.want, *tp); diff != "" {
			t.Errorf("Topic.UnmarshalJSON mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestTopic_UnmarshalJSON_unknownTopic(t *testing.T) {
	tp := new(Topic)
	err := tp.UnmarshalJSON([]byte(`"test"`))
	if err == nil {
		t.Errorf("Expected error to be returned")
	}

	if !errors.Is(err, ErrUnknownTopic) {
		t.Errorf("expected unknown topic error")
	}
}
