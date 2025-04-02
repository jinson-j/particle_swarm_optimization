package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bird struct {
	position rl.Vector2
	velocity rl.Vector2
	p_best   rl.Vector2
}

type Goal struct {
	position rl.Vector2
	g_best   rl.Vector2
}

const window_width = 500
const window_height = 500

const num_birds = 1000

const inertia float32 = 0.7
const cognitive float32 = 1.5
const social float32 = 1.5

const maxSpeed float32 = 2

func calculate_distance(a, b rl.Vector2) float32 {
	return float32(math.Sqrt(float64((a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y))))
}

func main() {
	rl.InitWindow(int32(window_width), int32(window_height), "Bird Simulator!")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	birds := make([]Bird, num_birds)
	goal := Goal{
		position: rl.Vector2{X: window_width / 2, Y: window_height / 2},
		g_best:   rl.Vector2{X: math.MaxFloat32, Y: math.MaxFloat32},
	}

	for i := 0; i < num_birds; i++ {
		x := (rand.Float32() * window_width)
		y := (rand.Float32() * window_height)

		x_prime := (-1.0 + (rand.Float32() * 2.0))
		y_prime := (-1.0 + (rand.Float32() * 2.0))

		birds[i] = Bird{
			position: rl.Vector2{X: x, Y: y},
			velocity: rl.Vector2{X: x_prime, Y: y_prime},
			p_best:   rl.Vector2{X: x, Y: y},
		}

		if calculate_distance(birds[i].position, goal.position) < calculate_distance(goal.g_best, goal.position) {
			goal.g_best.X = x
			goal.g_best.Y = y
		}
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			goal.position = rl.GetMousePosition()

			for i := 0; i < num_birds; i++ {
				x := birds[i].position.X
				y := birds[i].position.Y

				x_prime := (-1.0 + (rand.Float32() * 2.0))
				y_prime := (-1.0 + (rand.Float32() * 2.0))

				birds[i] = Bird{
					position: rl.Vector2{X: x, Y: y},
					velocity: rl.Vector2{X: x_prime, Y: y_prime},
					p_best:   rl.Vector2{X: x, Y: y},
				}

				if calculate_distance(birds[i].position, goal.position) < calculate_distance(goal.g_best, goal.position) {
					goal.g_best.X = x
					goal.g_best.Y = y
				}
			}
		}

		for i := range birds {
			r1, r2 := rand.Float32(), rand.Float32()

			vx := birds[i].velocity.X
			vy := birds[i].velocity.Y

			px := birds[i].position.X
			py := birds[i].position.Y

			pbx := birds[i].p_best.X
			pby := birds[i].p_best.Y

			gbx := goal.g_best.X
			gby := goal.g_best.Y

			birds[i].velocity.X =
				(inertia*vx + cognitive*r1*(pbx-px) +
					social*r2*(gbx-px))

			birds[i].velocity.Y =
				(inertia*vy + cognitive*r1*(pby-py) +
					social*r2*(gby-py))

			if rl.Vector2Length(birds[i].velocity) > maxSpeed {
				birds[i].velocity = rl.Vector2Scale(rl.Vector2Normalize(birds[i].velocity), maxSpeed)
			}

			birds[i].position.X += birds[i].velocity.X
			birds[i].position.Y += birds[i].velocity.Y
			rl.DrawCircleV(birds[i].position, 2, rl.Red)

			if calculate_distance(birds[i].position, goal.position) < calculate_distance(birds[i].p_best, goal.position) {
				birds[i].p_best.X = birds[i].position.X
				birds[i].p_best.Y = birds[i].position.Y
			}

			if calculate_distance(birds[i].position, goal.position) < calculate_distance(goal.g_best, goal.position) {
				goal.g_best.X = birds[i].position.X
				goal.g_best.Y = birds[i].position.Y
			}

		}

		rl.DrawText("Bird Simulator", 5, 5, 20, rl.Gray)
		rl.DrawFPS(5, window_height-20)

		rl.DrawCircleV(goal.position, 5, rl.Blue)
		rl.DrawCircleV(goal.g_best, 3, rl.Green)
		rl.EndDrawing()
	}
}
