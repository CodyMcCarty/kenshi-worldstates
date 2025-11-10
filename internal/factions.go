package internal

type Faction struct {
	Name  string
	Notes string
}

// a list of all the factions in play
var (
	F_AntiSlavers    = &Faction{Name: "Anti-Slavers"}
	F_Cannibal       = &Faction{Name: "Cannibal"}
	F_CrabRaiders    = &Faction{Name: "Crab Raiders"}
	F_Deadcat        = &Faction{Name: "Deadcat"}
	F_EmpirePeasants = &Faction{Name: "Empire Peasants"}
	F_Fishmen        = &Faction{Name: "Fishmen"}
	F_Fogmen         = &Faction{Name: "Fogmen"}
	F_FreeTraders    = &Faction{Name: "Free Traders"}
	F_HolyNation     = &Faction{Name: "Holy Nation"}
	F_Manhunters     = &Faction{Name: "Manhunters"}
	F_MercenaryGuild = &Faction{Name: "Mercenary Guild"}
	F_PreacherCult   = &Faction{Name: "Preacher Cult"}
	F_Reavers        = &Faction{Name: "Reavers"}
	F_RebelFarmers   = &Faction{Name: "Rebel Farmers"}
	F_RebelSwordsmen = &Faction{Name: "Rebel Swordsmen"}
	F_Nomads         = &Faction{Name: "Nomads"}
	F_SlaveHunters   = &Faction{Name: "Slave Hunters"}
	F_SlaveTraders   = &Faction{Name: "Slave Traders"}
	F_TechHunters    = &Faction{Name: "Tech Hunters"}
	F_TradersGuild   = &Faction{Name: "Trader's Guild"}
	F_UnitedCities   = &Faction{Name: "United Cities"}
	F_WesternHive    = &Faction{Name: "Western Hive"}
	F_YabutaOutlaws  = &Faction{Name: "Yabuta Outlaws"}
)
