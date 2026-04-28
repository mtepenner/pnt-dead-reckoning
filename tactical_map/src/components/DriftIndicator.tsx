import { NavigationState } from '../types';

type Props = {
  state: NavigationState | null;
};

export function DriftIndicator({ state }: Props) {
  const major = state?.drift_ellipse_m[0] ?? 0;
  const minor = state?.drift_ellipse_m[1] ?? 0;
  const width = Math.max(32, major * 38);
  const height = Math.max(24, minor * 42);

  return (
    <section className="panel stack-gap">
      <div className="section-heading">
        <span>Drift Indicator</span>
        <small>{state ? `${state.vision_quality.toFixed(2)} vision weight` : 'No state fix yet'}</small>
      </div>
      <svg viewBox="0 0 240 160" className="drift-chart" role="img" aria-label="Uncertainty ellipse">
        <line x1="20" y1="80" x2="220" y2="80" />
        <line x1="120" y1="20" x2="120" y2="140" />
        <ellipse cx="120" cy="80" rx={width} ry={height} />
        <circle cx="120" cy="80" r="4" />
      </svg>
      <div className="metric-row">
        <span>Major axis</span>
        <strong>{major.toFixed(2)} m</strong>
      </div>
      <div className="metric-row">
        <span>Minor axis</span>
        <strong>{minor.toFixed(2)} m</strong>
      </div>
    </section>
  );
}
