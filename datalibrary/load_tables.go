package datalib

import "fmt"

// LoadTables replaces the data-library tables with the given decrypted files, by the name they have in the game's
// data/game folder ("generated_entities.dl_bin", "dl_library.dl_typelib"). Call it before any Parse / Load function
// (several cache their first result). Names not given keep the embedded table (none with -tags noembed).
func LoadTables(files map[string][]byte) error {
	loaded = len(files) > 0
	for name, b := range files {
		switch name {
		case "generated_arc_settings.dl_bin":
			arcSettings = b
		case "generated_animation_event_trigger_settings.dl_bin":
			animationEventTriggerSettings = b
		case "generated_beam_settings.dl_bin":
			beamSettings = b
		case "generated_damage_settings.dl_bin":
			damageSettings = b
		case "generated_planet_data.dl_bin":
			planetData = b
		case "generated_planet_override_settings.dl_bin":
			planetOverrideSettings = b
		case "generated_planet_region_settings.dl_bin":
			planetRegionSettings = b
		case "generated_planet_types_settings.dl_bin":
			planetTypesSettings = b
		case "generated_region_settings.dl_bin":
			regionSettings = b
		case "generated_sky_settings.dl_bin":
			skySettings = b
		case "generated_projectile_settings.dl_bin":
			projectileSettings = b
		case "generated_environment_settings.dl_bin":
			environmentSettings = b
		case "generated_explosion_settings.dl_bin":
			explosionSettings = b
		case "generated_customization_armor_sets.dl_bin":
			customizationArmorSets = b
			pieceStride = 0 // (found again for this table)
		case "generated_customization_passive_bonuses.dl_bin":
			customizationPassiveBonuses = b
		case "generated_unit_customization_settings.dl_bin":
			unitCustomizationSettings = b
		case "generated_weapon_customization_settings.dl_bin":
			weaponCustomizationSettings = b
		case "generated_entities.dl_bin":
			entities = b
		case "generated_entity_deltas.dl_bin":
			entityDeltas = b
		case "generated_zone_data.dl_bin":
			zoneSettings = b
		case "dl_library.dl_typelib":
			typelib = b
		default:
			return fmt.Errorf("datalib: unknown table %q", name)
		}
	}
	return nil
}

var loaded bool

// HaveTables: the data-library tables are there (embedded, or given to LoadTables). Built with -tags noembed and
// without LoadTables, everything that needs them is unavailable (App opens without armor sets, skins, paint schemes,
// planet data; units export without their variants).
func HaveTables() bool { return embedded || loaded }

// TableNames lists the files LoadTables accepts.
func TableNames() []string {
	return []string{
		"generated_arc_settings.dl_bin",
		"generated_animation_event_trigger_settings.dl_bin",
		"generated_beam_settings.dl_bin",
		"generated_damage_settings.dl_bin",
		"generated_planet_data.dl_bin",
		"generated_planet_override_settings.dl_bin",
		"generated_planet_region_settings.dl_bin",
		"generated_planet_types_settings.dl_bin",
		"generated_region_settings.dl_bin",
		"generated_sky_settings.dl_bin",
		"generated_projectile_settings.dl_bin",
		"generated_environment_settings.dl_bin",
		"generated_explosion_settings.dl_bin",
		"generated_customization_armor_sets.dl_bin",
		"generated_customization_passive_bonuses.dl_bin",
		"generated_unit_customization_settings.dl_bin",
		"generated_weapon_customization_settings.dl_bin",
		"generated_entities.dl_bin",
		"generated_entity_deltas.dl_bin",
		"generated_zone_data.dl_bin",
		"dl_library.dl_typelib",
	}
}
