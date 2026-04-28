package filter

import "math"

type SquareRootInformationFilter struct {
	state    StateVector
	sqrtInfo [9]float64
}

func NewSquareRootInformationFilter() *SquareRootInformationFilter {
	filter := &SquareRootInformationFilter{
		sqrtInfo: [9]float64{0.65, 0.65, 0.8, 0.75, 0.75, 0.85, 0.9, 0.9, 0.95},
	}
	filter.updateDriftEllipse()
	return filter
}

func (filter *SquareRootInformationFilter) Predict(accelMS2 [3]float64, gyroRadS [3]float64, dt float64) StateVector {
	for axis := 0; axis < 3; axis++ {
		filter.state.PositionM[axis] += filter.state.VelocityMS[axis]*dt + 0.5*accelMS2[axis]*dt*dt
		filter.state.VelocityMS[axis] += accelMS2[axis] * dt
		filter.state.AttitudeDeg[axis] += gyroRadS[axis] * dt * 180 / math.Pi
	}

	for index, value := range filter.sqrtInfo {
		decay := 0.992 - 0.0005*float64(index)
		filter.sqrtInfo[index] = math.Max(0.18, value*decay)
	}

	filter.updateDriftEllipse()
	return filter.state
}

func (filter *SquareRootInformationFilter) UpdateWithVision(observation VisionObservation) StateVector {
	weight := clamp(0.25+0.6*observation.Quality, 0.15, 0.9)
	filter.state.PositionM[0] += observation.DeltaXM * weight
	filter.state.PositionM[1] += observation.DeltaYM * weight
	filter.state.VelocityMS[0] = blend(filter.state.VelocityMS[0], observation.VxMS, weight)
	filter.state.VelocityMS[1] = blend(filter.state.VelocityMS[1], observation.VyMS, weight)
	filter.state.FeatureCount = observation.FeatureCount
	filter.state.VisionQuality = observation.Quality
	filter.state.OpticalFlowMS = [2]float64{observation.VxMS, observation.VyMS}
	filter.state.TrackedFeatures = observation.Tracks
	filter.state.Timestamp = observation.Timestamp

	for _, index := range []int{0, 1, 3, 4} {
		filter.sqrtInfo[index] = math.Min(2.4, filter.sqrtInfo[index]*(1.08+0.35*observation.Quality))
	}
	filter.updateDriftEllipse()
	return filter.state
}

func (filter *SquareRootInformationFilter) Snapshot() StateVector {
	return filter.state
}

func (filter *SquareRootInformationFilter) updateDriftEllipse() {
	major := 1.0 / math.Max(filter.sqrtInfo[0], 0.18)
	minor := 1.0 / math.Max(filter.sqrtInfo[1], 0.18)
	filter.state.DriftEllipseM = [2]float64{major, minor}
}

func blend(current, observed, weight float64) float64 {
	return current*(1-weight) + observed*weight
}

func clamp(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}
