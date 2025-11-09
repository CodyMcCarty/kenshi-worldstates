package internal

import "fmt"

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
	L_BossSimion    = &Leader{Name: "Boss Simion", Faction: F_RebelFarmers, Notes: "Can disappear. Tinfist will send you to make an aliance"} // todo revisit
	L_CrabQueen     = &Leader{}
	L_LdInaba       = &Leader{}
	L_LdKana        = &Leader{}
	L_LdMerin       = &Leader{}
	L_LdNagata      = &Leader{Name: "Lord Nagata", Faction: F_UnitedCities, Home: T_ShoBattai}
	L_LdSanda       = &Leader{}
	L_LdShiro       = &Leader{}
	L_LdTsugi       = &Leader{}
	L_LdYoshinaga   = &Leader{Name: "Lord Yoshinaga", Faction: F_UnitedCities, Home: T_Heng, Notes: "No world States. if Heng is half-destroyed or handed over to Empire Peasants, Lord Yoshinaga will despawn"}
	L_Longen        = &Leader{Name: "Longen", Faction: F_TradersGuild, Home: T_TradersEdge, Notes: "Can be turned into Tinfist. Can reward for Tinfist. captured and/or killed Longen and talk to Tinfist, you will get +75 reputation with the anti-slavers instantly, without the -30 to the UC, Slave traders or Traders Guild. "}
	L_Preacher      = &Leader{}
	L_SMGrace       = &Leader{}
	L_SMGrande      = &Leader{Name: "Slave Master Grande", Faction: F_TradersGuild, Home: T_Eyesocket}
	L_SMHaga        = &Leader{}
	L_SMMaster      = &Leader{Name: "Slave Market Master", Faction: F_TradersGuild, Home: T_SlaveMarkets}
	L_SMRen         = &Leader{}
	L_SMRuben       = &Leader{}
	L_SMWada        = &Leader{}
	L_Tengu         = &Leader{Name: "Emperor Tengu", Faction: F_UnitedCities, Home: T_Heft}
	L_Tinfist       = &Leader{Name: "Tinfist", Faction: F_AntiSlavers, Home: T_Spring}
	L_Valamon       = &Leader{}
	L_Valtena       = &Leader{}
	L_WestHiveQueen = &Leader{}
	L_Yabuta        = &Leader{Status: Imprisoned}
)

func PrintLeaderStatus(ldrs []*Leader) {
	fmt.Println("***Leader Status***")
	for _, ldr := range ldrs {
		fmt.Println(ldr.Name, ldr.Status)
	}
	fmt.Println("---")
}

// IsAlive bit of a misnomer. means free
func (l *Leader) IsAlive() BoolExpr {
	return func() bool { return l.Status == Alive }
}

// IsNotAlive is dead. Used in a few of Tinfist's world states.
func (l *Leader) IsNotAlive() BoolExpr {
	return func() bool { return l.Status == Dead }
}

// Okay NPC is Alive. First seen with the west hive queen.
func (l *Leader) Okay() BoolExpr {
	return func() bool { return l.Status == Alive }
}

// Ok NPC is 1 (alive). The reason for reduant functions is A. that's how it is in the FCS 2. it deals with player involvment and other factors out of scope for this app.
func (l *Leader) Ok() BoolExpr {
	return func() bool { return l.Status == Alive }
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
