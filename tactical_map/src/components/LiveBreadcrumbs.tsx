import { CRS, LatLngExpression } from 'leaflet';
import { CircleMarker, MapContainer, Polygon, Polyline, Rectangle, useMap } from 'react-leaflet';
import { useEffect, useMemo } from 'react';

import { HistoryPoint, NavigationState } from '../types';

type Props = {
  state: NavigationState | null;
  history: HistoryPoint[];
};

export function LiveBreadcrumbs({ state, history }: Props) {
  const path = useMemo<LatLngExpression[]>(() => history.map((point) => [point.y, point.x]), [history]);
  const latest = history[history.length - 1];
  const ellipse = useMemo(() => buildEllipse(latest, state), [latest, state]);

  return (
    <div className="map-shell">
      <MapContainer
        center={latest ? [latest.y, latest.x] : [0, 0]}
        zoom={2}
        minZoom={-2}
        maxZoom={5}
        crs={CRS.Simple}
        className="map-canvas"
        attributionControl={false}
      >
        <Rectangle bounds={[[-80, -80], [80, 80]]} pathOptions={{ color: '#224f55', weight: 1 }} />
        <AutoCenter latest={latest ? [latest.y, latest.x] : null} />
        {path.length > 1 && <Polyline positions={path} pathOptions={{ color: '#7fe7be', weight: 3 }} />}
        {latest && <CircleMarker center={[latest.y, latest.x]} radius={6} pathOptions={{ color: '#ffcc74', fillOpacity: 0.8 }} />}
        {ellipse.length > 0 && <Polygon positions={ellipse} pathOptions={{ color: '#ff8b5c', weight: 2, fillOpacity: 0.12 }} />}
      </MapContainer>
    </div>
  );
}

function AutoCenter({ latest }: { latest: [number, number] | null }) {
  const map = useMap();

  useEffect(() => {
    if (latest) {
      map.setView(latest, map.getZoom(), { animate: false });
    }
  }, [latest, map]);

  return null;
}

function buildEllipse(latest: HistoryPoint | undefined, state: NavigationState | null): LatLngExpression[] {
  if (!latest || !state) {
    return [];
  }

  const [major, minor] = state.drift_ellipse_m;
  const points: LatLngExpression[] = [];
  for (let index = 0; index < 24; index += 1) {
    const theta = (Math.PI * 2 * index) / 24;
    points.push([latest.y + Math.sin(theta) * minor, latest.x + Math.cos(theta) * major]);
  }
  return points;
}
