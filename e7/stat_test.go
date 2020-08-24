package e7

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStat_String(t *testing.T) {
	tests := []struct {
		in   Stat
		want string
	}{
		{in: Attack, want: "att"},
		{in: AttackPercent, want: "att_rate"},
		{in: Defense, want: "def"},
		{in: DefensePercent, want: "def_rate"},
		{in: Health, want: "max_hp"},
		{in: HealthPercent, want: "max_hp_rate"},
		{in: Speed, want: "speed"},
		{in: CriticalHitChance, want: "cri"},
		{in: CriticalHitDamage, want: "cri_dmg"},
		{in: Effectiveness, want: "acc"},
		{in: EffectResistance, want: "res"},
		{in: DualAttackChance, want: "coop"},
	}

	for _, tt := range tests {
		got := tt.in.String()
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Stat.String mismatch (-want, +got):\n%s", diff)
		}
	}
}

func TestStat_MarshalJSON(t *testing.T) {
	tests := []struct {
		in   Stat
		want []byte
	}{
		{in: Attack, want: []byte(`"att"`)},
		{in: AttackPercent, want: []byte(`"att_rate"`)},
		{in: Defense, want: []byte(`"def"`)},
		{in: DefensePercent, want: []byte(`"def_rate"`)},
		{in: Health, want: []byte(`"max_hp"`)},
		{in: HealthPercent, want: []byte(`"max_hp_rate"`)},
		{in: Speed, want: []byte(`"speed"`)},
		{in: CriticalHitChance, want: []byte(`"cri"`)},
		{in: CriticalHitDamage, want: []byte(`"cri_dmg"`)},
		{in: Effectiveness, want: []byte(`"acc"`)},
		{in: EffectResistance, want: []byte(`"res"`)},
		{in: DualAttackChance, want: []byte(`"coop"`)},
	}

	for _, tt := range tests {
		got, err := tt.in.MarshalJSON()
		if err != nil {
			t.Errorf("Stat.MarshalJSON returned error: %v", err)
		}
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("Stat.MarshalJSON mismatch (-want, +got):\n%s", diff)
		}
	}
}

func TestStat_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		in   []byte
		want Stat
	}{
		{in: []byte(`"att"`), want: Attack},
		{in: []byte(`"att_rate"`), want: AttackPercent},
		{in: []byte(`"def"`), want: Defense},
		{in: []byte(`"def_rate"`), want: DefensePercent},
		{in: []byte(`"max_hp"`), want: Health},
		{in: []byte(`"max_hp_rate"`), want: HealthPercent},
		{in: []byte(`"speed"`), want: Speed},
		{in: []byte(`"cri"`), want: CriticalHitChance},
		{in: []byte(`"cri_dmg"`), want: CriticalHitDamage},
		{in: []byte(`"acc"`), want: Effectiveness},
		{in: []byte(`"res"`), want: EffectResistance},
		{in: []byte(`"coop"`), want: DualAttackChance},
	}

	for _, tt := range tests {
		s := new(Stat)
		err := s.UnmarshalJSON(tt.in)
		if err != nil {
			t.Errorf("Stat.UnmarshalJSON returned error: %v", err)
		}
		if diff := cmp.Diff(tt.want, *s); diff != "" {
			t.Errorf("Stat.UnmarshalJSON mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestStat_UnmarshalJSON_unknownStat(t *testing.T) {
	s := new(Stat)
	err := s.UnmarshalJSON([]byte(`"test"`))
	if err == nil {
		t.Errorf("Expected error to be returned")
	}

	if !errors.Is(err, ErrUnknownStat) {
		t.Errorf("expected unknown stat error")
	}
}
