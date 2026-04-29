'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '@/lib/api';
import Link from 'next/link';

export default function NewEventPage() {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [startTime, setStartTime] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleCreateEvent = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      await apiRequest('/api/v1/events', {
        method: 'POST',
        body: JSON.stringify({
          title,
          description,
          start_time: startTime ? new Date(startTime).toISOString() : new Date().toISOString(),
        }),
      });
      router.push('/dashboard');
    } catch (err: any) {
      alert(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50/30 text-gray-600 font-sans">
      <main className="max-w-3xl mx-auto px-8 py-16">
        <header className="mb-12">
          <Link href="/dashboard" className="text-xs font-bold text-gray-400 uppercase tracking-widest hover:text-blue-600 transition-colors mb-4 block">
            ← ダッシュボードに戻る
          </Link>
          <h1 className="text-3xl font-bold text-gray-900 tracking-tight">Create New Event</h1>
          <p className="text-sm text-gray-400 mt-2">イベントは「下書き」として作成されます。後で公開設定が可能です。</p>
        </header>

        <div className="bg-white p-10 rounded-2xl border border-gray-100 shadow-sm">
          <form onSubmit={handleCreateEvent} className="space-y-10">
            <div>
              <label className="block text-[10px] font-bold text-gray-400 uppercase tracking-[0.2em] mb-3">イベント名</label>
              <input
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full p-4 text-lg border border-gray-100 rounded-xl focus:ring-1 focus:ring-blue-600 outline-none transition-all"
                placeholder="例: 第1回 コミュニティミートアップ"
                required
              />
            </div>

            <div>
              <label className="block text-[10px] font-bold text-gray-400 uppercase tracking-[0.2em] mb-3">開催日時</label>
              <input
                type="datetime-local"
                value={startTime}
                onChange={(e) => setStartTime(e.target.value)}
                className="w-full p-4 text-base border border-gray-100 rounded-xl focus:ring-1 focus:ring-blue-600 outline-none transition-all"
                required
              />
            </div>

            <div>
              <label className="block text-[10px] font-bold text-gray-400 uppercase tracking-[0.2em] mb-3">イベント概要</label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={6}
                className="w-full p-4 text-base border border-gray-100 rounded-xl focus:ring-1 focus:ring-blue-600 outline-none transition-all resize-none"
                placeholder="イベントの詳細な説明を入力してください..."
              />
            </div>

            <div className="flex gap-6 pt-6 border-t border-gray-50">
              <button
                type="submit"
                disabled={isLoading}
                className="flex-1 py-4 bg-blue-600 text-white text-xs font-bold uppercase tracking-widest rounded-xl hover:bg-blue-700 transition-all shadow-xl shadow-blue-50 disabled:bg-blue-300"
              >
                {isLoading ? '作成中...' : '下書きとして保存'}
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
