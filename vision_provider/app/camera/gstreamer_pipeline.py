from __future__ import annotations

from dataclasses import dataclass
import math
import time

import cv2
import numpy as np


@dataclass(slots=True)
class SyntheticFrame:
    image: np.ndarray
    true_dx_pixels: float
    true_dy_pixels: float
    timestamp: float


def build_gstreamer_pipeline(width: int = 640, height: int = 480, fps: int = 20) -> str:
    return (
        f"videotestsrc is-live=true ! video/x-raw,width={width},height={height},framerate={fps}/1 "
        "! videoconvert ! appsink"
    )


class SyntheticGroundCamera:
    def __init__(self, width: int = 640, height: int = 480) -> None:
        self.width = width
        self.height = height
        self._step = 0

    def next_frame(self) -> SyntheticFrame:
        phase = self._step * 0.18
        dx = 4.0 * math.cos(phase * 0.7)
        dy = 2.5 * math.sin(phase * 0.5)

        canvas = np.zeros((self.height, self.width), dtype=np.uint8)
        for column in range(60, self.width, 90):
            for row in range(50, self.height, 70):
                x = int(column + dx + 12 * math.sin((row + phase * 20) / 55.0))
                y = int(row + dy + 10 * math.cos((column + phase * 15) / 80.0))
                cv2.circle(canvas, (x, y), 4, 220, -1)

        runway_points = np.array(
            [
                [110 + dx, 360 + dy],
                [280 + dx, 250 + dy],
                [520 + dx, 245 + dy],
                [610 + dx, 340 + dy],
            ],
            dtype=np.int32,
        )
        cv2.polylines(canvas, [runway_points], isClosed=False, color=180, thickness=2)
        cv2.putText(canvas, "PNT", (34, 42), cv2.FONT_HERSHEY_SIMPLEX, 0.9, 200, 2, cv2.LINE_AA)

        frame = SyntheticFrame(
            image=cv2.GaussianBlur(canvas, (5, 5), 0),
            true_dx_pixels=dx,
            true_dy_pixels=dy,
            timestamp=time.time(),
        )
        self._step += 1
        return frame
