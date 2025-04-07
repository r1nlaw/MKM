package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	imgWidth    = 800
	imgHeight   = 800
	lambda      = 500e-9  // Длина волны (500 нм)
	diskRadius  = 100e-6  // Радиус диска (100 мкм)
	distance    = 7.14e-3 // Расстояние до экрана мм
	samples     = 100000  // Количество точек на краю диска
	screenWidth = 0.5e-3  // Ширина экрана мм
)

type Point struct{ X, Y float64 }

func main() {
	start := time.Now()

	fmt.Println("Генерация точек...")
	startPoints := time.Now()
	edgePoints := generateDiskEdgePoints(samples, diskRadius)
	fmt.Printf("Генерация точек заняла: %v\n", time.Since(startPoints))

	fmt.Println("Создание изображения...")
	startImage := time.Now()
	createPoissonEffectImage(edgePoints, "poisson_effect.png")
	fmt.Printf("Создание изображения заняло: %v\n", time.Since(startImage))

	fmt.Println("Создание графика интенсивности...")
	startPlot := time.Now()
	createIntensityPlot(edgePoints, "intensity_plot.png")
	fmt.Printf("Создание графика заняло: %v\n", time.Since(startPlot))

	calculateFresnelZones(diskRadius, lambda, distance)

	centerRe, centerIm := calculateAmplitude(edgePoints, 0, 0)
	centerIntensity := centerRe*centerRe + centerIm*centerIm
	fmt.Printf("Интенсивность в центре экрана: %.6f\n", centerIntensity)

	fmt.Printf("Полное время выполнения программы: %v\n", time.Since(start))
}

func generateDiskEdgePoints(n int, r float64) []Point {
	points := make([]Point, n)
	numWorkers := runtime.NumCPU()
	chunk := (n + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		start := w * chunk
		end := start + chunk
		if end > n {
			end = n
		}
		go func(start, end int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(start)))
			for i := start; i < end; i++ {
				theta := rng.Float64() * 2 * math.Pi
				points[i] = Point{
					X: r * math.Cos(theta),
					Y: r * math.Sin(theta),
				}
			}
		}(start, end)
	}

	wg.Wait()
	return points
}

func calculateAmplitude(points []Point, x, y float64) (float64, float64) {
	var re, im float64
	k := 2 * math.Pi / lambda

	// Расчет зоны Френеля
	r0 := diskRadius                    // Радиус диска
	b := distance                       // Расстояние до экрана
	m := (r0 * r0 / lambda) * (1.0 / b) // Количество зон Френеля

	// Если количество зон Френеля больше, делаем интенсивность в центре более темной
	fresnelFactor := 1.0
	if m > 1 {
		fresnelFactor = 1 / math.Sqrt(m)
	}

	// Вычисление амплитуды
	for _, p := range points {
		dx := x - p.X
		dy := y - p.Y
		phase := (k / (2 * distance)) * (dx*dx + dy*dy)
		re += math.Cos(phase)
		im += math.Sin(phase) // Суммируем синусы и косинусы фаз от всех точек и получаем суммарную амплитуду в точках
	}

	n := float64(len(points))
	amplitude := fresnelFactor * (re / n) // Средняя амплитуда на точке
	return amplitude, fresnelFactor * (im / n)
}

func createPoissonEffectImage(points []Point, filename string) {
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	scale := screenWidth / float64(imgWidth)
	diskCenterX, diskCenterY := imgWidth/2, imgHeight/2
	diskRadiusPx := int(diskRadius / scale)

	// Вычисление количества зон Френеля
	m := (diskRadius * diskRadius / lambda) * (1.0 / distance)
	fresnelFactor := 1.0
	if m > 1 {
		fresnelFactor = 1 / math.Sqrt(m)
	}

	// Заполнение изображения черным цветом для пустой области
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			dx := x - diskCenterX
			dy := y - diskCenterY
			if dx*dx+dy*dy < diskRadiusPx*diskRadiusPx {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}
		}
	}

	// Массив интенсивности
	intensity := make([][]float64, imgHeight)
	for i := range intensity {
		intensity[i] = make([]float64, imgWidth)
	}

	// Многозадачность с использованием атомарного максимума интенсивности
	var maxIntensity uint64
	var wg sync.WaitGroup
	blockSize := imgHeight / 8
	if blockSize == 0 {
		blockSize = 1
	}

	for blockY := 0; blockY < 8 && blockY*blockSize < imgHeight; blockY++ {
		wg.Add(1)
		go func(blockY int) {
			defer wg.Done()
			startY := blockY * blockSize
			endY := startY + blockSize
			if endY > imgHeight {
				endY = imgHeight
			}
			for y := startY; y < endY; y++ {
				for x := 0; x < imgWidth; x++ {
					dx := x - diskCenterX
					dy := y - diskCenterY
					if dx*dx+dy*dy < diskRadiusPx*diskRadiusPx {
						continue
					}
					xPos := (float64(x) - float64(imgWidth)/2) * scale
					yPos := (float64(y) - float64(imgHeight)/2) * scale
					re, im := calculateAmplitude(points, xPos, yPos) // действительные и мнимые части амплитуды
					intensity[y][x] = re*re + im*im
					current := math.Float64bits(intensity[y][x])
					for { // находим максимальную интенсивность для нормализации всех пикселей [0....1]
						old := atomic.LoadUint64(&maxIntensity)
						if current <= old || atomic.CompareAndSwapUint64(&maxIntensity, old, current) {
							break
						}
					}
				}
			}
		}(blockY)
	}

	wg.Wait()

	// Нормализация интенсивности
	maxI := math.Float64frombits(maxIntensity)
	if maxI == 0 {
		maxI = 1
	}

	// Отображение Пуазона с учетом интенсивности и затемнения центра
	poissonRadius := 3
	for y := diskCenterY - poissonRadius; y <= diskCenterY+poissonRadius; y++ {
		for x := diskCenterX - poissonRadius; x <= diskCenterX+poissonRadius; x++ {
			if x >= 0 && x < imgWidth && y >= 0 && y < imgHeight {
				xPos := (float64(x) - float64(imgWidth)/2) * scale
				yPos := (float64(y) - float64(imgHeight)/2) * scale
				re, im := calculateAmplitude(points, xPos, yPos)
				intens := re*re + im*im

				// Уменьшаем интенсивность по центру, если много зон
				intens *= fresnelFactor

				normIntensity := intens / maxI
				img.Set(x, y, colorFromRingIntensity(normIntensity, x, y, diskCenterX, diskCenterY))
			}
		}
	}

	// Установка цвета пикселей с нормализованной интенсивностью
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			dx := x - diskCenterX
			dy := y - diskCenterY
			if dx*dx+dy*dy >= diskRadiusPx*diskRadiusPx {
				normIntensity := intensity[y][x] / maxI
				img.Set(x, y, colorFromRingIntensity(normIntensity, x, y, diskCenterX, diskCenterY))
			}
		}
	}

	saveImage(img, filename)
}

func colorFromRingIntensity(intensity float64, x, y, centerX, centerY int) color.RGBA {
	r := uint8(255 * math.Pow(intensity, 0.6))
	g := uint8(255 * math.Pow(intensity, 0.9))
	b := uint8(255 * math.Pow(intensity, 0.4))
	dx := float64(x - centerX)
	dy := float64(y - centerY)
	dist := math.Sqrt(dx*dx + dy*dy)

	// Увеличение синего цвета внутри радиуса
	if dist < 100 {
		b = uint8(math.Min(255, float64(b)+150*math.Pow(intensity, 0.8)))
	}

	return color.RGBA{r, g, b, 255}
}

func createIntensityPlot(points []Point, filename string) {
	p := plot.New()
	p.Title.Text = "Распределение интенсивности"
	p.X.Label.Text = "Расстояние от центра, мм"
	p.Y.Label.Text = "Абсолютная интенсивность"

	scale := screenWidth / float64(imgWidth)
	pts := make(plotter.XYs, imgWidth)

	var maxIntensity float64
	for x := 0; x < imgWidth; x++ {
		xPos := (float64(x) - float64(imgWidth)/2) * scale
		re, im := calculateAmplitude(points, xPos, 0)
		intensity := re*re + im*im
		pts[x].X = xPos * 1000
		pts[x].Y = intensity

		if intensity > maxIntensity {
			maxIntensity = intensity
		}
	}

	if maxIntensity == 0 {
		maxIntensity = 1
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatal(err)
	}
	p.Add(line)
	p.Add(plotter.NewGrid())

	p.X.Min = -screenWidth * 1000 / 2
	p.X.Max = screenWidth * 1000 / 2
	p.Y.Min = 0
	p.Y.Max = 1.1

	if err := p.Save(10*vg.Centimeter, 6*vg.Centimeter, filename); err != nil {
		log.Fatal(err)
	}
}

func saveImage(img *image.RGBA, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

func calculateFresnelZones(r0, lambda, b float64) {
	m := (r0 * r0 / lambda) * (1.0 / b)
	fmt.Printf("Количество открытых зон Френеля: m = %.2f\n", m)
}
