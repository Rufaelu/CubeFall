package objects

import (
	"path/filepath"

	"CubeFall/helper"
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type PoppedEnemies []mgl32.Vec3

type ExtraEnemyProperty struct {
	Id               uint32
	Position         mgl32.Vec3
	IsHit            bool
	HavePassedPlayer bool
	BaseY            float32 // original Y
	BounceAmp        float32 // height of bounce
	BounceSpd        float32 // speed of bounce
	Time             float32 // internal timer
}

type Enemey struct {
	ModelMatrix   mgl32.Mat4
	textures      []helper.TextureId
	Vertices      []float32
	Extras        []ExtraEnemyProperty
	VAO           helper.BufferId
	MovementSpeed float32
}

func (enemy *Enemey) New() {
	green_file_path := filepath.Join("assets", "diddy.png")
	green_texture := helper.LoadTextureAlphaPng(green_file_path)
	white_file_path := filepath.Join("assets", "magic-star.png")
	white_texture := helper.LoadTextureAlphaPng(white_file_path)

	enemy.textures = []helper.TextureId{green_texture, white_texture}

	stacks := 18
	sectors := 36
	radius := float32(0.5)

	enemy.Vertices = []float32{}

	for i := 0; i < stacks; i++ {
		lat1 := math.Pi*float64(i)/float64(stacks) - math.Pi/2
		lat2 := math.Pi*float64(i+1)/float64(stacks) - math.Pi/2

		y1 := radius * float32(math.Sin(lat1))
		y2 := radius * float32(math.Sin(lat2))

		r1 := radius * float32(math.Cos(lat1))
		r2 := radius * float32(math.Cos(lat2))

		for j := 0; j < sectors; j++ {
			lng1 := 2 * math.Pi * float64(j) / float64(sectors)
			lng2 := 2 * math.Pi * float64(j+1) / float64(sectors)

			x1 := r1 * float32(math.Cos(lng1))
			z1 := r1 * float32(math.Sin(lng1))
			x2 := r1 * float32(math.Cos(lng2))
			z2 := r1 * float32(math.Sin(lng2))

			x3 := r2 * float32(math.Cos(lng1))
			z3 := r2 * float32(math.Sin(lng1))
			x4 := r2 * float32(math.Cos(lng2))
			z4 := r2 * float32(math.Sin(lng2))

			u1 := float32(j) / float32(sectors)
			u2 := float32(j+1) / float32(sectors)
			v1 := float32(i) / float32(stacks)
			v2 := float32(i+1) / float32(stacks)

			// Triangle 1
			enemy.Vertices = append(enemy.Vertices,
				x1, y1, z1, u1, v1,
				x3, y2, z3, u1, v2,
				x2, y1, z2, u2, v1,
			)

			// Triangle 2
			enemy.Vertices = append(enemy.Vertices,
				x2, y1, z2, u2, v1,
				x3, y2, z3, u1, v2,
				x4, y2, z4, u2, v2,
			)
		}
	}

	enemy.Extras = []ExtraEnemyProperty{
		{Id: 1, IsHit: false, Position: mgl32.Vec3{-4.0, 0.6, -14.0}, HavePassedPlayer: false},
		{Id: 2, IsHit: false, Position: mgl32.Vec3{1.0, 0.6, -14.0}, HavePassedPlayer: false},
		{Id: 3, IsHit: false, Position: mgl32.Vec3{4.0, 0.6, -14.0}, HavePassedPlayer: false},
		{Id: 4, IsHit: false, Position: mgl32.Vec3{8.0, 0.6, -14.0}, HavePassedPlayer: false},
	}

	// enemy.Positions = []mgl32.Vec3{
	// 	{-4.0, 0.6, -7.0},
	// 	{1.0, 0.6, -7.0},
	// 	{4.0, 0.6, -7.0},
	// 	{8.0, 0.6, -7.0},
	// 	{-4.0, 0.6, -10.0},
	// 	{1.0, 0.6, -10.0},
	// 	{4.0, 0.6, -10.0},
	// 	{8.0, 0.6, -10.0},
	// }
}

func (enemy *Enemey) LoadVertexAttribs() {
	enemy.VAO = helper.GenBindVertexArray(2)
	helper.GenBindBuffer(gl.ARRAY_BUFFER, 2)
	helper.BufferDataFloat(gl.ARRAY_BUFFER, enemy.Vertices, gl.DYNAMIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)
	helper.UnbindVertexArray()
}

var increment float32 = 0.0444
var angle float32 = 0

func (enemy *Enemey) Renderer(shader_program *helper.Shader) {
	helper.BindVertextArray(enemy.VAO)
	for i, prop := range enemy.Extras {
		if i%2 == 0 {
			helper.BindTexture(enemy.textures[0])
		} else {
			helper.BindTexture(enemy.textures[1])
		}
		angle += increment
		enemy.ModelMatrix = mgl32.Ident4()
		// enemy.ModelMatrix = mgl32.HomogRotate3DX(mgl32.DegToRad(angle)).Mul4(enemy.ModelMatrix)
		enemy.ModelMatrix = mgl32.HomogRotate3DY(mgl32.DegToRad(angle)).Mul4(enemy.ModelMatrix)

		enemy.ModelMatrix = mgl32.Translate3D(prop.Position.X(), prop.Position.Y(), prop.Position.Z()).Mul4(enemy.ModelMatrix)
		shader_program.SetMat4("model", enemy.ModelMatrix)
		// gl.DrawArrays(gl.TRIANGLES, 0, 36)
		vertexCount := int32(len(enemy.Vertices) / 5)
		gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)

	}
}
