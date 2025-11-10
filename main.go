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

	// todo: what can be done with HN and Shek to improve their towns and life?

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
	// Yoshi and Ohta can disappear.
	w.Capture(L_LdYoshinaga)
	w.Capture(L_LdOhta)
	w.Capture(L_Longen) // turn Longen into Tinfist. Make alliance with Boss Simion for Tinfist.
	// PortNorth	Nagata & Kana		Simion
	// ShoBattai	Nagata				Simion
	// DistantHiveVillage				falls if !tinfist && Kana
	w.Capture(L_LdNagata)
	w.Capture(L_LdKana)
	w.Kill(L_LdKana) // to prevent distant hive village issues. Alternatively, this would be a good check to ensure releasing Tinfist works.

	// TraderEdge & Heng Simion Can move in.
	w.Capture(L_Tinfist)
	w.Capture(L_Tengu) // Simion disappears if tinfist is alive.

	// FreeCity	Free Tinfist
	w.Release(L_Tinfist)

	// Drin 		Inaba				HN
	// Stoat		Inaba				TechHunter
	w.Capture(L_LdInaba)

	//ldrs := []*Leader{L_LdSanda, L_SMGrace, L_SMMaster, L_LdMerin, L_SMWada}
	//for _, l := range ldrs {
	//	l.LogInfo()
	//}

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

	// todo: I initally set out to get all the Trader's guild
	fmt.Println("What's the status of the United Cities")
	for _, l := range UnitedCitiesLeaders {
		l.LogInfo()
	}
	fmt.Println("What's the status of the Trader's Guild")
	for _, l := range TradersGuildLeaders {
		l.LogInfo()
	}
	// UC
	//L_LdTsugi, L_LdShiro,, L_SMRen, L_SMHaga, L_SMRuben, L_SMGrande
	//w.Capture(L_LdTsugi) // Brink 	->	 Brink Takeover (Reavers)
	//w.Capture(L_LdShiro)  // Catun 	->	 Catun Fishman Takeover
	// Trader's Guild
	//w.Capture(L_SMRen)    // Slave Farm 	->	 Slave Farm Destroyed
	//w.Capture(L_SMHaga)   // Stone Camp 	->	 Stone Camp Destroyed
	//w.Capture(L_SMGrande) // Eyesocket 	->	 Eyesocket Destroyed
	//w.Capture(L_SMRuben)  // South Stone Camp 	->	 South Stone Camp Destroyed

	w.Kill(L_Longen) // Settled Nomads -> Settled Nomads Empty

	fmt.Println("kill tinfist at the end. I want his cpu.  Can I do it at the start for cult village?")
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
