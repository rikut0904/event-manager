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

  if (!user) return null;

  return (
    <div className="min-h-screen bg-gray-50/30 text-gray-600 font-sans">
      <main className="max-w-5xl mx-auto px-8 py-16">
        <header className="mb-16 flex justify-between items-end border-b border-gray-100 pb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 tracking-tight uppercase">My Events</h1>
            <p className="text-sm text-gray-400 mt-2 font-medium tracking-wide">作成済みのイベントを管理します。</p>
          </div>
          <Link
            href="/dashboard/events/new"
            className="px-6 py-2.5 bg-blue-600 text-white text-sm font-bold rounded-lg hover:bg-blue-700 transition-all shadow-md shadow-blue-50"
          >
            新規イベント作成
          </Link>
        </header>

        <div className="space-y-6">
          {events.length === 0 ? (
            <div className="bg-white p-24 rounded-2xl border border-gray-100 shadow-sm text-center">
              <p className="text-sm text-gray-400 tracking-wider italic mb-8">管理中のイベントはありません</p>
              <Link
                href="/dashboard/events/new"
                className="text-xs font-bold text-blue-600 uppercase tracking-widest hover:underline"
              >
                最初のイベントを作成する →
              </Link>
            </div>
          ) : (
            <div className="grid grid-cols-1 gap-4">
              {events.map((event) => (
                <div 
                  key={event.id} 
                  className="bg-white rounded-xl border border-gray-100 shadow-sm hover:shadow-md transition-all group overflow-hidden flex"
                >
                  {/* Status Indicator Color on the left */}
                  <div className={`w-1.5 h-full ${
                    event.status === 'published' ? 'bg-blue-600' : 'bg-gray-200'
                  }`} />
                  
                  <div className="flex-1 p-6 flex items-center justify-between">
                    <div>
                      <div className="flex items-center gap-3 mb-1">
                        <h3 className="text-base font-bold text-gray-900 group-hover:text-blue-600 transition-colors">
                          {event.title}
                        </h3>
                        <span className={`text-[9px] font-bold px-2 py-0.5 rounded uppercase tracking-tighter ${
                          event.status === 'published' ? 'bg-blue-50 text-blue-600' : 'bg-gray-100 text-gray-400'
                        }`}>
                          {event.status === 'published' ? '公開中' : '下書き'}
                        </span>
                      </div>
                      <p className="text-xs text-gray-400 font-mono italic">/{event.display_id}</p>
                    </div>
                    
                    <div className="flex items-center gap-4">
                      <span className="text-[10px] font-bold text-gray-300 uppercase tracking-widest group-hover:text-gray-400 transition-colors">
                        Edit Settings
                      </span>
                      <svg className="w-4 h-4 text-gray-300 group-hover:text-blue-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                      </svg>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
