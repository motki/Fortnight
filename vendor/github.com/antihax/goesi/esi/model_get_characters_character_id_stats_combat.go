/*
 * EVE Swagger Interface
 *
 * An OpenAPI for EVE Online
 *
 * OpenAPI spec version: 0.8.2
 *
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package esi

/* A list of GetCharactersCharacterIdStatsCombat. */
//easyjson:json
type GetCharactersCharacterIdStatsCombatList []GetCharactersCharacterIdStatsCombat

/* combat object */
//easyjson:json
type GetCharactersCharacterIdStatsCombat struct {
	CapDrainedbyNpc                        int64 `json:"cap_drainedby_npc,omitempty"`                            /* cap_drainedby_npc integer */
	CapDrainedbyPc                         int64 `json:"cap_drainedby_pc,omitempty"`                             /* cap_drainedby_pc integer */
	CapDrainingPc                          int64 `json:"cap_draining_pc,omitempty"`                              /* cap_draining_pc integer */
	CriminalFlagSet                        int64 `json:"criminal_flag_set,omitempty"`                            /* criminal_flag_set integer */
	DamageFromNpCsAmount                   int64 `json:"damage_from_np_cs_amount,omitempty"`                     /* damage_from_np_cs_amount integer */
	DamageFromNpCsNumShots                 int64 `json:"damage_from_np_cs_num_shots,omitempty"`                  /* damage_from_np_cs_num_shots integer */
	DamageFromPlayersBombAmount            int64 `json:"damage_from_players_bomb_amount,omitempty"`              /* damage_from_players_bomb_amount integer */
	DamageFromPlayersBombNumShots          int64 `json:"damage_from_players_bomb_num_shots,omitempty"`           /* damage_from_players_bomb_num_shots integer */
	DamageFromPlayersCombatDroneAmount     int64 `json:"damage_from_players_combat_drone_amount,omitempty"`      /* damage_from_players_combat_drone_amount integer */
	DamageFromPlayersCombatDroneNumShots   int64 `json:"damage_from_players_combat_drone_num_shots,omitempty"`   /* damage_from_players_combat_drone_num_shots integer */
	DamageFromPlayersEnergyAmount          int64 `json:"damage_from_players_energy_amount,omitempty"`            /* damage_from_players_energy_amount integer */
	DamageFromPlayersEnergyNumShots        int64 `json:"damage_from_players_energy_num_shots,omitempty"`         /* damage_from_players_energy_num_shots integer */
	DamageFromPlayersFighterBomberAmount   int64 `json:"damage_from_players_fighter_bomber_amount,omitempty"`    /* damage_from_players_fighter_bomber_amount integer */
	DamageFromPlayersFighterBomberNumShots int64 `json:"damage_from_players_fighter_bomber_num_shots,omitempty"` /* damage_from_players_fighter_bomber_num_shots integer */
	DamageFromPlayersFighterDroneAmount    int64 `json:"damage_from_players_fighter_drone_amount,omitempty"`     /* damage_from_players_fighter_drone_amount integer */
	DamageFromPlayersFighterDroneNumShots  int64 `json:"damage_from_players_fighter_drone_num_shots,omitempty"`  /* damage_from_players_fighter_drone_num_shots integer */
	DamageFromPlayersHybridAmount          int64 `json:"damage_from_players_hybrid_amount,omitempty"`            /* damage_from_players_hybrid_amount integer */
	DamageFromPlayersHybridNumShots        int64 `json:"damage_from_players_hybrid_num_shots,omitempty"`         /* damage_from_players_hybrid_num_shots integer */
	DamageFromPlayersMissileAmount         int64 `json:"damage_from_players_missile_amount,omitempty"`           /* damage_from_players_missile_amount integer */
	DamageFromPlayersMissileNumShots       int64 `json:"damage_from_players_missile_num_shots,omitempty"`        /* damage_from_players_missile_num_shots integer */
	DamageFromPlayersProjectileAmount      int64 `json:"damage_from_players_projectile_amount,omitempty"`        /* damage_from_players_projectile_amount integer */
	DamageFromPlayersProjectileNumShots    int64 `json:"damage_from_players_projectile_num_shots,omitempty"`     /* damage_from_players_projectile_num_shots integer */
	DamageFromPlayersSmartBombAmount       int64 `json:"damage_from_players_smart_bomb_amount,omitempty"`        /* damage_from_players_smart_bomb_amount integer */
	DamageFromPlayersSmartBombNumShots     int64 `json:"damage_from_players_smart_bomb_num_shots,omitempty"`     /* damage_from_players_smart_bomb_num_shots integer */
	DamageFromPlayersSuperAmount           int64 `json:"damage_from_players_super_amount,omitempty"`             /* damage_from_players_super_amount integer */
	DamageFromPlayersSuperNumShots         int64 `json:"damage_from_players_super_num_shots,omitempty"`          /* damage_from_players_super_num_shots integer */
	DamageFromStructuresTotalAmount        int64 `json:"damage_from_structures_total_amount,omitempty"`          /* damage_from_structures_total_amount integer */
	DamageFromStructuresTotalNumShots      int64 `json:"damage_from_structures_total_num_shots,omitempty"`       /* damage_from_structures_total_num_shots integer */
	DamageToPlayersBombAmount              int64 `json:"damage_to_players_bomb_amount,omitempty"`                /* damage_to_players_bomb_amount integer */
	DamageToPlayersBombNumShots            int64 `json:"damage_to_players_bomb_num_shots,omitempty"`             /* damage_to_players_bomb_num_shots integer */
	DamageToPlayersCombatDroneAmount       int64 `json:"damage_to_players_combat_drone_amount,omitempty"`        /* damage_to_players_combat_drone_amount integer */
	DamageToPlayersCombatDroneNumShots     int64 `json:"damage_to_players_combat_drone_num_shots,omitempty"`     /* damage_to_players_combat_drone_num_shots integer */
	DamageToPlayersEnergyAmount            int64 `json:"damage_to_players_energy_amount,omitempty"`              /* damage_to_players_energy_amount integer */
	DamageToPlayersEnergyNumShots          int64 `json:"damage_to_players_energy_num_shots,omitempty"`           /* damage_to_players_energy_num_shots integer */
	DamageToPlayersFighterBomberAmount     int64 `json:"damage_to_players_fighter_bomber_amount,omitempty"`      /* damage_to_players_fighter_bomber_amount integer */
	DamageToPlayersFighterBomberNumShots   int64 `json:"damage_to_players_fighter_bomber_num_shots,omitempty"`   /* damage_to_players_fighter_bomber_num_shots integer */
	DamageToPlayersFighterDroneAmount      int64 `json:"damage_to_players_fighter_drone_amount,omitempty"`       /* damage_to_players_fighter_drone_amount integer */
	DamageToPlayersFighterDroneNumShots    int64 `json:"damage_to_players_fighter_drone_num_shots,omitempty"`    /* damage_to_players_fighter_drone_num_shots integer */
	DamageToPlayersHybridAmount            int64 `json:"damage_to_players_hybrid_amount,omitempty"`              /* damage_to_players_hybrid_amount integer */
	DamageToPlayersHybridNumShots          int64 `json:"damage_to_players_hybrid_num_shots,omitempty"`           /* damage_to_players_hybrid_num_shots integer */
	DamageToPlayersMissileAmount           int64 `json:"damage_to_players_missile_amount,omitempty"`             /* damage_to_players_missile_amount integer */
	DamageToPlayersMissileNumShots         int64 `json:"damage_to_players_missile_num_shots,omitempty"`          /* damage_to_players_missile_num_shots integer */
	DamageToPlayersProjectileAmount        int64 `json:"damage_to_players_projectile_amount,omitempty"`          /* damage_to_players_projectile_amount integer */
	DamageToPlayersProjectileNumShots      int64 `json:"damage_to_players_projectile_num_shots,omitempty"`       /* damage_to_players_projectile_num_shots integer */
	DamageToPlayersSmartBombAmount         int64 `json:"damage_to_players_smart_bomb_amount,omitempty"`          /* damage_to_players_smart_bomb_amount integer */
	DamageToPlayersSmartBombNumShots       int64 `json:"damage_to_players_smart_bomb_num_shots,omitempty"`       /* damage_to_players_smart_bomb_num_shots integer */
	DamageToPlayersSuperAmount             int64 `json:"damage_to_players_super_amount,omitempty"`               /* damage_to_players_super_amount integer */
	DamageToPlayersSuperNumShots           int64 `json:"damage_to_players_super_num_shots,omitempty"`            /* damage_to_players_super_num_shots integer */
	DamageToStructuresTotalAmount          int64 `json:"damage_to_structures_total_amount,omitempty"`            /* damage_to_structures_total_amount integer */
	DamageToStructuresTotalNumShots        int64 `json:"damage_to_structures_total_num_shots,omitempty"`         /* damage_to_structures_total_num_shots integer */
	DeathsHighSec                          int64 `json:"deaths_high_sec,omitempty"`                              /* deaths_high_sec integer */
	DeathsLowSec                           int64 `json:"deaths_low_sec,omitempty"`                               /* deaths_low_sec integer */
	DeathsNullSec                          int64 `json:"deaths_null_sec,omitempty"`                              /* deaths_null_sec integer */
	DeathsPodHighSec                       int64 `json:"deaths_pod_high_sec,omitempty"`                          /* deaths_pod_high_sec integer */
	DeathsPodLowSec                        int64 `json:"deaths_pod_low_sec,omitempty"`                           /* deaths_pod_low_sec integer */
	DeathsPodNullSec                       int64 `json:"deaths_pod_null_sec,omitempty"`                          /* deaths_pod_null_sec integer */
	DeathsPodWormhole                      int64 `json:"deaths_pod_wormhole,omitempty"`                          /* deaths_pod_wormhole integer */
	DeathsWormhole                         int64 `json:"deaths_wormhole,omitempty"`                              /* deaths_wormhole integer */
	DroneEngage                            int64 `json:"drone_engage,omitempty"`                                 /* drone_engage integer */
	Dscans                                 int64 `json:"dscans,omitempty"`                                       /* dscans integer */
	DuelRequested                          int64 `json:"duel_requested,omitempty"`                               /* duel_requested integer */
	EngagementRegister                     int64 `json:"engagement_register,omitempty"`                          /* engagement_register integer */
	KillsAssists                           int64 `json:"kills_assists,omitempty"`                                /* kills_assists integer */
	KillsHighSec                           int64 `json:"kills_high_sec,omitempty"`                               /* kills_high_sec integer */
	KillsLowSec                            int64 `json:"kills_low_sec,omitempty"`                                /* kills_low_sec integer */
	KillsNullSec                           int64 `json:"kills_null_sec,omitempty"`                               /* kills_null_sec integer */
	KillsPodHighSec                        int64 `json:"kills_pod_high_sec,omitempty"`                           /* kills_pod_high_sec integer */
	KillsPodLowSec                         int64 `json:"kills_pod_low_sec,omitempty"`                            /* kills_pod_low_sec integer */
	KillsPodNullSec                        int64 `json:"kills_pod_null_sec,omitempty"`                           /* kills_pod_null_sec integer */
	KillsPodWormhole                       int64 `json:"kills_pod_wormhole,omitempty"`                           /* kills_pod_wormhole integer */
	KillsWormhole                          int64 `json:"kills_wormhole,omitempty"`                               /* kills_wormhole integer */
	NpcFlagSet                             int64 `json:"npc_flag_set,omitempty"`                                 /* npc_flag_set integer */
	ProbeScans                             int64 `json:"probe_scans,omitempty"`                                  /* probe_scans integer */
	PvpFlagSet                             int64 `json:"pvp_flag_set,omitempty"`                                 /* pvp_flag_set integer */
	RepairArmorByRemoteAmount              int64 `json:"repair_armor_by_remote_amount,omitempty"`                /* repair_armor_by_remote_amount integer */
	RepairArmorRemoteAmount                int64 `json:"repair_armor_remote_amount,omitempty"`                   /* repair_armor_remote_amount integer */
	RepairArmorSelfAmount                  int64 `json:"repair_armor_self_amount,omitempty"`                     /* repair_armor_self_amount integer */
	RepairCapacitorByRemoteAmount          int64 `json:"repair_capacitor_by_remote_amount,omitempty"`            /* repair_capacitor_by_remote_amount integer */
	RepairCapacitorRemoteAmount            int64 `json:"repair_capacitor_remote_amount,omitempty"`               /* repair_capacitor_remote_amount integer */
	RepairCapacitorSelfAmount              int64 `json:"repair_capacitor_self_amount,omitempty"`                 /* repair_capacitor_self_amount integer */
	RepairHullByRemoteAmount               int64 `json:"repair_hull_by_remote_amount,omitempty"`                 /* repair_hull_by_remote_amount integer */
	RepairHullRemoteAmount                 int64 `json:"repair_hull_remote_amount,omitempty"`                    /* repair_hull_remote_amount integer */
	RepairHullSelfAmount                   int64 `json:"repair_hull_self_amount,omitempty"`                      /* repair_hull_self_amount integer */
	RepairShieldByRemoteAmount             int64 `json:"repair_shield_by_remote_amount,omitempty"`               /* repair_shield_by_remote_amount integer */
	RepairShieldRemoteAmount               int64 `json:"repair_shield_remote_amount,omitempty"`                  /* repair_shield_remote_amount integer */
	RepairShieldSelfAmount                 int64 `json:"repair_shield_self_amount,omitempty"`                    /* repair_shield_self_amount integer */
	SelfDestructs                          int64 `json:"self_destructs,omitempty"`                               /* self_destructs integer */
	WarpScramblePc                         int64 `json:"warp_scramble_pc,omitempty"`                             /* warp_scramble_pc integer */
	WarpScrambledbyNpc                     int64 `json:"warp_scrambledby_npc,omitempty"`                         /* warp_scrambledby_npc integer */
	WarpScrambledbyPc                      int64 `json:"warp_scrambledby_pc,omitempty"`                          /* warp_scrambledby_pc integer */
	WeaponFlagSet                          int64 `json:"weapon_flag_set,omitempty"`                              /* weapon_flag_set integer */
	WebifiedbyNpc                          int64 `json:"webifiedby_npc,omitempty"`                               /* webifiedby_npc integer */
	WebifiedbyPc                           int64 `json:"webifiedby_pc,omitempty"`                                /* webifiedby_pc integer */
	WebifyingPc                            int64 `json:"webifying_pc,omitempty"`                                 /* webifying_pc integer */
}
