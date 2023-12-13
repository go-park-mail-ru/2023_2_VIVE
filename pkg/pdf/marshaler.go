package pdf

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

type FontFamilyName string

const (
	TitleSep                         = " "
	SectionContentSep                = "\n"
	fontDir                          = "./assets/"
	RobotoFont        FontFamilyName = "Roboto"
	PDFTagName                       = "pdf"
	Header1Height                    = 10
)

type PDFFieldType string

const (
	TitleType   PDFFieldType = "title"
	ContentType PDFFieldType = "content"
)

type PDFTagContent struct {
	Type       string
	HeaderName string
	// Level      int
}

func parsePDFTagContent(tagContent string) *PDFTagContent {
	elems := strings.Split(tagContent, ",")

	res := &PDFTagContent{}
	val := reflect.ValueOf(res).Elem()

	for i := 0; i < val.NumField() && i < len(elems); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		switch typeField.Type.Kind() {
		case reflect.Int:
			elemInt, _ := strconv.Atoi(elems[i])
			valueField.SetInt(int64(elemInt))
		case reflect.String:
			valueField.SetString(elems[i])
		}
		// typeField := val.Type().Field(i)

	}

	return res
}

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
	Header1Font     FontConfig
	Header2Font     FontConfig
	Header3Font     FontConfig
	RegularTextFont FontConfig
}

// type PDFSection struct {
// 	header  string
// 	content string
// }

type PDFSection map[string][]string

// TODO: define order of sections in pdf file
type PDFFileStruct struct {
	titleParts []string
	sections   PDFSection
}

func (f *PDFFileStruct) Build(pdf *gofpdf.Fpdf, config *PDFConfig) {
	title := strings.Join(f.titleParts, TitleSep)
	addTitle(pdf, config, title)

	for header, content := range f.sections {
		addSection(pdf, config, header, strings.Join(content, SectionContentSep))
	}
}

func getContentAreaSize(pdf *gofpdf.Fpdf) gofpdf.SizeType {

	res := gofpdf.SizeType{}
	pageWidth, pageHeight := pdf.GetPageSize()
	marginL, marginT, marginR, marginB := pdf.GetMargins()

	res.Wd = pageWidth - marginL - marginR
	res.Ht = pageHeight - marginB - marginT

	return res
}

func initFonts(pdf *gofpdf.Fpdf) {
	pdf.AddUTF8Font(string(RobotoFont), "", "Roboto-Regular.ttf")
	pdf.AddUTF8Font(string(RobotoFont), "B", "Roboto-Bold.ttf")
}

func initPDF(config *PDFConfig) *gofpdf.Fpdf {
	pdf := gofpdf.New(config.Orientation, config.Unit, config.PageFormat, fontDir)
	initFonts(pdf)
	pdf.SetMargins(config.LeftMargin, config.TopMargin, config.RightMargin)
	pdf.SetAutoPageBreak(true, config.BottomMargin)

	pdf.AddPage()

	return pdf
}

func insertTextLines(pdf *gofpdf.Fpdf, currFontConfig *FontConfig, lines []string, blankLine bool) {
	ht := pdf.PointConvert(currFontConfig.Size) * currFontConfig.HeightCoef
	for _, line := range lines {
		pdf.CellFormat(getContentAreaSize(pdf).Wd, ht, line, "", 1, "LT", false, 0, "")
	}
	if blankLine {
		pdf.Ln(ht)
	}
}

func addTitle(pdf *gofpdf.Fpdf, config *PDFConfig, title string) {
	pdf.SetFont(config.Header1Font.Family, config.Header1Font.Style, config.Header1Font.Size)
	contentWidth := getContentAreaSize(pdf).Wd

	lines := pdf.SplitText(title, contentWidth)
	insertTextLines(pdf, &config.Header1Font, lines, false)
}

func addSection(pdf *gofpdf.Fpdf, config *PDFConfig, headerName, content string) {
	pdf.SetFont(config.Header2Font.Family, config.Header2Font.Style, config.Header2Font.Size)

	headerLines := pdf.SplitText(headerName, getContentAreaSize(pdf).Wd)
	insertTextLines(pdf, &config.Header2Font, headerLines, false)

	pdf.SetFont(config.RegularTextFont.Family, config.RegularTextFont.Style, config.RegularTextFont.Size)
	contentLines := pdf.SplitText(content, getContentAreaSize(pdf).Wd)
	insertTextLines(pdf, &config.RegularTextFont, contentLines, true)
}

func parseData( /* pdf *gofpdf.Fpdf,  */ config *PDFConfig, data any) *PDFFileStruct {
	res := PDFFileStruct{
		titleParts: []string{},
		sections: PDFSection{},
	}
	val := reflect.ValueOf(data).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		if valueField.Kind() == reflect.Pointer && valueField.IsNil() {
			continue
		}
		if valueField.Kind() == reflect.Ptr && !valueField.IsNil() {
			valueField = valueField.Elem()
		}
		pdfTag := typeField.Tag.Get(PDFTagName)

		if pdfTag == "" {
			continue
		}

		tagContent := parsePDFTagContent(pdfTag)

		switch tagContent.Type {
		case string(TitleType):
			if strings.TrimSpace(valueField.String()) == "" {
				continue
			}
			res.titleParts = append(res.titleParts, strings.TrimSpace(valueField.String()))
			// addTitle(pdf, config, valueField.String())

		case string(ContentType):
			header := strings.TrimSpace(tagContent.HeaderName)
			content := strings.TrimSpace(valueField.String())
			if _, ok := res.sections[header]; !ok {
				res.sections[header] = []string{}
			}
			res.sections[header] = append(res.sections[header], content)
		}
	}
	return &res
}

func MarshalPDF(config *PDFConfig, data any) (*gofpdf.Fpdf, error) {
	pdf := initPDF(config)

	pdfFileStruct := parseData( /* pdf,  */ config, data)
	pdfFileStruct.Build(pdf, config)

	return pdf, nil
}
