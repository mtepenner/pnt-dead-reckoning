from __future__ import annotations

from dataclasses import dataclass

import cv2
import numpy as np


@dataclass(slots=True)
class FeatureTrack:
    x: float
    y: float
    dx: float
    dy: float


@dataclass(slots=True)
class MotionEstimate:
    dx_pixels: float
    dy_pixels: float
    vx_m_s: float
    vy_m_s: float
    quality: float
    tracks: list[FeatureTrack]


class OpticalFlowEstimator:
    def __init__(self, meters_per_pixel: float = 0.018, frame_period_s: float = 0.25) -> None:
        self._meters_per_pixel = meters_per_pixel
        self._frame_period_s = frame_period_s

    def estimate(
        self,
        previous_frame: np.ndarray,
        current_frame: np.ndarray,
        previous_points: list[tuple[float, float]],
    ) -> MotionEstimate:
        if len(previous_points) < 4:
            return MotionEstimate(0.0, 0.0, 0.0, 0.0, 0.0, [])

        previous_np = np.array(previous_points, dtype=np.float32).reshape(-1, 1, 2)
        next_points, status, _ = cv2.calcOpticalFlowPyrLK(
            previous_frame,
            current_frame,
            previous_np,
            None,
            winSize=(21, 21),
            maxLevel=3,
            criteria=(cv2.TERM_CRITERIA_EPS | cv2.TERM_CRITERIA_COUNT, 20, 0.03),
        )

        if next_points is None or status is None:
            return MotionEstimate(0.0, 0.0, 0.0, 0.0, 0.0, [])

        valid = status.flatten() == 1
        if not np.any(valid):
            return MotionEstimate(0.0, 0.0, 0.0, 0.0, 0.0, [])

        prev_valid = previous_np[valid][:, 0, :]
        next_valid = next_points[valid][:, 0, :]
        deltas = next_valid - prev_valid
        mean_delta = deltas.mean(axis=0)

        dx_pixels = float(mean_delta[0])
        dy_pixels = float(mean_delta[1])
        vx_m_s = dx_pixels * self._meters_per_pixel / self._frame_period_s
        vy_m_s = dy_pixels * self._meters_per_pixel / self._frame_period_s
        quality = min(1.0, float(len(prev_valid)) / 48.0)
        tracks = [
            FeatureTrack(
                x=float(point[0]),
                y=float(point[1]),
                dx=float(delta[0]),
                dy=float(delta[1]),
            )
            for point, delta in zip(next_valid[:32], deltas[:32], strict=False)
        ]
        return MotionEstimate(dx_pixels, dy_pixels, vx_m_s, vy_m_s, quality, tracks)
