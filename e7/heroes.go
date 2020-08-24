package e7

import (
	"context"
	"fmt"
	"net/http"
)

// HeroesService handles communication with the heroes related
// methods of the EpicSevenDB API.
type HeroesService service

// Hero represents an Epic Seven hero profile.
type Hero struct {
	UUID                string                                `json:"_id,omitempty"`
	ID                  string                                `json:"id,omitempty"`
	Name                string                                `json:"name,omitempty"`
	Moonlight           bool                                  `json:"moonlight,omitempty"`
	Rarity              uint                                  `json:"rarity,omitempty"`
	Attribute           Attribute                             `json:"attribute,omitempty"`
	Role                Role                                  `json:"role,omitempty"`
	Zodiac              string                                `json:"zodiac,omitempty"`
	Description         string                                `json:"description,omitempty"`
	Story               string                                `json:"story,omitempty"`
	GetLine             string                                `json:"get_line,omitempty"`
	Stats               BaseStats                             `json:"stats,omitempty"`
	Relationships       []Relationship                        `json:"relationships,omitempty"`
	SelfDevotion        SelfDevotion                          `json:"self_devotion,omitempty"`
	Devotion            Devotion                              `json:"devotion,omitempty"`
	Specialty           Specialty                             `json:"specialty,omitempty"`
	Camping             Camping                               `json:"camping,omitempty"`
	ZodiacTree          []ZodiacNode                          `json:"zodiac_tree,omitempty"`
	Skills              []Skill                               `json:"skills,omitempty"`
	SpecialtyChange     SpecialtyChange                       `json:"specialty_change,omitempty"`
	Assets              Assets                                `json:"assets,omitempty"`
	Buffs               []Buff                                `json:"buffs,omitempty"`
	Debuffs             []Debuff                              `json:"debuffs,omitempty"`
	Common              []Common                              `json:"common,omitempty"`
	ExclusiveEquipments []ExclusiveEquipment                  `json:"exclusiveEquipments,omitempty"`
	CalculatedStats     map[PreCalculatedState]CalculatedStat `json:"calculatedStatus,omitempty"`
}

// BaseStats represents an Epic Seven hero's base stats modifiers.
type BaseStats struct {
	Bra int `json:"bra,omitempty"`
	Int int `json:"int,omitempty"`
	Fai int `json:"fai,omitempty"`
	Des int `json:"des,omitempty"`
}

// Relationship represents an Epic Seven hero's relationship details with other
// heroes.
type Relationship struct {
	ID          string      `json:"id,omitempty"`
	Slot        int         `json:"slot,omitempty"`
	Description string      `json:"description,omitempty"`
	Relation    string      `json:"relation,omitempty"`
	Upgrade     interface{} `json:"upgrade,omitempty"`
	RelationID  string      `json:"relation_id,omitempty"`
}

// SelfDevotion represents an Epic Seven hero's self devotion details.
type SelfDevotion struct {
	Type   Stat           `json:"type,omitempty"`
	Grades DevotionGrades `json:"grades,omitempty"`
}

// Devotion represents an Epic Seven hero's devotion details. Devotion are
// stat multipliers added to other heroes based on slot.
//
// e.g. Aramintha has an ATK devotion for slots 1, 2, 3, 4. Heroes in those
// slots will gain additional ATK multipliers.
type Devotion struct {
	Type   Stat           `json:"type,omitempty"`
	Grades DevotionGrades `json:"grades,omitempty"`
	Slots  Slots          `json:"slots,omitempty"`
}

// DevotionGrades represents an Epic Seven hero's devotion grade multipliers.
// Devotion grades are tiered and increases a particular stat by some factor.
//
// e.g.
// B - 3.6%
// A - 5.4%
// S - 7.2%
// SS - 9.0%
// SSS - 1.08%
type DevotionGrades struct {
	B   float32 `json:"B,omitempty"`
	A   float32 `json:"A,omitempty"`
	S   float32 `json:"S,omitempty"`
	SS  float32 `json:"SS,omitempty"`
	SSS float32 `json:"SSS,omitempty"`
}

// Slots represents the position in a party.
type Slots struct {
	One   bool `json:"1,omitempty"`
	Two   bool `json:"2,omitempty"`
	Three bool `json:"3,omitempty"`
	Four  bool `json:"4,omitempty"`
}

// Specialty represents an Epic Seven hero's specialty details.
type Specialty struct {
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	EffectType  string        `json:"effect_type,omitempty"`
	EffectValue float32       `json:"effect_value,omitempty"`
	Command     int           `json:"command,omitempty"`
	Charm       int           `json:"charm,omitempty"`
	Politics    int           `json:"politics,omitempty"`
	Type        SpecialtyType `json:"type,omitempty"`
	Assets      Assets        `json:"assets,omitempty"`
}

// SpecialtyType represents an Epic Seven hero's specialty type.
type SpecialtyType struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Assets contains URLs to a hero's assets.
type Assets struct {
	Thumbnail string `json:"thumbnail,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Image     string `json:"image,omitempty"`
}

// Camping represents an Epic Seven hero's camping details.
type Camping struct {
	Personalities []string      `json:"personalities,omitempty"`
	Topics        []Topic       `json:"topics,omitempty"`
	Values        CampingValues `json:"values,omitempty"`
}

// CampingValues represents the morale points gained or lost when selecting
// the accompanying topic.
type CampingValues struct {
	Criticism        int `json:"Criticism,omitempty"`
	RealityCheck     int `json:"Reality Check,omitempty"`
	HeroicTale       int `json:"Heroic Tale,omitempty"`
	ComfortingCheer  int `json:"Comforting Cheer,omitempty"`
	CuteCheer        int `json:"Cute Cheer,omitempty"`
	HeroicCheer      int `json:"Heroic Cheer,omitempty"`
	SadMemory        int `json:"Sad Memory,omitempty"`
	JoyfulMemory     int `json:"Joyful Memory,omitempty"`
	HappyMemory      int `json:"Happy Memory,omitempty"`
	UniqueComment    int `json:"Unique Comment,omitempty"`
	SelfIndulgent    int `json:"Self-Indulgent,omitempty"`
	Occult           int `json:"Occult,omitempty"`
	Myth             int `json:"Myth,omitempty"`
	BizarreStory     int `json:"Bizarre Story,omitempty"`
	FoodStory        int `json:"Food Story,omitempty"`
	HorrorStory      int `json:"Horror Story,omitempty"`
	Gossip           int `json:"Gossip,omitempty"`
	Dream            int `json:"Dream,omitempty"`
	Advice           int `json:"Advice,omitempty"`
	Complain         int `json:"Complain,omitempty"`
	Belief           int `json:"Belief,omitempty"`
	InterestingStory int `json:"Interesting Story,omitempty"`
}

// ZodiacNode represents an Epic Seven hero's zodiac tree details.
type ZodiacNode struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	// SkillEnhanced value will be false when the name of the ZodiacNode
	// is "Potential Stone"; otherwise it will be a number when the name
	// of the ZodiacStone is "Ability Stone".
	SkillEnhanced interface{} `json:"skill_enhanced,omitempty"`
	Costs         []NodeCost  `json:"costs,omitempty"`
	Stats         []NodeStat  `json:"stats,omitempty"`
	UUID          string      `json:"_id,omitempty"`
}

// NodeCost represents the cost for unlocking a ZodiacNode.
type NodeCost struct {
	Item         string     `json:"item,omitempty"`
	Count        int        `json:"count,omitempty"`
	ID           string     `json:"_id,omitempty"`
	Identifier   string     `json:"identifier,omitempty"`
	Name         string     `json:"name,omitempty"`
	Description  string     `json:"description,omitempty"`
	Category     string     `json:"category,omitempty"`
	Attribute    *Attribute `json:"attribute,omitempty"`
	Grade        uint       `json:"grade,omitempty"`
	Type1        string     `json:"type1,omitempty"`
	Type2        *string    `json:"type2"`
	Assets       Assets     `json:"assets,omitempty"`
	RequestCount uint       `json:"request_count,omitempty"`
	SupportCount uint       `json:"support_count,omitempty"`
}

// NodeStat represents the stat increase when unlocking a ZodiacNode.
type NodeStat struct {
	Stat  Stat    `json:"stat,omitempty"`
	Value float32 `json:"value,omitempty"`
	Type  string  `json:"type,omitempty"`
}

// Skill represents an Epic Seven hero's skill.
type Skill struct {
	Name              string        `json:"name,omitempty"`
	CanEnhance        bool          `json:"can_enhance,omitempty"`
	Description       string        `json:"description,omitempty"`
	Values            []float32     `json:"values,omitempty"`
	Passive           bool          `json:"passive,omitempty"`
	Cooldown          uint          `json:"cooldown,omitempty"`
	SoulGain          uint          `json:"soul_gain,omitempty"`
	Pow               float32       `json:"pow,omitempty"`
	AttackPercent     float32       `json:"att_rate,omitempty"`
	Buff              []uint        `json:"buff,omitempty"`
	Debuff            []uint        `json:"debuff,omitempty"`
	Common            []uint        `json:"common,omitempty"`
	SoulDescription   string        `json:"soul_description,omitempty"`
	SoulRequirement   uint          `json:"soul_requirement,omitempty"`
	SoulPow           float32       `json:"soul_pow,omitempty"`
	SoulAttackPercent float32       `json:"soul_att_rate,omitempty"`
	Enhancements      []Enhancement `json:"enhancements,omitempty"`
}

// SpecialtyChange represents an Epic Seven hero's specialty change details.
type SpecialtyChange struct {
	ID           string    `json:"id,omitempty"`
	ChangedSkill uint      `json:"changedSkill,omitempty"`
	Quests       []Quest   `json:"quests,omitempty"`
	Tree         SkillTree `json:"tree,omitempty"`
}

// Quest represents the quests needed to complete for a hero's specialty change.
type Quest struct {
	Category           string `json:"category,omitempty"`
	MissionName        string `json:"mission_name,omitempty"`
	MissionDescription string `json:"mission_description,omitempty"`
}

// SkillTree represents an Epic Seven hero's specialty change skill tree.
type SkillTree []SkillBranch

// SkillBranch represents a branch of SkillNodes in a skill tree.
type SkillBranch []SkillNode

// SkillNode represents a node in a skill branch.
type SkillNode struct {
	ID           uint               `json:"id,omitempty"`
	Position     uint               `json:"position,omitempty"`
	RequireID    *uint              `json:"require_id,omitempty"`
	Enhancements []SkillEnhancement `json:"enhancements,omitempty"`
}

// SkillEnhancement represents the enhancements when upgrading a skill node.
type SkillEnhancement struct {
	Type        string  `json:"type,omitempty"`
	Stat        *Stat   `json:"stat,omitempty"`
	Value       float32 `json:"value,omitempty"`
	Description string  `json:"description,omitempty"`
	Upgrade     *string `json:"upgrade,omitempty"`
}

// Enhancement represents the enhancement details for a skill.
type Enhancement struct {
	Description string            `json:"string,omitempty"`
	Costs       []EnhancementCost `json:"costs,omitempty"`
	UUID        string            `json:"_id,omitempty"`
}

// EnhancementCost represents the cost for enhancing a skill.
type EnhancementCost struct {
	Item         string      `json:"item,omitempty"`
	Count        uint        `json:"count,omitempty"`
	UUID         string      `json:"_id,omitempty"`
	Identifier   string      `json:"identifier,omitempty"`
	Name         string      `json:"name,omitempty"`
	Description  string      `json:"description,omitempty"`
	Category     string      `json:"category,omitempty"`
	Attribute    interface{} `json:"attribute,omitempty"`
	Grade        uint        `json:"grade,omitempty"`
	Type1        string      `json:"type1,omitempty"`
	Type2        interface{} `json:"type2,omitempty"`
	Assets       Assets      `json:"assets,omitempty"`
	RequestCount uint        `json:"request_count,omitempty"`
	SupportCount uint        `json:"support_count,omitempty"`
}

// Common represents a common effect's detail.
type Common struct {
	UUID   string `json:"_id,omitempty"`
	ID     uint   `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Name   string `json:"name,omitempty"`
	Effect string `json:"effect,omitempty"`
	Assets Assets `json:"assets,omitempty"`
}

// Buff represents a buff effect's details.
type Buff struct {
	Common
}

// Debuff represents a debuff effect's details.
type Debuff struct {
	Common
}

// ExclusiveEquipment represents exclusive equipment details.
type ExclusiveEquipment struct {
	UUID        string                    `json:"_id,omitempty"`
	ID          string                    `json:"id,omitempty"`
	Name        string                    `json:"name,omitempty"`
	Description string                    `json:"description,omitempty"`
	Unit        string                    `json:"unit,omitempty"`
	Role        Role                      `json:"role,omitempty"`
	Rarity      uint                      `json:"rarity,omitempty"`
	Stat        ExclusiveEquipmentStat    `json:"stat,omitempty"`
	Skills      []ExclusiveEquipmentSkill `json:"skills,omitempty"`
	Assets      Assets                    `json:"assets,omitempty"`
}

// ExclusiveEquipmentStat represents the stat increase when equipping an Exclusive Equipment.
type ExclusiveEquipmentStat struct {
	Type  Stat    `json:"type,omitempty"`
	Value float32 `json:"value,omitempty"`
}

// ExclusiveEquipmentSkill represents the skills provided by an Exclusive Equipment.
type ExclusiveEquipmentSkill struct {
	Skill            uint   `json:"skill,omitempty"`
	Description      string `json:"description,omitempty"`
	SkillDescription string `json:"skill_description,omitempty"`
	Values           []uint `json:"values,omitempty"`
	UUID             uint   `json:"_id,omitempty"`
}

// PreCalculatedState represents a state identifier of precalculated stats
// for a hero.
type PreCalculatedState string

const (
	// Level50FiveStarNoAwaken is a precalculated state representing
	// a level 50 five starred hero that has no awakenings.
	Level50FiveStarNoAwaken PreCalculatedState = "lv50FiveStarNoAwaken"

	// Level50FiveStarFullyAwakened is a precalculated state representing
	// a level 50 five starred hero that is fully awakened.
	Level50FiveStarFullyAwakened PreCalculatedState = "lv50FiveStarFullyAwakened"

	// Level60SixStarNoAwaken is a precalculated state representing
	// a level 60 six starred hero that has no awakenings.
	Level60SixStarNoAwaken PreCalculatedState = "lv60SixStarNoAwaken"

	// Level60SixStarFullyAwakened is a precalculated state representing
	// a level 60 six starred hero that is fully awakend.
	Level60SixStarFullyAwakened PreCalculatedState = "lv60SixStarFullyAwakened"
)

// CalculatedStat represents the total stats at a given precalculate state.
type CalculatedStat struct {
	CombatPoints      uint    `json:"cp,omitempty"`
	Attack            uint    `json:"atk,omitempty"`
	Health            uint    `json:"hp,omitempty"`
	Speed             uint    `json:"spd,omitempty"`
	Defense           uint    `json:"def,omitempty"`
	CriticalHitChance float32 `json:"chc,omitempty"`
	CriticalHitDamage float32 `json:"chd,omitempty"`
	DualAttackChance  float32 `json:"dac,omitempty"`
	Effectiveness     float32 `json:"eff,omitempty"`
	EffectResistance  float32 `json:"efr,omitempty"`
}

// HeroesResponse is an EpicSevenDB API response that contains a list
// of heroes.
type HeroesResponse struct {
	Results  []*Hero  `json:"results,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// GetByID fetches a hero by ID. The ID is the hero's name in lowercase. Heroes with space in their names
// must be hyphenated. e.g. Little Queen Charlotte would be little-queen-charlotte.
func (s *HeroesService) GetByID(ctx context.Context, hero string) (*Hero, *http.Response, error) {
	u := fmt.Sprintf("hero/%v", hero)
	req, err := s.client.NewRequest(http.MethodGet, u)
	if err != nil {
		return nil, nil, err
	}

	response := new(HeroesResponse)
	resp, err := s.client.Do(ctx, req, response)
	if err != nil {
		return nil, resp, err
	}

	return response.Results[0], resp, nil
}
