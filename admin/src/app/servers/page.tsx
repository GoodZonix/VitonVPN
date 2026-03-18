'use client';

import useSWR from 'swr';

const fetcher = (url: string) => fetch(url).then((r) => r.json());

export default function ServersPage() {
  const api = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const { data, error } = useSWR<{ servers: Array<{
    id: string;
    name: string;
    region: string;
    host: string;
    port: number;
    type: string;
    is_active: boolean;
  }> }>(`${api}/api/servers`, fetcher);

  if (error) return <div>Failed to load servers</div>;
  if (!data) return <div>Loading...</div>;

  return (
    <div style={{ padding: 24 }}>
      <h1>Servers</h1>
      <p style={{ color: '#666' }}>Manage VPN nodes. Add/edit via database or separate admin API.</p>
      <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: 16 }}>
        <thead>
          <tr style={{ borderBottom: '1px solid #ddd' }}>
            <th style={{ textAlign: 'left', padding: 8 }}>Name</th>
            <th style={{ textAlign: 'left', padding: 8 }}>Region</th>
            <th style={{ textAlign: 'left', padding: 8 }}>Host</th>
            <th style={{ textAlign: 'left', padding: 8 }}>Port</th>
            <th style={{ textAlign: 'left', padding: 8 }}>Type</th>
            <th style={{ textAlign: 'left', padding: 8 }}>Active</th>
          </tr>
        </thead>
        <tbody>
          {data.servers.map((s) => (
            <tr key={s.id} style={{ borderBottom: '1px solid #eee' }}>
              <td style={{ padding: 8 }}>{s.name}</td>
              <td style={{ padding: 8 }}>{s.region}</td>
              <td style={{ padding: 8 }}>{s.host}</td>
              <td style={{ padding: 8 }}>{s.port}</td>
              <td style={{ padding: 8 }}>{s.type}</td>
              <td style={{ padding: 8 }}>{s.is_active ? 'Yes' : 'No'}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
