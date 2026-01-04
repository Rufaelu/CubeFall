package Menu

import (
	"path/filepath"
	"time"

	"CubeFall/helper"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth = 1280
const winHeight = 730

func RunMainMenu(window *sdl.Window) bool {
	menuShader, err := helper.NewShader(
		filepath.Join("shaders", "menu.vert"), filepath.Join("shaders", "menu.frag"),
	)
	if err != nil {
		panic(err)
	}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {

			case *sdl.QuitEvent:
				return false

			case *sdl.MouseButtonEvent:
				if e.Type == sdl.MOUSEBUTTONDOWN && e.Button == sdl.BUTTON_LEFT {
					mx := float32(e.X)
					my := float32(e.Y)

					if inside(mx, my, 490, 260, 300, 70) {
						time.Sleep(150 * time.Millisecond)
						return true
					}

					if inside(mx, my, 490, 360, 300, 70) {
						return false
					}
				}
			}
		}

		gl.ClearColor(0.08, 0.08, 0.12, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		DrawRect(menuShader, 490, 260, 300, 70, 0.2, 0.7, 0.2)
		DrawRect(menuShader, 490, 360, 300, 70, 0.7, 0.2, 0.2)

		window.GLSwap()
		sdl.Delay(16)
	}
}

func DrawRect(shader *helper.Shader, x, y, w, h float32, r, g, b float32) {
	nx := func(px float32) float32 {
		return (px/float32(winWidth))*2 - 1
	}
	ny := func(py float32) float32 {
		return 1 - (py/float32(winHeight))*2
	}

	vertices := []float32{
		nx(x), ny(y),
		nx(x + w), ny(y),
		nx(x + w), ny(y + h),
		nx(x), ny(y + h),
	}

	indices := []uint32{0, 1, 2, 2, 3, 0}

	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	shader.Use()
	shader.SetVec3("color", r, g, b)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
}

func inside(mx, my, x, y, w, h float32) bool {
	return mx >= x && mx <= x+w && my >= y && my <= y+h
}
