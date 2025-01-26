package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"math/rand"
)

import (
	"ecs/lib"
)

type DebounceUpdateComponent struct {
	DebounceTime float32
	CurrentTime  float32
}

type TextComponent struct {
	Text     string
	FontSize int32
}

type PositionComponent struct {
	X, Y int32
}

type ColorComponent struct {
	R, G, B, A uint8
}

type ParticleSpawnComponent struct {
	LifeTime float32
}
type ParticleComponent struct {
	Radius float32
}
type LifetimeComponent struct {
	LifeTime    float32
	CurrentTime float32
}

type FPSComponent struct{}
type FrameTimeComponent struct{}
type ShouldUpdateComponent struct{}
type EntityCounterComponent struct{}
type VisibleComponent struct{}
type SpeedComponent struct {
	Speed float32
}

var textComponentID = lib.RegisterComponent[TextComponent]()
var positionComponentID = lib.RegisterComponent[PositionComponent]()
var fpsComponentID = lib.RegisterComponent[FPSComponent]()
var frameTimeComponentID = lib.RegisterComponent[FrameTimeComponent]()
var debounceUpdateComponentID = lib.RegisterComponent[DebounceUpdateComponent]()
var shouldUpdateComponentID = lib.RegisterComponent[ShouldUpdateComponent]()
var colorComponentID = lib.RegisterComponent[ColorComponent]()
var particleSpawnComponentID = lib.RegisterComponent[ParticleSpawnComponent]()
var particleComponentID = lib.RegisterComponent[ParticleComponent]()
var lifetimeComponentID = lib.RegisterComponent[LifetimeComponent]()
var entityCounterComponentID = lib.RegisterComponent[EntityCounterComponent]()
var visibleComponentID = lib.RegisterComponent[VisibleComponent]()
var speedComponentID = lib.RegisterComponent[SpeedComponent]()

func main() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint | rl.FlagWindowHighdpi | rl.FlagWindowResizable)
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	postprocessingShader := rl.LoadShader("", "shaders/postprocessing.glsl")
	fishEyeShader := rl.LoadShader("", "shaders/fisheye.glsl")
	strengthLoc := rl.GetShaderLocation(fishEyeShader, "strength")
	timeLoc := rl.GetShaderLocation(postprocessingShader, "time")

	world := lib.NewWorld()

	fpsCounter := world.CreateEntity()
	world.AddComponents(fpsCounter,
		TextComponent{Text: "FPS: %d", FontSize: toScaled(20)},
		PositionComponent{X: toScaled(10), Y: toScaled(10)},
		FPSComponent{},
		DebounceUpdateComponent{DebounceTime: 0.1, CurrentTime: 0},
		ShouldUpdateComponent{},
		VisibleComponent{},
	)
	frameTime := world.CreateEntity()
	world.AddComponents(frameTime,
		TextComponent{Text: "Frame time: %d", FontSize: toScaled(14)},
		PositionComponent{X: toScaled(10), Y: toScaled(30)},
		FrameTimeComponent{},
		DebounceUpdateComponent{DebounceTime: 1, CurrentTime: 0},
		ShouldUpdateComponent{},
		VisibleComponent{},
	)

	particleSpawn := world.CreateEntity()
	world.AddComponents(particleSpawn,
		ParticleSpawnComponent{LifeTime: 20},
		PositionComponent{X: 0, Y: 0},
		ColorComponent{R: 1, G: 0, B: 0, A: 1},
	)

	entitiesCounter := world.CreateEntity()
	world.AddComponents(entitiesCounter,
		TextComponent{Text: "Entities: %d", FontSize: toScaled(14)},
		PositionComponent{X: toScaled(10), Y: toScaled(50)},
		EntityCounterComponent{},
		ShouldUpdateComponent{},
		DebounceUpdateComponent{DebounceTime: 0.1, CurrentTime: 0},
		EntityCounterComponent{},
		VisibleComponent{},
	)

	//rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		world.Query().With(lib.GetComponentID[ShouldUpdateComponent]()).Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
			world.DisableComponent(id, shouldUpdateComponentID)
		})

		world.Query().
			With(lib.GetComponentID[DebounceUpdateComponent]()).
			Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
				debounceUpdate := components[debounceUpdateComponentID].(DebounceUpdateComponent)
				debounceUpdate.CurrentTime += rl.GetFrameTime()
				if debounceUpdate.CurrentTime >= debounceUpdate.DebounceTime {
					debounceUpdate.CurrentTime = 0
					world.EnableComponent(id, shouldUpdateComponentID)
				}
				world.AddComponents(id, debounceUpdate)
			})

		world.Query().
			With(fpsComponentID, positionComponentID).
			Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
				screenPosition := components[positionComponentID].(PositionComponent)
				mousePosition := scaledMousePosition()
				screenPosition.X = int32(mousePosition.X)
				screenPosition.Y = int32(mousePosition.Y)
				world.AddComponents(id, screenPosition)
			})

		world.Query().With(lifetimeComponentID).Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
			lifetime := components[lifetimeComponentID].(LifetimeComponent)
			lifetime.CurrentTime += rl.GetFrameTime()
			if lifetime.CurrentTime >= lifetime.LifeTime {
				world.DestroyEntity(id)
			} else {
				world.AddComponents(id, lifetime)
			}
		})

		world.Query().
			With(lifetimeComponentID, particleComponentID, colorComponentID).
			Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
				lifetime := components[lifetimeComponentID].(LifetimeComponent)
				colorComponent := components[colorComponentID].(ColorComponent)
				colorComponent.G = uint8(255.0 * (lifetime.CurrentTime / lifetime.LifeTime))
				colorComponent.B = uint8(1 - lifetime.CurrentTime/lifetime.LifeTime)
				colorComponent.A = uint8(max(255.0*(0.7-lifetime.CurrentTime/lifetime.LifeTime), 0.0))
				world.AddComponents(id, colorComponent)
			})

		world.Query().With(particleSpawnComponentID).Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
			spawner := components[particleSpawnComponentID].(ParticleSpawnComponent)
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				for i := 0; i < 300; i++ {
					x := int32(rand.Float32()*float32(rl.GetRenderWidth())*2 - float32(rl.GetRenderWidth())*0.5)
					y := int32(rand.Float32()*float32(rl.GetRenderHeight())*1.5 - float32(rl.GetRenderHeight())*0.5)
					if rl.Vector2Distance(scaledMousePosition(), rl.Vector2{X: float32(x), Y: float32(y)}) < 200 {
						continue
					}
					particleEntity := world.CreateEntity()
					world.AddComponents(particleEntity,
						ParticleComponent{},
						PositionComponent{X: x, Y: y},
						ColorComponent{R: 255, G: 0, B: 0, A: 255},
						LifetimeComponent{LifeTime: spawner.LifeTime, CurrentTime: 0},
						VisibleComponent{},
						SpeedComponent{Speed: rand.Float32() * 3},
					)
				}
			}
		})

		world.Query().
			With(fpsComponentID, textComponentID, shouldUpdateComponentID).
			Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
				text := components[textComponentID].(TextComponent)
				text.Text = fmt.Sprintf("FPS: %d", rl.GetFPS())
				world.AddComponents(id, text)
			})

		world.Query().
			With(frameTimeComponentID, textComponentID, shouldUpdateComponentID).
			Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
				text := components[textComponentID].(TextComponent)
				text.Text = fmt.Sprintf("Frame time: %f", rl.GetFrameTime())
				world.AddComponents(id, text)
			})

		world.Query().
			With(entityCounterComponentID, textComponentID).
			Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
				text := components[textComponentID].(TextComponent)
				text.Text = fmt.Sprintf("Entities: %d", world.GetEntityCount())
				world.AddComponents(id, text)
			})

		world.Query().
			With(particleComponentID, positionComponentID, speedComponentID).
			Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
				position := components[positionComponentID].(PositionComponent)
				speed := components[speedComponentID].(SpeedComponent)
				particlePosition := rl.Vector2{X: float32(position.X), Y: float32(position.Y)}
				mousePosition := scaledMousePosition()
				movedPosition := rl.Vector2Lerp(particlePosition, mousePosition, speed.Speed*rl.GetFrameTime())
				radius := float32(toScaled(10)) * rl.Vector2Distance(mousePosition, movedPosition) / float32(rl.GetRenderWidth())
				position.X = int32(movedPosition.X)
				position.Y = int32(movedPosition.Y)
				world.AddComponents(id, position)
				world.AddComponents(id, ParticleComponent{Radius: radius})
			})

		// Drawing
		firstRenderPass := rl.LoadRenderTexture(int32(rl.GetRenderWidth()), int32(rl.GetRenderHeight()))
		rl.BeginDrawing()
		rl.BeginTextureMode(firstRenderPass)
		rl.ClearBackground(rl.RayWhite)

		world.Query().
			With(particleComponentID, positionComponentID, colorComponentID, visibleComponentID).
			Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
				particle := components[particleComponentID].(ParticleComponent)
				position := components[positionComponentID].(PositionComponent)
				color := components[colorComponentID].(ColorComponent)

				pos := rl.Vector2{X: float32(position.X), Y: float32(position.Y)}
				col := rl.NewColor(color.R, color.G, color.B, color.A)

				rl.DrawCircleV(pos, particle.Radius, col)

				c := rl.Red
				c.A = col.A / 10
				rl.DrawLineV(pos, scaledMousePosition(), c)
			})

		rl.EndTextureMode()

		secondRenderPass := rl.LoadRenderTexture(int32(rl.GetRenderWidth()), int32(rl.GetRenderHeight()))
		rl.BeginTextureMode(secondRenderPass)
		rl.BeginShaderMode(fishEyeShader)
		mp := scaledMousePosition()
		middle := rl.Vector2{X: float32(rl.GetRenderWidth()) / 2, Y: float32(rl.GetRenderHeight() / 2)}
		distance := rl.Vector2Distance(middle, mp)
		normalized := distance / float32(rl.GetRenderWidth())
		rl.SetShaderValue(fishEyeShader, strengthLoc, []float32{normalized}, rl.ShaderUniformFloat)
		rl.DrawTexturePro(firstRenderPass.Texture,
			rl.NewRectangle(0, 0, float32(firstRenderPass.Texture.Width), float32(-firstRenderPass.Texture.Height)),
			rl.NewRectangle(0, 0, float32(firstRenderPass.Texture.Width), float32(firstRenderPass.Texture.Height)),
			rl.NewVector2(0, 0), 0, rl.White,
		)
		rl.EndShaderMode()
		world.Query().
			With(textComponentID, positionComponentID, visibleComponentID).
			Each(func(id lib.EntityID, components map[lib.ComponentID]interface{}) {
				positions := components[positionComponentID].(PositionComponent)
				text := components[textComponentID].(TextComponent)
				rl.DrawText(text.Text, positions.X, positions.Y, text.FontSize, rl.Black)
			})
		rl.EndTextureMode()

		rl.SetShaderValue(postprocessingShader, timeLoc, []float32{rl.GetFrameTime()}, rl.ShaderUniformFloat)
		rl.BeginShaderMode(postprocessingShader)
		rl.DrawTexturePro(secondRenderPass.Texture,
			rl.NewRectangle(0, 0, float32(firstRenderPass.Texture.Width), float32(-firstRenderPass.Texture.Height)),
			rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
			rl.NewVector2(0, 0), 0, rl.White,
		)
		rl.EndShaderMode()

		rl.EndDrawing()
		rl.UnloadRenderTexture(firstRenderPass)
		rl.UnloadRenderTexture(secondRenderPass)
	}
}

func toScaled(value int32) int32 {
	return int32(math.Floor(float64(value) * float64(rl.GetWindowScaleDPI().X)))
}

func toScaledV(vector rl.Vector2) rl.Vector2 {
	return rl.Vector2Scale(vector, rl.GetWindowScaleDPI().X)
}

func scaledMousePosition() rl.Vector2 {
	return toScaledV(rl.GetMousePosition())
}
