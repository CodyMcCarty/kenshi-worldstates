package internal

type Faction struct {
	Name       string
	HeadLeader *Leader
	Notes      string
}

// a list of all the factions in play
var (
	F_AntiSlavers    = &Faction{Name: "Anti-Slavers"}
	F_EmpirePeasants = &Faction{Name: "Empire Peasants"}
	F_FreeTraders    = &Faction{Name: "Free Traders"}
	F_Manhunters     = &Faction{Name: "Manhunters"}
	F_MercenaryGuild = &Faction{Name: "Mercenary Guild"}
	F_SlaveHunters   = &Faction{Name: "Slave Hunters"}
	F_TradersGuild   = &Faction{Name: "Trader's Guild"}
	F_UnitedCities   = &Faction{Name: "United Cities"}
)
