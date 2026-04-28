import { NavigationState } from '../types';

type Props = {
  state: NavigationState | null;
};

export function VisionFeed({ state }: Props) {
  const tracks = state?.tracked_features ?? [];

  return (
    <section className="panel stack-gap">
      <div className="section-heading">
        <span>Vision Feed</span>
        <small>{state ? `${state.feature_count} tracked features` : 'No feature data yet'}</small>
      </div>
      <svg viewBox="0 0 640 480" className="vision-overlay" role="img" aria-label="Feature track overlay">
        <rect x="0" y="0" width="640" height="480" rx="20" />
        {tracks.map((track, index) => (
          <g key={`${track.x}-${track.y}-${index}`}>
            <line x1={track.x - track.dx} y1={track.y - track.dy} x2={track.x} y2={track.y} />
            <circle cx={track.x} cy={track.y} r="4" />
          </g>
        ))}
      </svg>
      <div className="metric-row">
        <span>Optical flow</span>
        <strong>
          {state ? `${state.optical_flow_m_s[0].toFixed(2)} / ${state.optical_flow_m_s[1].toFixed(2)} m/s` : 'Pending'}
        </strong>
      </div>
    </section>
  );
}
