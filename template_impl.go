package gofpdf

//
//  GoFPDI
//
//    Copyright 2015 Marcus Downing
//
//  FPDI
//
//    Copyright 2004-2014 Setasign - Jan Slabon
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

import (
// "bytes"
)

// newTpl creates a template, copying graphics settings from a template if one is given
func newTpl(corner PointType, size SizeType, unitStr, fontDirStr string, fn func(*Tpl), copyFrom *Fpdf) Template {
	orientationStr := "p"
	if size.Wd > size.Ht {
		orientationStr = "l"
	}
	sizeStr := ""

	fpdf := fpdfNew(orientationStr, unitStr, sizeStr, fontDirStr, size)
	tpl := Tpl{*fpdf}
	if copyFrom != nil {
		tpl.loadParamsFromFpdf(copyFrom)
	}
	tpl.Fpdf.SetAutoPageBreak(false, 0)
	tpl.Fpdf.AddPage()
	fn(&tpl)
	bytes := tpl.Fpdf.pages[tpl.Fpdf.page].Bytes()
	templates := make([]Template, 0, len(tpl.Fpdf.templates))
	for _, t := range tpl.Fpdf.templates {
		templates = append(templates, t)
	}
	images := tpl.Fpdf.images

	id := GenerateTemplateID()
	template := FpdfTpl{id, corner, size, bytes, images, templates}
	return &template
}

// FpdfTpl is a concrete implementation of the Template interface.
type FpdfTpl struct {
	id        int64
	corner    PointType
	size      SizeType
	bytes     []byte
	images    map[string]*ImageInfoType
	templates []Template
}

// ID returns the global template identifier
func (t *FpdfTpl) ID() int64 {
	return t.id
}

// Size gives the bounding dimensions of this template
func (t *FpdfTpl) Size() (corner PointType, size SizeType) {
	return t.corner, t.size
}

// Bytes returns the actual template data, not including resources
func (t *FpdfTpl) Bytes() []byte {
	return t.bytes
}

// Images returns a list of the images used in this template
func (t *FpdfTpl) Images() map[string]*ImageInfoType {
	return t.images
}

// Templates returns a list of templates used in this template
func (t *FpdfTpl) Templates() []Template {
	return t.templates
}

// Tpl is an Fpdf used for writing a template.
// It has most of the facilities of an Fpdf,but cannot add more pages.
// Tpl is used directly only during the limited time a template is writable.
type Tpl struct {
	Fpdf
}

func (t *Tpl) loadParamsFromFpdf(f *Fpdf) {
	t.Fpdf.compress = false

	t.Fpdf.k = f.k
	t.Fpdf.x = f.x
	t.Fpdf.y = f.y
	t.Fpdf.lineWidth = f.lineWidth
	t.Fpdf.capStyle = f.capStyle
	t.Fpdf.joinStyle = f.joinStyle

	t.Fpdf.color.draw = f.color.draw
	t.Fpdf.color.fill = f.color.fill
	t.Fpdf.color.text = f.color.text

	t.Fpdf.currentFont = f.currentFont
	t.Fpdf.fontFamily = f.fontFamily
	t.Fpdf.fontSize = f.fontSize
	t.Fpdf.fontSizePt = f.fontSizePt
	t.Fpdf.fontStyle = f.fontStyle
	t.Fpdf.ws = f.ws
}

// AddPage does nothing because you cannot add pages to a template
func (t *Tpl) AddPage() {
}

// AddPageFormat does nothign becasue you cannot add pages to a template
func (t *Tpl) AddPageFormat(orientationStr string, size SizeType) {
}

// SetAutoPageBreak does nothing because you cannot add pages to a template
func (t *Tpl) SetAutoPageBreak(auto bool, margin float64) {
}
