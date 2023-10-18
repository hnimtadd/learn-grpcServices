package sample

import (
	"grpcCource/pb"
	"math/rand"
)

func RandomLaptopBrand() string {
	return randomStringFromSet("DELL", "ASUS", "APPLE")
}

func RandomLaptopName(brand string) string {
	switch brand {
	case "DELL":
		return randomStringFromSet(
			"XPS 9560",
			"Latitude 1600",
			"Vostro 9559",
		)
	case "ASUS":
		return randomStringFromSet(
			"Vivobook 1600",
			"Vivobook 1400",
		)
	case "APPLE":
		return randomStringFromSet(
			"Macbook Pro M1 2022",
			"Macbook Air M1 2021",
		)
	default:
		return "unknown"

	}

}
func RandomResolution() *pb.Screen_Resolution {
	height := randomInt(624, 4080)
	width := height * 16 / 9
	return &pb.Screen_Resolution{
		Width:  uint32(width),
		Height: uint32(height),
	}
}

func NewScreen() *pb.Screen {
	panel := pb.Screen_Panel(randomInt(0, 1))
	screen := &pb.Screen{
		SizeInch:   randomFloat32(10, 16),
		Resolution: RandomResolution(),
		Multitouch: randomBool(),
		Panel:      panel,
	}
	return screen
}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet(
			"Xeon E-3556",
			"Core i9-9980HK",
			"Core i7 9700HQ",
			"Core i5-9560f",
		)
	}
	return randomStringFromSet(
		"Ryzen 7 PRO 2780U",
		"Ryzen 5 PRO 3580U",
		"Ryzen 3 PRO 3280GE",
	)
}

func randomStringFromSet(a ...string) string {
	n := len(a)
	return a[randomInt(0, n-1)]
}

func randomKeyboardLayout() pb.Keyboard_Layout {
	return pb.Keyboard_Layout(randomInt(0, 2))
}

func randomInt(start int, end int) int {
	return rand.Intn(end-start+1) + start

}
func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomFloat64(min float64, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func randomFloat32(min float32, max float32) float32 {
	return rand.Float32()*(max-min) + min
}
