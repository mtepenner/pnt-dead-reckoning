package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/mtepenner/pnt-dead-reckoning/navigation_core/internal/filter"
	"github.com/mtepenner/pnt-dead-reckoning/navigation_core/internal/sensors"
	"github.com/mtepenner/pnt-dead-reckoning/navigation_core/internal/telemetry"
)

type latestObservation struct {
	mu          sync.RWMutex
	observation filter.VisionObservation
	available   bool
}

func main() {
	hub := telemetry.NewHub()
	imuDriver := sensors.NewDriver()
	fusion := filter.NewSquareRootInformationFilter()
	vision := &latestObservation{}

	go serveVisionIngress(vision)
	go runFusionLoop(fusion, imuDriver, hub, vision)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"service":"navigation-core","status":"ok"}`))
	})
	mux.HandleFunc("/api/state", hub.HandleState)
	mux.HandleFunc("/api/history", hub.HandleHistory)
	mux.HandleFunc("/ws", hub.HandleWebSocket)

	log.Println("navigation core listening on http://127.0.0.1:8081")
	if err := http.ListenAndServe(":8081", withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func runFusionLoop(fusion *filter.SquareRootInformationFilter, imuDriver *sensors.Driver, hub *telemetry.Hub, vision *latestObservation) {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	history := make([]filter.HistoryPoint, 0, 256)
	for range ticker.C {
		sample := imuDriver.Next()
		state := fusion.Predict(sample.AccelMS2, sample.GyroRadS, 0.25)
		state.Timestamp = sample.Timestamp

		if observation, ok := vision.get(); ok && time.Since(observation.Timestamp) < 2*time.Second {
			state = fusion.UpdateWithVision(observation)
		}

		history = append(history, filter.HistoryPoint{
			Timestamp:    state.Timestamp,
			X:            state.PositionM[0],
			Y:            state.PositionM[1],
			DriftMajorM:  state.DriftEllipseM[0],
			DriftMinorM:  state.DriftEllipseM[1],
			HeadingDeg:   state.AttitudeDeg[2],
			VisionWeight: state.VisionQuality,
		})
		if len(history) > 180 {
			history = append([]filter.HistoryPoint(nil), history[len(history)-180:]...)
		}

		hub.Publish(telemetry.Snapshot{State: state, History: append([]filter.HistoryPoint(nil), history...)})
	}
}

func serveVisionIngress(latest *latestObservation) {
	transport, address := defaultVisionTransport()
	transport = envOrDefault("VO_TRANSPORT", transport)
	address = envOrDefault("VO_ADDRESS", address)

	listener, err := listen(transport, address)
	if err != nil {
		log.Printf("vision ingress unavailable: %v", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("vision accept error: %v", err)
			continue
		}
		go handleVisionConn(conn, latest)
	}
}

func handleVisionConn(conn net.Conn, latest *latestObservation) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var observation filter.VisionObservation
		if err := json.Unmarshal(scanner.Bytes(), &observation); err != nil {
			continue
		}
		if observation.Timestamp.IsZero() {
			observation.Timestamp = time.Now().UTC()
		}
		latest.set(observation)
	}
}

func defaultVisionTransport() (string, string) {
	if runtime.GOOS == "windows" {
		return "tcp", "127.0.0.1:9101"
	}
	return "unix", "/tmp/pnt_vo.sock"
}

func listen(transport, address string) (net.Listener, error) {
	if transport == "unix" {
		_ = os.Remove(address)
	}
	return net.Listen(transport, address)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func (latest *latestObservation) set(observation filter.VisionObservation) {
	latest.mu.Lock()
	defer latest.mu.Unlock()
	latest.observation = observation
	latest.available = true
}

func (latest *latestObservation) get() (filter.VisionObservation, bool) {
	latest.mu.RLock()
	defer latest.mu.RUnlock()
	return latest.observation, latest.available
}
