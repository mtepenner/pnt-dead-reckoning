export type FeatureTrack = {
  x: number;
  y: number;
  dx: number;
  dy: number;
};

export type NavigationState = {
  timestamp: string;
  position_m: [number, number, number];
  velocity_m_s: [number, number, number];
  attitude_deg: [number, number, number];
  drift_ellipse_m: [number, number];
  feature_count: number;
  vision_quality: number;
  optical_flow_m_s: [number, number];
  tracked_features: FeatureTrack[];
};

export type HistoryPoint = {
  timestamp: string;
  x: number;
  y: number;
  drift_major_m: number;
  drift_minor_m: number;
  heading_deg: number;
  vision_weight: number;
};

export type NavigationSnapshot = {
  state: NavigationState;
  history: HistoryPoint[];
};
