package sensors

import (
	"math"
	"time"
)

type Sample struct {
	Timestamp time.Time  `json:"timestamp"`
	AccelMS2  [3]float64 `json:"accel_m_s2"`
	GyroRadS  [3]float64 `json:"gyro_rad_s"`
}

type Driver struct {
	step int
}

func NewDriver() *Driver {
	return &Driver{}
}

func (driver *Driver) Next() Sample {
	t := float64(driver.step) * 0.25
	driver.step++
	return Sample{
		Timestamp: time.Now().UTC(),
		AccelMS2: [3]float64{
			0.18*math.Cos(t/2.8) + 0.02,
			0.14*math.Sin(t/3.5) + 0.015,
			0.0,
		},
		GyroRadS: [3]float64{
			0.008 * math.Sin(t/4.0),
			0.010 * math.Cos(t/5.0),
			0.04 + 0.003*math.Sin(t/6.0),
		},
	}
}
