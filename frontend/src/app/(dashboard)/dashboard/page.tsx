'use client';

import { useAuth } from '@/contexts/AuthContext';

export default function DashboardPage() {
  const { user } = useAuth();

  if (!user) return null;

  return (
    <div className="min-h-screen bg-gray-50/30 text-gray-600 font-sans">
      <main className="max-w-5xl mx-auto px-8 py-16">
        <header className="mb-16">
          <h1 className="text-3xl font-bold text-gray-900 tracking-tight uppercase">Dashboard</h1>
          <p className="text-sm text-gray-400 mt-2 font-medium tracking-wide">{user.email}</p>
        </header>

        <div className="bg-white p-12 rounded-2xl border border-gray-100 shadow-sm text-center">
          <div className="max-w-md mx-auto">
            <p className="text-sm text-gray-500 leading-relaxed mb-8">
              Event Managerへようこそ。イベントの作成や管理を行うことができます。
            </p>
          </div>
        </div>
      </main>
    </div>
  );
}
