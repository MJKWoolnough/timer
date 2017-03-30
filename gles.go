// +build gles

package main

import (
	"unsafe"

	_ "github.com/MJKWoolnough/engine/graphics/gles2"
	"github.com/go-gl/gl/v3.1/gles2"
)

func setup() {
	gles2.ClearColor(0, 0, 0, 1)
	pid := CreateProgram(vs, fs)
	oid = gles2.GetUniformLocation(pid, &offsetName[0])
	sid = gles2.GetUniformLocation(pid, &scaleName[0])
	cid = gles2.GetAttribLocation(pid, &posName[0])
	gles2.UseProgram(pid)
}

var (
	vs = []byte("uniform vec2 offset;" +
		"uniform vec2 scale;" +
		"attribute vec2 pos;" +
		"void main() {" +
		"	gl_Position = vec4((offset + pos) * scale, 0, 1);" +
		"}")
	fs = []byte("void main() {" +
		"	gl_FragColor = vec4(1, 1, 1, 1);" +
		"}")
	offsetName = []byte("offset\x00")
	scaleName  = []byte("scale\x00")
	posName    = []byte("pos\x00")
)

func CreateShader(typ uint32, source []byte) uint32 {
	id := gles2.CreateShader(typ)
	sourceP := &source[0]
	sourceL := int32(len(source))
	gles2.ShaderSource(id, 1, &sourceP, &sourceL)
	gles2.CompileShader(id)
	return id
}

func CreateProgram(vertexShader, fragmentShader []byte) uint32 {
	vs := CreateShader(gles2.VERTEX_SHADER, vertexShader)
	gles2.GetError()
	fs := CreateShader(gles2.FRAGMENT_SHADER, fragmentShader)
	pid := gles2.CreateProgram()
	gles2.AttachShader(pid, vs)
	gles2.AttachShader(pid, fs)
	gles2.LinkProgram(pid)
	return pid
}

var sid, cid, oid int32

const (
	BarTop = 1 << iota
	BarTopLeft
	BarTopRight
	BarMiddle
	BarBottomLeft
	BarBottomRight
	BarBottom
)

var (
	digits = [10]byte{
		BarTop | BarTopLeft | BarTopRight | BarBottomLeft | BarBottomRight | BarBottom,
		BarTopRight | BarBottomRight,
		BarTop | BarTopRight | BarMiddle | BarBottomLeft | BarBottom,
		BarTop | BarTopRight | BarMiddle | BarBottomRight | BarBottom,
		BarTopLeft | BarTopRight | BarMiddle | BarBottomRight,
		BarTop | BarTopLeft | BarMiddle | BarBottomRight | BarBottom,
		BarTop | BarTopLeft | BarMiddle | BarBottomLeft | BarBottomRight | BarBottom,
		BarTop | BarTopRight | BarBottomRight,
		BarTop | BarTopLeft | BarTopRight | BarMiddle | BarBottomLeft | BarBottomRight | BarBottom,
		BarTop | BarTopLeft | BarTopRight | BarMiddle | BarBottomRight | BarBottom,
	}
	horizontalBar = [6]XY{
		{bHeight * tWidth, 0},              //bottom-left
		{0, bHeight / 2},                   //middle-left
		{bWidth - bHeight*tWidth, 0},       //bottom-right
		{bHeight * tWidth, bHeight},        //top-left
		{bWidth, bHeight / 2},              //middle-right
		{bWidth - bHeight*tWidth, bHeight}, //top-right
	}
	verticalBar = [6]XY{
		{0, bHeight * tWidth},
		{bHeight / 2, 0},
		{0, bWidth - bHeight*tWidth},
		{bHeight, bHeight * tWidth},
		{bHeight / 2, bWidth},
		{bHeight, bWidth - bHeight*tWidth},
	}
	bars = [7]Bar{
		{XY{tWidth, 2 * bWidth}, &horizontalBar},    //top
		{XY{0, bWidth + tWidth}, &verticalBar},      //top-left
		{XY{bWidth, bWidth + tWidth}, &verticalBar}, //top-right
		{XY{tWidth, bWidth}, &horizontalBar},        //middle
		{XY{0, tWidth}, &verticalBar},               //bottom-left
		{XY{bWidth, tWidth}, &verticalBar},          //bottom-right
		{XY{tWidth, 0}, &horizontalBar},             //bottom
	}
	scale  = XY{0.02, 0.03}
	offset = XY{0, 0}
)

const (
	tWidth  float32 = bHeight / 2
	bHeight float32 = 1
	bWidth  float32 = 5
)

type Bar struct {
	Offset XY
	Bar    *[6]XY
}

type XY [2]float32

func render(w, h int, t float64) {
	gles2.Clear(gles2.COLOR_BUFFER_BIT)

	gles2.Uniform2f(sid, scale[0], scale[1]*r)

	gap := bWidth + tWidth + 1

	sec := int(t) % (360000)

	hours := sec / 3600
	minutes := sec / 60 % 60
	seconds := sec % 60

	displayDigit(hours/10%10, -20-2*gap)
	displayDigit(hours%10, -20-gap)

	displayDigit(minutes/10%10, -gap)
	displayDigit(minutes%10, 0)

	displayDigit(seconds/10%10, 20)
	displayDigit(seconds%10, 20+gap)

}

func displayDigit(p int, offsetX float32) {
	for n, b := range bars {
		if digits[p]&(1<<uint(n)) > 0 {
			gles2.Uniform2f(oid, offset[0]+offsetX+b.Offset[0], offset[1]+b.Offset[1])
			vertices := b.Bar
			gles2.VertexAttribPointer(uint32(cid), 2, gles2.FLOAT, false, 2*4, unsafe.Pointer(vertices))
			gles2.EnableVertexAttribArray(0)

			gles2.DrawArrays(gles2.TRIANGLE_STRIP, 0, int32(len(vertices)))
		}
	}
}
