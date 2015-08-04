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

// CreateTemplate defines a new template using the current page size.
func (f *Fpdf) CreateTemplate(fn func(*Tpl)) Template {
	return newTpl(PointType{0, 0}, f.curPageSize, f.unitStr, f.fontDirStr, fn, f)
}

// CreateTemplateCustom starts a template, using the given bounds.
func (f *Fpdf) CreateTemplateCustom(corner PointType, size SizeType, fn func(*Tpl)) Template {
	return newTpl(corner, size, f.unitStr, f.fontDirStr, fn, f)
}

// CreateTemplate creates a template not attached to any document
func CreateTemplate(corner PointType, size SizeType, unitStr, fontDirStr string, fn func(*Tpl)) Template {
	return newTpl(corner, size, unitStr, fontDirStr, fn, nil)
}

// UseTemplate adds a template to the current page or another template,
// using the size and position at which it was originally written.
func (f *Fpdf) UseTemplate(t Template) {
	if t == nil {
		f.SetErrorf("Template is nil")
		return
	}
	corner, size := t.Size()
	f.UseTemplateScaled(t, corner, size)
}

// UseTemplateScaled adds a template to the current page or another template,
// using the given page coordinates.
func (f *Fpdf) UseTemplateScaled(t Template, corner PointType, size SizeType) {
	if t == nil {
		f.SetErrorf("Template is nil")
		return
	}

	// You have to add at least a page first
	if f.page <= 0 {
		f.SetErrorf("Cannot use a template without first adding a page")
		return
	}

	// make a note of the fact that we actually use this template, as well as any other templates,
	// images or fonts it uses
	f.templates[t.ID()] = t
	for _, tt := range t.Templates() {
		f.templates[tt.ID()] = tt
	}
	for name, ti := range t.Images() {
		name = sprintf("t%d-%s", t.ID(), name)
		f.images[name] = ti
	}

	// template data
	_, templateSize := t.Size()
	scaleX := size.Wd / templateSize.Wd
	scaleY := size.Ht / templateSize.Ht
	tx := corner.X * f.k
	ty := (f.curPageSize.Ht - corner.Y - size.Ht) * f.k

	f.outf("q %.4F 0 0 %.4F %.4F %.4F cm", scaleX, scaleY, tx, ty) // Translate
	f.outf("/TPL%d Do Q", t.ID())
}

var nextTemplateIDChannel = func() chan int64 {
	ch := make(chan int64)
	go func() {
		var nextID int64 = 1
		for {
			ch <- nextID
			nextID++
		}
	}()
	return ch
}()

// GenerateTemplateID gives the next template ID. These numbers are global so that they can never clash.
func GenerateTemplateID() int64 {
	return <-nextTemplateIDChannel
}

// Template is an object that can be written to, then used and re-used any number of times within a document.
type Template interface {
	ID() int64
	Size() (PointType, SizeType)
	Bytes() []byte
	Images() map[string]*ImageInfoType
	Templates() []Template
}

// putTemplates writes the templates to the PDF
func (f *Fpdf) putTemplates() {
	filter := ""
	if f.compress {
		filter = "/Filter /FlateDecode "
	}

	templates := sortTemplates(f.templates)
	var t Template
	for _, t = range templates {
		corner, size := t.Size()

		f.newobj()
		f.templateObjects[t.ID()] = f.n
		f.outf("<<%s/Type /XObject", filter)
		f.out("/Subtype /Form")
		f.out("/Formtype 1")
		f.outf("/BBox [%.2F %.2F %.2F %.2F]", corner.X*f.k, corner.Y*f.k, (corner.X+size.Wd)*f.k, (corner.Y+size.Ht)*f.k)
		if corner.X != 0 || corner.Y != 0 {
			f.outf("/Matrix [1 0 0 1 %.5F %.5F]", -corner.X*f.k*2, corner.Y*f.k*2)
		}

		// Template's resource dictionary
		f.out("/Resources ")
		f.out("<</ProcSet [/PDF /Text /ImageB /ImageC /ImageI]")

		tImages := t.Images()
		tTemplates := t.Templates()
		if len(tImages) > 0 || len(tTemplates) > 0 {
			f.out("/XObject <<")
			for _, ti := range tImages {
				f.outf("/I%d %d 0 R", ti.i, ti.n)
			}
			for _, tt := range tTemplates {
				id := tt.ID()
				if objID, ok := f.templateObjects[id]; ok {
					f.outf("/TPL%d %d 0 R", id, objID)
				}
			}
			f.out(">>")
		}

		f.out(">>")

		//  Write the template's byte stream
		buffer := t.Bytes()
		// fmt.Println("Put template bytes", string(buffer[:]))
		if f.compress {
			buffer = sliceCompress(buffer)
		}
		f.outf("/Length %d >>", len(buffer))
		f.putstream(buffer)
		f.out("endobj")
	}
}

// sortTemplates puts templates in a suitable order based on dependices
func sortTemplates(templates map[int64]Template) []Template {
	chain := make([]Template, 0, len(templates)*2)

	// build a full set of dependency chains
	for _, t := range templates {
		tlist := templateChainDependencies(t)
		for _, tt := range tlist {
			if tt != nil {
				chain = append(chain, tt)
			}
		}
	}

	// reduce that to make a simple list
	sorted := make([]Template, 0, len(templates))
chain:
	for _, t := range chain {
		for _, already := range sorted {
			if t == already {
				continue chain
			}
		}
		sorted = append(sorted, t)
	}

	return sorted
}

//  templateChainDependencies is a recursive function for determining the full chain of template dependencies
func templateChainDependencies(template Template) []Template {
	requires := template.Templates()
	chain := make([]Template, len(requires)*2)
	for _, req := range requires {
		for _, sub := range templateChainDependencies(req) {
			chain = append(chain, sub)
		}
	}
	chain = append(chain, template)
	return chain
}
