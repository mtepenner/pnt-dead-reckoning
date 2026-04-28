import { useEffect, useMemo, useState } from 'react';

import { NavigationSnapshot } from '../types';

const apiBaseUrl = import.meta.env.VITE_NAV_API_URL ?? 'http://127.0.0.1:8081';

export function useNavigationState() {
  const [snapshot, setSnapshot] = useState<NavigationSnapshot | null>(null);
  const [status, setStatus] = useState('Waiting for telemetry');

  useEffect(() => {
    let socket: WebSocket | null = null;
    let reconnectTimer: number | null = null;
    let cancelled = false;

    const connect = () => {
      const wsUrl = `${apiBaseUrl.replace('http://', 'ws://').replace('https://', 'wss://')}/ws`;
      socket = new WebSocket(wsUrl);

      socket.onopen = () => {
        if (!cancelled) {
          setStatus('Connected to navigation core');
        }
      };

      socket.onmessage = (event) => {
        if (cancelled) {
          return;
        }
        const nextSnapshot = JSON.parse(event.data) as NavigationSnapshot;
        setSnapshot(nextSnapshot);
        setStatus(`Last fix ${new Date(nextSnapshot.state.timestamp).toLocaleTimeString()}`);
      };

      socket.onerror = () => {
        if (!cancelled) {
          setStatus('Telemetry stream interrupted');
        }
      };

      socket.onclose = () => {
        if (!cancelled) {
          reconnectTimer = window.setTimeout(connect, 1000);
        }
      };
    };

    connect();
    return () => {
      cancelled = true;
      socket?.close();
      if (reconnectTimer !== null) {
        window.clearTimeout(reconnectTimer);
      }
    };
  }, []);

  return useMemo(() => ({
    snapshot,
    status,
  }), [snapshot, status]);
}
