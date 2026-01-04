package main

import (
	"fmt"
	"path/filepath"
	"time"

	"CubeFall/Menu"
	"CubeFall/gamelogic"
	"CubeFall/helper"
	"CubeFall/objects"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth = 1280
const winHeight = 730

// var keyStates = sdl.GetKeyboardState()

func main() {
	// ---------------- SDL INIT ----------------
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	scoreTrack := gamelogic.NewScore(0, 0)

	window, err := sdl.CreateWindow(
		"CubeFall",
		50, 30,
		winWidth, winHeight,
		sdl.WINDOW_OPENGL,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	// sdl.SetRelativeMouseMode(true)

	window.GLCreateContext()
	gl.Init()

	Menu.RunMainMenu(window)
	gl.Enable(gl.DEPTH_TEST)

	// ---------------- SHADER ----------------
	shader, err := helper.NewShader(
		filepath.Join("shaders", "first.vert"),
		filepath.Join("shaders", "quad_tex.frag"),
	)
	if err != nil {
		panic(err)
	}

	// ---------------- CAMERA ----------------
	camera := helper.NewCamera(
		mgl32.Vec3{0, 0, 3},
		mgl32.Vec3{0, 1, 0},
		-90, 0,
		2.7, 0.4,
	)

	// ---------------- OBJECTS ----------------
	var land objects.Land
	var enemy objects.Enemey
	var player objects.Player
	// var bullet objects.Bullet
	var bullets []objects.Bullet

	land.New()
	enemy.New()
	player.New()
	// bullet.New()

	land.LoadVertexAttribs()
	enemy.LoadVertexAttribs()
	player.LoadVertexAttribs()
	// bullet.LoadVertexAttribs()

	enemySpeed := float32(0.01111)

	prevMouseX, prevMouseY, _ := sdl.GetMouseState()
	previousTime := time.Now()
	fireRate := float32(0.001) // 10 bullets/sec
	fireTimer := float32(0.01)

	// ================= GAME LOOP =================
	for {
		// frameStart := time.Now()
		currentTime := time.Now()
		elapsedTime := float32(currentTime.Sub(previousTime).Seconds())
		previousTime = currentTime

		// -------- EVENTS --------
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				gamelogic.PrintGameStats(scoreTrack)
				return
			}
		}

		// -------- DELTA TIME --------
		// elapsedTime := float32(time.Since(frameStart).Seconds())
		sdl.PumpEvents()
		keyStates := sdl.GetKeyboardState()

		// -------- INPUT --------
		mouseX, mouseY, _ := sdl.GetMouseState()
		direction := helper.Nowhere

		if keyStates[sdl.SCANCODE_A] != 0 {

			direction = helper.Left
		}
		if keyStates[sdl.SCANCODE_D] != 0 {
			direction = helper.Right
		}
		if keyStates[sdl.SCANCODE_W] != 0 {
			direction = helper.Forward
		}
		if keyStates[sdl.SCANCODE_S] != 0 {
			direction = helper.Backward
		}

		// -------- SHOOT --------
		_, _, mouseButtons := sdl.GetMouseState()
		leftClick := mouseButtons&sdl.Button(sdl.BUTTON_LEFT) != 0
		//  := mouseButtons&sdl.Button(sdl.BUTTON_LEFT) != 0

		fireTimer -= elapsedTime

		if (keyStates[sdl.SCANCODE_SPACE] != 0 || leftClick) && fireTimer <= 0 {

			var b objects.Bullet

			// 🔴 THESE GO HERE
			b.New()
			b.LoadVertexAttribs()
			b.Fire(camera.Position, camera.Front)

			bullets = append(bullets, b)

			fireTimer = fireRate
		}

		for i := 0; i < len(bullets); i++ {
			bullets[i].Update(elapsedTime)

			gamelogic.BulletHitsEnemy(&bullets[i], &enemy.Extras, scoreTrack)

			if bullets[i].Position.Sub(bullets[i].StartPosition).Len() > 50 {
				bullets[i].Alive = false
			}

			if !bullets[i].Alive {
				bullets = append(bullets[:i], bullets[i+1:]...)
				i--
			}
		}
		// for i := range bullets {
		// 	bullets[i].Render(shader)
		// }

		// -------- UPDATE --------
		camera.UpdateCamera(
			direction,
			elapsedTime,
			camera.MovementSpeed,
			float32(mouseX-prevMouseX),
			-float32(mouseY-prevMouseY),
		)

		prevMouseX = mouseX
		prevMouseY = mouseY

		if camera.Position.Y() <= 0.3 {
			camera.Position[1] = 0.4
		}

		gamelogic.MoveEnemies(&enemy.Extras, &enemySpeed, elapsedTime)

		if gamelogic.AllEnemiesAreHit(&enemy.Extras) || len(enemy.Extras) < 1 {
			gamelogic.LevelUp(scoreTrack, &enemySpeed, &enemy.Extras)
		}

		// -------- RENDER --------
		gl.ClearColor(0.53, 0.81, 0.92, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shader.Use()

		projection := mgl32.Perspective(
			mgl32.DegToRad(45.0),
			float32(winWidth)/float32(winHeight),
			0.1, 200.0,
		)

		shader.SetMat4("projection", projection)
		shader.SetMat4("view", camera.GetViewMatrix())

		land.Renderer(shader)
		enemy.Renderer(shader)
		player.Renderer(camera, shader)

		for i := range bullets {
			bullets[i].Render(shader)
		}

		updateWindowTitle(window, scoreTrack)

		window.GLSwap()
		shader.CheckForShaderChanges()
	}
}

// ---------------- UI ----------------
func updateWindowTitle(window *sdl.Window, score *gamelogic.ScoreTrack) {
	window.SetTitle(fmt.Sprintf(
		"CubeFall | Score: %d | Kills: %d | Max Level: %d",
		score.Points,
		score.KillCount,
		score.PassedMaxLevel,
	))
}
