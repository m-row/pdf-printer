package main

import (
	"bytes"
	"io"
	"log/slog"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/labstack/echo/v4"
)

func (c *Controller) Print(ctx echo.Context) error {
	v, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	pdf, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// Set global options
	pdf.Dpi.Set(300)
	pdf.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdf.Grayscale.Set(false)

	page := wkhtmltopdf.NewPageReader(bytes.NewReader(v))
	page.EnableLocalFileAccess.Set(true)
	page.FooterRight.Set("[page] / [toPage]")
	page.FooterLeft.Set("[date] [time]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(1)
	// Add to document
	pdf.AddPage(page)

	if err := pdf.Create(); err != nil {
		slog.Error(err.Error())
		return err
	}

	ctx.Response().Header().Set("Content-Type", "application/pdf")
	ctx.Response().
		Header().
		Add("Content-Disposition", "attachment;filename=generated.pdf")

	_, err = ctx.Response().Write(pdf.Buffer().Bytes())
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}
