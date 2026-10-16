//go:build !noembed

package datalib

import _ "embed"

const embedded = true

// The data-library tables, embedded (filediver's decrypted copies). Build with -tags noembed to leave them out
// and supply them at run time with LoadTables.

//go:embed generated_arc_settings.dl_bin.gz
var arcSettingsCompressed []byte

//go:embed generated_animation_event_trigger_settings.dl_bin.gz
var animationEventTriggerSettingsCompressed []byte

//go:embed generated_beam_settings.dl_bin.gz
var beamSettingsCompressed []byte

//go:embed generated_damage_settings.dl_bin.gz
var damageSettingsCompressed []byte

//go:embed generated_planet_data.dl_bin.gz
var planetDataCompressed []byte

//go:embed generated_planet_override_settings.dl_bin.gz
var planetOverrideSettingsCompressed []byte

//go:embed generated_planet_region_settings.dl_bin.gz
var planetRegionSettingsCompressed []byte

//go:embed generated_planet_types_settings.dl_bin.gz
var planetTypesSettingsCompressed []byte

//go:embed generated_region_settings.dl_bin.gz
var regionSettingsCompressed []byte

//go:embed generated_sky_settings.dl_bin.gz
var skySettingsCompressed []byte

//go:embed generated_projectile_settings.dl_bin.gz
var projectileSettingsCompressed []byte

//go:embed generated_environment_settings.dl_bin.gz
var environmentSettingsCompressed []byte

//go:embed generated_explosion_settings.dl_bin.gz
var explosionSettingsCompressed []byte

//go:embed generated_customization_armor_sets.dl_bin.gz
var customizationArmorSetsCompressed []byte

//go:embed generated_customization_passive_bonuses.dl_bin.gz
var customizationPassiveBonusesCompressed []byte

//go:embed generated_unit_customization_settings.dl_bin.gz
var unitCustomizationSettingsCompressed []byte

//go:embed generated_weapon_customization_settings.dl_bin.gz
var weaponCustomizationSettingsCompressed []byte

//go:embed generated_entities.dl_bin.gz
var entitiesCompressed []byte

//go:embed generated_entity_deltas.dl_bin.gz
var entityDeltasCompressed []byte

//go:embed generated_zone_data.dl_bin.gz
var zoneSettingsCompressed []byte

//go:embed dl_library.dl_typelib.gz
var typelibCompressed []byte
