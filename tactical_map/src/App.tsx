import { DriftIndicator } from './components/DriftIndicator';
import { LiveBreadcrumbs } from './components/LiveBreadcrumbs';
import { VisionFeed } from './components/VisionFeed';
import { useNavigationState } from './hooks/useNavigationState';

export default function App() {
  const { snapshot, status } = useNavigationState();
  const state = snapshot?.state ?? null;
  const history = snapshot?.history ?? [];

  return (
    <main className="shell">
      <section className="hero panel">
        <div>
          <p className="eyebrow">GPS-Denied Navigation</p>
          <h1>Dead reckoning console for fused IMU and visual odometry tracks.</h1>
        </div>
        <div className="hero-meta">
          <span>{status}</span>
          <span>{state ? `${state.position_m[0].toFixed(2)}, ${state.position_m[1].toFixed(2)} m` : 'Awaiting initial path'}</span>
        </div>
      </section>

      <section className="dashboard-grid">
        <section className="panel stack-gap wide-panel">
          <div className="section-heading">
            <span>Live Breadcrumb Trail</span>
            <small>{history.length} fused fixes</small>
          </div>
          <LiveBreadcrumbs state={state} history={history} />
        </section>

        <div className="sidebar-stack">
          <DriftIndicator state={state} />
          <VisionFeed state={state} />
          <section className="panel stack-gap">
            <div className="section-heading">
              <span>State Vector</span>
              <small>SRIF output</small>
            </div>
            <div className="metric-row">
              <span>Velocity</span>
              <strong>{state ? `${state.velocity_m_s[0].toFixed(2)}, ${state.velocity_m_s[1].toFixed(2)} m/s` : 'Pending'}</strong>
            </div>
            <div className="metric-row">
              <span>Attitude</span>
              <strong>{state ? `${state.attitude_deg[0].toFixed(1)}, ${state.attitude_deg[1].toFixed(1)}, ${state.attitude_deg[2].toFixed(1)} deg` : 'Pending'}</strong>
            </div>
          </section>
        </div>
      </section>
    </main>
  );
}
