package internal

import "fmt"

type DesiredTownLogType uint8

const (
	AlwaysLog DesiredTownLogType = iota
	LogIfNotFirstPick
	EitherIsFine
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
	w.verify()
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

// UpdateWorldStates Shouldn't be called in main(). It should be called by functions that can change world state.
// also handles special events like Simon or Yoshinaga disappearing
func (w *World) UpdateWorldStates() {
	w.handleSpecialEvents()
	bWordChanged := false
	for i, t := range w.Towns {
		var overrideMatches []*Town
		// check if all override conditions are true
		for _, o := range t.Overrides {
			if allTrue(o.WorldState) {
				overrideMatches = append(overrideMatches, o)
			}
		}

		// if all overrides are true update the town's override
		var overrideMatch *Town
		if len(overrideMatches) == 1 {
			overrideMatch = overrideMatches[0]
		}

		// resolve town in the event there's more than one matching override
		if len(overrideMatches) > 1 {
			// log that the conflict should be verified.
			fmt.Println("WARN Multiple overrides possible for", t.Name+". Attempting to resolve. Override with the greatest num of world stats wins.")
			largest := 0
			// find the town with the greater number of worldstates
			for _, o := range overrideMatches {
				fmt.Print(o.Name+".NumWorldStates=", len(o.WorldState), ", ")
				l := len(o.WorldState)
				if l > largest {
					largest = l
					overrideMatch = o
				}
			}
			// if there is a tie with greatest num of worldstate, panic.
			for _, o := range overrideMatches {
				if o != overrideMatch {
					l := len(o.WorldState)
					if l >= largest && o != overrideMatch {
						panic("Too many town overrides for " + t.Name)
					}
				}
			}
			fmt.Println("Override conflict resolved. The winner is", overrideMatch.Name, "Verify result in game")
		}

		if overrideMatch == w.Towns[i] {
			panic("OverrideMatch == the current town. How?")
		}

		if overrideMatch != nil {
			w.Towns[i] = overrideMatch
			// log the town change
			{
				fmt.Print("[Town] ", overrideMatch.MapLoc.Name+": "+t.Name+" -> "+overrideMatch.Name+". want: ")
				dtv, ok := w.DesiredTownMap[t.MapLoc]
				if ok {
					for _, dt := range dtv.Towns {
						fmt.Print(dt.Name + ", ")
					}
				}
				fmt.Print("\t", overrideMatch.Notes+".", t.MapLoc.Notes)
				fmt.Println()
			}
			bWordChanged = true
			w.handleSpecialEvents()
		}
	}

	// if there were changes, then check again in the event a town can override more than once
	if bWordChanged {
		w.UpdateWorldStates()
	}
}

// hard-coded & one-off events like a person disappearing
func (w *World) handleSpecialEvents() {
	// todo add to world log.
	// Lord Yoshinaga can disappear depending on Heng or Longen
	if L_LdYoshinaga != nil && L_LdYoshinaga.Status == Alive {
		bHengChanged := false
		for _, t := range w.Towns {
			if t.MapLoc == T_Heng {
				if t != T_Heng {
					bHengChanged = true
				}
				break
			}
		}
		if L_Longen.Status != Alive || bHengChanged {
			name := L_LdYoshinaga.Name
			L_LdYoshinaga = nil
			fmt.Println("WARN", name, "has disappeared")
		}
	}
	// todo boss simion disappeared?
	if L_BossSimion != nil {
		if L_Tengu.Status != Alive {
			if L_Longen.Status != Alive && L_Tinfist.Status != Alive {

			} else {
				fmt.Println("WARN", L_BossSimion.Name, "has disappeared")
				L_BossSimion = nil
			}
		}
	}
}

// CompareDesiredWorldStates is a debug to see check the world state.
func (w *World) CompareDesiredWorldStates() {
	fmt.Println()
	fmt.Println("*** List of current overrides ***")
	for _, t := range w.Towns {
		dv, ok := w.DesiredTownMap[t.MapLoc]
		if ok {
			dtStr := ""
			bDesiredTown := false
			for _, dt := range dv.Towns {
				if dt == t {
					fmt.Println("✓", "["+t.MapLoc.Name+"]", t.Name+". \t"+dv.Notes, t.Notes, t.MapLoc.Notes)
					// todo: can it override? has overrides and they can trigger look at settled nomads for example
					bDesiredTown = true
					dtStr += dt.Name
				}
			}
			if bDesiredTown {
				// print a list of other desired towns if there's multiple
				if len(dv.Towns) > 1 {
					str := ""
					for i, dt := range dv.Towns {
						str += dt.Name
						if i != len(dv.Towns)-1 {
							str += ", "
						}
					}
					fmt.Println("\tDesiredTowns=[" + str + "]")
				}
				// can this town override?
				if len(t.Overrides) > 1 {
					str := ""
					for i, dt := range t.Overrides {
						str += dt.Name
						if i != len(t.Overrides)-1 {
							str += ", "
						}
					}
					fmt.Println("\tCan still override to [" + str + "]")
				}
			}
			if !bDesiredTown {
				fmt.Println("x", "["+t.MapLoc.Name+"]", t.Name+". want:")
				dtMsgs := ""
				fmt.Print("\t ")
				for _, dt := range dv.Towns {
					fmt.Print(dt.Name + ",")
					dtMsgs += dt.Notes
				}
				fmt.Println()
				if dv.Notes != "" || t.MapLoc.Notes != "" || dtMsgs != "" {
					fmt.Println("\t", dv.Notes, t.MapLoc.Notes, dtMsgs)
				}

				// print list of override towns
				if len(t.Overrides) > 0 {
					str := ""
					for i, dt := range t.Overrides {
						str += dt.Name
						if i != len(t.Overrides)-1 {
							str += ", "
						}
					}
					fmt.Println("\t", t.Name+".Overrides=["+str+"]")
				} else {
					fmt.Println("\t", t.Name, "cannot override.")
				}
			}
		} else {
			fmt.Println("✓", "["+t.MapLoc.Name+"]", t.Name, "No desired town.", dv.Notes, t.Notes, t.MapLoc.Notes)
		}
	}
}

func (w *World) GetTowns() {
	fmt.Println()
	fmt.Println("***Towns***")
	for _, t := range w.Towns {
		overridesStr := ""
		for i, o := range t.Overrides {
			overridesStr += o.Name
			if i < len(t.Overrides)-1 {
				overridesStr += ", "
			}
		}
		fmt.Println(t.Name)
		fmt.Println("\tOverrides=[" + overridesStr + "]")
		v, ok := w.DesiredTownMap[t.MapLoc]
		if ok {
			str := ""
			for i, dt := range v.Towns {
				str += dt.Name
				if i < len(v.Towns)-1 {
					str += ", "
				}
			}
			fmt.Println("\tPrefered=[" + str + "]")
			if v.Notes != "" {
				fmt.Println("\t" + v.Notes)
			}
		} else {
			fmt.Println("\tPrefered=[]")
		}
		if t.Notes != "" {
			fmt.Println("\t" + t.Notes)
		}
	}
	fmt.Println()
}

func (w *World) PrintTown(t *Town) {
	fmt.Println("***")
	if t == nil {
		fmt.Println("Town is nil")
		return
	}

	var currentTown *Town
	for _, town := range w.Towns {
		if town.MapLoc == t {
			currentTown = town
		}
	}

	fmt.Println("["+t.MapLoc.Name+"]", currentTown.Name)
	v, ok := w.DesiredTownMap[t.MapLoc]
	if ok {
		str := ""
		for i, dt := range v.Towns {
			str += dt.Name
			if i < len(v.Towns)-1 {
				str += ", "
			}
		}
		fmt.Println("\tDesiredTowns=[" + str + "]")
	} else {
		fmt.Println("\tNo desired towns")
	}

	// list the current towns overrides
	str := ""
	for i, o := range currentTown.Overrides {
		str += o.Name
		if i < len(currentTown.Overrides)-1 {
			str += ", "
		}
	}
	fmt.Println("\t", currentTown.Name+".Overrides=["+str+"]")

	fmt.Println("***")
}

// verify that the towns are set up as expected to avoid issues down the line.
func (w *World) verify() {
	// todo check that non of the overrides have the same world states
	for _, t := range w.Towns {
		if t.Faction == nil {
			panic("Faction not set")
		}
		if t.MapLoc != nil {
			for _, o := range t.Overrides {
				if o.Faction == nil {
					panic("Faction not set")
				}
				if o.MapLoc != t.MapLoc {
					panic("MapLoc mismatch.")
				}
				if o == t {
					panic("Override == the current town.")
				}
				if t.Name == o.Name {
					panic("Override is using the town name.")
				}
				for _, o2 := range o.Overrides {
					if o2 == o {
						panic("Override == the current override town.")
					}
					if o2.Name == o.Name {
						panic("Override is using the same name.")
					}
				}
			}
		} else {
			panic("MapLoc is nil")
		}
	}
}

// Only add main towns that can be locations at the start.
func (w *World) initTown(t *Town) {
	w.Towns = append(w.Towns, t)
	t.MapLoc = t
}

// GuidedStep2 seed the town.
// SeedTowns fleshes out the towns' info
func (w *World) SeedTowns() {
	var towns []*Town

	w.initTown(T_Bark)
	towns = bark()
	w.DesiredTownMap[T_Bark] = DesiredTown{Towns: towns, LogType: EitherIsFine, Notes: "the two prosperous are 99% the same town."}

	w.initTown(T_Brink)
	towns = brink()
	w.DesiredTownMap[T_Brink] = DesiredTown{Towns: towns, LogType: AlwaysLog, Notes: "No override is prolly best. Reaver takeover is interesting. Malnourished is not bad."}

	w.initTown(T_Catun)
	towns = catun()
	w.DesiredTownMap[T_Catun] = DesiredTown{Towns: towns}

	w.initTown(T_ClownSteady)
	towns = clownSteady()
	w.DesiredTownMap[T_ClownSteady] = DesiredTown{Towns: towns}

	w.initTown(T_CultVillage)
	towns = cultVillage()
	w.DesiredTownMap[T_CultVillage] = DesiredTown{Towns: towns, Notes: "Slaver takeover sounds better than cultist. but I would have to capture Tinfist while Longen and Grande are free."}

	w.initTown(T_DistantHiveVillage)
	towns = distantHiveVillage()
	w.DesiredTownMap[T_DistantHiveVillage] = DesiredTown{Towns: towns}

	w.initTown(T_DriftersLast)
	towns = driftersLast()
	w.DesiredTownMap[T_DriftersLast] = DesiredTown{Towns: towns}

	w.initTown(T_Drin)
	towns = drin()
	w.DesiredTownMap[T_Drin] = DesiredTown{Towns: towns, Notes: "At least in one of my play through, the HN version was more than I expected."}

	w.initTown(T_Eyesocket)
	towns = eyesocket()
	w.DesiredTownMap[T_Eyesocket] = DesiredTown{Towns: towns}

	w.initTown(T_FreeSettlement)
	towns = freeSettlement()
	w.DesiredTownMap[T_FreeSettlement] = DesiredTown{Towns: towns}

	w.initTown(T_FishingVillage)
	towns = fishingVillage()
	w.DesiredTownMap[T_FishingVillage] = DesiredTown{Towns: towns, Notes: "There's two Fishing Villages in the FCS. There's more info on the wiki."}

	w.initTown(T_FortSimion)
	towns = fortSimion()
	w.DesiredTownMap[T_FortSimion] = DesiredTown{Towns: towns, Notes: "Can I use the fort if it's abandoned? If not, It doesn't matter what happens to the fort."}

	w.initTown(T_Heft)
	towns = heft()
	w.DesiredTownMap[T_Heft] = DesiredTown{Towns: towns, LogType: LogIfNotFirstPick, Notes: "The original version of the city has many shops. The Rebel override is in the aftermath of the revolution."}

	w.initTown(T_Heng)
	towns = heng()
	w.DesiredTownMap[T_Heng] = DesiredTown{Towns: towns, LogType: LogIfNotFirstPick, Notes: "Tinfist has to die for Override."}

	w.initTown(T_ManhunterBase)
	towns = manhunterBase()
	w.DesiredTownMap[T_ManhunterBase] = DesiredTown{Towns: towns, Notes: "Override has some unique shops worth checking out."}

	w.initTown(T_Mourn)
	T_Mourn.Notes = "Mourn has no overrides. It does have the Super adventure shop. It needs protection from beak things. Kill camps and bring down some mercs"

	w.initTown(T_PortNorth)
	towns = portNorth()
	w.DesiredTownMap[T_PortNorth] = DesiredTown{Towns: towns, Notes: "Destroying spawns the most Anti-Slaver presence."}

	w.initTown(T_PortSouth)
	towns = portSouth()
	w.DesiredTownMap[T_PortSouth] = DesiredTown{Towns: towns, Notes: "the overrides are bugged and should be modded."}

	w.initTown(T_SettledNomads)
	towns = settledNomads()
	w.DesiredTownMap[T_SettledNomads] = DesiredTown{Towns: towns, Notes: "Prevent death by clearing nearby animal camps"}

	w.initTown(T_ShoBattai)
	towns = shoBattai()
	w.DesiredTownMap[T_ShoBattai] = DesiredTown{Towns: towns, Notes: "Has unique hat shop. The Peasants' version is almost as nice as the original."}

	w.initTown(T_SlaveFarm)
	towns = slaveFarm()
	w.DesiredTownMap[T_SlaveFarm] = DesiredTown{Towns: towns}

	w.initTown(T_SlaveFarmSouth)
	towns = slaveFarmSouth()
	w.DesiredTownMap[T_SlaveFarmSouth] = DesiredTown{Towns: towns}

	w.initTown(T_SlaveMarkets) // give up for sm master?
	towns = slaveMarkets()
	w.DesiredTownMap[T_SlaveMarkets] = DesiredTown{Towns: towns, Notes: "I'm ok letting this one become destroyed. I never use it."}

	w.initTown(T_SouthStoneCamp)
	towns = southStoneCamp()
	w.DesiredTownMap[T_SouthStoneCamp] = DesiredTown{Towns: towns, Notes: "this one may have to be destroyed."}

	w.initTown(T_Spring)
	T_Spring.Notes = "Spring doesn't have overrides"

	w.initTown(T_Stoat)
	towns = stoat()
	w.DesiredTownMap[T_Stoat] = DesiredTown{Towns: towns, Notes: "the Tech hunter take over is good albeit a bit wild, but not as nice as the og town. I should clear out some camps nearby bc the tech hunters have broken gates."}

	w.initTown(T_StoneCamp)
	towns = stoneCamp()
	w.DesiredTownMap[T_StoneCamp] = DesiredTown{Towns: towns, Notes: "Destroyed is the only other option"}

	w.initTown(T_TradersEdge)
	towns = tradersEdge()
	w.DesiredTownMap[T_TradersEdge] = DesiredTown{Towns: towns, Notes: "Traders Edge Prosperous needs Tinfist alive."}
}

func bark() (desiredTowns []*Town) {
	BarkProsperous0 := &Town{Name: "Bark Prosperous 0", MapLoc: T_Bark, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), false),
		Equal(L_LdSanda.IsAlive(), false),
	}}
	BarkHalfDestroyed := &Town{Name: "Bark Half Destroyed", MapLoc: T_Bark, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_LdSanda.IsAlive(), false),
	}}
	BarkMalnourished := &Town{Name: "Bark Malnourished", MapLoc: T_Bark, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_LdSanda.IsAlive(), true),
		Equal(L_SMWada.IsAlive(), false),
		Equal(L_LdKana.IsAlive(), false),
	}}
	BarkProsperous2 := &Town{Name: "Bark Prosperous 2", MapLoc: T_Bark, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), true),
		Equal(L_LdSanda.IsAlive(), false),
		Equal(L_SMWada.IsAlive(), false),
		Equal(L_LdKana.IsAlive(), false),
	}}
	T_Bark.Overrides = append(T_Bark.Overrides, BarkProsperous0, BarkHalfDestroyed, BarkMalnourished, BarkProsperous2)
	BarkHalfDestroyed.Overrides = append(BarkHalfDestroyed.Overrides, BarkProsperous0)
	BarkMalnourished.Overrides = append(BarkMalnourished.Overrides, BarkProsperous0, BarkHalfDestroyed, BarkProsperous2)
	BarkProsperous2.Overrides = append(BarkProsperous2.Overrides)
	return []*Town{BarkProsperous0, BarkProsperous2}
}

func brink() (desiredTowns []*Town) {
	BrinkTakeover := &Town{Name: "Brink Takeover", MapLoc: T_Brink, Faction: F_Reavers, WorldState: []Cond{
		Equal(L_LdTsugi.IsNotAlive(), true),
		Equal(L_Valamon.IsAlive(), true),
	}}
	BrinkDestroyed := &Town{Name: "Brink Destroyed", MapLoc: T_Brink, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_LdTsugi.IsAlive(), false),
		Equal(L_Valamon.IsAlive(), false),
	}}
	BrinkMalnourished := &Town{Name: "Brink Malnourished", MapLoc: T_Brink, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_SMGrande.IsAlive(), false),
		Equal(L_LdTsugi.IsAlive(), true),
	}}
	T_Brink.Overrides = append(T_Brink.Overrides, BrinkTakeover, BrinkDestroyed, BrinkMalnourished)
	BrinkTakeover.Overrides = append(BrinkTakeover.Overrides, BrinkDestroyed)
	BrinkMalnourished.Overrides = append(BrinkMalnourished.Overrides, BrinkDestroyed, BrinkTakeover)
	return []*Town{T_Brink, BrinkTakeover, BrinkMalnourished}
}

func catun() (desiredTowns []*Town) {
	CatunFishmanTakeover := &Town{Name: "Catun Fishman Takeover", MapLoc: T_Catun, Faction: F_Fishmen, WorldState: []Cond{
		Equal(L_LdShiro.IsAlive(), false),
	}}
	T_Catun.Overrides = append(T_Catun.Overrides, CatunFishmanTakeover)
	return []*Town{T_Catun}
}

func clownSteady() (desiredTowns []*Town) {
	ClownSteadyMalnourished1 := &Town{Name: "Clown Steady Malnourished 1", MapLoc: T_ClownSteady, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_SMGrace.IsAlive(), false),
	}}
	ClownSteadyProsperous := &Town{Name: "Clown Steady Prosperous", MapLoc: T_ClownSteady, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_SMMaster.IsAlive(), false),
		Equal(L_SMGrace.IsAlive(), false),
	}}
	ClownSteadyMalnourished2 := &Town{Name: "Clown Steady Malnourished 2", MapLoc: T_ClownSteady, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_SMMaster.IsAlive(), false),
	}}
	T_ClownSteady.Overrides = append(T_ClownSteady.Overrides, ClownSteadyMalnourished1, ClownSteadyProsperous, ClownSteadyMalnourished2)
	ClownSteadyMalnourished1.Overrides = append(ClownSteadyMalnourished1.Overrides, ClownSteadyProsperous)
	ClownSteadyMalnourished2.Overrides = append(ClownSteadyMalnourished2.Overrides, ClownSteadyProsperous)
	return []*Town{ClownSteadyProsperous}
}

func cultVillage() (desiredTowns []*Town) {
	CultVillageSlaverTakeover := &Town{Name: "Cult Village Slaver Takeover", MapLoc: T_CultVillage, Faction: F_PreacherCult, WorldState: []Cond{
		Equal(L_Tinfist.IsAlive(), false),
		Equal(L_Longen.IsAlive(), true),
		Equal(L_SMGrande.IsAlive(), true),
	}}
	T_CultVillage.Overrides = append(T_CultVillage.Overrides, CultVillageSlaverTakeover)
	CultVillageSlaverTakeover.Overrides = append(CultVillageSlaverTakeover.Overrides)
	return []*Town{CultVillageSlaverTakeover, T_CultVillage}
}

func distantHiveVillage() (desiredTowns []*Town) {
	DeadhiveOverrunWest := &Town{Name: "Deadhive Overrun West", MapLoc: T_DistantHiveVillage, Faction: F_Fogmen, WorldState: []Cond{
		Equal(L_WestHiveQueen.Okay(), false),
	}}
	HiveVillageNEmpty := &Town{Name: "Hive Village N Empty", MapLoc: T_DistantHiveVillage, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_Tinfist.IsAlive(), false),
		Equal(L_LdKana.IsAlive(), true),
	}}
	T_DistantHiveVillage.Overrides = append(T_DistantHiveVillage.Overrides, DeadhiveOverrunWest, HiveVillageNEmpty)
	DeadhiveOverrunWest.Overrides = append(DeadhiveOverrunWest.Overrides, HiveVillageNEmpty)
	return []*Town{T_DistantHiveVillage}
}

func driftersLast() (desiredTowns []*Town) {
	DriftersLastProsperous := &Town{Name: "Drifters Last Prosperous", MapLoc: T_DriftersLast, Faction: F_FreeTraders, WorldState: []Cond{
		Equal(L_LdMerin.IsAlive(), false),
		Equal(L_SMGrace.IsAlive(), false),
	}}
	DriftersLastHalfDestroyed := &Town{Name: "Drifters Last Half Destroyed", MapLoc: T_DriftersLast, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_LdMerin.IsAlive(), false),
	}}
	T_DriftersLast.Overrides = append(T_DriftersLast.Overrides, DriftersLastProsperous, DriftersLastHalfDestroyed)
	DriftersLastHalfDestroyed.Overrides = append(DriftersLastHalfDestroyed.Overrides, DriftersLastProsperous)
	return []*Town{DriftersLastProsperous}
}

func drin() (desiredTowns []*Town) {
	DrinOverride := &Town{Name: "Drin Override", MapLoc: T_Drin, Faction: F_HolyNation, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), false),
		Equal(L_LdInaba.IsAlive(), false),
		Equal(L_Valtena.IsAlive(), true),
	}}
	DrinDestroyed := &Town{Name: "Drin Destroyed", MapLoc: T_Drin, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), false),
		Equal(L_LdInaba.IsAlive(), false),
		Equal(L_Valtena.IsAlive(), false),
	}}
	T_Drin.Overrides = append(T_Drin.Overrides, DrinOverride, DrinDestroyed)
	DrinOverride.Overrides = append(DrinOverride.Overrides, DrinDestroyed)
	return []*Town{DrinOverride, T_Drin}
}

func eyesocket() (desiredTowns []*Town) {
	EyesocketDestroyed := &Town{Name: "Eyesocket Destroyed", MapLoc: T_Eyesocket, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_SMGrande.IsAlive(), false),
	}}
	T_Eyesocket.Overrides = append(T_Eyesocket.Overrides, EyesocketDestroyed)
	return []*Town{T_Eyesocket}
}

func freeSettlement() (desiredTowns []*Town) {
	FreeCity := &Town{Name: "Free City", MapLoc: T_FreeSettlement, Faction: F_AntiSlavers, WorldState: []Cond{
		Equal(L_Tinfist.IsNotAlive(), false),
		Equal(L_Tengu.Ok(), false),
		Equal(L_Longen.IsAlive(), false),
	}}
	T_FreeSettlement.Overrides = append(T_FreeSettlement.Overrides, FreeCity)
	return []*Town{FreeCity}
}

func fishingVillage() (desiredTowns []*Town) {
	FishingVillageSlaverTakeover := &Town{Name: "Fishing Village Slaver Takeover", MapLoc: T_FishingVillage, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_Tinfist.IsAlive(), false),
		Equal(L_Longen.IsAlive(), true),
		Equal(L_Tengu.IsAlive(), true),
	}}
	T_FishingVillage.Overrides = append(T_FishingVillage.Overrides, FishingVillageSlaverTakeover)
	CultTownDestroyed := &Town{Name: "Cult Town Destroyed", MapLoc: T_FishingVillage, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_SMGrande.IsAlive(), false),
		Equal(L_Longen.IsAlive(), false),
	}}
	FishingVillageSlaverTakeover.Overrides = append(FishingVillageSlaverTakeover.Overrides, CultTownDestroyed)
	return []*Town{T_FishingVillage}
}

func fortSimion() (desiredTowns []*Town) {
	FortYabuta := &Town{Name: "Fort Yabuta", MapLoc: T_FortSimion, Faction: F_YabutaOutlaws, WorldState: []Cond{
		Equal(L_Yabuta.IsAlive(), true),
		// (note) the actual world state is not allied with rebel farmers. that's outside the scope of this app considering I'll never go down this path.
		Equal(L_BossSimion.IsAlive(), false),
	}}
	FortSimionDestroyed := &Town{Name: "Fort Simion Destroyed", MapLoc: T_FortSimion, Faction: F_RebelSwordsmen, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), false),
	}}
	T_FortSimion.Overrides = append(T_FortSimion.Overrides, FortYabuta, FortSimionDestroyed)
	return []*Town{FortSimionDestroyed, FortYabuta, T_FortSimion}
}

func heft() (desiredTowns []*Town) {
	HeftOverride := &Town{Name: "Heft Override", MapLoc: T_Heft, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), false),
	}, Notes: "keeps best armor shop"}
	HeftMalnourished0 := &Town{Name: "Heft Malnourished 0", MapLoc: T_Heft, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), true),
		Equal(L_SMRen.IsAlive(), false),
		Equal(L_SMWada.IsAlive(), false),
	}}
	HeftMalnourished2 := &Town{Name: "Heft Malnourished 2", MapLoc: T_Heft, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), true),
		Equal(L_SMWada.IsAlive(), false),
		Equal(L_Longen.IsAlive(), false),
	}}
	HeftMalnourished3 := &Town{Name: "Heft Malnourished 3", MapLoc: T_Heft, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), true),
		Equal(L_SMWada.IsAlive(), false),
		Equal(L_SMRuben.IsAlive(), false),
	}}
	T_Heft.Overrides = append(T_Heft.Overrides, HeftOverride, HeftMalnourished0, HeftMalnourished2, HeftMalnourished3)
	HeftMalnourished0.Overrides = append(HeftMalnourished0.Overrides, HeftOverride)
	HeftMalnourished2.Overrides = append(HeftMalnourished2.Overrides, HeftOverride)
	HeftMalnourished3.Overrides = append(HeftMalnourished3.Overrides, HeftOverride)
	return []*Town{T_Heft, HeftOverride}
}

func heng() (desiredTowns []*Town) {
	HengProsperous := &Town{Name: "Heng Prosperous", MapLoc: T_Heng, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), false),
		Equal(L_Tinfist.IsNotAlive(), false),
		Equal(L_Tengu.IsAlive(), false),
	}}
	HengOverride := &Town{Name: "Heng Override", MapLoc: T_Heng, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), false),
		Equal(L_Tinfist.IsNotAlive(), true),
		Equal(L_Tengu.IsAlive(), false),
	}}
	HengHalfDestroyed43 := &Town{Name: "Heng Half Destroyed 43", MapLoc: T_Heng, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), false),
		Equal(L_Tengu.IsAlive(), true),
	}}
	HengHalfDestroyed44 := &Town{Name: "Heng Half Destroyed 44", MapLoc: T_Heng, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), true),
		Equal(L_Tengu.IsAlive(), false),
	}}
	HengMalnourished := &Town{Name: "HengMalnourished", MapLoc: T_Heng, Faction: F_UnitedCities, WorldState: []Cond{
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

func manhunterBase() (desiredTowns []*Town) {
	ManhunterBaseOverride := &Town{Name: "Manhunter Base Override", MapLoc: T_ManhunterBase, Faction: F_MercenaryGuild, WorldState: []Cond{
		Equal(L_SMMaster.IsAlive(), false),
		Equal(L_SMGrace.IsAlive(), false),
	}}
	T_ManhunterBase.Overrides = append(T_ManhunterBase.Overrides, ManhunterBaseOverride)
	return []*Town{ManhunterBaseOverride}
}

func portNorth() (desiredTowns []*Town) {
	PortNorthDestroyed := &Town{Name: "Port North Destroyed", MapLoc: T_PortNorth, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_LdKana.IsAlive(), false),
	}}
	PortNorthFarmTown := &Town{Name: "Port North Farm Town", MapLoc: T_PortNorth, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_LdNagata.IsAlive(), false),
		Equal(L_LdKana.IsAlive(), false),
		Equal(L_Tengu.Ok(), false),
		Equal(L_BossSimion.IsAlive(), true),
	}}
	T_PortNorth.Overrides = append(T_PortNorth.Overrides, PortNorthDestroyed, PortNorthFarmTown)
	PortNorthDestroyed.Overrides = append(PortNorthDestroyed.Overrides, PortNorthFarmTown)
	return []*Town{PortNorthFarmTown}
}

func portSouth() (desiredTowns []*Town) {
	PortSouthDestroyed := &Town{Name: "Port South Destroyed", MapLoc: T_PortSouth, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_SMWada.IsAlive(), false),
	}}
	PortSouthFarmTown := &Town{Name: "Port South Farm Town", MapLoc: T_PortSouth, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_Tengu.Ok(), false),
		Equal(L_SMWada.IsAlive(), false),
		Equal(L_Longen.IsAlive(), false),
	}}
	T_PortSouth.Overrides = append(T_PortSouth.Overrides, PortSouthDestroyed, PortSouthFarmTown)

	bHasModFix := false
	if bHasModFix {
		PortSouthDestroyed.Overrides = append(PortSouthDestroyed.Overrides, PortSouthFarmTown)
	}

	return []*Town{PortSouthFarmTown}
}

func settledNomads() (desiredTowns []*Town) {
	SettledNomadsEmpty := &Town{Name: "Settled Nomads Empty", MapLoc: T_SettledNomads, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_Tinfist.IsAlive(), false),
		Equal(L_Longen.IsAlive(), true),
	}}
	T_SettledNomads.Overrides = append(T_SettledNomads.Overrides, SettledNomadsEmpty)
	return []*Town{T_SettledNomads}
}

func shoBattai() (desiredTowns []*Town) {
	ShoBattaiHalfDestroyed := &Town{Name: "Sho-Battai Half Destroyed", MapLoc: T_ShoBattai, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_LdNagata.IsAlive(), false),
		Equal(L_Tengu.IsAlive(), true),
	}}
	ShoBattaiCannibals := &Town{Name: "Sho-Battai Cannibals", MapLoc: T_ShoBattai, Faction: F_Cannibal, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), false),
		Equal(L_LdNagata.IsAlive(), false),
		Equal(L_BossSimion.IsAlive(), false),
	}}
	ShoBattaiProsperous := &Town{Name: "Sho-Battai Prosperous", MapLoc: T_ShoBattai, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_Tengu.Ok(), false),
		Equal(L_LdNagata.IsAlive(), false),
		Equal(L_BossSimion.IsAlive(), true),
	}}
	ShoBattaiMalnourished := &Town{Name: "Sho-Battai Malnourished", MapLoc: T_ShoBattai, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), true),
		Equal(L_LdNagata.IsAlive(), true),
		Equal(L_LdKana.IsNotAlive(), true),
	}}
	ShoBattaiMalnourishedHalfDestroyed0 := &Town{Name: "Sho-Battai Malnourished & Half Destroyed 0", MapLoc: T_ShoBattai, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), true),
		Equal(L_LdNagata.IsAlive(), false),
		Equal(L_LdKana.IsNotAlive(), true),
	}}
	ShoBattaiMalnourishedHalfDestroyed2 := &Town{Name: "Sho-Battai Malnourished & Half Destroyed 2", MapLoc: T_ShoBattai, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), false),
		Equal(L_LdNagata.IsAlive(), true),
		Equal(L_LdKana.IsNotAlive(), true),
	}}
	T_ShoBattai.Overrides = append(T_ShoBattai.Overrides, ShoBattaiHalfDestroyed, ShoBattaiCannibals, ShoBattaiProsperous, ShoBattaiMalnourished, ShoBattaiMalnourishedHalfDestroyed0, ShoBattaiMalnourishedHalfDestroyed2)
	ShoBattaiHalfDestroyed.Overrides = append(ShoBattaiHalfDestroyed.Overrides, ShoBattaiCannibals, ShoBattaiProsperous, ShoBattaiMalnourishedHalfDestroyed0)
	ShoBattaiMalnourished.Overrides = append(ShoBattaiMalnourished.Overrides, ShoBattaiCannibals, ShoBattaiProsperous, ShoBattaiMalnourishedHalfDestroyed0, ShoBattaiMalnourishedHalfDestroyed2)
	ShoBattaiMalnourishedHalfDestroyed0.Overrides = append(ShoBattaiMalnourishedHalfDestroyed0.Overrides, ShoBattaiCannibals, ShoBattaiProsperous)
	ShoBattaiMalnourishedHalfDestroyed2.Overrides = append(ShoBattaiMalnourishedHalfDestroyed2.Overrides, ShoBattaiCannibals, ShoBattaiProsperous)

	bAmIOkWithOgTown := false
	if bAmIOkWithOgTown {
		return []*Town{T_ShoBattai, ShoBattaiProsperous}
	}
	return []*Town{ShoBattaiProsperous}
}

func slaveFarm() (desiredTowns []*Town) {
	SlaveFarmDestroyed := &Town{Name: "Slave Farm Destroyed", MapLoc: T_SlaveFarm, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_SMRen.IsAlive(), false),
	}}
	T_SlaveFarm.Overrides = append(T_SlaveFarm.Overrides, SlaveFarmDestroyed)
	return []*Town{T_SlaveFarm}
}

func slaveFarmSouth() (desiredTowns []*Town) {
	SlaveFarmSDestroyed := &Town{Name: "Slave Farm South Destroyed", MapLoc: T_SlaveFarmSouth, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_SMGrace.IsAlive(), false),
	}}
	SlaveFarmSFarmTown := &Town{Name: "Slave Farm South Farm Town", MapLoc: T_SlaveFarmSouth, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_SMMaster.IsAlive(), false),
		Equal(L_SMGrace.IsAlive(), false),
	}}
	T_SlaveFarmSouth.Overrides = append(T_SlaveFarmSouth.Overrides, SlaveFarmSDestroyed, SlaveFarmSFarmTown)
	SlaveFarmSDestroyed.Overrides = append(SlaveFarmSDestroyed.Overrides, SlaveFarmSFarmTown)
	return []*Town{SlaveFarmSFarmTown}
}

func slaveMarkets() (desiredTowns []*Town) {
	SlaveMarketDestroyed := &Town{Name: "Slave market Destroyed", MapLoc: T_SlaveMarkets, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_SMMaster.IsAlive(), false),
	}}
	T_SlaveMarkets.Overrides = append(T_SlaveMarkets.Overrides, SlaveMarketDestroyed)
	return []*Town{T_SlaveMarkets}
}

func southStoneCamp() (desiredTowns []*Town) {
	SouthStoneCampDestroyed := &Town{Name: "South Stone Camp Destroyed", MapLoc: T_SouthStoneCamp, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_SMRuben.IsAlive(), false),
	}}
	T_SouthStoneCamp.Overrides = append(T_SouthStoneCamp.Overrides, SouthStoneCampDestroyed)
	return []*Town{T_SouthStoneCamp}
}

func stoat() (desiredTowns []*Town) {
	StoatHalfDesttroyed29 := &Town{Name: "Stoat Half Desttroyed 29", MapLoc: T_Stoat, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_LdInaba.IsAlive(), false),
		Equal(L_Tengu.IsAlive(), true),
	}}
	StoatHalfDesttroyed30 := &Town{Name: "Stoat Half Desttroyed 30", MapLoc: T_Stoat, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), false),
		Equal(L_LdInaba.IsAlive(), true),
	}}
	StoatYabuta := &Town{Name: "Stoat Yabuta", MapLoc: T_Stoat, Faction: F_YabutaOutlaws, WorldState: []Cond{
		Equal(L_LdInaba.IsAlive(), false),
		Equal(L_Yabuta.IsAlive(), true),
		Equal(L_Tengu.IsAlive(), false),
	}}
	StoatTakeover := &Town{Name: "Stoat Takeover", MapLoc: T_Stoat, Faction: F_TechHunters, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), false),
		Equal(L_LdInaba.IsAlive(), false),
		Equal(L_Yabuta.IsAlive(), false),
	}}
	StoatMalnourished := &Town{Name: "Stoat Malnourished", MapLoc: T_Stoat, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), true),
		Equal(L_LdInaba.IsAlive(), true),
		Equal(L_SMHaga.IsAlive(), false),
		Equal(L_SMRen.IsAlive(), false),
	}}
	StoatMalnourishedHalfDestroyed0 := &Town{Name: "Stoat Malnourished & Half Destroyed 0", MapLoc: T_Stoat, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), false),
		Equal(L_LdInaba.IsAlive(), true),
		Equal(L_SMHaga.IsAlive(), false),
		Equal(L_SMRen.IsAlive(), false),
	}}
	StoatMalnourishedHalfDestroyed2 := &Town{Name: "Stoat Malnourished & Half Destroyed 2", MapLoc: T_Stoat, Faction: F_UnitedCities, WorldState: []Cond{
		Equal(L_Tengu.IsAlive(), true),
		Equal(L_LdInaba.IsAlive(), false),
		Equal(L_SMHaga.IsAlive(), false),
		Equal(L_SMRen.IsAlive(), false),
	}}
	T_Stoat.Overrides = append(T_Stoat.Overrides, StoatHalfDesttroyed29, StoatHalfDesttroyed30, StoatYabuta, StoatTakeover, StoatMalnourished, StoatMalnourishedHalfDestroyed0, StoatMalnourishedHalfDestroyed2)
	StoatHalfDesttroyed29.Overrides = append(StoatHalfDesttroyed29.Overrides, StoatYabuta, StoatTakeover, StoatMalnourishedHalfDestroyed0, StoatMalnourishedHalfDestroyed2)
	StoatHalfDesttroyed30.Overrides = append(StoatHalfDesttroyed30.Overrides, StoatYabuta, StoatTakeover, StoatMalnourishedHalfDestroyed0, StoatMalnourishedHalfDestroyed2)
	StoatMalnourished.Overrides = append(StoatMalnourished.Overrides, StoatYabuta, StoatTakeover, StoatMalnourishedHalfDestroyed0, StoatMalnourishedHalfDestroyed2)
	StoatMalnourishedHalfDestroyed0.Overrides = append(StoatMalnourishedHalfDestroyed0.Overrides, StoatTakeover, StoatYabuta)
	StoatMalnourishedHalfDestroyed2.Overrides = append(StoatMalnourishedHalfDestroyed2.Overrides, StoatTakeover, StoatYabuta)
	return []*Town{T_Stoat, StoatTakeover}
}

func stoneCamp() (desiredTowns []*Town) {
	StoneCampDestroyed := &Town{Name: "Stone Camp Destroyed", MapLoc: T_StoneCamp, Faction: F_SlaveTraders, WorldState: []Cond{
		Equal(L_SMHaga.IsAlive(), false),
	}}
	T_StoneCamp.Overrides = append(T_StoneCamp.Overrides, StoneCampDestroyed)
	return []*Town{T_StoneCamp}
}

func tradersEdge() (desiredTowns []*Town) {
	TradersEdgeProsperous := &Town{Name: "Trade's Edge Prosperous", MapLoc: T_TradersEdge, Faction: F_AntiSlavers, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), false),
		Equal(L_Tinfist.IsNotAlive(), false),
		Equal(L_Tengu.IsAlive(), false),
	}}
	TradersEdgeOverride := &Town{Name: "Trade's Edge Override", MapLoc: T_TradersEdge, Faction: F_EmpirePeasants, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), false),
		Equal(L_Tinfist.IsNotAlive(), true),
		Equal(L_Tengu.IsAlive(), false),
	}}
	TradersEdgeDestroyed := &Town{Name: "Trade's Edge Destroyed", MapLoc: T_TradersEdge, Faction: F_TradersGuild, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), false),
		Equal(L_Tengu.IsAlive(), true),
	}}
	TradersEdgeMalnourished0 := &Town{Name: "Trade's Edge Malnourished 0", MapLoc: T_TradersEdge, Faction: F_TradersGuild, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), true),
		Equal(L_SMGrande.IsAlive(), false),
		Equal(L_SMRuben.IsAlive(), false),
		Equal(L_SMRen.IsAlive(), false),
	}}
	TradersEdgeMalnourished1 := &Town{Name: "Trade's Edge Malnourished 1", MapLoc: T_TradersEdge, Faction: F_TradersGuild, WorldState: []Cond{
		Equal(L_Longen.IsAlive(), false),
		Equal(L_SMGrande.IsAlive(), false),
		Equal(L_SMRuben.IsAlive(), false),
		Equal(L_SMRen.IsAlive(), false),
	}}
	T_TradersEdge.Overrides = append(T_TradersEdge.Overrides, TradersEdgeProsperous, TradersEdgeOverride, TradersEdgeDestroyed, TradersEdgeMalnourished0, TradersEdgeMalnourished1)
	TradersEdgeDestroyed.Overrides = append(TradersEdgeDestroyed.Overrides, TradersEdgeProsperous, TradersEdgeOverride)
	TradersEdgeMalnourished0.Overrides = append(TradersEdgeMalnourished0.Overrides, TradersEdgeProsperous, TradersEdgeOverride)
	TradersEdgeMalnourished1.Overrides = append(TradersEdgeMalnourished1.Overrides, TradersEdgeProsperous, TradersEdgeOverride)
	return []*Town{TradersEdgeProsperous}
}
