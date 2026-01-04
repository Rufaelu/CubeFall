package gamelogic

import (
	// "errors"
	"math"
	"math/rand"
	"time"

	"CubeFall/helper"
	"CubeFall/objects"

	"github.com/go-gl/mathgl/mgl32"
)

func BulletHitsEnemy(
	bullet *objects.Bullet,
	enemies *[]objects.ExtraEnemyProperty,
	score *ScoreTrack,
) {
	for i := range *enemies {
		enemy := &(*enemies)[i]

		if enemy.IsHit {
			continue
		}

		distance := bullet.Position.Sub(enemy.Position).Len()

		if distance < bullet.Radius+0.5 {
			enemy.IsHit = true
			bullet.Alive = false

			score.KillCount++
			score.Points++

			return // one bullet = one kill
		}
	}
}

func MoveEnemies(
	enemies *[]objects.ExtraEnemyProperty,
	speed *float32,
	deltaTime float32,
) {
	for i := range *enemies {
		enemy := &(*enemies)[i]

		if enemy.IsHit {
			enemy.Position[1] += 0.025
			enemy.Position[2] -= 0.025
			continue
		}

		// Move forward (Z axis)
		enemy.Position[2] += *speed

		// Bounce (Y axis)
		enemy.Time += deltaTime
		enemy.Position[1] =
			enemy.BaseY +
				float32(math.Sin(float64(enemy.Time*enemy.BounceSpd)))*
					enemy.BounceAmp
	}
}

func KillEnemy(
	hit_enemy *objects.ExtraEnemyProperty,
	enemies *[]objects.ExtraEnemyProperty,
	score *ScoreTrack,
) {
	score.KillCount += 1
	score.Points += 1
	for i, enemy := range *enemies {
		if enemy.Id == hit_enemy.Id {
			(*enemies)[i].IsHit = true
			break
		}
	}
}

func LevelUp(score *ScoreTrack, enemySpeed *float32, enemies *[]objects.ExtraEnemyProperty) {
	score.PassedMaxLevel += 1
	*enemySpeed *= 1.2
	spawnMore(enemies)
}

var init_count int = 4

func spawnMore(enemies *[]objects.ExtraEnemyProperty) {
	*enemies = (*enemies)[:0]

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	minX, maxX := float32(-12.0), float32(12.0)
	minZ, maxZ := float32(-30.0), float32(-10.0)
	y := float32(0.6)

	for i := 0; i < init_count*2; i++ {
		x := minX + rng.Float32()*(maxX-minX)
		z := minZ + rng.Float32()*(maxZ-minZ)

		(*enemies) = append(*enemies,
			objects.ExtraEnemyProperty{
				Id:        uint32(i + 1),
				IsHit:     false,
				Position:  mgl32.Vec3{x, y, z},
				BaseY:     y,
				BounceAmp: 0.4 + rng.Float32()*0.3, // random height
				BounceSpd: 2.0 + rng.Float32()*2.0, // random speed
				Time:      rng.Float32() * 10.0,    // desync enemies
			},
		)
	}

	init_count *= 2
}

func AllEnemiesAreHit(enemies *[]objects.ExtraEnemyProperty) bool {
	for _, enemy := range *enemies {
		if !enemy.IsHit || enemy.Position.Y() <= 5 {
			return false
		}
	}
	return true
}

func HandlePassedEnemies(enemies *[]objects.ExtraEnemyProperty, player_position *helper.Camera, scoreTrack *ScoreTrack) {
	for i, enemy := range *enemies {
		if enemy.Position.Z() > player_position.Position.Z() {
			scoreTrack.Points -= 1
			// passedID := enemy.Id
			*enemies = append((*enemies)[:i], (*enemies)[i+1:]...)
			// fmt.Printf("Enemy with ID %d passed you!\n", passedID)
		}
	}
}
