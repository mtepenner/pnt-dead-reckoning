from __future__ import annotations

from dataclasses import dataclass

import cv2
import numpy as np


@dataclass(slots=True)
class FeatureSet:
    points: list[tuple[float, float]]
    strengths: list[float]


class FeatureExtractor:
    def __init__(self, max_features: int = 160) -> None:
        self._orb = cv2.ORB_create(nfeatures=max_features)

    def extract(self, frame: np.ndarray) -> FeatureSet:
        keypoints = self._orb.detect(frame, None)
        keypoints = sorted(keypoints, key=lambda keypoint: keypoint.response, reverse=True)[:80]
        points = [(float(keypoint.pt[0]), float(keypoint.pt[1])) for keypoint in keypoints]
        strengths = [float(keypoint.response) for keypoint in keypoints]
        return FeatureSet(points=points, strengths=strengths)
