package pdfgen

import (
	"fmt"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"tilank/dto"
	"time"
)

func GeneratePDF(viol *dto.Violation, truck *dto.Truck, rules *dto.Rules) error {
	if rules == nil {
		rules = &dto.Rules{
			ID:          primitive.ObjectID{},
			Description: "Diberikan teguran.",
		}
	}

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetCompression(true)
	m.SetPageMargins(20, 10, 20)

	err := buildHeading(m)
	if err != nil {
		return err
	}
	buildBody(m, viol, truck, rules)
	err = buildImage(m, viol)
	if err != nil {
		return err
	}
	buildSignature(m, viol)

	err = m.OutputFileAndClose(fmt.Sprintf("static/pdf/%s.pdf", viol.ID.Hex()))
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

func buildBody(m pdf.Maroto, viol *dto.Violation, truck *dto.Truck, rules *dto.Rules) {
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
			textBody(m, viol.NoIdentity, 0)
			textBody(m, viol.NoPol, 5)
			textBody(m, viol.Location, 10)
			textBody(m, time.Unix(viol.CreatedAt, 0).Format("02-01-2006 15:04"), 15)
			textBody(m, strconv.Itoa(truck.Score), 20)
			textBody(m, viol.DetailViolation, 25)
		})
		m.ColSpace(1)
		m.Col(3, func() {
			textBody(m, rules.Description, 0)
		})
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
		m.Col(4, func() {
			textBodyCenter(m, fmt.Sprintf("Banjarmain %s", approveAt), 0)
		})
		m.ColSpace(8)
	})
	m.Row(25, func() {
		m.Col(4, func() {
			m.QrCode(data.ID.Hex(), props.Rect{
				Top:     0,
				Percent: 100,
				Center:  true,
			})
		})
		m.ColSpace(8)
	})
	m.Row(10, func() {
		m.Col(4, func() {
			textBodyCenter(m, "HSSE TPKB", 5)
		})
		m.ColSpace(8)
	})
}
