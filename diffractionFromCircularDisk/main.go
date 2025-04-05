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
	samples     = 10000   // Количество точек на краю диска
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
	fmt.Printf("Создание графика заняла: %v\n", time.Since(startPlot))

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
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(start))) // отдельный генератор на поток
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

	for _, p := range points { // вычисляем амплитуду для каждой точки (расстояние до точки экрана и фазе волны)
		dx := x - p.X
		dy := y - p.Y
		r := math.Sqrt(dx*dx + dy*dy + distance*distance)
		phase := k * r
		amplitude := 1.0 / r

		re += amplitude * math.Cos(phase)
		im += amplitude * math.Sin(phase)
	}

	n := float64(len(points))
	return re / n, im / n // действительная и мнимая части амплитуды для точки экрана
}

func createPoissonEffectImage(points []Point, filename string) {
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	scale := screenWidth / float64(imgWidth)

	// Закрашиваем диск черным цветом
	diskCenterX, diskCenterY := imgWidth/2, imgHeight/2
	diskRadiusPx := int(diskRadius / scale)

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

	// Инициализация структуры для интенсивности
	intensity := make([][]float64, imgHeight)
	for i := range intensity {
		intensity[i] = make([]float64, imgWidth)
	}

	// Используем atomic для поиска максимальной интенсивности
	var maxIntensity uint64

	// Параллельная обработка блоков изображения
	var wg sync.WaitGroup
	blockSize := imgHeight / 1024 // Разделим на 128 блоков

	for blockY := 0; blockY < 1024; blockY++ {
		wg.Add(1)
		go func(blockY int) {
			defer wg.Done()
			startY := blockY * blockSize
			endY := startY + blockSize
			if blockY == 7 {
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
					re, im := calculateAmplitude(points, xPos, yPos)
					intensity[y][x] = re*re + im*im

					// Использование atomic для обновления максимальной интенсивности
					current := math.Float64bits(intensity[y][x])
					for {
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

	// Находим максимальную интенсивность
	maxI := math.Float64frombits(maxIntensity)
	if maxI == 0 {
		maxI = 1
	}

	// Добавляем яркую точку в центре (эффект Пуассона)
	poissonRadius := 3
	for y := diskCenterY - poissonRadius; y <= diskCenterY+poissonRadius; y++ {
		for x := diskCenterX - poissonRadius; x <= diskCenterX+poissonRadius; x++ {
			if x >= 0 && x < imgWidth && y >= 0 && y < imgHeight {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	// Рисуем дифракционные кольца
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
	// Цвета для дифракционных колец
	r := uint8(255 * math.Pow(intensity, 0.4))
	g := uint8(255 * math.Pow(intensity, 0.8))
	b := uint8(255 * math.Pow(intensity, 0.8))

	// Вычисляем расстояние от центра с правильным приведением типов
	dx := float64(x - centerX)
	dy := float64(y - centerY)
	dist := math.Sqrt(dx*dx + dy*dy)

	// Добавляем синий оттенок для первых колец
	if dist < 100 {
		b = uint8(math.Min(255, float64(b)+150*math.Pow(intensity, 0.8)))
	}

	return color.RGBA{r, g, b, 255}
}
func createIntensityPlot(points []Point, filename string) {
	p := plot.New()
	p.Title.Text = "Распределение интенсивности"
	p.X.Label.Text = "Расстояние от центра, мм"
	p.Y.Label.Text = "Относительная интенсивность"

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
	for i := range pts {
		pts[i].Y /= maxIntensity
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
