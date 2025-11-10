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

func main() {
	w := &World{DesiredTownMap: make(map[*Town]DesiredTown)}
	w.Seed()

	TryToGetToProsperousQuickly(w)
	PrintAllTownInfo(w)

	// todo: add UC TG Slaver only bounties i.e. Blue Eyes.
	// todo: I may want to keep heft as is
	// todo what's the earliest I can taken longen and tengu?
	// todo: what can be done with HN and Shek to improve their towns and life?
	// todo: I initally set out to get all the Trader's guild
	// todo: double check: what happens if all captured leaders get free? die?
}

// Goal is Make the world fun, better towns, go after the trade guild.
func TryToGetToProsperousQuickly(w *World) {
	w.Capture(L_LdYoshinaga)
	w.Capture(L_Longen)
	//w.Capture(L_Tengu)
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
