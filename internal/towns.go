package internal

type Town struct {
	Name string
	// where on the map. i.e. Heng override prosperous is at Heng
	Location   *Town
	Faction    *Faction
	Overrides  []*Town
	WorldState []Cond
	Notes      string
}

// list of towns at play
var (
	T_ClownSteady    = &Town{}
	T_CultVillage    = &Town{}
	T_Eyesocket      = &Town{}
	T_FreeSettlement = &Town{}
	T_FishingVillage = &Town{}
	T_Heft           = &Town{Name: "Heft"}
	T_Heng           = &Town{Name: "Heng", Faction: F_UnitedCities, Notes: "if Heng is half-destroyed or handed over to Empire Peasants, Lord Yoshinaga will despawn"}
	T_ManhunterBase  = &Town{}
	T_SettledNomads  = &Town{}
	T_SlaveFarmSouth = &Town{}
	T_SlaveMarkets   = &Town{}
	T_ShoBattai      = &Town{Name: "Sho-Battai"}
	T_Spring         = &Town{Name: "Spring"}
	T_TradersEdge    = &Town{Name: "Traders-Edge"}
)
