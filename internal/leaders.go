package internal

import "fmt"

type Leader struct {
	Name    string
	Faction *Faction
	Home    *Town
	// Alive&Free, Dead, or Imprisoned
	Status LeaderStatus
	Notes  string
}

// list of people in play
var (
	L_BossSimion    = &Leader{Name: "Boss Simion", Faction: F_RebelFarmers, Notes: "Can disappear. Tinfist will send you to make an aliance with anti slavers"} // todo revisit
	L_CrabQueen     = &Leader{Name: "nil"}
	L_LdInaba       = &Leader{Name: "nil"}
	L_LdKana        = &Leader{Name: "nil"}
	L_LdMerin       = &Leader{Name: "nil"}
	L_LdNagata      = &Leader{Name: "Lord Nagata", Faction: F_UnitedCities, Home: T_ShoBattai}
	L_LdSanda       = &Leader{Name: "Sanda"}
	L_LdShiro       = &Leader{Name: "nil"}
	L_LdTsugi       = &Leader{Name: "nil"}
	L_LdYoshinaga   = &Leader{Name: "Lord Yoshinaga", Faction: F_UnitedCities, Home: T_Heng, Notes: "No world States. if Heng is half-destroyed or handed over to Empire Peasants, Lord Yoshinaga will despawn"}
	L_Longen        = &Leader{Name: "Longen", Faction: F_TradersGuild, Home: T_TradersEdge, Notes: "Can be turned into Tinfist. Can reward for Tinfist. captured and/or killed Longen and talk to Tinfist, you will get +75 reputation with the anti-slavers instantly, without the -30 to the UC, Slave traders or Traders Guild. "}
	L_Preacher      = &Leader{Name: "nil"}
	L_SMGrace       = &Leader{Name: "nil"}
	L_SMGrande      = &Leader{Name: "Slave Master Grande", Faction: F_TradersGuild, Home: T_Eyesocket}
	L_SMHaga        = &Leader{Name: "nil"}
	L_SMMaster      = &Leader{Name: "Slave Market Master", Faction: F_TradersGuild, Home: T_SlaveMarkets}
	L_SMRen         = &Leader{Name: "nil"}
	L_SMRuben       = &Leader{Name: "nil"}
	L_SMWada        = &Leader{Name: "nil"}
	L_Tengu         = &Leader{Name: "Emperor Tengu", Faction: F_UnitedCities, Home: T_Heft}
	L_Tinfist       = &Leader{Name: "Tinfist", Faction: F_AntiSlavers, Home: T_Spring}
	L_Valamon       = &Leader{Name: "nil"}
	L_Valtena       = &Leader{Name: "nil"}
	L_WestHiveQueen = &Leader{Name: "nil"}
	L_Yabuta        = &Leader{Status: Imprisoned}
)

var Leaders = []*Leader{
	L_BossSimion,
	L_CrabQueen,
	L_LdInaba,
	L_LdKana,
	L_LdMerin,
	L_LdNagata,
	L_LdSanda,
	L_LdShiro,
	L_LdTsugi,
	L_LdYoshinaga,
	L_Longen,
	L_Preacher,
	L_SMGrace,
	L_SMGrande,
	L_SMHaga,
	L_SMMaster,
	L_SMRen,
	L_SMRuben,
	L_SMWada,
	L_Tengu,
	L_Tinfist,
	L_Valamon,
	L_Valtena,
	L_WestHiveQueen,
	L_Yabuta,
}

func (ldr *Leader) LogInfo() {
	info := ldr.GetInfo()
	fmt.Println("Log [Info]", info)
}

func (ldr *Leader) GetInfo() string {
	if ldr == nil {
		info := fmt.Sprint("Leader NOT found while logging leader info")
		return info
	}

	f := "nil"
	if ldr.Faction != nil {
		f = ldr.Faction.Name
	}
	h := "nil"
	if ldr.Home != nil {
		h = ldr.Home.Name
	}
	info := fmt.Sprint(ldr.Name+"("+f+")", " from ", h+".", " is ", ldr.Status.String()+". ", ldr.Notes)
	return info
}

type Cond struct {
	Label  string
	Want   bool
	Eval   BoolExpr
	Leader *Leader
}

func (l *Leader) IsAlive(want bool) Cond {
	return Cond{
		Label:  "is alive",
		Want:   want,
		Eval:   l.npcIs1(),
		Leader: l,
	}
}

func (l *Leader) IsNotAlive(want bool) Cond {
	return Cond{
		Label:  "is not alive",
		Want:   want,
		Eval:   l.npcIsNot1(),
		Leader: l,
	}
}

func (l *Leader) Okay(want bool) Cond {
	return Cond{
		Label:  "okay",
		Want:   want,
		Eval:   l.npcIs1(),
		Leader: l,
	}
}

func (l *Leader) Ok(want bool) Cond {
	return Cond{
		Label:  "okay",
		Want:   want,
		Eval:   l.npcIs1(),
		Leader: l,
	}
}

//// IsAlive bit of a misnomer. means free   Ok
//func (l *Leader) IsAlive() BoolExpr {
//	return func() bool { return l.Status == Alive }
//}
//
//// IsNotAlive is dead. Used in a few of Tinfist's world states.
//func (l *Leader) IsNotAlive() BoolExpr {
//	return func() bool { return l.Status == Dead }
//}
//
//// Okay NPC is Alive. First seen with the west hive queen.
//func (l *Leader) Okay() BoolExpr {
//	return func() bool { return l.Status == Alive }
//}
//
//// Ok NPC is 1 (alive). The reason for reduant functions is A. that's how it is in the FCS 2. it deals with player involvment and other factors out of scope for this app.
//func (l *Leader) Ok() BoolExpr {
//	return func() bool { return l.Status == Alive }
//}

func (l *Leader) npcIs1() BoolExpr {
	return func() bool { return l.Status == Alive }
}

func (l *Leader) npcIsNot1() BoolExpr {
	return func() bool { return l.Status != Alive }
}

// Matched to the FCS' status of people
type LeaderStatus uint8

const (
	// bit of a misnomer.  it means free
	Alive LeaderStatus = iota
	Dead
	Imprisoned
)

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
