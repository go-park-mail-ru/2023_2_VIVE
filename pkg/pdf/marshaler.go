package pdf

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

type FontFamilyName string

const (
	HeaderSep                        = " "
	SectionContentSep                = "\n"
	SkipTag                          = "-"
	fontDir                          = "./assets/"
	RobotoFont        FontFamilyName = "Roboto"
	PDFTagName                       = "pdf"
	// Header1Height                    = 10
)

type PDFFieldType string

const (
	HeaderType  PDFFieldType = "header"
	ContentType PDFFieldType = "content"
	// SectionType PDFFieldType = "section"
)

type PDFTagContent struct {
	Type       string
	HeaderName string
	Prefix     string
	MapperType string
}

func parsePDFTagContent(tagContent string) *PDFTagContent {
	elems := strings.Split(tagContent, ",")

	res := &PDFTagContent{}
	val := reflect.ValueOf(res).Elem()

	for i := 0; i < val.NumField() && i < len(elems); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		switch typeField.Type.Kind() {
		case reflect.String:
			valueField.SetString(elems[i])
		}
	}

	return res
}

type PDFFileStruct struct {
	headerParts []string
	content     string
	sections    []*PDFFileStruct
}

func (f *PDFFileStruct) Build(pdf *gofpdf.Fpdf, config *PDFConfig) {
	f.build(pdf, config, 0)
}

func (f *PDFFileStruct) build(pdf *gofpdf.Fpdf, config *PDFConfig, level int) {
	if level >= len(config.HeaderFonts) {
		return
	}

	f.addTitle(pdf, config, level)
	f.addContent(pdf, config)

	for _, section := range f.sections {
		section.build(pdf, config, level+1)
		// addSection(pdf, config, header, strings.Join(section, SectionContentSep))
	}
}

func (f *PDFFileStruct) addTitle(pdf *gofpdf.Fpdf, config *PDFConfig /* title string, */, level int) {
	title := strings.Join(f.headerParts, HeaderSep)
	headerFont := config.HeaderFonts[level]
	pdf.SetFont(headerFont.Family, headerFont.Style, headerFont.Size)
	contentWidth := getContentAreaSize(pdf).Wd

	lines := pdf.SplitText(title, contentWidth)
	insertTextLines(pdf, &headerFont, lines, false)
}

func (f *PDFFileStruct) addContent(pdf *gofpdf.Fpdf, config *PDFConfig) {
	if content := strings.TrimSpace(f.content); content == "" {
		return
	}
	pdf.SetFont(config.RegularTextFont.Family, config.RegularTextFont.Style, config.RegularTextFont.Size)

	lines := pdf.SplitText(strings.TrimSpace(f.content), getContentAreaSize(pdf).Wd)
	insertTextLines(pdf, &config.RegularTextFont, lines, true)
}

func (f *PDFFileStruct) appendSection(innerSection *PDFFileStruct) {
	if innerSection == nil || (strings.TrimSpace(innerSection.content) == "" && len(innerSection.sections) == 0) {
		return
	}
	innerSectionHeader := strings.Join(innerSection.headerParts, HeaderSep)
	for _, section := range f.sections {
		sectionHeader := strings.Join(section.headerParts, HeaderSep)
		if sectionHeader == innerSectionHeader {
			section.content += SectionContentSep + innerSection.content
			section.sections = append(section.sections, innerSection.sections...)
			return
		}
	}
	f.sections = append(f.sections, innerSection)
}

func getContentAreaSize(pdf *gofpdf.Fpdf) gofpdf.SizeType {
	// res := gofpdf.SizeType{}
	pageWidth, pageHeight := pdf.GetPageSize()
	marginL, marginT, marginR, marginB := pdf.GetMargins()

	// res.Wd = pageWidth - marginL - marginR
	// res.Ht = pageHeight - marginB - marginT

	return gofpdf.SizeType{
		Wd: pageWidth - marginL - marginR,
		Ht: pageHeight - marginB - marginT,
	}
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

func parseData(val reflect.Value, mappers map[string]Mapper) *PDFFileStruct {
	// val := reflect.ValueOf(data)
	if val.Kind() != reflect.Pointer {
		return nil
	}

	res := PDFFileStruct{}
	val = val.Elem()

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
		case string(HeaderType):
			if strings.TrimSpace(valueField.String()) == "" {
				continue
			}
			res.headerParts = append(res.headerParts, strings.TrimSpace(valueField.String()))

		case string(ContentType):
			innerSection := PDFFileStruct{}

			headerTag := strings.TrimSpace(tagContent.HeaderName)
			if headerTag != SkipTag {
				innerSection.headerParts = append(innerSection.headerParts, strings.TrimSpace(tagContent.HeaderName))
			}

			switch valueField.Kind() {
			case reflect.String:
				innerSection.content = Map(mappers, strings.TrimSpace(valueField.String()), strings.TrimSpace(tagContent.MapperType))

			case reflect.Slice:
				for j := 0; j < valueField.Len(); j++ {
					switch valueField.Index(j).Kind() {
					case reflect.String:
						innerSection.content += strings.TrimSpace(valueField.Index(j).String()) + " "

					case reflect.Pointer:
						innerInnerSection := parseData(valueField.Index(j), mappers)

						innerSection.appendSection(innerInnerSection)
					}
				}
			}

			if prefix := strings.TrimSpace(tagContent.Prefix); prefix != "" && innerSection.content != "" {
				innerSection.content = fmt.Sprintf("%s: %s", prefix, innerSection.content)
			}

			res.appendSection(&innerSection)

		}
	}
	return &res
}

func ParseData(data any, mappers map[string]Mapper) *PDFFileStruct {
	val := reflect.ValueOf(data)
	return parseData(val, mappers)
}

func MarshalPDF(config *PDFConfig, data any) (*gofpdf.Fpdf, error) {
	pdf := initPDF(config)

	pdfFileStruct := ParseData(data, config.Mappers)
	pdfFileStruct.Build(pdf, config)

	return pdf, nil
}
