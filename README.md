# 🧭 PNT Dead Reckoning

## Description
PNT (Position, Navigation, and Timing) Dead Reckoning is a robust sensor fusion pipeline designed for accurate navigation in GPS-denied environments. The system relies on a Python/OpenCV visual odometry provider and an IMU driver to gather raw spatial data. A high-performance Go-based navigation core ingests this data, applying a Square-Root Information Filter (SRIF) to calculate coordinates, velocity, and orientation. Operators can monitor the calculated path, visual tracking features, and drift uncertainty in real-time via a React/TypeScript tactical dashboard.

## 📑 Table of Contents
- [Features](#-features)
- [Technologies Used](#-technologies-used)
- [Installation](#-installation)
- [Usage](#-usage)
- [Project Structure](#-project-structure)
- [Contributing](#-contributing)
- [License](#-license)

## 🚀 Features
* **Visual Odometry (VO):** A Python application utilizing OpenCV to detect ORB/SIFT keypoints and calculate optical flow, streaming velocity vectors via Unix sockets.
* **SRIF Sensor Fusion Brain:** A Go-based navigation core that reads raw IMU (Accel/Gyro) data and fuses it with the VO feed using an advanced Square-Root Information Filter (SRIF) to maintain a highly accurate 9-DOF state vector.
* **Tactical Navigation HUD:** A React and TypeScript interface featuring a Mapbox/Leaflet live breadcrumb trail, real-time CV feature tracking overlays, and a visual drift indicator (Uncertainty Ellipse).
* **Local Simulation:** Includes Python scripts to generate "True" versus "Noisy" drone path data, allowing for rapid local testing of the filter's stability without physical hardware.

## 🛠️ Technologies Used
* **Navigation Core:** Go
* **Vision Provider:** Python, OpenCV, GStreamer
* **Frontend Dashboard:** React, TypeScript, Mapbox/Leaflet
* **Inter-Process Communication:** Unix Sockets, WebSockets
* **Infrastructure:** Docker, Docker Compose

## ⚙️ Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/pnt-dead-reckoning.git
   cd pnt-dead-reckoning
   ```

2. Boot the full stack (Vision Provider, Navigation Core, and Tactical UI) using Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. *Optional:* If you do not have camera hardware attached, run the local simulator to generate sample telemetry:
   ```bash
   python simulation/drone_path_generator.py
   ```

## 💻 Usage

* **Access the Tactical Map:** Open your browser and navigate to `http://localhost:3000` (or your configured port) to view the Live Breadcrumbs and Drift Indicator.
* **Monitor the Vision Feed:** Toggle the Vision Feed overlay in the UI to see real-time optical flow and keypoint extraction processing from the camera.
* **Analyze State History:** The Go daemon serves coordinate history over WebSockets, allowing the UI to instantly trace the historical path of the agent.

## 📂 Project Structure
* `/vision_provider`: Python/OpenCV container handling raw frame capture and visual odometry calculations.
* `/navigation_core`: Go application executing the primary SRIF math and state vector management.
* `/tactical_map`: React frontend providing the real-time PNT mapping interface.
* `/simulation`: Tools for generating noisy test data for filter calibration.
* `/.github/workflows`: CI/CD pipelines including rigorous validation for the SRIF matrix math.

## 🤝 Contributing
Contributions are highly encouraged! When submitting a Pull Request, especially involving the `/filter` directory, please ensure all automated tests for filter stability and SRIF matrix math pass successfully.

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
