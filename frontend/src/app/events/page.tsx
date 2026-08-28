'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { apiRequest } from '@/lib/api';
import Link from 'next/link';

export default function EventsPage() {
  const { user } = useAuth();
  const [events, setEvents] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const data = await apiRequest('/api/v1/events/me');
        setEvents(data || []);
      } catch (err) {
        console.error('Failed to fetch events:', err);
      } finally {
        setIsLoading(false);
      }
    };
    if (user) fetchEvents();
  }, [user]);

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleString('ja-JP', {
      month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
    });
  };

  if (!user) return null;

  return (
    <div className="min-h-screen bg-gray-50/30 text-gray-600 font-sans p-8 lg:p-16">
      <main className="max-w-6xl mx-auto">
        <header className="mb-16 flex justify-between items-end border-b border-gray-100 pb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 tracking-tight">イベント管理</h1>
            <p className="text-sm text-gray-400 mt-2 font-medium tracking-wide">作成済みのイベントを管理します。</p>
          </div>
          <Link href="/events/new" className="px-6 py-2.5 bg-blue-600 text-white text-sm font-bold rounded-xl hover:bg-blue-700 shadow-lg shadow-blue-100 transition-all">
            ＋ 新規作成
          </Link>
        </header>

        <div className="space-y-6">
          {isLoading ? (
            [1, 2].map(i => <div key={i} className="h-32 bg-white rounded-3xl animate-pulse border border-gray-100" />)
          ) : events.length === 0 ? (
            <div className="bg-white p-24 rounded-[3rem] border border-gray-100 shadow-sm text-center">
              <p className="text-sm text-gray-400 italic mb-8">管理中のイベントはありません</p>
              <Link href="/events/new" className="px-8 py-3 bg-blue-50 text-blue-600 rounded-xl text-sm font-bold tracking-widest hover:bg-blue-100 transition-all">
                最初のイベントを作成する →
              </Link>
            </div>
          ) : (
            events.map((event) => (
              <div key={event.id} className="bg-white rounded-[2rem] border border-gray-100 shadow-sm hover:shadow-md transition-all flex flex-col md:flex-row overflow-hidden group">
                <div className="w-full md:w-48 h-32 md:h-auto bg-gray-100 shrink-0 relative">
                  {event.thumbnail_url ? (
                    <img src={event.thumbnail_url} alt={event.title} className="w-full h-full object-cover" />
                  ) : (
                    <div className="w-full h-full bg-gradient-to-br from-gray-50 to-gray-100" />
                  )}
                  <div className="absolute top-3 left-3 flex gap-1.5">
                    <span className={`px-2.5 py-1 rounded text-xs font-black tracking-widest ${event.status === 'published' ? 'bg-blue-600 text-white' : 'bg-gray-400 text-white'}`}>
                      {event.status === 'published' ? '公開中' : '下書き'}
                    </span>
                    <span className={`px-2.5 py-1 rounded text-xs font-black tracking-widest bg-white/90 backdrop-blur-sm text-gray-600`}>
                      {event.is_online ? 'オンライン' : '会場'}
                    </span>
                  </div>
                </div>
                <div className="flex-1 p-6 md:p-8 flex flex-col md:flex-row md:items-center justify-between gap-6">
                  <div className="flex-1 min-w-0">
                    <Link href={`/events/${event.id}`} className="block group-hover:text-blue-600 transition-colors">
                      <h3 className="text-2xl font-bold text-gray-900 truncate">{event.title}</h3>
                    </Link>
                    <div className="mt-2 flex items-center gap-4 text-sm font-bold text-gray-500 tracking-widest">
                      <span>{formatDate(event.start_time)}</span>
                      {event.capacity > 0 && <span>• 定員 {event.capacity}名</span>}
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <Link href={`/events/${event.id}`} className="px-5 py-2.5 bg-white text-gray-600 border border-gray-100 text-sm font-bold tracking-widest rounded-xl hover:bg-gray-50 transition-all shadow-sm">
                      閲覧
                    </Link>
                    <Link href={`/events/${event.id}/edit`} className="px-5 py-2.5 bg-blue-50 text-blue-600 text-sm font-bold tracking-widest rounded-xl hover:bg-blue-100 transition-all">
                      編集
                    </Link>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </main>
    </div>
  );
}
