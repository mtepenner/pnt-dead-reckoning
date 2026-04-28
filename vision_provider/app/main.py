from __future__ import annotations

import json
import os
import platform
import socket
import time
from datetime import datetime, timezone
from pathlib import Path

from camera.gstreamer_pipeline import SyntheticGroundCamera
from cv.feature_extractor import FeatureExtractor
from cv.optical_flow import OpticalFlowEstimator


def default_transport() -> tuple[str, str]:
    if platform.system().lower().startswith("win"):
        return "tcp", "127.0.0.1:9101"
    return "unix", "/tmp/pnt_vo.sock"


def connect_stream(transport: str, address: str) -> socket.socket:
    if transport == "tcp":
        host, port = address.split(":", 1)
        return socket.create_connection((host, int(port)))

    client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    client.connect(address)
    return client


def stream_visual_odometry() -> None:
    transport, address = default_transport()
    transport = os.getenv("VO_TRANSPORT", transport)
    address = os.getenv("VO_ADDRESS", address)
    max_frames = int(os.getenv("PROVIDER_MAX_FRAMES", "0"))
    sleep_seconds = float(os.getenv("PROVIDER_FRAME_PERIOD_S", "0.25"))

    camera = SyntheticGroundCamera()
    extractor = FeatureExtractor()
    estimator = OpticalFlowEstimator(frame_period_s=sleep_seconds)
    previous_frame = None
    previous_points: list[tuple[float, float]] = []
    sent_frames = 0

    while True:
        try:
            stream = connect_stream(transport, address)
            break
        except OSError:
            time.sleep(0.5)

    with stream:
        while True:
            frame = camera.next_frame()
            features = extractor.extract(frame.image)
            if previous_frame is None:
                previous_frame = frame.image
                previous_points = features.points
                time.sleep(sleep_seconds)
                continue

            motion = estimator.estimate(previous_frame, frame.image, previous_points)
            payload = {
                "timestamp": datetime.fromtimestamp(frame.timestamp, tz=timezone.utc).isoformat().replace("+00:00", "Z"),
                "delta_x_m": motion.dx_pixels * 0.018,
                "delta_y_m": motion.dy_pixels * 0.018,
                "vx_m_s": motion.vx_m_s,
                "vy_m_s": motion.vy_m_s,
                "feature_count": len(features.points),
                "quality": motion.quality,
                "tracks": [
                    {"x": track.x, "y": track.y, "dx": track.dx, "dy": track.dy}
                    for track in motion.tracks
                ],
                "preview_path": str(Path("frames") / "synthetic-preview"),
            }

            try:
                stream.sendall((json.dumps(payload) + "\n").encode("utf-8"))
            except OSError:
                stream.close()
                stream = connect_stream(transport, address)
                continue

            previous_frame = frame.image
            previous_points = features.points
            sent_frames += 1
            if max_frames and sent_frames >= max_frames:
                break
            time.sleep(sleep_seconds)


if __name__ == "__main__":
    stream_visual_odometry()
