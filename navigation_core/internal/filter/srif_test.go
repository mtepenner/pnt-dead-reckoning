package filter

import (
	"math"
	"testing"
	"time"
)

func TestPredictIntegratesAcceleration(t *testing.T) {
	filter := NewSquareRootInformationFilter()
	state := filter.Predict([3]float64{1, 0, 0}, [3]float64{}, 1.0)

	if math.Abs(state.PositionM[0]-0.5) > 0.05 {
		t.Fatalf("expected position near 0.5 m, got %.3f", state.PositionM[0])
	}
	if math.Abs(state.VelocityMS[0]-1.0) > 0.05 {
		t.Fatalf("expected velocity near 1.0 m/s, got %.3f", state.VelocityMS[0])
	}
}

func TestVisionUpdateShrinksDrift(t *testing.T) {
	filter := NewSquareRootInformationFilter()
	filter.sqrtInfo[0] = 0.25
	filter.sqrtInfo[1] = 0.3
	filter.updateDriftEllipse()
	before := filter.Snapshot().DriftEllipseM

	state := filter.UpdateWithVision(VisionObservation{
		Timestamp:    time.Now().UTC(),
		DeltaXM:      0.4,
		DeltaYM:      -0.1,
		VxMS:         1.2,
		VyMS:         -0.2,
		FeatureCount: 48,
		Quality:      0.9,
	})

	if state.DriftEllipseM[0] >= before[0] || state.DriftEllipseM[1] >= before[1] {
		t.Fatalf("expected drift ellipse to shrink, before=%v after=%v", before, state.DriftEllipseM)
	}
	if state.FeatureCount != 48 {
		t.Fatalf("expected feature count to update, got %d", state.FeatureCount)
	}
	if math.Abs(state.VelocityMS[0]-1.2) > 0.35 {
		t.Fatalf("expected x velocity to move toward VO measurement, got %.3f", state.VelocityMS[0])
	}
}
