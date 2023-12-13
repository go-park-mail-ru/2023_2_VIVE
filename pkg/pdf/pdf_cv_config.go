package pdf

var (
	Header1Config = FontConfig{
		Family: "Arial",
		Style:  "B",
		Size:   25,
		Height: 30,
	}
	Header2Config = FontConfig{
		Family: "Arial",
		Style:  "B",
		Size:   18,
		Height: 20,
	}
	Header3Config = FontConfig{
		Family: "Arial",
		Style:  "",
		Size:   12,
		Height: 15,
	}
	RegularTextConfig = FontConfig{
		Family: "Arial",
		Style:  "",
		Size:   9,
		Height: 11,
	}

	CVConfig = PDFConfig{
		Orientation:     "P",
		Unit:            "mm",
		PageFormat:      "A4",
		FontDir:         "",
		LeftMargin:      20,
		TopMargin:       20,
		RightMargin:     20,
		BottomMargin:    20,
		Header1Font:     Header1Config,
		Header2Font:     Header2Config,
		Header3Font:     Header3Config,
		RegularTextFont: RegularTextConfig,
	}
)
