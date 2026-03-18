'use client';

import Link from 'next/link';

export default function AdminHome() {
  return (
    <div style={{ padding: 24, maxWidth: 800 }}>
      <h1>VPN Startup — Admin</h1>
      <nav style={{ display: 'flex', gap: 16, marginTop: 24 }}>
        <Link href="/servers">Servers</Link>
        <Link href="/users">Users</Link>
        <Link href="/subscriptions">Subscriptions</Link>
        <Link href="/stats">Traffic stats</Link>
      </nav>
      <p style={{ marginTop: 24, color: '#666' }}>
        Backend API base URL: <code>{process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}</code>
      </p>
    </div>
  );
}
