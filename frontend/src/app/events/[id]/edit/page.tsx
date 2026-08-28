'use client';

import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { apiRequest } from '@/lib/api';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';

export default function EditEventPage() {
  const { id } = useParams();
  const { user } = useAuth();
  
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [location, setLocation] = useState('');
  const [isOnline, setIsOnline] = useState(false);
  const [capacity, setCapacity] = useState(0);
  const [sourceURL, setSourceURL] = useState('');
  const [thumbnailURL, setThumbnailURL] = useState('');
  const [status, setStatus] = useState('draft');
  
  const [startDate, setStartDate] = useState('');
  const [startTime, setStartTime] = useState('');
  const [endDate, setEndDate] = useState('');
  const [endTime, setEndTime] = useState('');
  
  const [isEndTimeManual, setIsEndTimeManual] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isFetching, setIsFetching] = useState(true);
  const [saveSuccess, setSaveSuccess] = useState(false);
  const [error, setError] = useState('');
  const router = useRouter();

  useEffect(() => {
    if (saveSuccess) {
      const timer = setTimeout(() => setSaveSuccess(false), 3000);
      return () => clearTimeout(timer);
    }
  }, [saveSuccess]);

  useEffect(() => {
    const fetchEvent = async () => {
      try {
        const data = await apiRequest(`/api/v1/events/${id}`);
        if (user && data.creator_id !== user.id) {
          router.replace('/events');
          return;
        }
        setTitle(data.title);
        setDescription(data.description || '');
        setLocation(data.location || '');
        setIsOnline(data.is_online);
        setCapacity(data.capacity || 0);
        setSourceURL(data.source_url || '');
        setThumbnailURL(data.thumbnail_url || '');
        setStatus(data.status);
        
        if (data.start_time) {
          const dt = new Date(data.start_time);
          setStartDate(dt.toISOString().split('T')[0]);
          setStartTime(dt.toTimeString().slice(0, 5));
        }
        if (data.end_time) {
          const dt = new Date(data.end_time);
          setEndDate(dt.toISOString().split('T')[0]);
          setEndTime(dt.toTimeString().slice(0, 5));
          setIsEndTimeManual(true);
        }
      } catch (err) {
        console.error('Failed to fetch event:', err);
        router.replace('/events');
      } finally {
        setIsFetching(false);
      }
    };
    if (user && id) fetchEvent();
  }, [id, user, router]);

  useEffect(() => {
    if (!startDate || !startTime || isEndTimeManual || isFetching) return;
    const start = new Date(`${startDate}T${startTime}`);
    if (isNaN(start.getTime())) return;
    const end = new Date(start.getTime() + 60 * 60 * 1000);
    setEndDate(end.toISOString().split('T')[0]);
    setEndTime(end.toTimeString().slice(0, 5));
  }, [startDate, startTime, isEndTimeManual, isFetching]);

  const handleUpdateEvent = async (e: React.SyntheticEvent, finalStatus: string) => {
    e.preventDefault();
    setError('');

    const start = new Date(`${startDate}T${startTime}`);
    const end = endDate && endTime ? new Date(`${endDate}T${endTime}`) : null;

    if (end && end <= start) {
      setError('終了時間は開始時間より後である必要があります。');
      return;
    }

    setIsLoading(true);
    try {
      await apiRequest(`/api/v1/events/${id}`, {
        method: 'PATCH',
        body: JSON.stringify({
          title, description, location,
          status: finalStatus,
          is_online: isOnline,
          capacity: Number(capacity),
          source_url: sourceURL,
          thumbnail_url: thumbnailURL,
          start_time: start.toISOString(),
          end_time: end ? end.toISOString() : null,
        }),
      });
      setSaveSuccess(true);
      if (status !== finalStatus) setStatus(finalStatus);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  if (isFetching) return null;

  const isPublished = status === 'published';

  return (
    <div className="min-h-screen bg-gray-50/30 text-gray-900 font-sans pb-32">
      {/* Fixed Top Actions */}
      <div className="sticky top-0 z-30 bg-white/90 backdrop-blur-md border-b border-gray-100 px-8 h-16 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link href="/events" className="p-2 hover:bg-gray-100 rounded-full transition-colors text-gray-400">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M15 19l-7-7 7-7" /></svg>
          </Link>
          <h1 className="text-sm font-black tracking-[0.2em] text-gray-500">
            {isPublished ? 'イベントを更新' : '下書きを編集'}
          </h1>
        </div>
        <div className="flex gap-3 items-center">
          {saveSuccess && (
            <span className="text-sm font-bold text-green-500 tracking-widest animate-pulse mr-2">
              ✓ 保存しました
            </span>
          )}
          <Link href={`/events/${id}`} className="px-5 py-2 text-sm font-black tracking-widest text-gray-500 hover:text-gray-900 transition-colors">
            プレビュー
          </Link>
          <button
            onClick={(e) => handleUpdateEvent(e, isPublished ? 'published' : 'draft')}
            disabled={isLoading}
            className="px-6 py-2.5 bg-white text-gray-600 text-sm font-black tracking-widest rounded-full border border-gray-200 hover:bg-gray-50 transition-all disabled:opacity-50"
          >
            {isLoading ? '保存中...' : '変更を保存'}
          </button>
          {!isPublished && (
            <button
              onClick={(e) => handleUpdateEvent(e, 'published')}
              disabled={isLoading}
              className="px-8 py-2.5 bg-blue-600 text-white text-sm font-black tracking-widest rounded-full hover:bg-blue-700 transition-all shadow-xl shadow-blue-500/20"
            >
              今すぐ公開
            </button>
          )}
        </div>
      </div>

      <main className="max-w-[1400px] mx-auto px-8 pt-12">
        {error && (
          <div className="mb-12 p-4 bg-red-50 border border-red-100 text-red-600 text-sm font-bold tracking-widest rounded-xl flex items-center gap-3">
            <span>⚠️</span>
            <span>{error}</span>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-12">
          {/* Left Column: Editor (2/3) */}
          <div className="lg:col-span-2 space-y-12">
            <div className="bg-white p-10 rounded-[2.5rem] border border-gray-100 shadow-sm min-h-[700px] flex flex-col">
              <label className="block text-base font-black text-gray-500 tracking-[0.2em] mb-6 ml-1">イベントの説明</label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="flex-1 w-full p-8 bg-gray-50/30 border border-gray-100 rounded-[2rem] focus:ring-2 focus:ring-blue-600 focus:bg-white outline-none transition-all resize-none font-light leading-relaxed text-gray-600 text-lg"
                placeholder="詳細情報を記入..."
              />
            </div>
          </div>

          {/* Right Column: Overview & Settings (1/3) */}
          <div className="space-y-8">
            {/* 1. Overview Section: Title, Format, Capacity */}
            <div className="bg-white p-8 rounded-[2.5rem] border border-gray-100 shadow-sm space-y-8">
              <div>
                <label className="block text-sm font-black text-gray-400 tracking-[0.2em] mb-6 ml-1 flex items-center gap-2">
                  <span className="w-1.5 h-1.5 rounded-full bg-blue-600"></span> 基本情報
                </label>
                <div className="space-y-6">
                  <div>
                    <label className="block text-sm font-black text-gray-500 tracking-widest mb-3 ml-1">タイトル</label>
                    <input
                      type="text"
                      value={title}
                      onChange={(e) => setTitle(e.target.value)}
                      className="w-full text-xl font-black text-gray-900 bg-gray-50/50 border border-gray-100 rounded-xl p-4 focus:ring-2 focus:ring-blue-600 focus:bg-white outline-none transition-all placeholder:text-gray-200"
                      placeholder="タイトルを入力"
                      required
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-black text-gray-500 tracking-widest mb-3 ml-1">開催形式</label>
                      <div className="grid grid-cols-2 p-1 bg-gray-50 border border-gray-100 rounded-xl">
                        <button type="button" onClick={() => setIsOnline(false)} className={`py-2 rounded-lg text-[9px] font-black tracking-tighter transition-all ${!isOnline ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-400 hover:text-gray-900'}`}>会場</button>
                        <button type="button" onClick={() => setIsOnline(true)} className={`py-2 rounded-lg text-[9px] font-black tracking-tighter transition-all ${isOnline ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-400 hover:text-gray-900'}`}>オンライン</button>
                      </div>
                    </div>
                    <div>
                      <label className="block text-sm font-black text-gray-500 tracking-widest mb-3 ml-1">定員</label>
                      <div className="flex items-center gap-2">
                        <input
                          type="number"
                          min="0"
                          value={capacity}
                          onChange={(e) => setCapacity(parseInt(e.target.value) || 0)}
                          className="w-full p-2.5 text-lg font-bold bg-gray-50/50 border border-gray-100 rounded-lg focus:ring-2 focus:ring-blue-600 focus:bg-white outline-none transition-all font-mono"
                        />
                        <span className="text-sm font-bold text-gray-500 tracking-tighter">名</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* 2. Logistics Section: Location & Dates */}
            <div className="bg-white p-8 rounded-[2.5rem] border border-gray-100 shadow-sm space-y-8">
              <div>
                <label className="block text-sm font-black text-gray-400 tracking-[0.2em] mb-6 ml-1 flex items-center gap-2">
                  <span className="w-1.5 h-1.5 rounded-full bg-gray-300"></span> 開催情報
                </label>
                <div className="space-y-6">
                  <div>
                    <label className="block text-sm font-black text-gray-500 tracking-widest mb-3 ml-1">
                      {isOnline ? '配信URL' : '開催場所'}
                    </label>
                    <input
                      type="text"
                      value={location}
                      onChange={(e) => setLocation(e.target.value)}
                      className="w-full p-4 bg-gray-50/50 border border-gray-100 rounded-xl focus:ring-2 focus:ring-blue-600 focus:bg-white outline-none transition-all text-sm"
                      placeholder={isOnline ? 'URLを入力...' : '開催場所を入力...'}
                    />
                  </div>

                  <div className="space-y-4 pt-2">
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-black text-gray-500 tracking-widest mb-2 ml-1">開始日</label>
                        <input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} className="w-full p-3 bg-gray-50/50 border border-gray-100 rounded-lg focus:ring-2 focus:ring-blue-600 focus:bg-white outline-none font-mono text-sm" />
                      </div>
                      <div>
                        <label className="block text-sm font-black text-gray-500 tracking-widest mb-2 ml-1">開始時刻</label>
                        <input type="time" value={startTime} onChange={(e) => setStartTime(e.target.value)} className="w-full p-3 bg-gray-50/50 border border-gray-100 rounded-lg focus:ring-2 focus:ring-blue-600 focus:bg-white outline-none font-mono text-sm" />
                      </div>
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-black text-gray-500 tracking-widest mb-2 ml-1">終了日</label>
                        <input type="date" value={endDate} onChange={(e) => { setEndDate(e.target.value); setIsEndTimeManual(true); }} className="w-full p-3 bg-gray-50/50 border border-gray-100 rounded-lg focus:ring-2 focus:ring-blue-600 focus:bg-white outline-none font-mono text-sm" />
                      </div>
                      <div>
                        <label className="block text-sm font-black text-gray-500 tracking-widest mb-2 ml-1">終了時刻</label>
                        <input type="time" value={endTime} onChange={(e) => { setEndTime(e.target.value); setIsEndTimeManual(true); }} className="w-full p-3 bg-gray-50/50 border border-gray-100 rounded-lg focus:ring-2 focus:ring-blue-600 focus:bg-white outline-none font-mono text-sm" />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* 3. Media & Links Section: Thumbnail & External URL */}
            <div className="bg-white p-8 rounded-[2.5rem] border border-gray-100 shadow-sm space-y-8">
              <div>
                <label className="block text-sm font-black text-gray-400 tracking-[0.2em] mb-6 ml-1 flex items-center gap-2">
                <span className="w-1.5 h-1.5 rounded-full bg-gray-300"></span> 画像・リンク
                </label>
                <div className="space-y-6">
                  <div>
                    <label className="block text-sm font-black text-gray-500 tracking-widest mb-3 ml-1">サムネイルプレビュー</label>
                    <div className="aspect-video rounded-xl overflow-hidden bg-gray-50 border border-gray-100 relative mb-4">
                      {thumbnailURL ? (
                        <img src={thumbnailURL} alt="プレビュー" className="w-full h-full object-cover" onError={(e) => (e.currentTarget.style.display = 'none')} />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center text-sm font-bold text-gray-400 tracking-widest text-center px-8">サムネイルなし</div>
                      )}
                    </div>
                    <input
                      type="url"
                      value={thumbnailURL}
                      onChange={(e) => setThumbnailURL(e.target.value)}
                      className="w-full p-3 text-sm bg-gray-50/50 border border-gray-100 rounded-lg focus:ring-2 focus:ring-blue-600 focus:bg-white outline-none transition-all font-mono"
                      placeholder="画像URLを入力..."
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-black text-gray-500 tracking-widest mb-3 ml-1">外部リンク</label>
                    <input
                      type="url"
                      value={sourceURL}
                      onChange={(e) => setSourceURL(e.target.value)}
                      className="w-full p-3 text-sm bg-gray-50/50 border border-gray-100 rounded-lg focus:ring-2 focus:ring-blue-600 focus:bg-white outline-none transition-all font-mono"
                      placeholder="https://..."
                    />
                  </div>
                </div>
              </div>
            </div>

            {/* Danger Zone */}
            <div className="pt-12">
              <button 
                onClick={() => { if(window.confirm('本当に削除しますか？\nこの操作は元に戻せません。')) apiRequest(`/api/v1/events/${id}`, { method: 'DELETE' }).then(() => router.push('/events')); }} 
                className="w-full py-4 bg-white text-red-600 border border-red-100 rounded-2xl text-sm font-black tracking-widest hover:bg-red-600 hover:text-white transition-all shadow-sm active:scale-95"
              >
                イベントを完全に削除
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
