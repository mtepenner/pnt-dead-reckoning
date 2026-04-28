package filter

import "time"

type FeatureTrack struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	DX float64 `json:"dx"`
	DY float64 `json:"dy"`
}

type VisionObservation struct {
	Timestamp    time.Time      `json:"timestamp"`
	DeltaXM      float64        `json:"delta_x_m"`
	DeltaYM      float64        `json:"delta_y_m"`
	VxMS         float64        `json:"vx_m_s"`
	VyMS         float64        `json:"vy_m_s"`
	FeatureCount int            `json:"feature_count"`
	Quality      float64        `json:"quality"`
	Tracks       []FeatureTrack `json:"tracks"`
}

type StateVector struct {
	Timestamp       time.Time      `json:"timestamp"`
	PositionM       [3]float64     `json:"position_m"`
	VelocityMS      [3]float64     `json:"velocity_m_s"`
	AttitudeDeg     [3]float64     `json:"attitude_deg"`
	DriftEllipseM   [2]float64     `json:"drift_ellipse_m"`
	FeatureCount    int            `json:"feature_count"`
	VisionQuality   float64        `json:"vision_quality"`
	OpticalFlowMS   [2]float64     `json:"optical_flow_m_s"`
	TrackedFeatures []FeatureTrack `json:"tracked_features"`
}

type HistoryPoint struct {
	Timestamp    time.Time `json:"timestamp"`
	X            float64   `json:"x"`
	Y            float64   `json:"y"`
	DriftMajorM  float64   `json:"drift_major_m"`
	DriftMinorM  float64   `json:"drift_minor_m"`
	HeadingDeg   float64   `json:"heading_deg"`
	VisionWeight float64   `json:"vision_weight"`
}
