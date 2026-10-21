package datalib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/xypwn/filediver/datalibrary/enum"
	"github.com/xypwn/filediver/stingray"
	"github.com/xypwn/filediver/util"
)

type WeatherColorSet struct {
	ID                             uint32     `json:"id"`
	CloudColor                     mgl32.Vec3 `json:"cloud_color"`
	VistaCloudColor                mgl32.Vec3 `json:"vista_cloud_color"`
	FogColor                       mgl32.Vec3 `json:"fog_color"`
	WeatherEffectColor             mgl32.Vec3 `json:"weather_effect_color"`
	WeatherGroundFogVolumeColor    mgl32.Vec3 `json:"weather_ground_fog_volume_color"`
	NearGroundFogColor             mgl32.Vec3 `json:"near_ground_fog_color"`
	ZoneFogTint                    mgl32.Vec3 `json:"zone_fog_tint"`
	WeatherSunColorMult            mgl32.Vec3 `json:"weather_sun_color_mult"`
	WeatherSunIntensityMult        float32    `json:"weather_sun_intensity_mult"`
	ZoneFogMainColor               mgl32.Vec3 `json:"zone_fog_main_color"`
	ZoneFogSecondaryColor          mgl32.Vec3 `json:"zone_fog_secondary_color"`
	ZoneFogBugsColor               mgl32.Vec3 `json:"zone_fog_bugs_color"`
	ZoneFogBotsColor               mgl32.Vec3 `json:"zone_fog_bots_color"`
	ZoneFogIlluminateColor         mgl32.Vec3 `json:"zone_fog_illuminate_color"`
	AmbientDiffuseTint             mgl32.Vec3 `json:"ambient_diffuse_tint"`
	AmbientDiffuseGroundBounceTint mgl32.Vec3 `json:"ambient_diffuse_ground_bounce_tint"`
	FogAmbientTint                 mgl32.Vec3 `json:"fog_ambient_tint"`
	FogAmbientGroundBounceTint     mgl32.Vec3 `json:"fog_ambient_ground_bounce_tint"`
}

// PlanetColorSet: the planet's own zone fog colors, split off WeatherColorSet in game build 25988413 (member names
// by length: the zone_fog_* colors that left WeatherColorSet).
type PlanetColorSet struct {
	ID                     uint32     `json:"id"`
	ZoneFogMainColor       mgl32.Vec3 `json:"zone_fog_main_color"`
	ZoneFogSecondaryColor  mgl32.Vec3 `json:"zone_fog_secondary_color"`
	ZoneFogBugsColor       mgl32.Vec3 `json:"zone_fog_bugs_color"`
	ZoneFogBotsColor       mgl32.Vec3 `json:"zone_fog_bots_color"`
	ZoneFogIlluminateColor mgl32.Vec3 `json:"zone_fog_illuminate_color"`
}

type LevelGenerationPaletteGroup struct {
	//Palette                 stingray.Hash
	AssetGrading            stingray.Hash
	SkySettingsGroup        stingray.Hash
	DayGrading              stingray.Hash
	NightGrading            stingray.Hash
	SunsetGrading           stingray.Hash
	WeatherColorSet         stingray.ThinHash // up to game build 25480438
	WeatherColorSetInternal WeatherColorSet   // up to game build 25480438
	PlanetColorSet          stingray.ThinHash // since game build 25988413
	PlanetColorSetInternal  PlanetColorSet    // since game build 25988413
}

// rawLevelGenerationPaletteGroup: the palette group up to game build 25480438 (256 bytes).
type rawLevelGenerationPaletteGroup struct {
	AssetGrading            stingray.Hash
	SkySettingsGroup        stingray.Hash
	DayGrading              stingray.Hash
	NightGrading            stingray.Hash
	SunsetGrading           stingray.Hash
	WeatherColorSet         stingray.ThinHash
	WeatherColorSetInternal WeatherColorSet
}

func (p rawLevelGenerationPaletteGroup) Resolve() LevelGenerationPaletteGroup {
	return LevelGenerationPaletteGroup{
		AssetGrading: p.AssetGrading, SkySettingsGroup: p.SkySettingsGroup, DayGrading: p.DayGrading,
		NightGrading: p.NightGrading, SunsetGrading: p.SunsetGrading,
		WeatherColorSet: p.WeatherColorSet, WeatherColorSetInternal: p.WeatherColorSetInternal,
	}
}

// rawLevelGenerationPaletteGroup25988413: the palette group since game build 25988413 (112 bytes): a planet color
// set in place of the weather color set.
type rawLevelGenerationPaletteGroup25988413 struct {
	AssetGrading           stingray.Hash
	SkySettingsGroup       stingray.Hash
	DayGrading             stingray.Hash
	NightGrading           stingray.Hash
	SunsetGrading          stingray.Hash
	PlanetColorSet         stingray.ThinHash
	PlanetColorSetInternal PlanetColorSet
	_                      [4]uint8
}

func (p rawLevelGenerationPaletteGroup25988413) Resolve() LevelGenerationPaletteGroup {
	return LevelGenerationPaletteGroup{
		AssetGrading: p.AssetGrading, SkySettingsGroup: p.SkySettingsGroup, DayGrading: p.DayGrading,
		NightGrading: p.NightGrading, SunsetGrading: p.SunsetGrading,
		PlanetColorSet: p.PlanetColorSet, PlanetColorSetInternal: p.PlanetColorSetInternal,
	}
}

type rawLevelGenerationRegion struct {
	Name         DLString
	Id           stingray.ThinHash
	Type         enum.LevelGenerationRegionType
	Variants     DLPtr // Seems to be unused
	Weight       float32
	UnknownFloat float32
	UnknownInt   uint32
	_            [4]uint8
}

type LevelGenerationRegion struct {
	Name         *string
	Id           stingray.ThinHash
	Type         enum.LevelGenerationRegionType
	Variants     *GenerationRegionVariantList
	Weight       float32
	UnknownFloat float32
	UnknownInt   uint32
}

type SimpleLevelGenerationRegion struct {
	Name         *string                            `json:"name"`
	Id           string                             `json:"id"`
	Type         enum.LevelGenerationRegionType     `json:"region"`
	Variants     *SimpleGenerationRegionVariantList `json:"variants"` // Seems to be unused
	Weight       float32                            `json:"weight"`
	UnknownFloat float32                            `json:"unknown_float"`
	UnknownInt   uint32                             `json:"unknown_int"`
}

func (r LevelGenerationRegion) ToSimple(lookupHash HashLookup, lookupThinHash ThinHashLookup, lookupStrings StringsLookup) SimpleLevelGenerationRegion {
	var variants *SimpleGenerationRegionVariantList
	if r.Variants != nil {
		result := r.Variants.ToSimple(lookupHash, lookupThinHash, lookupStrings)
		variants = &result
	}

	return SimpleLevelGenerationRegion{
		Name:         r.Name,
		Id:           lookupThinHash(r.Id),
		Type:         r.Type,
		Variants:     variants,
		Weight:       r.Weight,
		UnknownFloat: r.UnknownFloat,
		UnknownInt:   r.UnknownInt,
	}
}

// rawPlanetData: a planet in either layout (rawPlanetData25480438, rawPlanetData25988413), palettes resolved.
type rawPlanetData struct {
	rawPlanetDataHead
	PaletteGroupLowland      LevelGenerationPaletteGroup
	PaletteGroupHighland     LevelGenerationPaletteGroup
	ScenarioSettingsLowland  stingray.Hash
	ScenarioSettingsHighland stingray.Hash
	UnknownEnum              uint8 // since game build 25988413 (an enum of 4 values, name length 20)
	rawPlanetDataTail
}

// rawPlanetData25480438: LevelGenerationPlanetData up to game build 25480438 (904 bytes).
type rawPlanetData25480438 struct {
	rawPlanetDataHead
	PaletteGroupLowland      rawLevelGenerationPaletteGroup
	PaletteGroupHighland     rawLevelGenerationPaletteGroup
	ScenarioSettingsLowland  stingray.Hash
	ScenarioSettingsHighland stingray.Hash
	rawPlanetDataTail
}

// rawPlanetData25988413: LevelGenerationPlanetData since game build 25988413 (624 bytes): smaller palette groups,
// and a one-byte enum after the scenario settings.
type rawPlanetData25988413 struct {
	rawPlanetDataHead
	PaletteGroupLowland      rawLevelGenerationPaletteGroup25988413
	PaletteGroupHighland     rawLevelGenerationPaletteGroup25988413
	ScenarioSettingsLowland  stingray.Hash
	ScenarioSettingsHighland stingray.Hash
	UnknownEnum              uint8
	_                        [7]uint8
	rawPlanetDataTail
}

type rawPlanetDataHead struct {
	InheritsOffset            int64
	PlanetNameLoc             uint32
	PlanetDescriptionLoc      uint32
	PlanetDescriptionShortLoc uint32
	PlanetSystemNameLoc       uint32
	PlanetLayoutId            stingray.ThinHash
	_                         [4]uint8
	MaterialLookupUnit1       stingray.Hash // name len 32
	MaterialLookupUnit2       stingray.Hash // name len 33
	MaterialLookupUnit3       stingray.Hash // name len 28
	OperationTagsOffset       int64
	OperationTagsCount        int64
	ResourceOverridesOffset   int64
	ResourceOverridesCount    int64
	DebugNameOffset           int64
	RegionLowland             rawLevelGenerationRegion
	RegionHighland            rawLevelGenerationRegion
}

type rawPlanetDataTail struct {
	GameplayModifiersOffset          int64
	GameplayModifiersCount           uint64
	PlanetType                       enum.PlanetType
	Unknown                          stingray.ThinHash
	NatureLocationTagOffset          uint64
	NatureLocationTagCount           uint64
	ScatterSettings                  stingray.ThinHash
	_                                [4]uint8
	MissionPlanetUnit                stingray.Hash
	MissionPlanetHologramUnit        stingray.Hash
	MissionPlanetUnitPackage         stingray.Hash
	MissionPlanetHologramUnitPackage stingray.Hash
	SolarSystemSettings              stingray.Hash
	SolarSystemIdSelectionOffset     int64
	SolarSystemIdSelectionCount      uint64
	SampleTypesOffset                int64
	SampleTypesCount                 uint64
	AmbienceSoundIdStart             uint32
	AmbienceSoundIdStop              uint32
	HologramPlanetMaterial           stingray.Hash
	PlanetPreviewImage               stingray.Hash
	ResourceRegionOverridesOffset    int64
	ResourceRegionOverridesCount     uint64
	ShadingEnvironmentEntity         stingray.Hash
	WaterEntity                      stingray.Hash
	HeathazeEntity                   stingray.Hash
	PackagePath                      stingray.Hash
	UnknownFloat                     float32
	_                                [4]uint8
}

type ResourceRegionOverride struct {
	ID         stingray.ThinHash
	RegionFlag enum.RegionFlag
}

type ResourceOverride struct {
	Type        stingray.Hash
	Replace     stingray.Hash
	ReplaceWith stingray.Hash
}

type PlanetData struct {
	Inherits                         string
	PlanetNameLoc                    string
	PlanetDescriptionLoc             string
	PlanetDescriptionShortLoc        string
	PlanetSystemNameLoc              string
	PlanetLayoutId                   stingray.ThinHash
	MaterialLookupUnit1              stingray.Hash
	MaterialLookupUnit2              stingray.Hash
	MaterialLookupUnit3              stingray.Hash
	OperationTags                    []enum.OperationTag
	ResourceOverrides                []ResourceOverride
	DebugName                        string
	RegionLowland                    LevelGenerationRegion
	RegionHighland                   LevelGenerationRegion
	PaletteGroupLowland              LevelGenerationPaletteGroup
	PaletteGroupHighland             LevelGenerationPaletteGroup
	ScenarioSettingsLowland          stingray.Hash
	ScenarioSettingsHighland         stingray.Hash
	UnknownEnum                      uint8 // since game build 25988413
	GameplayModifiers                []stingray.Hash
	PlanetType                       enum.PlanetType
	Unknown                          stingray.ThinHash
	NatureLocationTags               []enum.NatureLocationTag
	ScatterSettings                  stingray.ThinHash
	MissionPlanetUnit                stingray.Hash
	MissionPlanetHologramUnit        stingray.Hash
	MissionPlanetUnitPackage         stingray.Hash
	MissionPlanetHologramUnitPackage stingray.Hash
	SolarSystemSettings              stingray.Hash
	SolarSystemIdSelections          []stingray.Hash
	SampleTypes                      []enum.SampleType
	AmbienceSoundIdStart             uint32
	AmbienceSoundIdStop              uint32
	HologramPlanetMaterial           stingray.Hash
	PlanetPreviewImage               stingray.Hash
	ResourceRegionOverrides          []ResourceRegionOverride
	ShadingEnvironmentEntity         stingray.Hash
	WaterEntity                      stingray.Hash
	HeathazeEntity                   stingray.Hash
	PackagePath                      stingray.Hash
	UnknownFloat                     float32
}

func (a rawPlanetData) Resolve(lookupHash HashLookup, lookupThinHash ThinHashLookup, lookupStrings StringsLookup) PlanetData {
	return PlanetData{
		PlanetNameLoc:                    lookupStrings(a.PlanetNameLoc),
		PlanetDescriptionLoc:             lookupStrings(a.PlanetDescriptionLoc),
		PlanetDescriptionShortLoc:        lookupStrings(a.PlanetDescriptionShortLoc),
		PlanetSystemNameLoc:              lookupStrings(a.PlanetSystemNameLoc),
		PlanetLayoutId:                   a.PlanetLayoutId,
		MaterialLookupUnit1:              a.MaterialLookupUnit1,
		MaterialLookupUnit2:              a.MaterialLookupUnit2,
		MaterialLookupUnit3:              a.MaterialLookupUnit3,
		PaletteGroupLowland:              a.PaletteGroupLowland,
		PaletteGroupHighland:             a.PaletteGroupHighland,
		ScenarioSettingsLowland:          a.ScenarioSettingsLowland,
		ScenarioSettingsHighland:         a.ScenarioSettingsHighland,
		UnknownEnum:                      a.UnknownEnum,
		PlanetType:                       a.PlanetType,
		Unknown:                          a.Unknown,
		ScatterSettings:                  a.ScatterSettings,
		MissionPlanetUnit:                a.MissionPlanetUnit,
		MissionPlanetHologramUnit:        a.MissionPlanetHologramUnit,
		MissionPlanetUnitPackage:         a.MissionPlanetUnitPackage,
		MissionPlanetHologramUnitPackage: a.MissionPlanetHologramUnitPackage,
		SolarSystemSettings:              a.SolarSystemSettings,
		AmbienceSoundIdStart:             a.AmbienceSoundIdStart,
		AmbienceSoundIdStop:              a.AmbienceSoundIdStop,
		HologramPlanetMaterial:           a.HologramPlanetMaterial,
		PlanetPreviewImage:               a.PlanetPreviewImage,
		ShadingEnvironmentEntity:         a.ShadingEnvironmentEntity,
		WaterEntity:                      a.WaterEntity,
		HeathazeEntity:                   a.HeathazeEntity,
		PackagePath:                      a.PackagePath,
		UnknownFloat:                     a.UnknownFloat,
	}
}

func LoadPlanetData(lookupHash HashLookup, lookupThinHash ThinHashLookup, lookupStrings StringsLookup) ([]PlanetData, error) {
	r := bytes.NewReader(planetData)

	infos := make([]PlanetData, 0)
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, fmt.Errorf("reading count: %v", err)
	}
	for i := uint32(0); i < count; i++ {
		var header DLSubdataHeader
		if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
			return nil, fmt.Errorf("reading item %v: %v", i, err)
		}

		if header.Type != Sum("LevelGenerationPlanetData") {
			return nil, fmt.Errorf("invalid planet data file")
		}

		base, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("finding planet data base: %v", err)
		}

		rawSetting, err := readRawPlanetData(r)
		if err != nil {
			return nil, fmt.Errorf("reading planet data: %v", err)
		}

		setting := rawSetting.Resolve(lookupHash, lookupThinHash, lookupStrings)
		if _, err := r.Seek(int64(base+rawSetting.InheritsOffset), io.SeekStart); err != nil {
			return nil, err
		}
		inherits, err := util.ReadCString(r)
		if err != nil {
			return nil, err
		}
		setting.Inherits = inherits

		operationTags := make([]enum.OperationTag, rawSetting.OperationTagsCount)
		if rawSetting.OperationTagsOffset > 0 {
			if _, err := r.Seek(int64(base+rawSetting.OperationTagsOffset), io.SeekStart); err != nil {
				return nil, err
			}
			if err := binary.Read(r, binary.LittleEndian, &operationTags); err != nil {
				return nil, fmt.Errorf("reading operation tags: %v", err)
			}
		}

		resourceOverrides := make([]ResourceOverride, rawSetting.ResourceOverridesCount)
		if rawSetting.ResourceOverridesOffset > 0 {
			if _, err := r.Seek(int64(base+rawSetting.ResourceOverridesOffset), io.SeekStart); err != nil {
				return nil, err
			}
			if err := binary.Read(r, binary.LittleEndian, &resourceOverrides); err != nil {
				return nil, fmt.Errorf("reading resource overrides: %v", err)
			}
		}
		setting.ResourceOverrides = resourceOverrides

		if _, err := r.Seek(int64(base+rawSetting.DebugNameOffset), io.SeekStart); err != nil {
			return nil, err
		}
		debugName, err := util.ReadCString(r)
		if err != nil {
			return nil, err
		}
		setting.DebugName = debugName

		name, err := rawSetting.RegionLowland.Name.Resolve(r, base)
		if err != nil {
			return nil, err
		}

		rawVariants, err := ResolveDLPtr[rawGenerationRegionVariantList](rawSetting.RegionLowland.Variants, r, base)
		if err != nil {
			return nil, err
		}
		var variants *GenerationRegionVariantList
		if rawVariants != nil {
			variants, err = rawVariants.Deserialize(r, base)
			if err != nil {
				return nil, err
			}
		}

		setting.RegionLowland = LevelGenerationRegion{
			Name:         name,
			Id:           rawSetting.RegionLowland.Id,
			Type:         rawSetting.RegionLowland.Type,
			Variants:     variants,
			Weight:       rawSetting.RegionLowland.Weight,
			UnknownFloat: rawSetting.RegionLowland.UnknownFloat,
			UnknownInt:   rawSetting.RegionLowland.UnknownInt,
		}

		name, err = rawSetting.RegionHighland.Name.Resolve(r, base)
		if err != nil {
			return nil, err
		}

		rawVariants, err = ResolveDLPtr[rawGenerationRegionVariantList](rawSetting.RegionHighland.Variants, r, base)
		if err != nil {
			return nil, err
		}
		variants = nil
		if rawVariants != nil {
			variants, err = rawVariants.Deserialize(r, base)
			if err != nil {
				return nil, err
			}
		}
		setting.RegionHighland = LevelGenerationRegion{
			Name:         name,
			Id:           rawSetting.RegionHighland.Id,
			Type:         rawSetting.RegionHighland.Type,
			Variants:     variants,
			Weight:       rawSetting.RegionHighland.Weight,
			UnknownFloat: rawSetting.RegionHighland.UnknownFloat,
			UnknownInt:   rawSetting.RegionHighland.UnknownInt,
		}

		if _, err := r.Seek(int64(base+rawSetting.GameplayModifiersOffset), io.SeekStart); err != nil {
			return nil, err
		}
		gameplayModifiers := make([]stingray.Hash, rawSetting.GameplayModifiersCount)
		if err := binary.Read(r, binary.LittleEndian, gameplayModifiers); err != nil {
			return nil, err
		}

		if _, err := r.Seek(int64(base+int64(rawSetting.NatureLocationTagOffset)), io.SeekStart); err != nil {
			return nil, err
		}
		natureLocationTags := make([]enum.NatureLocationTag, rawSetting.NatureLocationTagCount)
		if err := binary.Read(r, binary.LittleEndian, natureLocationTags); err != nil {
			return nil, err
		}

		if _, err := r.Seek(int64(base+rawSetting.SolarSystemIdSelectionOffset), io.SeekStart); err != nil {
			return nil, err
		}
		solarSystemIdSelections := make([]stingray.Hash, rawSetting.SolarSystemIdSelectionCount)
		if err := binary.Read(r, binary.LittleEndian, solarSystemIdSelections); err != nil {
			return nil, err
		}

		if _, err := r.Seek(int64(base+rawSetting.SampleTypesOffset), io.SeekStart); err != nil {
			return nil, err
		}
		sampleTypes := make([]enum.SampleType, rawSetting.SampleTypesCount)
		if err := binary.Read(r, binary.LittleEndian, sampleTypes); err != nil {
			return nil, err
		}

		if _, err := r.Seek(int64(base+rawSetting.ResourceRegionOverridesOffset), io.SeekStart); err != nil {
			return nil, err
		}
		resourceRegionOverrides := make([]ResourceRegionOverride, rawSetting.ResourceRegionOverridesCount)
		if err := binary.Read(r, binary.LittleEndian, resourceRegionOverrides); err != nil {
			return nil, err
		}
		setting.GameplayModifiers = gameplayModifiers
		setting.NatureLocationTags = natureLocationTags
		setting.SolarSystemIdSelections = solarSystemIdSelections
		setting.SampleTypes = sampleTypes
		setting.ResourceRegionOverrides = resourceRegionOverrides
		infos = append(infos, setting)
		if _, err := r.Seek(base+int64(header.Size), io.SeekStart); err != nil {
			return nil, err
		}
	}
	return infos, nil
}

// readRawPlanetData reads a planet in the layout of its game build. The first string (inherits) follows the struct
// directly in every planet of both builds, so its offset tells the struct's size: 904 bytes up to 25480438, 624 since
// 25988413.
func readRawPlanetData(r *bytes.Reader) (rawPlanetData, error) {
	base, _ := r.Seek(0, io.SeekCurrent)
	var inherits int64
	if err := binary.Read(r, binary.LittleEndian, &inherits); err != nil {
		return rawPlanetData{}, err
	}
	r.Seek(base, io.SeekStart)
	if inherits >= int64(binary.Size(rawPlanetData25480438{})) {
		var p rawPlanetData25480438
		if err := binary.Read(r, binary.LittleEndian, &p); err != nil {
			return rawPlanetData{}, err
		}
		return rawPlanetData{rawPlanetDataHead: p.rawPlanetDataHead, PaletteGroupLowland: p.PaletteGroupLowland.Resolve(),
			PaletteGroupHighland: p.PaletteGroupHighland.Resolve(), ScenarioSettingsLowland: p.ScenarioSettingsLowland,
			ScenarioSettingsHighland: p.ScenarioSettingsHighland, rawPlanetDataTail: p.rawPlanetDataTail}, nil
	}
	var p rawPlanetData25988413
	if err := binary.Read(r, binary.LittleEndian, &p); err != nil {
		return rawPlanetData{}, err
	}
	return rawPlanetData{rawPlanetDataHead: p.rawPlanetDataHead, PaletteGroupLowland: p.PaletteGroupLowland.Resolve(),
		PaletteGroupHighland: p.PaletteGroupHighland.Resolve(), ScenarioSettingsLowland: p.ScenarioSettingsLowland,
		ScenarioSettingsHighland: p.ScenarioSettingsHighland, UnknownEnum: p.UnknownEnum,
		rawPlanetDataTail: p.rawPlanetDataTail}, nil
}
