from __future__ import annotations

import argparse
import json
import math
import random


def generate_samples(steps: int, seed: int) -> list[dict[str, float]]:
    random.seed(seed)
    samples: list[dict[str, float]] = []
    for index in range(steps):
        t = index * 0.25
        true_x = 12 * math.cos(t / 3.0) + 0.4 * index
        true_y = 9 * math.sin(t / 4.0) + 0.2 * index
        noisy_x = true_x + random.gauss(0, 0.45)
        noisy_y = true_y + random.gauss(0, 0.45)
        samples.append(
            {
                "time_s": round(t, 3),
                "true_x_m": round(true_x, 3),
                "true_y_m": round(true_y, 3),
                "noisy_x_m": round(noisy_x, 3),
                "noisy_y_m": round(noisy_y, 3),
            }
        )
    return samples


def main() -> None:
    parser = argparse.ArgumentParser(description="Generate deterministic true/noisy path samples.")
    parser.add_argument("--steps", type=int, default=20)
    parser.add_argument("--seed", type=int, default=42)
    args = parser.parse_args()
    print(json.dumps(generate_samples(args.steps, args.seed), indent=2))


if __name__ == "__main__":
    main()
