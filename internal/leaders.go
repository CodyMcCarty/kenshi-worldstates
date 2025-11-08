package internal

// Matched to the FCS' status of people
type LeaderStatus uint8

const (
	// bit of a misnomer.  it means free
	Alive LeaderStatus = iota
	Dead
	Imprisoned
)

type Leader struct {
	Name    string
	Faction *Faction
	Home    *Town
	// Free, Dead, or Imprisoned
	Status LeaderStatus
	Notes  string
}

// list of people in play
var (
	L_CrabQueen   = &Leader{}
	L_Longen      = &Leader{Name: "Longen", Faction: F_TradersGuild, Home: T_TradersEdge, Notes: "Can be turned into Tinfist. Can reward for Tinfist. captured and/or killed Longen and talk to Tinfist, you will get +75 reputation with the anti-slavers instantly, without the -30 to the UC, Slave traders or Traders Guild. "}
	L_Preacher    = &Leader{}
	L_LdInaba     = &Leader{}
	L_LdKana      = &Leader{}
	L_LdMerin     = &Leader{}
	L_LdNagata    = &Leader{Name: "Lord Nagata", Faction: F_UnitedCities, Home: T_ShoBattai}
	L_LdSanda     = &Leader{}
	L_LdTsugi     = &Leader{}
	L_LdYoshinaga = &Leader{Name: "Lord Yoshinaga", Faction: F_UnitedCities, Home: T_Heng, Notes: "No world States. if Heng is half-destroyed or handed over to Empire Peasants, Lord Yoshinaga will despawn"}
	L_SMGrace     = &Leader{}
	L_SMGrande    = &Leader{Name: "Slave Master Grande", Faction: F_TradersGuild, Home: T_Eyesocket}
	L_SMMaster    = &Leader{Name: "Slave Market Master", Faction: F_TradersGuild, Home: T_SlaveMarkets}
	L_SMRen       = &Leader{}
	L_SMRuben     = &Leader{}
	L_SMWada      = &Leader{}
	L_Tengu       = &Leader{Name: "Emperor Tengu", Faction: F_UnitedCities, Home: T_Heft}
	L_Tinfist     = &Leader{Name: "Tinfist", Faction: F_AntiSlavers, Home: T_Spring}
	L_Valamon     = &Leader{}
)

// bit of a misnomer. means free
func (l *Leader) IsAlive() BoolExpr {
	return func() bool { return l.Status == Alive }
}

// is dead
func (l *Leader) IsNotAlive() BoolExpr {
	return func() bool { return l.Status == Dead }
}

func (s LeaderStatus) String() string {
	switch s {
	case Alive:
		return "Alive"
	case Dead:
		return "Dead"
	case Imprisoned:
		return "Imprisoned"
	default:
		panic("Unknown leader status")
		return "Unknown"
	}
}
