package pdf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

const (
	PDFTagName    = "pdf"
	Header1Height = 10
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

	// "B"- bold;
	// "I" - italic;
	// "U" - underscore;
	// "S" - strike-out
	Style  string
	Size   float64
	Height float64
}

type PDFConfig struct {
	// "P" or "Portrait"; "L" or "Landscape"
	Orientation string
	//"pt" for point, "mm" for millimeter, "cm" for centimeter, or "in" for inch
	Unit string
	// "A3", "A4", "A5", "Letter", "Legal", or "Tabloid". An empty string will be replaced with "A4"
	PageFormat      string
	FontDir         string
	LeftMargin      float64
	TopMargin       float64
	RightMargin     float64
	BottomMargin    float64
	Header1Font     FontConfig
	Header2Font     FontConfig
	Header3Font     FontConfig
	RegularTextFont FontConfig
}

func getContentAreaSize(pdf *gofpdf.Fpdf) gofpdf.SizeType {
	fmt.Println(2.521)

	res := gofpdf.SizeType{}
	pageWidth, pageHeight := pdf.GetPageSize()
	marginL, marginT, marginR, marginB := pdf.GetMargins()
	fmt.Println(2.522)

	res.Wd = pageWidth - marginL - marginR
	res.Ht = pageHeight - marginB - marginT
	fmt.Println(2.523)

	return res
}

func initPDF(config *PDFConfig) *gofpdf.Fpdf {
	pdf := gofpdf.New(config.Orientation, config.Unit, config.PageFormat, config.FontDir)
	pdf.SetMargins(config.LeftMargin, config.TopMargin, config.RightMargin)
	pdf.SetAutoPageBreak(true, config.BottomMargin)

	pdf.AddPage()

	return pdf
}

func insertTextLines(pdf *gofpdf.Fpdf, currFontConfig *FontConfig, lines []string) {
	fmt.Printf("lines: %v\n", lines)
	for _, line := range lines {
		fmt.Println(2.531)
		pdf.Cell(getContentAreaSize(pdf).Wd, currFontConfig.Height, line)
		fmt.Println(2.532)
	}
}

func addTitle(pdf *gofpdf.Fpdf, config *PDFConfig, title string) {
	fmt.Println(2.51)
	pdf.SetFont(config.Header1Font.Family, config.Header1Font.Style, config.Header1Font.Size)
	fmt.Println(2.52)
	fmt.Printf("title: %s\n", title)
	contentWidth := getContentAreaSize(pdf).Wd
	fmt.Printf("contentWidth: %v\n", contentWidth)

	lines := pdf.SplitText(title, contentWidth)
	fmt.Printf("lines: %v\n", lines)
	fmt.Println(2.53)
	insertTextLines(pdf, &config.Header1Font, lines)
	fmt.Println(2.54)
}

func addSection(pdf *gofpdf.Fpdf, config *PDFConfig, headerName, content string) {
	pdf.SetFont(config.Header2Font.Family, config.Header2Font.Style, config.Header2Font.Size)

	headerLines := pdf.SplitText(headerName, getContentAreaSize(pdf).Wd)
	insertTextLines(pdf, &config.Header1Font, headerLines)

	contentLines := pdf.SplitText(content, getContentAreaSize(pdf).Wd)
	insertTextLines(pdf, &config.Header1Font, contentLines)
}

func parseData(pdf *gofpdf.Fpdf, config *PDFConfig, data any) {
	fmt.Println(2.1)
	fmt.Printf("data: %#v\n", data)
	val := reflect.ValueOf(data).Elem()
	fmt.Println(2.2)

	fmt.Printf("%T have %d fields:\n", data, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		pdfTag := typeField.Tag.Get(PDFTagName)

		fmt.Printf(
			"\tname=%v, type=%v, value=%v, tag=`%v`\n",
			typeField.Name,
			typeField.Type.Kind(),
			valueField,
			pdfTag,
		)

		if pdfTag == "" {
			continue
		}
		fmt.Println(2.3)

		tagContent := parsePDFTagContent(pdfTag)
		fmt.Printf("tagContent: %#v\n", tagContent)
		fmt.Println(2.4)

		switch tagContent.Type {
		case string(TitleType):
			fmt.Println(2.5)
			addTitle(pdf, config, valueField.String())
			fmt.Println(2.6)

		case string(ContentType):
			fmt.Println(2.7)
			if strings.TrimSpace(tagContent.HeaderName) == "" {
				continue
			}
			fmt.Println(2.8)
			addSection(pdf, config, tagContent.HeaderName, valueField.String())
		}
		fmt.Println(2.9)
	}
}

func MarshalPDF(config *PDFConfig, data any) (*gofpdf.Fpdf, error) {
	fmt.Println(1)
	pdf := initPDF(config)
	fmt.Println(2)

	parseData(pdf, config, data)
	fmt.Println(3)

	return pdf, nil
}
