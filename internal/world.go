package internal

import "fmt"

type DesiredTownLogType uint8

const (
	AlwaysLog DesiredTownLogType = iota
	LogIfNotFirstPick
)

type DesiredTown struct {
	// will be just one town in most cases. Can accept multiple desired towns
	Towns   []*Town
	LogType DesiredTownLogType
	Notes   string
}

type World struct {
	//Leaders []*Leader
	// current towns in the game. doesn't include non-existent overrides
	Towns []*Town
	// This is where I'm storing the overrides that I want.
	DesiredTownMap map[*Town]DesiredTown
	Notes          []string
	//RegionNotes map[RegionID][]string
}

func (w *World) Seed() {
	w.SeedTowns()
}

func (w *World) Capture(ldr *Leader) {
	if ldr.Status == Alive {
		prev := ldr.Status
		ldr.Status = Imprisoned

		{ // log
			var loc string
			if ldr.Home != nil {
				loc = ldr.Home.Name
			} else {
				loc = "??"
			}
			fmt.Println("[Leader]", ldr.Name+"("+loc+")"+":", prev.String(), "->", ldr.Status.String()+". \t", ldr.Notes)
		}

		w.UpdateWorldStates()
	} else {
		fmt.Println("WARN", ldr.Name, "is", ldr.Status.String()+".", "Only people that are alive can be Imprisoned")
	}
}

// you shouldn't have to call this in main(). It should be called by functions that can change world state.
// handles special events like Simon or Yoshinaga disapearing
func (w *World) UpdateWorldStates() {
	bChanged := false
	for i, t := range w.Towns {
		var matches []*Town
		// check if all override conditions are true
		for _, o := range t.Overrides {
			if allTrue(o.WorldState) {
				matches = append(matches, o)
			}
		}
		// if all overrides are true update the town's override

		if len(matches) > 1 {
			panic("Too many town overrides for " + t.Name)
		}

		if len(matches) == 1 {
			nextTown := matches[0]
			if w.Towns[i] != nextTown {
				{ // update next town in the event it's missing some data
					if nextTown.Location == nil {
						nextTown.Location = t.Location
					}
					if nextTown.Faction == nil {
						nextTown.Faction = t.Location.Faction
					}
				}

				w.Towns[i] = nextTown
				{ // log the town change
					fmt.Print("[Town] ", nextTown.Location.Name+": "+t.Name+" -> "+nextTown.Name+". want: ")
					dtv, ok := w.DesiredTownMap[t.Location]
					if ok {
						for _, dt := range dtv.Towns {
							fmt.Print(dt.Name + ", ")
						}
					}
					fmt.Print("\t", nextTown.Notes+".", t.Location.Notes)
					fmt.Println()
				}
				bChanged = true
			}
		}
	}

	// handle special events
	if L_LdYoshinaga.Status == Alive {
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
			L_LdYoshinaga.Status = Dead
			fmt.Println("WARN", L_LdYoshinaga.Name, "has disappeared")
		}
	}
	// todo boss simion dissapear?

	if bChanged {
		w.UpdateWorldStates()
	}
}

// debug to see check the world state.
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

func (w *World) SeedTowns() {
	var towns []*Town

	w.Towns = append(w.Towns, T_Heng)
	T_Heng.Location = T_Heng
	towns = heng()
	w.DesiredTownMap[T_Heng] = DesiredTown{Towns: towns, LogType: LogIfNotFirstPick, Notes: "Tinfist has to die for Override"}

}

func heng() (desiredTowns []*Town) {
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

	return []*Town{HengProsperous, HengOverride}
}
