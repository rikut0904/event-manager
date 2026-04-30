'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { apiRequest } from '@/lib/api';
import Link from 'next/link';

export default function EventsPage() {
  const { user } = useAuth();
  const [events, setEvents] = useState<any[]>([]);

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const data = await apiRequest('/api/v1/events/me');
        setEvents(data || []);
      } catch (err) {
        console.error('Failed to fetch events:', err);
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
      <main className="max-w-5xl mx-auto">
        <header className="mb-16 flex justify-between items-end border-b border-gray-100 pb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 tracking-tight uppercase">Event Management</h1>
            <p className="text-sm text-gray-400 mt-2 font-medium tracking-wide">作成済みのイベントを管理します。</p>
          </div>
          <Link href="/events/new" className="px-6 py-2.5 bg-blue-600 text-white text-sm font-bold rounded-lg hover:bg-blue-700 transition-all">
            ＋ 新規作成
          </Link>
        </header>

        <div className="space-y-4">
          {events.length === 0 ? (
            <div className="bg-white p-24 rounded-2xl border border-gray-100 shadow-sm text-center">
              <p className="text-sm text-gray-400 italic mb-8">管理中のイベントはありません</p>
              <Link href="/events/new" className="text-xs font-bold text-blue-600 hover:underline">
                最初のイベントを作成する →
              </Link>
            </div>
          ) : (
            events.map((event) => (
              <div key={event.id} className="bg-white rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-all flex overflow-hidden">
                <div className={`w-2 h-auto ${event.status === 'published' ? 'bg-blue-600' : 'bg-gray-200'}`} />
                <div className="flex-1 p-8 flex flex-col md:flex-row md:items-center justify-between gap-6">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-3 mb-2">
                      <span className={`text-[9px] font-bold px-2 py-0.5 rounded uppercase ${event.status === 'published' ? 'bg-blue-50 text-blue-600 border border-blue-100' : 'bg-gray-50 text-gray-400 border border-gray-100'}`}>
                        {event.status === 'published' ? '公開中' : '下書き'}
                      </span>
                    </div>
                    <Link href={`/events/${event.id}`} className="block group-hover:text-blue-600 transition-colors">
                      <h3 className="text-xl font-bold text-gray-900 truncate">{event.title}</h3>
                    </Link>
                    <p className="mt-2 text-xs text-gray-500">{formatDate(event.start_time)}</p>
                  </div>
                  <div className="flex items-center gap-3">
                    <Link href={`/events/${event.id}`} className="px-4 py-2 bg-white text-gray-600 border border-gray-100 text-[15px] font-bold uppercase rounded-lg hover:bg-gray-50 transition-all shadow-sm">
                      閲覧
                    </Link>
                    <Link href={`/events/${event.id}/edit`} className="px-4 py-2 bg-blue-50 text-blue-600 text-[15px] font-bold uppercase rounded-lg hover:bg-blue-100 transition-all font-mono">
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
