package pdfgen

import (
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func textH1(m pdf.Maroto, text string) {
	m.Text(text, props.Text{
		Top:         3,
		Style:       consts.Bold,
		Size:        18,
		Align:       consts.Center,
		Extrapolate: true,
		Color:       getDarkPurpleColor(),
	})

}

func textH2(m pdf.Maroto, text string, top float64) {
	m.Text(text, props.Text{
		Top:         top,
		Extrapolate: false,
		Style:       consts.Bold,
		Size:        14,
		Color:       getDarkPurpleColor(),
	})

}

func textBody(m pdf.Maroto, text string, top float64) {
	m.Text(text, props.Text{
		Top:         top,
		Extrapolate: false,
		Color:       getDarkPurpleColor(),
	})

}

func textBodyCenter(m pdf.Maroto, text string, top float64) {
	m.Text(text, props.Text{
		Top:         top,
		Extrapolate: false,
		Align:       consts.Center,
		Color:       getDarkPurpleColor(),
	})

}

func getDarkPurpleColor() color.Color {
	return color.Color{
		Red:   88,
		Green: 80,
		Blue:  99,
	}
}
