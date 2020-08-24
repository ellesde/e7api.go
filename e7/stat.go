package e7

import "errors"

// Stat represents a hero's stat.
type Stat int

// Hero stat.
const (
	Attack Stat = iota
	AttackPercent
	Defense
	DefensePercent
	Health
	HealthPercent
	Speed
	CriticalHitChance
	CriticalHitDamage
	Effectiveness
	EffectResistance
	DualAttackChance
)

func (s Stat) String() string {
	return statStrings[s]
}

var statStrings = map[Stat]string{
	Attack:            "att",
	AttackPercent:     "att_rate",
	Defense:           "def",
	DefensePercent:    "def_rate",
	Health:            "max_hp",
	HealthPercent:     "max_hp_rate",
	Speed:             "speed",
	CriticalHitChance: "cri",
	CriticalHitDamage: "cri_dmg",
	Effectiveness:     "acc",
	EffectResistance:  "res",
	DualAttackChance:  "coop",
}

var stats = map[string]Stat{
	"att":         Attack,
	"att_rate":    AttackPercent,
	"def":         Defense,
	"def_rate":    DefensePercent,
	"max_hp":      Health,
	"max_hp_rate": HealthPercent,
	"speed":       Speed,
	"cri":         CriticalHitChance,
	"cri_dmg":     CriticalHitDamage,
	"acc":         Effectiveness,
	"res":         EffectResistance,
	"coop":        DualAttackChance,
}

// MarshalJSON marshals s to a quoted JSON string.
func (s Stat) MarshalJSON() ([]byte, error) {
	buf := writeStringBuffer(statStrings[s])
	return buf.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted JSON string to s.
func (s *Stat) UnmarshalJSON(b []byte) error {
	str, err := unmarshalJSON(b)
	if err != nil {
		return err
	}

	val, ok := stats[str]
	if !ok {
		return ErrUnknownStat
	}
	*s = val
	return nil
}

// ErrUnknownStat is returned when unmarshalling a quoted JSON string whose
// value is not in the list of defined stats.
var ErrUnknownStat = errors.New("unknown stat")
