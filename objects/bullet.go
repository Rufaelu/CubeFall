package objects

import (
	"math"
	"path/filepath"

	"CubeFall/helper"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Bullet struct {
	texture       helper.TextureId
	Vertices      []float32
	Position      mgl32.Vec3
	StartPosition mgl32.Vec3 // ✅ ADD THIS
	Direction     mgl32.Vec3
	Velocity      mgl32.Vec3
	ShotSpeed     float32
	VAO           helper.BufferId
	Alive         bool
	Radius        float32
}

// type BulletInMotion struct {
// 	PosX float32
// 	PosY float32
// 	PosZ float32
// }

func (b *Bullet) New() {
	b.texture = helper.LoadTextureAlphaPng(
		filepath.Join("assets", "diddy.png"),
	)
	b.Radius = 0.15

	b.ShotSpeed = 10.0
	b.Alive = false

	segments := 16
	radius := float32(0.05)

	b.Vertices = []float32{}

	for i := 0; i < segments; i++ {
		angle1 := float32(i) * 2.0 * math.Pi / float32(segments)
		angle2 := float32(i+1) * 2.0 * math.Pi / float32(segments)

		// center
		b.Vertices = append(b.Vertices,
			0, 0, 0, 0.5, 0.5,
		)

		// first point
		b.Vertices = append(b.Vertices,
			radius*float32(math.Cos(float64(angle1))),
			radius*float32(math.Sin(float64(angle1))),
			0,
			1, 0,
		)

		// second point
		b.Vertices = append(b.Vertices,
			radius*float32(math.Cos(float64(angle2))),
			radius*float32(math.Sin(float64(angle2))),
			0,
			0, 1,
		)
	}

}

func (b *Bullet) Fire(start, direction mgl32.Vec3) {
	b.Alive = true
	b.Position = start
	b.StartPosition = start // ✅ ADD THIS
	b.Direction = direction.Normalize()
	b.Velocity = b.Direction.Mul(b.ShotSpeed)
}

func (b *Bullet) Update(dt float32) {
	if !b.Alive {
		return
	}

	b.Position = b.Position.Add(
		b.Velocity.Mul(dt),
	)
}

func (b *Bullet) Render(shader *helper.Shader) {
	if !b.Alive {
		return
	}

	helper.BindVertextArray(b.VAO)
	helper.BindTexture(b.texture)

	model := mgl32.Translate3D(
		b.Position.X(),
		b.Position.Y(),
		b.Position.Z(),
	)

	shader.SetMat4("model", model)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(b.Vertices)/5))
}

func (b *Bullet) LoadVertexAttribs() {
	b.VAO = helper.GenBindVertexArray(3)
	helper.GenBindBuffer(gl.ARRAY_BUFFER, 1)
	helper.BufferDataFloat(gl.ARRAY_BUFFER, b.Vertices, gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	helper.UnbindVertexArray()
}
