'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { apiRequest } from '@/lib/api';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';

export default function EventViewPage() {
  const { id } = useParams();
  const { user } = useAuth();
  const [event, setEvent] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchEvent = async () => {
      try {
        const data = await apiRequest(`/api/v1/events/view/${id}`);
        setEvent(data);
      } catch (err: any) {
        setError('イベントが見つからないか、非公開になっています。');
      } finally {
        setIsLoading(false);
      }
    };
    if (id) fetchEvent();
  }, [id]);

  if (isLoading) return <div className="min-h-screen bg-white flex items-center justify-center font-mono text-xs uppercase tracking-widest text-gray-400 animate-pulse">Fetching...</div>;

  if (error || !event) return (
    <div className="min-h-screen bg-white flex flex-col items-center justify-center p-8 text-center">
      <h1 className="text-2xl font-bold text-gray-900 mb-4">404 - Not Found</h1>
      <p className="text-gray-500 mb-8 max-w-xs leading-relaxed text-sm">{error}</p>
      <Link href="/" className="px-6 py-2 border border-gray-200 text-gray-500 rounded-full text-xs font-bold uppercase tracking-widest hover:bg-gray-50 transition-colors">Back to Home</Link>
    </div>
  );

  const isOwner = user && user.id === event.creator_id;
  const googleMapsUrl = event.location ? `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(event.location)}` : null;

  return (
    <div className="min-h-screen bg-white text-gray-900 font-sans pb-24">
      <header className="border-b border-gray-50 h-16 flex items-center px-8 justify-between sticky top-0 bg-white/90 backdrop-blur-sm z-10">
        <Link href="/" className="font-bold tracking-tight text-lg">Event Manager</Link>
        {isOwner && (
          <Link href={`/events/${event.id}/edit`} className="text-[10px] font-bold text-blue-600 bg-blue-50 px-4 py-2 rounded-lg hover:bg-blue-100 transition-colors uppercase tracking-widest">
            Edit Event
          </Link>
        )}
      </header>

      <main className="max-w-3xl mx-auto px-8 pt-20">
        <div className="mb-20">
          <div className="flex items-center gap-3 mb-6">
            <time className="text-[10px] font-bold text-blue-600 uppercase tracking-[0.2em] bg-blue-50 px-2 py-1 rounded">
              {new Date(event.start_time).toLocaleString('ja-JP', { dateStyle: 'full', timeStyle: 'short' })}
            </time>
            {event.status === 'draft' && (
              <span className="text-[10px] font-bold bg-gray-100 text-gray-400 px-2 py-1 rounded uppercase tracking-widest">Preview Mode</span>
            )}
          </div>
          <h1 className="text-4xl lg:text-6xl font-extrabold text-gray-900 mb-10 tracking-tight leading-[1.1]">{event.title}</h1>
          
          <div className="flex items-center gap-4 border-t border-gray-50 pt-10">
            <div className="w-12 h-12 rounded-full bg-gray-900 flex items-center justify-center text-xs font-bold text-white">
              {event.creator_id.slice(0, 2).toUpperCase()}
            </div>
            <div>
              <p className="text-[10px] text-gray-400 font-bold uppercase tracking-widest mb-0.5">Host</p>
              <p className="text-sm font-bold text-gray-700">User #{event.creator_id.slice(0, 8)}</p>
            </div>
          </div>
        </div>

        {event.location && (
          <div className="mb-16 p-8 bg-gray-50 rounded-[2rem] border border-gray-100 flex flex-col sm:flex-row sm:items-center justify-between gap-6">
            <div className="flex items-start gap-4">
              <div className="w-10 h-10 rounded-full bg-white flex items-center justify-center shadow-sm text-gray-400 shrink-0">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
              </div>
              <div>
                <p className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-1">Location</p>
                <p className="text-base font-medium text-gray-800">{event.location}</p>
              </div>
            </div>
            {googleMapsUrl && (
              <a href={googleMapsUrl} target="_blank" rel="noopener noreferrer" className="px-6 py-3 bg-white text-gray-900 border border-gray-200 rounded-2xl text-[10px] font-bold uppercase tracking-widest hover:bg-gray-900 hover:text-white hover:border-gray-900 transition-all shadow-sm text-center">
                Google Maps で見る
              </a>
            )}
          </div>
        )}

        <div className="border-t border-gray-50 pt-16">
          <h2 className="text-[10px] font-bold text-gray-300 uppercase tracking-[0.3em] mb-10">Event Description</h2>
          <div className="whitespace-pre-wrap text-lg text-gray-600 leading-[1.8] font-light tracking-wide">
            {event.description || 'このイベントには詳細な説明がありません。'}
          </div>
        </div>

        <div className="mt-32 p-16 bg-gray-900 rounded-[3rem] text-center text-white overflow-hidden relative group">
          <div className="relative z-10">
            <h3 className="text-2xl font-bold mb-6">参加申し込み</h3>
            <p className="text-gray-400 mb-10 text-sm max-w-sm mx-auto leading-relaxed">
              このイベントへの参加を希望される方は、受付開始をお待ちください。
            </p>
            <button disabled className="px-12 py-4 bg-white/10 text-gray-500 font-bold rounded-2xl cursor-not-allowed uppercase text-xs tracking-widest">
              Registration Opening Soon
            </button>
          </div>
          {/* Decorative background element */}
          <div className="absolute top-0 right-0 w-64 h-64 bg-blue-600 rounded-full blur-[120px] opacity-20 -mr-32 -mt-32 transition-all group-hover:opacity-30"></div>
        </div>
      </main>

      <footer className="mt-24 pt-20 border-t border-gray-50 text-center">
        <p className="text-[10px] text-gray-300 font-bold uppercase tracking-[0.4em]">Event Manager</p>
      </footer>
    </div>
  );
}
