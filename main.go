package main

import (
	"fmt"

	. "github.com/CodyMcCarty/kenshi-worldstates/internal"
)

/* Todos
[x] reintroduce seedtowns
[x] rename things
[x] test out logs and debugers herlpers
[x] finish leader infos and looking at the wiki
[] figure out order
[] add head leaders to factions
[] add HN and Shek towns
[] add bounties
*/

func PrintAllTownInfo(w *World) {
	fmt.Println("*** Town info ***")
	for _, t := range Towns {
		w.LogTownInfo(t)
	}
	fmt.Println()
}

var TradersGuildLeaders = []*Leader{
	L_Longen, L_LdKana, L_SMRen, L_SMHaga, L_SMRuben, L_SMWada, L_SMGrande, L_SMGrace, L_SMMaster,
}

var UnitedCitiesLeaders = []*Leader{
	L_Tengu, L_LdInaba, L_LdYoshinaga, L_LdNagata, L_LdOhta, L_LdSanda, L_LdTsugi, L_LdShiro, L_LdMerin,
}

func ReleaseAllLeaders(w *World) {
	for _, l := range Leaders {
		if l.Status == Imprisoned {
			if l != L_Yabuta {
				w.Release(l)
			}
		}
	}
}

func KillAllImprisoned(w *World) {
	for _, l := range Leaders {
		if l.Status == Imprisoned {
			if l != L_Yabuta {
				w.Kill(l)
			}
		}
	}
}

func main() {
	w := &World{DesiredTownMap: make(map[*Town]DesiredTown)}
	w.Seed()

	// todo: add UC TG Slaver only bounties i.e. Blue Eyes.
	// todo: I may want to keep heft as is
	TryToGetToProsperousQuickly(w)

	// todo: can I ally with the shek and HN?

	// todo: what can be done with HN and Shek to improve their towns and life?
	// maybe capture seto. check in game what shops are in admag vs the fcs.
	/* og admag : ovr
	armour		armor
	Bar			bar
	Bar			bar
	General		trade
	Thieves		Thieves
	Travel		adventure
	Weapon		weapon
				advMechanical
	HN and Skek spawn more and fight more places.








	*/

	// todo: anything with minor factions Deadcat, Flotsam Ninjas, Mongrol, etc.
	// todo: the swamps
	// todo: world states not tied to towns

	PrintAllTownInfo(w)

	fmt.Println("\ndouble check: what happens if all captured leaders get free? die?")
	if true {
		fmt.Println("Releasing All Leaders")
		ReleaseAllLeaders(w)
	} else {
		fmt.Println("Killing All Imprisoned Leaders")
		KillAllImprisoned(w)
	}
}

// todo what's the earliest I can take longen and tengu?
// TryToGetToProsperousQuickly Goal is Make the world fun, better towns, go after the trade guild.
func TryToGetToProsperousQuickly(w *World) {
	fmt.Println("Yoshi and Ohta can disappear.")
	w.Capture(L_LdYoshinaga)
	w.Capture(L_LdOhta)
	w.Capture(L_Longen) // turn Longen into Tinfist. Make alliance with Boss Simion for Tinfist.
	// PortNorth	Nagata & Kana		Simion
	// ShoBattai	Nagata				Simion
	// DistantHiveVillage				falls if !tinfist && Kana
	w.Capture(L_LdNagata)
	w.Capture(L_LdKana)
	fmt.Println("kill Kana to prevent distant hive village override. Alternatively, this would be a good check to ensure releasing Tinfist works.")
	w.Kill(L_LdKana)

	fmt.Println("The order is important. Tinfist Imprisoned, Tengu Imprisoned. Then Tinfist Release and latter death. verify BossSimion at Heng")
	// TraderEdge & Heng Simion Can move in.
	fmt.Println("Capture Tinfist so Boss Simion moves into Heng. Otherwise Boss Simion disappears when Tengu is captured.")
	w.Capture(L_Tinfist)
	w.Capture(L_Tengu) // Simion disappears if tinfist is alive.

	fmt.Println("Release Tinfist so Free Settlement -> Free City")
	w.Release(L_Tinfist)
	fmt.Println("Verify Boss Simion at Heng.")

	// Drin 		Inaba				HN
	// Stoat		Inaba				TechHunter
	w.Capture(L_LdInaba)

	// Bark 		sanda
	// PortSouth 	Wada
	w.Capture(L_SMWada)
	w.Capture(L_LdSanda)
	// ClownSteady 	Master & Grace
	// DriftersLast	Merin & Grace
	// ManHBase		Master & Grace
	// SlaveFarmS	Master & Grace
	w.Capture(L_SMMaster)
	w.Capture(L_SMGrace)
	w.Capture(L_LdMerin)

	{
		fmt.Println("I initally set out to get all the Trader's guild")
		fmt.Println("What's the status of the United Cities")
		for _, l := range UnitedCitiesLeaders {
			l.LogInfo()
		}
		fmt.Println(`
		if capture Tsugi. Brink -> Brink Takeover (Reavers). This is a possibility
		if capture Shiro. Catun -> Catun Fishman Takeover
		`)
		fmt.Println("\nWhat's the status of the Trader's Guild")
		for _, l := range TradersGuildLeaders {
			l.LogInfo()
		}
		fmt.Println(`
			if capture Ren Slave Farm -> Slave Farm Destroyed
			if capture Haga Stone Camp -> Stone Camp Destroyed
			if capture Grande Eyesocket -> Eyesocket Destroyed
			if capture Ruben South Stone Camp -> South Stone Camp Destroyed
		`)
	}

	fmt.Println("Kill Longen to prevent Settled Nomads -> Settled Nomads Empty")
	w.Kill(L_Longen)

	fmt.Println("kill tinfist at the end. He's going a little nuts. The last time skeles tried to rule people, it didn't turn out well. I want his cpu.  Can I do it at the start for cult village? No, if cult village turns then settled Nomads turns.")
	w.Kill(L_Tinfist)

	fmt.Println("Double check that Boss Simion didn't disappear.")
}

func fromOldApp(w *World) {
	fmt.Println("Goal is to have better town overrides and avoid destroyed towns while going after the trader's guild.")
	fmt.Println("Start in the north with a bit of prep work.")
	// Capture Preacher?
	w.Capture(L_LdYoshinaga)
	w.Capture(L_Longen)
	fmt.Println("Maintain positive/non-hostile relations with UC, TG, ST, and all factions. Relations can be increase in a number of ways, such as, repeatable events like donations (see the wiki), increase relationship with nearly every faction by rapidly throwing them into slavery and then buying them out of it, healing them, turning in bounties.")
	fmt.Println("Talk to Tinfist while carrying Longen. Turn Longen into Grey. Ask Boss Simion to join Anti-Slavers. That order gives the best Rep in my testing. -25 with ST&UC !TG; +100AS +75Rebel.")
	// todo kill longen to prevent empty settled nomads. Unlocking his cage is enough for AS to kill him?
	// todo what's the earliest I can taken longen and tengu?

	fmt.Println("South UC & TG")
	w.Capture(L_SMMaster)
	w.Capture(L_SMGrace)
	w.Capture(L_LdMerin)

	fmt.Println("North UC and TG")
	w.Capture(L_LdTsugi)
	// release Tsugi. if she dies in prison, Reavers take Brink. Is that true? I don't think that's how that works. She could die outside of prison.
	w.Capture(L_LdInaba)
	w.Capture(L_LdNagata)
	w.Capture(L_LdSanda)
	w.Capture(L_LdKana)
	// kill kana. // Dead to prevent Distant HiveSlaver. N_PortN Order Maters for Bark. Has to come after Sanda or Tengu. Killing Kana to prevent Distant Hive village override if she escapes prison
	w.Capture(L_SMWada)

	fmt.Println("decide on the fate of heng, Tinfist, and Boss Simion")
	// todo: boss simion
	w.Capture(L_Tinfist)
	// todo lord ohta
	w.Capture(L_Tengu)
	fmt.Println("turn tengu into the", F_AntiSlavers.Name)
	fmt.Println("Deal with Tinfist. The order is important. Tinfist Imprisoned, Tengu Imprisoned, verify BossSimion at Heng, then Tinfist Dead.")
	// todo what happens if I releaes Tinfist in stead of kill him?
	// todo after my inital look at the overrides, I think I will leave tinfist be.
	fmt.Println("double check that Reavers didn't spawn near", L_LdTsugi.Home, "due to", L_LdTsugi.Name)
	fmt.Println("Send UC and TG Leaders to Rebirth for Redemption")
}
