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

  if (isLoading) return <div className="min-h-screen bg-white flex items-center justify-center font-mono text-base tracking-widest text-gray-400 animate-pulse">読み込み中...</div>;

  if (error || !event) return (
    <div className="min-h-screen bg-white flex flex-col items-center justify-center p-8 text-center">
      <h1 className="text-2xl font-bold text-gray-900 mb-4">404 - ページが見つかりません</h1>
      <p className="text-gray-500 mb-8 max-w-xs leading-relaxed text-sm">{error}</p>
      <Link href="/" className="px-6 py-2 border border-gray-200 text-gray-500 rounded-full text-xs font-bold tracking-widest hover:bg-gray-50 transition-colors">ホームへ戻る</Link>
    </div>
  );

  const isOwner = user && user.id === event.creator_id;
  const googleMapsUrl = event.location ? `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(event.location)}` : null;
  const mapEmbedUrl = !event.is_online && event.location ? `https://maps.google.com/maps?q=${encodeURIComponent(event.location)}&output=embed` : null;

  return (
    <div className="min-h-screen bg-white text-gray-900 font-sans pb-40">
      {/* Hero Section with Thumbnail */}
      <div className="relative w-full h-[40vh] lg:h-[60vh] bg-gray-900 overflow-hidden">
        {event.thumbnail_url ? (
          <img src={event.thumbnail_url} alt={event.title} className="w-full h-full object-cover opacity-60" />
        ) : (
          <div className="w-full h-full bg-gradient-to-br from-blue-900 to-black opacity-60" />
        )}
        <div className="absolute inset-0 bg-gradient-to-t from-white via-transparent to-transparent" />
        
        <header className="absolute top-0 w-full h-16 flex items-center px-8 justify-between z-20">
          <Link href="/" className="font-bold tracking-tight text-lg text-white mix-blend-difference">イベント管理</Link>
          {isOwner && (
            <Link href={`/events/${event.id}/edit`} className="text-sm font-bold text-blue-600 bg-white/90 backdrop-blur-sm px-5 py-2.5 rounded-xl hover:bg-white transition-all tracking-widest shadow-xl">
              イベントを編集
            </Link>
          )}
        </header>

        <div className="absolute bottom-0 left-0 w-full p-8 lg:p-20 z-10">
          <div className="max-w-4xl mx-auto">
            <div className="flex flex-wrap items-center gap-3 mb-6">
              <span className={`px-3 py-1.5 rounded-full text-sm font-black tracking-widest shadow-sm ${event.is_online ? 'bg-green-500 text-white' : 'bg-blue-600 text-white'}`}>
                {event.is_online ? 'オンライン' : '会場開催'}
              </span>
              <time className="text-sm font-bold text-gray-900 tracking-[0.1em] bg-white/80 backdrop-blur-sm px-3 py-1.5 rounded-full">
                {new Date(event.start_time).toLocaleString('ja-JP', { dateStyle: 'full', timeStyle: 'short' })}
              </time>
              {event.capacity > 0 && (
                <span className="text-sm font-bold text-gray-900 tracking-[0.1em] bg-white/80 backdrop-blur-sm px-3 py-1.5 rounded-full">
                  定員: {event.capacity}名
                </span>
              )}
            </div>
            <h1 className="max-w-4xl text-4xl lg:text-7xl font-black text-gray-950 tracking-tighter leading-tight drop-shadow-sm">
              <span className="box-decoration-clone bg-white/85 px-3 py-2 rounded-xl">{event.title}</span>
            </h1>
          </div>
        </div>
      </div>

      <main className="max-w-4xl mx-auto px-8 pt-16">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-16">
          <div className="lg:col-span-2 space-y-16">
            {/* Description */}
            <section>
              <h2 className="text-sm font-bold text-gray-400 tracking-[0.2em] mb-8 flex items-center gap-2">
                <span className="w-4 h-px bg-gray-200"></span> イベント概要
              </h2>
              <div className="whitespace-pre-wrap text-lg text-gray-600 leading-[1.8] font-light tracking-wide">
                {event.description || 'このイベントには詳細な説明がありません。'}
              </div>
            </section>

            {/* Map (Venue Only) */}
            {!event.is_online && mapEmbedUrl && (
              <section>
                <h2 className="text-sm font-bold text-gray-400 tracking-[0.2em] mb-8 flex items-center gap-2">
                  <span className="w-4 h-px bg-gray-200"></span> 会場・アクセス
                </h2>
                <div className="rounded-[2.5rem] overflow-hidden border border-gray-100 shadow-sm aspect-video relative group">
                  <iframe
                    width="100%"
                    height="100%"
                    title="イベント会場の地図"
                    style={{ border: 0 }}
                    loading="lazy"
                    allowFullScreen
                    src={mapEmbedUrl}
                  />
                  <div className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
                    <a href={googleMapsUrl || '#'} target="_blank" rel="noopener noreferrer" className="px-4 py-2 bg-white/90 backdrop-blur shadow-xl rounded-xl text-sm font-bold tracking-widest text-gray-900 hover:bg-white">
                      Google Maps で開く
                    </a>
                  </div>
                </div>
                <p className="mt-4 text-sm font-medium text-gray-500 flex items-center gap-2">
                  <svg className="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                  {event.location}
                </p>
              </section>
            )}

            {/* Online Info */}
            {event.is_online && (
              <section className="p-10 bg-green-50/50 border border-green-100 rounded-[2.5rem]">
                <h2 className="text-green-600 text-sm font-bold tracking-[0.2em] mb-4">オンラインイベント</h2>
                <p className="text-gray-600 text-sm mb-6">このイベントはオンラインで開催されます。配信先や参加URLは以下の通りです。</p>
                <div className="bg-white p-4 rounded-2xl border border-green-100 flex items-center justify-between">
                  <span className="text-xs font-mono text-gray-500 truncate mr-4">{event.location}</span>
                  <a href={event.location.startsWith('http') ? event.location : '#'} target="_blank" rel="noopener noreferrer" className="px-6 py-2 bg-green-600 text-white text-[10px] font-bold uppercase tracking-widest rounded-xl hover:bg-green-700 transition-colors shrink-0">
                    URLを開く
                  </a>
                </div>
              </section>
            )}
          </div>

          <aside className="space-y-8">
            {/* Host Info */}
            <div className="p-8 border border-gray-100 rounded-[2rem] bg-gray-50/30">
              <p className="text-sm font-bold text-gray-400 tracking-widest mb-4">主催者</p>
              <div className="flex items-center gap-4">
                <div className="w-10 h-10 rounded-full bg-gray-900 flex items-center justify-center text-[10px] font-bold text-white">
                  {event.creator_id.slice(0, 2).toUpperCase()}
                </div>
                <div>
                  <p className="text-sm font-bold text-gray-700 leading-none mb-1">ユーザー #{event.creator_id.slice(0, 8)}</p>
                  <p className="text-sm text-gray-400 font-medium tracking-tight">コミュニティメンバー</p>
                </div>
              </div>
            </div>
          </aside>
        </div>
      </main>

      {/* Fixed Bottom Registration Bar */}
      <div className="fixed bottom-0 left-0 w-full bg-white/80 backdrop-blur-xl border-t border-gray-100 p-6 z-40">
        <div className="max-w-4xl mx-auto flex items-center justify-between gap-8">
          <div className="hidden md:block">
            <p className="text-sm font-bold text-gray-500 tracking-widest mb-1">参加の準備はできましたか？</p>
            <h4 className="text-base font-bold text-gray-900 line-clamp-1">{event.title}</h4>
          </div>
          
          <div className="flex-1 md:flex-none">
            {event.source_url ? (
              <a 
                href={event.source_url} 
                target="_blank" 
                rel="noopener noreferrer"
                className="block w-full md:w-64 py-4 bg-gray-900 text-white text-center font-black rounded-2xl text-sm tracking-widest hover:bg-black transition-all shadow-xl shadow-gray-200 active:scale-95"
              >
                申し込みページへ移動
              </a>
            ) : (
                <button disabled className="w-full md:w-64 py-4 bg-gray-100 text-gray-400 font-bold rounded-2xl cursor-not-allowed text-sm tracking-widest">
                受付終了
              </button>
            )}
          </div>
        </div>
      </div>

      <footer className="mt-32 pt-20 border-t border-gray-50 text-center">
        <p className="text-[10px] text-gray-200 font-bold tracking-[0.4em]">イベント管理プラットフォーム</p>
      </footer>
    </div>
  );
}
