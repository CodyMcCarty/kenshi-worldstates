package main

import (
	"fmt"

	. "github.com/CodyMcCarty/kenshi-worldstates/internal"
)

func main() {
	w := &World{DesiredTownMap: make(map[*Town]DesiredTown)}
	w.Seed()

	fmt.Println("Goal is to have better town overrides and avoid destroyed towns while going after the trader's guild")
	fmt.Println("Start in the north with a bit of prep work.")
	// Capture Preacher?
	w.Capture(L_LdYoshinaga)
	w.Capture(L_Longen)
	fmt.Println("Maintain positive/non-hostile relations with UC, TG, ST, and all factions.")
	fmt.Println("Talk to Tinfist while carrying Longen. Turn Longen into Grey. Ask Boss Simion to join Anti-Slavers. That order gives the best Rep in my testing. -25 with ST&UC !TG; +100AS +75Rebel.")
	// kill longen to prevent empty settled nomads. Unlocking his cage is enough for AS to kill him?

	fmt.Println("South UC & TG")
	w.Capture(L_SMMaster)
	w.Capture(L_SMGrace)
	w.Capture(L_LdMerin)

	fmt.Println("North UC and TG")
	w.Capture(L_LdTsugi)
	// release Tsugi. if she dies in prison, Reavers take Brink
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
	fmt.Println("double check that Reavers didn't spawn near", L_LdTsugi.Home, "due to", L_LdTsugi.Name)
	fmt.Println("Send UC and TG Leaders to Rebirth for Redemption")

	// todo: double check: what happens if all captured leaders get free? die?
}
