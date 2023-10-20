package sample

import (
	"grpcCource/pkg/pb"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewLaptop() *pb.Laptop {
	id := uuid.NewString()
	brand := RandomLaptopBrand()
	name := RandomLaptopName(brand)
	cpu := NewCPU()
	gpu := newGPU()
	ram := NewRAM()
	storage := NewHDD()
	screen := NewScreen()
	keyboard := NewKeyboard()
	return &pb.Laptop{
		Id:       id,
		Brand:    brand,
		Name:     name,
		Cpu:      cpu,
		Ram:      ram,
		Gpus:     []*pb.GPU{gpu},
		Storages: []*pb.Storage{storage},
		Screen:   screen,
		Keyboard: keyboard,
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0, 3.0),
		},
		PriceUsd:    randomFloat64(1500, 3000),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		UpdateAt:    timestamppb.Now(),
	}
}

func NewKeyboard() *pb.Keyboard {
	keyboard := &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
	return keyboard
}

func NewCPU() *pb.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)
	numberCores := randomInt(2, 8)
	numberThreds := randomInt(2, 16)
	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.8)
	return &pb.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(numberCores),
		NumberThreads: uint32(numberThreds),
		MinGhz:        minGhz,
		MaxGhz:        maxGhz,
	}
}

func newGPU() *pb.GPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)
	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.8)
	memory := &pb.Memory{
		Value: uint64(randomInt(2, 6)),
		Unit:  pb.Memory_GIGABYTE,
	}
	return &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: memory,
	}
}

func NewRAM() *pb.Memory {
	ram := &pb.Memory{
		Unit:  pb.Memory_GIGABYTE,
		Value: uint64(randomInt(4, 64)),
	}
	return ram
}

func NewSSD() *pb.Storage {
	ssd := &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(128, 1024)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}
	return ssd
}

func NewHDD() *pb.Storage {
	hdd := &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit:  pb.Memory_TETRABYTE,
		},
	}
	return hdd
}

func NewScore() float64 {
	return randomFloat64(0, 10)

}
