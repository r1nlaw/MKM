package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const (
	width      = 1840
	height     = 2160
	sunRadius  = 30 // Радиус Солнца в пикселях
	bhRadius   = 15 // радиус горизонта событий
	pixelScale = 1e12
	G          = 6.67430e-11 // Гравитационная постоянная
	c          = 299792458.0 // Скорость света
	solarMass  = 1.98847e30  // Масса Солнца
)

func main() {
	var (
		blackHoleDist int
		rayCount      int
		massSolar     float64
	)

	fmt.Print("Введите массу чёрной дыры (в массах Солнца) например черная дыра Стрелец А* (4.3e6 масс Солнца): ")
	fmt.Scanln(&massSolar)

	fmt.Print("Введите расстояние до чёрной дыры (в пикселях) (1 пиксель = 1e12 метров): ")
	fmt.Scanln(&blackHoleDist)

	fmt.Print("Введите количество световых лучей: ")
	fmt.Scanln(&rayCount)

	mass := massSolar * solarMass
	sunX, sunY := 100, height/2
	bhX, bhY := sunX+blackHoleDist, height/2

	fmt.Println("Создание изображения...")
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	fillBackground(img, color.Black)

	drawSun(img, sunX, sunY, sunRadius)
	drawBlackHole(img, bhX, bhY, bhRadius)
	drawLightRays(img, sunX, sunY, sunRadius, bhX, bhY, mass, rayCount)

	fmt.Println("Сохранение изображения...")
	if err := saveImage(img, "black_hole_lensing.png"); err != nil {
		fmt.Println("Ошибка при сохранении:", err)
	} else {
		fmt.Println("Изображение успешно сохранено как black_hole_lensing.png")
	}
}

func fillBackground(img *image.RGBA, col color.Color) {
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, col)
		}
	}
}

func drawSun(img *image.RGBA, cx, cy, r int) {
	for y := -r * 3; y <= r*3; y++ {
		for x := -r * 3; x <= r*3; x++ {
			dist := math.Hypot(float64(x), float64(y))
			if dist <= float64(r) {
				intensity := 1 - dist/float64(r)
				col := color.RGBA{
					R: 255,
					G: uint8(255 * intensity * 0.6),
					B: uint8(80 * intensity * 0.6),
					A: 255,
				}
				img.Set(cx+x, cy+y, col)
			} else if dist <= float64(r*3) {
				alpha := uint8(100 * (1 - math.Pow(dist/float64(r*3), 3)))
				base := img.RGBAAt(cx+x, cy+y)
				img.Set(cx+x, cy+y, blendColors(base, color.RGBA{255, 200, 50, alpha}))
			}
		}
	}
}

func drawBlackHole(img *image.RGBA, cx, cy, r int) {
	diskOuter := r * 3
	for y := -diskOuter; y <= diskOuter; y++ {
		for x := -diskOuter; x <= diskOuter; x++ {
			dist := math.Hypot(float64(x), float64(y))
			if dist <= float64(diskOuter) && dist > float64(r) {
				intensity := 1 - dist/float64(diskOuter)
				col := color.RGBA{
					R: uint8(255 * intensity),
					G: uint8(150 * intensity),
					B: 0,
					A: uint8(255 * intensity),
				}
				img.Set(cx+x, cy+y, col)
			}
		}
	}
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				img.Set(cx+x, cy+y, color.Black)
			}
		}
	}
}

func drawLightRays(img *image.RGBA, sunX, sunY, sunR, bhX, bhY int, mass float64, rayCount int) {
	for i := 0; i < rayCount; i++ {
		angle := 2 * math.Pi * float64(i) / float64(rayCount)
		dx := math.Cos(angle)
		dy := math.Sin(angle)

		startX := sunX + int(float64(sunR)*dx)
		startY := sunY + int(float64(sunR)*dy)
		endX := startX + int(2000*dx)
		endY := startY + int(2000*dy)

		points := calculateDeflectedPath(startX, startY, endX, endY, bhX, bhY, mass)

		gradFactor := float64(i) / float64(rayCount)
		col := color.RGBA{
			R: 255,
			G: 255,
			B: uint8(gradFactor * 255),
			A: 255,
		}

		for j := 1; j < len(points); j++ {
			drawLine(img, points[j-1].X, points[j-1].Y, points[j].X, points[j].Y, col)
		}
	}
}

type Point struct {
	X, Y int
}

func calculateDeflectedPath(x0, y0, x1, y1, bhX, bhY int, mass float64) []Point { // Гравитационное линзирование (Эффект Эйнштейна)
	var points []Point
	steps := 2000
	bhXf := float64(bhX)
	bhYf := float64(bhY)

	points = append(points, Point{x0, y0})

	cx := float64(x0)
	cy := float64(y0)
	dirX := float64(x1 - x0)
	dirY := float64(y1 - y0)
	dirLen := math.Hypot(dirX, dirY)
	dirX /= dirLen
	dirY /= dirLen

	for i := 0; i < steps; i++ {
		dx := (bhXf - cx) * pixelScale
		dy := (bhYf - cy) * pixelScale
		r := math.Hypot(dx, dy)
		if r < 1e9 {
			r = 1e9
		}

		deflect := (4 * G * mass) / (c * c * r) // Гравитационное отклонение света

		deflectBoost := 1.0 + 1e4/(r/pixelScale) // Увеличение отклонения на близком расстоянии к черной дыре
		deflect *= deflectBoost

		orthoX := -dy / r
		orthoY := dx / r

		dirX += deflect * orthoX // траектории луча для обновления позиции луча с учетом дефлекции
		dirY += deflect * orthoY // траектории луча для обновления позиции луча с учетом дефлекции

		lenDir := math.Hypot(dirX, dirY)
		dirX /= lenDir
		dirY /= lenDir

		cx += dirX
		cy += dirY

		points = append(points, Point{int(cx), int(cy)})
	}

	return points
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, col color.RGBA) {
	dx := int(math.Abs(float64(x1 - x0)))
	dy := int(math.Abs(float64(y1 - y0)))
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	for {
		if x0 >= 0 && x0 < width && y0 >= 0 && y0 < height {
			img.Set(x0, y0, col)
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func blendColors(a, b color.RGBA) color.RGBA {
	alpha := float64(b.A) / 255
	return color.RGBA{
		R: uint8(float64(a.R)*(1-alpha) + float64(b.R)*alpha),
		G: uint8(float64(a.G)*(1-alpha) + float64(b.G)*alpha),
		B: uint8(float64(a.B)*(1-alpha) + float64(b.B)*alpha),
		A: 255,
	}
}

func saveImage(img image.Image, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
