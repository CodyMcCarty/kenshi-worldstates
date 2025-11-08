package main

import "fmt"

type Faction struct {
	Name       string
	HeadLeader *Leader
	Notes      string
}

var (
	F_AntiSlavers    = &Faction{Name: "Anti-Slavers"}
	F_EmpirePeasants = &Faction{Name: "Empire Peasants"}
	F_TradersGuild   = &Faction{Name: "Trader's Guild"}
	F_UnitedCities   = &Faction{Name: "United Cities"}
)

type LeaderStatus uint8

const (
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
		return "Unknown"
	}
}

type Leader struct {
	Name     string
	Faction  *Faction
	Location *Town
	Status   LeaderStatus
	Notes    string
}

type BoolExpr func() bool // lazy boolean (re-evaluated each time)
type Cond func() bool     // condition (also lazy)

func (l *Leader) IsAlive() BoolExpr {
	return func() bool { return l.Status == Alive }
}
func (l *Leader) IsNotAlive() BoolExpr {
	return func() bool { return l.Status == Dead }
}

func Equal(b BoolExpr, v bool) Cond { return func() bool { return b() == v } }

//	func And(cs ...Cond) Cond {
//		return func() bool {
//			for _, c := range cs {
//				if !c() {
//					return false
//				}
//			}
//			return true
//		}
//	}
func allTrue(cs []Cond) bool {
	for _, c := range cs {
		if !c() {
			return false
		}
	}
	return true
}

var (
	L_CrabQueen = &Leader{}
	L_Longen    = &Leader{Name: "Longen", Faction: F_TradersGuild, Location: T_TradersEdge, Notes: "Can be turned into Tinfist. Can reward for Tinfist. captured and/or killed Longen and talk to Tinfist, you will get +75 reputation with the anti-slavers instantly, without the -30 to the UC, Slave traders or Traders Guild. "}
	L_Nagata    = &Leader{Name: "Lord Nagata", Faction: F_UnitedCities, Location: T_ShoBattai}
	L_Preacher  = &Leader{}
	L_SMGrande  = &Leader{Name: "Slave Master Grande", Faction: F_TradersGuild, Location: T_Eyesocket}
	L_SMMaster  = &Leader{}
	L_SMRen     = &Leader{}
	L_SMRuben   = &Leader{}
	L_SMWada    = &Leader{}
	L_Tengu     = &Leader{Name: "Emperor Tengu", Faction: F_UnitedCities, Location: T_Heft}
	L_Tinfist   = &Leader{Name: "Tinfist", Faction: F_AntiSlavers, Location: T_Spring}
	L_Tsugi     = &Leader{}
	L_Valamon   = &Leader{}
	L_Wada      = &Leader{}
	L_Yoshinaga = &Leader{Name: "Lord Yoshinaga", Faction: F_UnitedCities, Location: T_Heng, Notes: "No world States. if Heng is half-destroyed or handed over to Empire Peasants, Lord Yoshinaga will despawn"}
)

type Town struct {
	Name       string
	Location   *Town
	Faction    *Faction
	Overrides  []*Town
	WorldState []Cond
	Notes      string
}

var (
	T_CultVillage    = &Town{}
	T_Eyesocket      = &Town{}
	T_FreeSettlement = &Town{}
	T_FishingVillage = &Town{}
	T_Heft           = &Town{Name: "Heft"}
	T_Heng           = &Town{Name: "Heng", Faction: F_UnitedCities, Notes: "if Heng is half-destroyed or handed over to Empire Peasants, Lord Yoshinaga will despawn"}
	T_SettledNomads  = &Town{}
	T_ShoBattai      = &Town{Name: "Sho-Battai"}
	T_Spring         = &Town{Name: "Spring"}
	T_TradersEdge    = &Town{Name: "Traders-Edge"}
)

type DesiredTownLogType uint8

const (
	AlwaysLog DesiredTownLogType = iota
	LogIfNotFirstPick
)

type DesiredTown struct {
	Towns   []*Town
	LogType DesiredTownLogType
	Notes   string
}

type World struct {
	//Leaders []*Leader
	Towns          []*Town
	DesiredTownMap map[*Town]DesiredTown
	Notes          []string
	//RegionNotes map[RegionID][]string
}

func (w *World) Seed() {
	w.SeedTowns()
}

func (w *World) SeedTowns() {
	w.Towns = append(w.Towns, T_Heng)
	T_Heng.Location = T_Heng
	HengProsperous := &Town{Name: "Heng Prosperous", Location: T_Heng, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), false),
		Equal(L_Tinfist.IsNotAlive(), false),
		Equal(L_Tengu.IsAlive(), false),
	}}
	HengOverride := &Town{Name: "Heng Override", Location: T_Heng, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), false),
		Equal(L_Tinfist.IsNotAlive(), true),
		Equal(L_Tengu.IsAlive(), false),
	}}
	HengHalfDestroyed43 := &Town{Name: "Heng Half Destroyed 43", Location: T_Heng, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), false),
		Equal(L_Tengu.IsAlive(), true),
	}}
	HengHalfDestroyed44 := &Town{Name: "Heng Half Destroyed 44", Location: T_Heng, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), true),
		Equal(L_Tengu.IsAlive(), false),
	}}
	HengMalnourished := &Town{Name: "HengMalnourished", Location: T_Heng, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_SMGrande.IsAlive(), false),
		Equal(L_SMRuben.IsAlive(), false),
		Equal(L_SMRen.IsAlive(), false),
	}}
	T_Heng.Overrides = append(T_Heng.Overrides, HengProsperous, HengOverride, HengHalfDestroyed43, HengHalfDestroyed44, HengMalnourished)
	HengHalfDestroyed43.Overrides = append(HengHalfDestroyed43.Overrides, HengProsperous, HengOverride, HengMalnourished)
	HengHalfDestroyed44.Overrides = append(HengHalfDestroyed44.Overrides, HengOverride, HengProsperous, HengMalnourished)
	HengMalnourished.Overrides = append(HengMalnourished.Overrides, HengProsperous, HengOverride)
	w.DesiredTownMap[T_Heng] = DesiredTown{Towns: []*Town{HengProsperous, HengOverride}, LogType: LogIfNotFirstPick, Notes: "Tinfist has to die for Override"}

}

func (w *World) UpdateWorldStates() {
	bChanged := false
	for i, t := range w.Towns {
		var matches []*Town
		for _, o := range t.Overrides {
			// check if all override conditions are true
			if allTrue(o.WorldState) {
				matches = append(matches, o)
			}
			// if all overrides are true update the town's override
		}

		if len(matches) > 1 {
			panic("Too many town overrides for " + t.Name)
		}

		if len(matches) == 1 {
			nextTown := matches[0]
			if w.Towns[i] != nextTown {
				if nextTown.Location == nil {
					nextTown.Location = t.Location
				}
				if nextTown.Faction == nil {
					nextTown.Faction = t.Location.Faction
				}

				w.Towns[i] = nextTown
				fmt.Print(nextTown.Location.Name + ": " + t.Name + " -> " + nextTown.Name + ". want: ")
				dtv, ok := w.DesiredTownMap[t.Location]
				if ok {
					for _, dt := range dtv.Towns {
						fmt.Print(dt.Name + ", ")
					}
				}
				fmt.Print("\t", nextTown.Notes+".", t.Location.Notes)
				fmt.Println()
				bChanged = true
			}
		}
	}

	// handle special events
	if L_Yoshinaga.Status == Alive {
		bHengChanged := false
		for _, t := range w.Towns {
			if t.Location == T_Heng {
				if t != T_Heng {
					bHengChanged = true
				}
				break
			}
		}
		if L_Longen.Status != Alive || bHengChanged {
			L_Yoshinaga.Status = Dead
			fmt.Println("WARN", L_Yoshinaga.Name, "has disappeared")
		}
	}

	if bChanged {
		w.UpdateWorldStates()
	}
}

func (w *World) CompareDesiredWorldStates() {
	for _, t := range w.Towns {
		dv, ok := w.DesiredTownMap[t.Location]
		if ok {
			bDesiredTown := false
			for _, dt := range dv.Towns {
				if dt == t {
					fmt.Println("✓ "+t.Name, dv.Notes, t.Notes, t.Location.Notes)
					bDesiredTown = true
				}
			}
			if !bDesiredTown {
				fmt.Println(t.Location.Name + ": is " + t.Name + " want:")
				dtMsgs := ""
				for _, dt := range dv.Towns {
					fmt.Print(" " + dt.Name + ",")
					dtMsgs += dt.Notes
				}
				fmt.Println()
				fmt.Println(dv.Notes, t.Location.Notes, dtMsgs)
				//bCanBecomeDesired := false todo
			}
		} else {
			fmt.Println(t.Name + " Doesn't have a desired town.")
		}
	}
}

func (w *World) Capture(ldr *Leader) {
	if ldr.Status == Alive {
		prev := ldr.Status
		ldr.Status = Imprisoned
		fmt.Println(ldr.Name+":", prev.String(), "->", ldr.Status.String()+". \t", ldr.Notes)
		w.UpdateWorldStates()
	} else {
		fmt.Println("WARN", ldr.Name, "is", ldr.Status.String()+".", "Only people that are alive can be Imprisoned")
	}
}

func main() {
	w := &World{DesiredTownMap: make(map[*Town]DesiredTown)}
	w.Seed()

	// start up. Capture Preacher?
	w.Capture(L_Yoshinaga)
	w.Capture(L_Longen)
	fmt.Println("Maintain positive/non-hostile relations with UC, TG, ST, and all factions.")
	fmt.Println("Talk to Tinfist while carrying Longen. Turn Longen into Grey. Ask Boss Simion to join Anti-Slavers. That order gives the best Rep in my testing. -25 with ST&UC !TG; +100AS +75Rebel.")
	// kill longen to prevent empty settled nomads. Unlocking his cage is enough for AS to kill him?

	w.Capture(L_SMMaster)

}
