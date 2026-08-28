'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { apiRequest } from '@/lib/api';
import Link from 'next/link';

export default function HomePage() {
  const { user } = useAuth();
  const [publishedEvents, setPublishedEvents] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchPublishedEvents = async () => {
      try {
        const data = await apiRequest('/api/v1/events/published');
        setPublishedEvents(data || []);
      } catch (err) {
        console.error('Failed to fetch published events:', err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchPublishedEvents();
  }, []);

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleString('ja-JP', {
      month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
    });
  };

  if (!user) return null;

  return (
    <div className="min-h-screen bg-gray-50/30 text-gray-600 font-sans p-8 lg:p-16">
      <main className="max-w-6xl mx-auto">
        <section className="space-y-8">
          <div className="flex items-center justify-between">
            <h2 className="text-sm font-bold text-gray-400 tracking-[0.2em]">公開中のイベント</h2>
          </div>

          {isLoading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              {[1, 2, 3].map(i => (
                <div key={i} className="bg-white aspect-video rounded-3xl border border-gray-100 animate-pulse" />
              ))}
            </div>
          ) : publishedEvents.length === 0 ? (
            <div className="bg-white p-24 rounded-[3rem] border border-gray-100 shadow-sm text-center text-sm italic text-gray-400">
              現在公開中のイベントはありません
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              {publishedEvents.map((event) => (
                <Link 
                  key={event.id} 
                  href={`/events/${event.id}`}
                  className="bg-white rounded-3xl border border-gray-100 shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all overflow-hidden group"
                >
                  <div className="aspect-video bg-gray-100 relative overflow-hidden">
                    {event.thumbnail_url ? (
                      <img src={event.thumbnail_url} alt={event.title} className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
                    ) : (
                      <div className="w-full h-full bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center">
                        <span className="text-[10px] font-bold text-gray-300 tracking-widest">画像なし</span>
                      </div>
                    )}
                    <div className="absolute top-4 left-4 flex gap-2">
                      <span className={`px-2.5 py-1 rounded text-xs font-black tracking-widest ${event.is_online ? 'bg-green-500 text-white' : 'bg-blue-600 text-white'}`}>
                        {event.is_online ? 'オンライン' : '会場開催'}
                      </span>
                    </div>
                  </div>
                  <div className="p-6">
                    <h3 className="text-xl font-bold text-gray-900 group-hover:text-blue-600 transition-colors mb-4 line-clamp-2">{event.title}</h3>
                    <div className="flex items-center justify-between text-sm font-bold">
                      <div className="flex items-center gap-1.5 text-gray-500 tracking-widest">
                        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-12 0 9 9 0 0112 0z" /></svg>
                        {formatDate(event.start_time)}
                      </div>
                      {event.capacity > 0 && (
                        <div className="text-blue-600/50">
                          定員 {event.capacity}名
                        </div>
                      )}
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}
