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
  const [status, setStatus] = useState('draft');
  
  const [startDate, setStartDate] = useState('');
  const [startTime, setStartTime] = useState('');
  const [endDate, setEndDate] = useState('');
  const [endTime, setEndTime] = useState('');
  
  const [isEndTimeManual, setIsEndTimeManual] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isFetching, setIsFetching] = useState(true);
  const [error, setError] = useState('');
  const router = useRouter();

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

  const handleUpdateEvent = async (e: React.FormEvent, finalStatus: string) => {
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
          start_time: start.toISOString(),
          end_time: end ? end.toISOString() : null,
        }),
      });
      router.push('/events');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  if (isFetching) return null;

  const isPublished = status === 'published';
  const googleMapsUrl = location ? `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(location)}` : null;

  return (
    <div className="min-h-screen bg-gray-50/30 text-gray-600 font-sans p-8 lg:p-16">
      <main className="max-w-4xl mx-auto">
        <header className="mb-12 flex justify-between items-start">
          <div>
            <Link href="/events" className="text-xs font-bold text-gray-400 uppercase tracking-widest hover:text-blue-600 mb-4 block">← Back to List</Link>
            <h1 className="text-3xl font-bold text-gray-900 tracking-tight uppercase tracking-widest">
              {isPublished ? 'Update Event' : 'Draft Event'}
            </h1>
          </div>
          <Link href={`/events/${id}`} className="px-6 py-2 bg-white text-blue-600 border border-blue-600 text-xs font-bold uppercase rounded-lg hover:bg-blue-50 transition-all shadow-sm">View Page</Link>
        </header>

        {error && (
          <div className="mb-8 p-4 bg-red-50 border border-red-100 text-red-600 text-xs font-bold rounded-xl flex items-center gap-3">
            <span>⚠️</span>
            <span>{error}</span>
          </div>
        )}

        <div className="bg-white p-10 rounded-3xl border border-gray-100 shadow-sm">
          <form className="space-y-12">
            <div className="space-y-8">
              <div>
                <label className="block text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-3">Title</label>
                <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} className="w-full p-4 text-xl font-bold border border-gray-100 rounded-2xl focus:ring-1 focus:ring-blue-600 outline-none transition-all" required />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
                <div className="space-y-4">
                  <label className="block text-[10px] font-bold text-gray-400 uppercase tracking-widest">Start Date & Time</label>
                  <div className="flex gap-2">
                    <input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} className="flex-1 p-3 border border-gray-100 rounded-xl focus:ring-1 focus:ring-blue-600 outline-none" required />
                    <input type="time" value={startTime} onChange={(e) => setStartTime(e.target.value)} className="w-32 p-3 border border-gray-100 rounded-xl focus:ring-1 focus:ring-blue-600 outline-none" required />
                  </div>
                </div>

                <div className="space-y-4">
                  <label className="block text-[10px] font-bold text-gray-400 uppercase tracking-widest">End Date & Time</label>
                  <div className="flex gap-2">
                    <input type="date" value={endDate} onChange={(e) => { setEndDate(e.target.value); setIsEndTimeManual(true); }} className="flex-1 p-3 border border-gray-100 rounded-xl focus:ring-1 focus:ring-blue-600 outline-none" />
                    <input type="time" value={endTime} onChange={(e) => { setEndTime(e.target.value); setIsEndTimeManual(true); }} className="w-32 p-3 border border-gray-100 rounded-xl focus:ring-1 focus:ring-blue-600 outline-none" />
                  </div>
                </div>
              </div>

              <div>
                <label className="block text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-3">Location</label>
                <div className="relative">
                  <input type="text" value={location} onChange={(e) => setLocation(e.target.value)} className="w-full p-4 border border-gray-100 rounded-2xl focus:ring-1 focus:ring-blue-600 outline-none pr-32" placeholder="住所やURLなど" />
                  {googleMapsUrl && (
                    <a href={googleMapsUrl} target="_blank" rel="noopener noreferrer" className="absolute right-4 top-1/2 -translate-y-1/2 text-[9px] font-bold text-blue-600 bg-blue-50 px-3 py-1 rounded-md hover:bg-blue-100 transition-colors">
                      Google Maps
                    </a>
                  )}
                </div>
              </div>

              <div>
                <label className="block text-[10px] font-bold text-gray-400 uppercase tracking-widest mb-3">Description</label>
                <textarea value={description} onChange={(e) => setDescription(e.target.value)} className="w-full p-4 border border-gray-100 rounded-2xl focus:ring-1 focus:ring-blue-600 outline-none resize-y min-h-[200px] font-light leading-relaxed" />
              </div>
            </div>

            <div className="flex gap-4">
              {!isPublished && (
                <button
                  type="button"
                  onClick={(e) => handleUpdateEvent(e, 'draft')}
                  disabled={isLoading}
                  className="flex-1 py-5 bg-gray-50 text-gray-400 text-xs font-bold uppercase tracking-[0.2em] rounded-2xl hover:bg-gray-100 transition-all"
                >
                  下書きとして保存
                </button>
              )}
              <button
                type="button"
                onClick={(e) => handleUpdateEvent(e, 'published')}
                disabled={isLoading}
                className="flex-[2] py-5 bg-blue-600 text-white text-xs font-bold uppercase tracking-[0.2em] rounded-2xl hover:bg-blue-700 shadow-xl shadow-blue-50 transition-all disabled:bg-blue-300"
              >
                {isLoading ? 'Processing...' : isPublished ? '内容を更新する' : 'イベントを公開する'}
              </button>
            </div>
          </form>

          <div className="mt-24 pt-16 border-t border-gray-100">
            <div className="relative overflow-hidden bg-gradient-to-br from-red-50 to-white p-10 rounded-[2.5rem] border border-red-100/50 shadow-sm flex flex-col md:flex-row md:items-center justify-between gap-10 group">
              {/* Decorative background pattern */}
              <div className="absolute -right-10 -bottom-10 w-40 h-40 bg-red-100/30 rounded-full blur-3xl transition-all group-hover:bg-red-200/40"></div>
              
              <div className="relative z-10 max-w-xl">
                <div className="flex items-center gap-3 mb-4">
                  <div className="w-8 h-8 rounded-full bg-red-500 flex items-center justify-center text-white shadow-lg shadow-red-200">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                  </div>
                  <h3 className="text-red-600 text-[11px] font-black uppercase tracking-[0.4em]">Danger Zone</h3>
                </div>
                <h4 className="text-xl font-bold text-gray-900 mb-2 tracking-tight">イベントの削除</h4>
                <p className="text-sm text-gray-500 font-medium leading-relaxed">
                  この操作は取り消せません。
                  <br />このイベントを削除すると、参加者リスト、統計データ、
                  <br />およびすべての関連コンテンツが<span className="text-red-500 font-bold italic">永久に消去</span>されます。
                </p>
              </div>
              
              <button 
                onClick={() => {
                  if(window.confirm('本当に削除しますか？\nこの操作は元に戻せません。')) {
                    apiRequest(`/api/v1/events/${id}`, { method: 'DELETE' }).then(() => router.push('/events'));
                  }
                }} 
                className="relative z-10 px-10 py-4 bg-white text-red-600 border-2 border-red-100 rounded-2xl text-xs font-black uppercase tracking-widest hover:bg-red-600 hover:text-white hover:border-red-600 transition-all duration-300 shadow-xl shadow-red-100/20 active:scale-95 whitespace-nowrap"
              >
                Delete Event
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
