package gofpdf

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)

/*
 * Copyright (c) 2015 Kurt Jung (Gmail: kurt.w.jung),
 *   Marcus Downing, Jan Slabon (Setasign)
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

// newTpl creates a template, copying graphics settings from a template if one is given
func newTpl(corner PointType, size SizeType, orientationStr, unitStr, fontDirStr string, fn func(*Tpl), copyFrom *Fpdf) Template {
	sizeStr := ""

	fpdf := fpdfNew(orientationStr, unitStr, sizeStr, fontDirStr, size)
	tpl := Tpl{*fpdf}
	if copyFrom != nil {
		tpl.loadParamsFromFpdf(copyFrom)
	}
	tpl.Fpdf.AddPage()
	fn(&tpl)

	bytes := make([][]byte, len(tpl.Fpdf.pages))
	// skip the first page as it will always be empty
	for x := 1; x < len(bytes); x++ {
		bytes[x] = tpl.Fpdf.pages[x].Bytes()
	}

	templates := make([]Template, 0, len(tpl.Fpdf.templates))
	for _, key := range templateKeyList(tpl.Fpdf.templates, true) {
		templates = append(templates, tpl.Fpdf.templates[key])
	}
	images := tpl.Fpdf.images

	id := GenerateTemplateID()
	template := FpdfTpl{id, corner, size, bytes, images, templates, tpl.Fpdf.page}
	return &template
}

// FpdfTpl is a concrete implementation of the Template interface.
type FpdfTpl struct {
	id        int64
	corner    PointType
	size      SizeType
	bytes     [][]byte
	images    map[string]*ImageInfoType
	templates []Template
	page      int
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
	return t.bytes[t.page]
}

// FromPage creates a new template from a specific Page
func (t *FpdfTpl) FromPage(page int) (Template, error) {
	// pages start at 1
	if page == 0 {
		return nil, errors.New("Pages start at 1 No template will have a page 0")
	}

	if page > t.NumPages() {
		return nil, fmt.Errorf("The template does not have a page %d", page)
	}
	// if it is already pointing to the correct page
	// there is no need to create a new template
	if t.page == page {
		return t, nil
	}

	t2 := *t
	t2.id = GenerateTemplateID()
	t2.page = page
	return &t2, nil
}

// FromPages creates a template slice with all the pages within a template.
func (t *FpdfTpl) FromPages() []Template {
	p := make([]Template, t.NumPages())
	for x := 1; x <= t.NumPages(); x++ {
		// the only error is when accessing a
		// non existing template... that can't happen
		// here
		p[x-1], _ = t.FromPage(x)
	}

	return p
}

// Images returns a list of the images used in this template
func (t *FpdfTpl) Images() map[string]*ImageInfoType {
	return t.images
}

// Templates returns a list of templates used in this template
func (t *FpdfTpl) Templates() []Template {
	return t.templates
}

// NumPages returns the number of available pages within the template. Look at FromPage and FromPages on access to that content.
func (t *FpdfTpl) NumPages() int {
	// the first page is empty to
	// make the pages begin at one
	return len(t.bytes) - 1
}

// Serialize turns a template into a byte string for later deserialization
func (t *FpdfTpl) Serialize() ([]byte, error) {
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	err := enc.Encode(t)

	return b.Bytes(), err
}

// DeserializeTemplate creaties a template from a previously serialized
// template
func DeserializeTemplate(b []byte) (Template, error) {
	tpl := new(FpdfTpl)
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err := dec.Decode(tpl)
	return tpl, err
}

// childrenImages returns the next layer of children images, it doesn't dig into
// children of children. Applies template namespace to keys to ensure
// no collisions. See UseTemplateScaled
func (t *FpdfTpl) childrenImages() map[string]*ImageInfoType {
	childrenImgs := make(map[string]*ImageInfoType)

	for x := 0; x < len(t.templates); x++ {
		imgs := t.templates[x].Images()
		for key, val := range imgs {
			name := sprintf("t%d-%s", t.templates[x].ID(), key)
			childrenImgs[name] = val
		}
	}

	return childrenImgs
}

// GobEncode encodes the receiving template into a byte buffer. Use GobDecode
// to decode the byte buffer back to a template.
func (t *FpdfTpl) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)

	err := encoder.Encode(t.templates)
	childrenImgs := t.childrenImages()

	if err == nil {
		encoder.Encode(len(t.images))
	}

	for key, img := range t.images {
		// if the image has already been saved as a child, then
		// save nil so we don't duplicate data
		err = encoder.Encode(key)

		if err != nil {
			break
		}

		if _, ok := childrenImgs[key]; ok {
			err = encoder.Encode("p")
		} else {
			err = encoder.Encode("o")
			if err == nil {
				err = encoder.Encode(img)
			}
		}
	}
	if err == nil {
		err = encoder.Encode(t.id)
	}
	if err == nil {
		err = encoder.Encode(t.corner)
	}
	if err == nil {
		err = encoder.Encode(t.size)
	}
	if err == nil {
		err = encoder.Encode(t.bytes)
	}
	if err == nil {
		err = encoder.Encode(t.page)
	}

	return w.Bytes(), err
}

// GobDecode decodes the specified byte buffer into the receiving template.
func (t *FpdfTpl) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)

	templates := make([]*FpdfTpl, 0)
	err := decoder.Decode(&templates)
	t.templates = make([]Template, len(templates))

	for x := 0; x < len(templates); x++ {
		t.templates[x] = templates[x]
	}

	var numImgs int
	if err == nil {
		err = decoder.Decode(&numImgs)
	}

	t.images = make(map[string]*ImageInfoType)
	childrenImgs := t.childrenImages()

	for x := 0; x < numImgs; x++ {
		var key string
		var tpe string

		if err == nil {
			err = decoder.Decode(&key)
		}

		if err == nil {
			err = decoder.Decode(&tpe)
		}

		if err == nil {
			switch tpe {
			case "p":
				if _, ok := childrenImgs[key]; !ok {
					err = fmt.Errorf("Encoded template is corrupt, could not find image %s", key)
				} else {
					t.images[key] = childrenImgs[key]
				}
			case "o":
				var img *ImageInfoType
				err = decoder.Decode(&img)

				if err == nil {
					t.images[key] = img
				}
			}
		}
	}
	if err == nil {
		err = decoder.Decode(&t.id)
	}
	if err == nil {
		err = decoder.Decode(&t.corner)
	}
	if err == nil {
		err = decoder.Decode(&t.size)
	}
	if err == nil {
		err = decoder.Decode(&t.bytes)
	}
	if err == nil {
		err = decoder.Decode(&t.page)
	}

	return err
}

// Tpl is an Fpdf used for writing a template. It has most of the facilities of
// an Fpdf, but cannot add more pages. Tpl is used directly only during the
// limited time a template is writable.
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

	t.Fpdf.fonts = f.fonts
	t.Fpdf.currentFont = f.currentFont
	t.Fpdf.fontFamily = f.fontFamily
	t.Fpdf.fontSize = f.fontSize
	t.Fpdf.fontSizePt = f.fontSizePt
	t.Fpdf.fontStyle = f.fontStyle
	t.Fpdf.ws = f.ws
}
