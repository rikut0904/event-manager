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
      <main className="max-w-5xl mx-auto">
        <header className="mb-16 border-b border-gray-100 pb-8">
          <h1 className="text-3xl font-bold text-gray-900 tracking-tight uppercase tracking-[0.05em]">Home</h1>
          <p className="text-sm text-gray-400 mt-2 font-medium tracking-wide">
            ようこそ、<span className="text-gray-900 font-bold">{user.email}</span> さん
          </p>
        </header>

        <section className="space-y-8">
          <div className="flex items-center justify-between">
            <h2 className="text-sm font-bold text-gray-400 uppercase tracking-[0.2em]">Events</h2>
          </div>

          {isLoading ? (
            <div className="bg-white p-20 rounded-2xl border border-gray-100 shadow-sm text-center">
              <p className="text-xs text-gray-300 animate-pulse font-mono tracking-widest">LOADING_STREAM...</p>
            </div>
          ) : publishedEvents.length === 0 ? (
            <div className="bg-white p-20 rounded-2xl border border-gray-100 shadow-sm text-center text-sm italic text-gray-400">
              現在公開中のイベントはありません
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {publishedEvents.map((event) => (
                <Link 
                  key={event.id} 
                  href={`/events/${event.id}`}
                  className="bg-white p-8 rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-all group"
                >
                  <div className="flex items-center gap-2 mb-4">
                    <span className="w-1.5 h-1.5 rounded-full bg-blue-600"></span>
                    <span className="text-[10px] font-bold text-blue-600 uppercase tracking-widest">Published</span>
                  </div>
                  <h3 className="text-xl font-bold text-gray-900 group-hover:text-blue-600 transition-colors mb-4 truncate">{event.title}</h3>
                  <div className="flex items-center gap-4 text-xs text-gray-500 font-medium">
                    <div className="flex items-center gap-1.5">
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-12 0 9 9 0 0112 0z" /></svg>
                      {formatDate(event.start_time)}
                    </div>
                    {event.location && (
                      <div className="flex items-center gap-1.5 truncate">
                        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                        <span className="truncate">{event.location}</span>
                      </div>
                    )}
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
