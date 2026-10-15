//go:build noembed

package datalib

// Built with -tags noembed: no tables in the binary; LoadTables supplies them.

var (
	arcSettingsCompressed                   []byte
	animationEventTriggerSettingsCompressed []byte
	beamSettingsCompressed                  []byte
	damageSettingsCompressed                []byte
	planetDataCompressed                    []byte
	planetOverrideSettingsCompressed        []byte
	planetRegionSettingsCompressed          []byte
	planetTypesSettingsCompressed           []byte
	regionSettingsCompressed                []byte
	skySettingsCompressed                   []byte
	projectileSettingsCompressed            []byte
	environmentSettingsCompressed           []byte
	explosionSettingsCompressed             []byte
	customizationArmorSetsCompressed        []byte
	customizationPassiveBonusesCompressed   []byte
	unitCustomizationSettingsCompressed     []byte
	weaponCustomizationSettingsCompressed   []byte
	entitiesCompressed                      []byte
	entityDeltasCompressed                  []byte
	zoneSettingsCompressed                  []byte
	typelibCompressed                       []byte
)
