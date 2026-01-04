// package main

// import (
// 	"path/filepath"
// 	"time"

// 	"CubeFall/gamelogic"
// 	"CubeFall/helper"
// 	"CubeFall/objects"

// 	"github.com/go-gl/gl/v3.3-core/gl"
// 	"github.com/go-gl/mathgl/mgl32"
// 	"github.com/veandco/go-sdl2/sdl"
// )

// const winWidth = 1280
// const winHeight = 730

// var keyStates = sdl.GetKeyboardState()

// var bullets []*objects.Bullet

// type GameState int

// const (
// 	StateMenu GameState = iota
// 	StatePlaying
// 	StateExit
// )

// var currentState GameState = StateMenu

// func RenderMenu() {
// 	gl.ClearColor(0.8, 0.1, 0.1, 1.0) // RED = menu
// 	gl.Clear(gl.COLOR_BUFFER_BIT)
// }

// func main() {
// 	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
// 		panic(err)
// 	}
// 	defer sdl.Quit()

// 	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
// 	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
// 	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

// 	scoreTrack := gamelogic.NewScore(0, 0)

// 	window, err := sdl.CreateWindow(
// 		"CubeFall",
// 		50,
// 		30,
// 		winWidth,
// 		winHeight,
// 		sdl.WINDOW_OPENGL,
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer window.Destroy()

// 	window.GLCreateContext()
// 	if err := gl.Init(); err != nil {
// 		panic(err)
// 	}

// 	gl.Enable(gl.DEPTH_TEST)

// 	vertexShaderPath := filepath.Join("shaders", "first.vert")
// 	fragmentShaderPath := filepath.Join("shaders", "quad_tex.frag")

// 	shaderProgram, err := helper.NewShader(vertexShaderPath, fragmentShaderPath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	worldUp := mgl32.Vec3{0, 1, 0}
// 	cameraPos := mgl32.Vec3{0, 0, 3}
// 	camera := helper.NewCamera(cameraPos, worldUp, -90, 0, 0.01, 0.4)

// 	var land objects.Land
// 	var enemy objects.Enemey
// 	var player objects.Player
// 	var bullet objects.Bullet

// 	land.New()
// 	enemy.New()
// 	player.New()
// 	bullet.New()

// 	land.LoadVertexAttribs()
// 	enemy.LoadVertexAttribs()
// 	player.LoadVertexAttribs()
// 	bullet.LoadVertexAttribs()

// 	var elapsedTime float32
// 	enemySpeed := float32(0.01111)

// 	bim := objects.BulletInMotion{
// 		PosX: camera.Position.X(),
// 		PosY: camera.Position.Y(),
// 		PosZ: camera.Position.Z(),
// 	}

// 	prevMouseX, prevMouseY, _ := sdl.GetMouseState()

// 	running := true
// 	for running {
// 		frameStart := time.Now()

// 		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
// 			switch e := event.(type) {

// 			case *sdl.QuitEvent:
// 				currentState = StateExit

// 			case *sdl.KeyboardEvent:
// 				if e.Type == sdl.KEYDOWN {
// 					switch e.Keysym.Sym {

// 					case sdl.K_RETURN:
// 						if currentState == StateMenu {
// 							currentState = StatePlaying
// 						}

// 					case sdl.K_ESCAPE:
// 						if currentState == StateMenu {
// 							currentState = StateExit
// 						}
// 					}
// 				}
// 			}
// 		}

// 		switch currentState {

// 		case StateMenu:
// 			RenderMenu()
// 			window.GLSwap()
// 			continue

// 		case StatePlaying:
// 			if keyStates[sdl.SCANCODE_SPACE] != 0 {
// 				bullet.IsFired = true
// 			} else {
// 				bullet.IsFired = false
// 			}

// 			mouseX, mouseY, _ := sdl.GetMouseState()
// 			direction := helper.Nowhere

// 			if keyStates[sdl.SCANCODE_A] != 0 {
// 				bim.PosX = camera.Position.X()
// 				direction = helper.Left
// 			}
// 			if keyStates[sdl.SCANCODE_D] != 0 {
// 				bim.PosX = camera.Position.X()
// 				direction = helper.Right
// 			}
// 			if keyStates[sdl.SCANCODE_W] != 0 {
// 				bim.PosZ = camera.Position.Z()
// 				direction = helper.Forward
// 			}
// 			if keyStates[sdl.SCANCODE_S] != 0 {
// 				bim.PosZ = camera.Position.Z()
// 				direction = helper.Backward
// 			}

// 			if camera.Position.Y() <= 0.3 {
// 				camera.Position = mgl32.Vec3{
// 					camera.Position.X(),
// 					0.4,
// 					camera.Position.Z(),
// 				}
// 			}

// 			camera.UpdateCamera(
// 				direction,
// 				elapsedTime,
// 				camera.MovementSpeed,
// 				float32(mouseX-prevMouseX),
// 				-float32(mouseY-prevMouseY),
// 			)

// 			prevMouseX = mouseX
// 			prevMouseY = mouseY

// 			gl.ClearColor(0.53, 0.81, 0.92, 1.0)
// 			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

// 			shaderProgram.Use()

// 			projection := mgl32.Perspective(
// 				mgl32.DegToRad(45),
// 				float32(winWidth)/float32(winHeight),
// 				0.1,
// 				200,
// 			)
// 			view := camera.GetViewMatrix()

// 			shaderProgram.SetMat4("projection", projection)
// 			shaderProgram.SetMat4("view", view)

// 			land.Renderer(shaderProgram)
// 			enemy.Renderer(shaderProgram)

// 			gamelogic.MoveEnemies(&enemy.Extras, &enemySpeed)

// 			if gamelogic.AllEnemiesAreHit(&enemy.Extras) || len(enemy.Extras) == 0 {
// 				gamelogic.LevelUp(scoreTrack, &enemySpeed, &enemy.Extras)
// 			}

// 			player.Renderer(camera, shaderProgram)
// 			bullet.Renderer(camera, shaderProgram, &bim)

// 			hitEnemy, err := gamelogic.GetHitEnemy(camera, &bim, &enemy.Extras)
// 			if err == nil {
// 				gamelogic.KillEnemy(&hitEnemy, &enemy.Extras, scoreTrack)
// 			}

// 			gamelogic.HandlePassedEnemies(&enemy.Extras, camera, scoreTrack)

// 			window.GLSwap()
// 			shaderProgram.CheckForShaderChanges()

// 			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)

// 		case StateExit:
// 			gamelogic.PrintGameStats(scoreTrack)
// 			running = false
// 		}
// 	}
// }

// func Shoot(camera *helper.Camera) {
// 	b := &objects.Bullet{}
// 	b.New()
// 	b.LoadVertexAttribs()
// 	b.Fire(camera.Position)
// 	bullets = append(bullets, b)
// }
