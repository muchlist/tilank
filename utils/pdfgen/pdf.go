package pdfgen

import (
	"fmt"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"tilank/dto"
	"time"
)

func GeneratePDF(data *dto.Violation) error {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 10, 20)

	// pelanggaran ke berapa
	// sangsinya apa

	err := buildHeading(m)
	if err != nil {
		return err
	}
	buildBody(m, data)
	err = buildImage(m, data)
	if err != nil {
		return err
	}
	buildSignature(m, data)

	err = m.OutputFileAndClose(fmt.Sprintf("static/pdf/%s.pdf", data.ID.Hex()))
	if err != nil {
		return err
	}
	return nil
}

func buildHeading(m pdf.Maroto) error {
	var errHead error

	m.Row(10, func() {

	})
	m.Row(20, func() {
		m.Col(2, func() {
			errHead = m.FileImage("static/image/logo/pelindo3.png", props.Rect{
				Percent: 100,
				Center:  false,
				Top:     3,
			})
		})
		m.Col(8, func() {
			textH1(m, "Surat Tilang Elektronik TPKB")
			textBodyCenter(m, "Aplikasi ETI Pelindo III Banjarmasin", 12)
		})
		m.ColSpace(2)
	})

	return errHead
}

func buildBody(m pdf.Maroto, data *dto.Violation) {
	m.Row(5, func() {

	})
	m.Row(15, func() {
		m.Col(8, func() {
			textH2(m, "Identifikasi Pelanggaran", 3)
		})
		m.Col(4, func() {
			textH2(m, "     [sangsi]", 3)
		})
	})
	m.Row(40, func() {
		m.Col(4, func() {
			textBody(m, "ID Truck / Nomer Lambung", 0)
			textBody(m, "Nomor Polisi", 5)
			textBody(m, "Lokasi", 10)
			textBody(m, "Tanggal", 15)
			textBody(m, "Pelanggaran ke-", 20)
			textBody(m, "Detail Pelanggaran", 25)
		})
		m.Col(4, func() {
			textBody(m, data.NoIdentity, 0)
			textBody(m, data.NoPol, 5)
			textBody(m, data.Location, 10)
			textBody(m, time.Unix(data.CreatedAt, 0).Format("02-01-2006 15:04"), 15)
			textBody(m, "2", 20) // todo pelanggaran keberapa
			textBody(m, data.DetailViolation, 25)
		})
		m.ColSpace(1)
		m.Col(3, func() {
			textBody(m, "Diberikan teguran ke dua (2), apabila melakukan pelanggaran 1 kali lagi di tahun yang sama maka akan dilakukan pemblokiran saat gate in pada truck dengan tersenbut", 0)
		}) // todo sangsinya apa
	})
}

func buildImage(m pdf.Maroto, data *dto.Violation) error {
	var err error
	m.Row(5, func() {

	})
	m.Row(10, func() {
		m.Col(12, func() {
			textH2(m, "Lampiran", 3)
		})
	})

	images := data.Images
	if len(images) == 0 {
		return nil
	}
	// memastikan hanya mencetak foto 3 pertama saja
	if len(images) > 3 {
		images = data.Images[:len(images)-3]
	}

	m.Row(60, func() {
		for _, image := range images {
			m.Col(4, func() {
				err = m.FileImage(fmt.Sprintf("static/%s", image), props.Rect{
					Percent: 90,
					Center:  true,
				})
			})
		}
	})

	return err
}

func buildSignature(m pdf.Maroto, data *dto.Violation) {
	approveAt := time.Unix(data.ApprovedAt, 0).Format("02-01-2006")
	m.Row(10, func() {

	})
	m.Row(8, func() {
		m.ColSpace(8)
		m.Col(4, func() {
			textBodyCenter(m, fmt.Sprintf("Banjarmain %s", approveAt), 0)
		})
	})
	m.Row(25, func() {
		m.ColSpace(8)
		m.Col(4, func() {
			m.QrCode(data.ID.Hex(), props.Rect{
				Top:     0,
				Percent: 100,
				Center:  true,
			})
		})
	})
	m.Row(10, func() {
		m.ColSpace(8)
		m.Col(4, func() {
			textBodyCenter(m, "HSSE TPKB", 5)
		})
	})
}
