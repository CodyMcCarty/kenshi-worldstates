package internal

type Town struct {
	Name string
	// where on the map. i.e. Heng override prosperous is at Heng
	MapLoc     *Town
	Faction    *Faction
	Overrides  []*Town
	WorldState []Cond
	Notes      string
}

// GuidedStep1 add a Town like T_SomeTown = ...

// list of towns at play. Leave empty until they have all their info populated. The towns here play the role of a town just like any override, and double duty as the map location.
var (
	T_Bark               = &Town{Name: "Bark", Faction: F_UnitedCities}
	T_Brink              = &Town{Name: "Brink", Faction: F_UnitedCities}
	T_Catun              = &Town{Name: "Catun", Faction: F_UnitedCities}
	T_ClownSteady        = &Town{Name: "ClownSteady", Faction: F_UnitedCities}
	T_CultVillage        = &Town{Name: "CultVillage", Faction: F_PreacherCult}
	T_DistantHiveVillage = &Town{Name: "Distant Hive Village", Faction: F_WesternHive}
	T_DriftersLast       = &Town{Name: "Drifters Last", Faction: F_UnitedCities}
	T_Drin               = &Town{Name: "Drin", Faction: F_UnitedCities}
	T_Eyesocket          = &Town{Name: "Eyesocket", Faction: F_SlaveTraders}
	T_FreeSettlement     = &Town{Name: "Free Settlement", Faction: F_UnitedCities}
	T_FishingVillage     = &Town{Name: "Fishing Village", Faction: F_Deadcat}
	T_FortSimion         = &Town{Name: "Fort Simion", Faction: F_RebelFarmers}
	T_Heft               = &Town{Name: "Heft", Faction: F_UnitedCities}
	T_Heng               = &Town{Name: "Heng", Faction: F_UnitedCities}
	T_ManhunterBase      = &Town{Name: "ManhunterBase", Faction: F_Manhunters}
	T_Mourn              = &Town{Name: "Morn", Faction: F_TechHunters}
	T_PortNorth          = &Town{Name: "Port North", Faction: F_SlaveTraders}
	T_PortSouth          = &Town{Name: "Port South", Faction: F_SlaveTraders}
	T_SettledNomads      = &Town{Name: "Settled Nomads", Faction: F_Nomads}
	T_ShoBattai          = &Town{Name: "Sho-Battai", Faction: F_UnitedCities}
	T_SlaveFarm          = &Town{Name: "Slave Farm", Faction: F_SlaveTraders}
	T_SlaveFarmSouth     = &Town{Name: "Slave Farm South", Faction: F_SlaveTraders}
	T_SlaveMarkets       = &Town{Name: "Slave Markets", Faction: F_SlaveTraders}
	T_SouthStoneCamp     = &Town{Name: "South Stone Camp", Faction: F_SlaveTraders}
	T_Spring             = &Town{Name: "Spring", Faction: F_AntiSlavers}
	T_Stoat              = &Town{Name: "Stoat", Faction: F_UnitedCities}
	T_StoneCamp          = &Town{Name: "Stone Camp", Faction: F_SlaveTraders}
	T_TradersEdge        = &Town{Name: "Traders-Edge", Faction: F_TradersGuild}
)
