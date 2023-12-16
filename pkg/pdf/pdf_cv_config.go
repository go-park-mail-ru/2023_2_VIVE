package pdf

var (
	Header1Config = FontConfig{
		Family:     string(RobotoFont),
		Style:      "B",
		Size:       25,
		HeightCoef: 1.5,
	}
	Header2Config = FontConfig{
		Family:     string(RobotoFont),
		Style:      "B",
		Size:       18,
		HeightCoef: 1.5,
	}
	Header3Config = FontConfig{
		Family:     string(RobotoFont),
		Style:      "B",
		Size:       12,
		HeightCoef: 1.5,
	}
	Header4Config = FontConfig{
		Family:     string(RobotoFont),
		Style:      "",
		Size:       10,
		HeightCoef: 1.5,
	}
	RegularTextConfig = FontConfig{
		Family:     string(RobotoFont),
		Style:      "",
		Size:       9,
		HeightCoef: 1.5,
	}

	CVConfig = PDFConfig{
		Orientation: "P",
		Unit:        "mm",
		PageFormat:  "A4",
		// FontDir:         "",
		LeftMargin:   20,
		TopMargin:    20,
		RightMargin:  20,
		BottomMargin: 20,
		HeaderFonts: []FontConfig{
			Header1Config,
			Header2Config,
			Header3Config,
			Header4Config,
		},
		// Header1Font:     Header1Config,
		// Header2Font:     Header2Config,
		// Header3Font:     Header3Config,
		RegularTextFont: RegularTextConfig,
		Mappers: map[string]Mapper{
			"gender":     MapGender,
			"dd.mm.yyyy": MapToDDMMYYYY,
			"mm.yyyy":    MapToMMYYYY,
		},
	}
)

type FontConfig struct {
	Family string

	// "B"- bold; "I" - italic; "U" - underscore; "S" - strike-out
	Style      string
	Size       float64
	HeightCoef float64
}

type PDFConfig struct {
	// "P" or "Portrait"; "L" or "Landscape"
	Orientation string
	//"pt" for point, "mm" for millimeter, "cm" for centimeter, or "in" for inch
	Unit string
	// "A3", "A4", "A5", "Letter", "Legal", or "Tabloid". An empty string will be replaced with "A4"
	PageFormat string
	// FontDir         string
	LeftMargin      float64
	TopMargin       float64
	RightMargin     float64
	BottomMargin    float64
	HeaderFonts     []FontConfig
	RegularTextFont FontConfig
	Mappers         map[string]Mapper
}
