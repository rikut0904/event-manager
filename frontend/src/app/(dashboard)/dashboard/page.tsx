'use client';

import { useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

export default function DashboardPage() {
  const { user } = useAuth();

  if (!user) return null;

  return (
    <div className="min-h-screen bg-gray-50/30 text-gray-600">
      <main className="max-w-5xl mx-auto px-8 py-16">
        <header className="mb-16 border-b border-gray-100 pb-8">
          <h1 className="text-3xl font-bold text-gray-900 tracking-tight uppercase">Dashboard</h1>
        </header>

        <div className="grid grid-cols-1 gap-12">
          <div className="bg-white p-12 rounded-2xl border border-gray-100 shadow-sm text-center">
            <div className="max-w-md mx-auto">
              <p className="text-sm text-gray-500 leading-relaxed">
                ダッシュボードへようこそ。ここでは、あなたが参加する予定のイベントを一覧で確認できるようになる予定です。
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
