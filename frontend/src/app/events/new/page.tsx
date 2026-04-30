'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '@/lib/api';
import Link from 'next/link';

export default function NewEventPage() {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [location, setLocation] = useState('');
  const [startDate, setStartDate] = useState('');
  const [startTime, setStartTime] = useState('');
  const [endDate, setEndDate] = useState('');
  const [endTime, setEndTime] = useState('');
  
  const [isEndTimeManual, setIsEndTimeManual] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const router = useRouter();

  const today = new Date().toISOString().split('T')[0];

  useEffect(() => {
    if (isEndTimeManual || !startDate || !startTime) return;
    try {
      const start = new Date(`${startDate}T${startTime}`);
      if (isNaN(start.getTime())) return;
      const end = new Date(start.getTime() + 60 * 60 * 1000);
      const endYear = end.getFullYear();
      const endMonth = String(end.getMonth() + 1).padStart(2, '0');
      const endDateStr = String(end.getDate()).padStart(2, '0');
      const endHours = String(end.getHours()).padStart(2, '0');
      const endMins = String(end.getMinutes()).padStart(2, '0');
      setEndDate(`${endYear}-${endMonth}-${endDateStr}`);
      setEndTime(`${endHours}:${endMins}`);
    } catch (e) {}
  }, [startDate, startTime, isEndTimeManual]);

  const handleCreateEvent = async (e: React.FormEvent, status: string) => {
    e.preventDefault();
    setError('');

    if (!title.trim()) {
      setError('イベント名は必須です。');
      return;
    }

    const start = new Date(`${startDate}T${startTime}`);
    const now = new Date();

    if (start < now) {
      setError('開始時間は現在時刻より後である必要があります。');
      return;
    }

    const end = endDate && endTime ? new Date(`${endDate}T${endTime}`) : null;

    if (end && end <= start) {
      setError('終了時間は開始時間より後である必要があります。');
      return;
    }

    setIsLoading(true);
    try {
      await apiRequest('/api/v1/events', {
        method: 'POST',
        body: JSON.stringify({
          title, description, location,
          status,
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

  const googleMapsUrl = location ? `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(location)}` : null;

  return (
    <div className="min-h-screen bg-gray-50/30 text-gray-600 font-sans p-8 lg:p-16">
      <main className="max-w-3xl mx-auto">
        <header className="mb-12">
          <Link href="/events" className="text-xs font-bold text-gray-400 uppercase tracking-widest hover:text-blue-600 transition-colors mb-4 block">← Back to List</Link>
          <h1 className="text-3xl font-bold text-gray-900 tracking-tight uppercase tracking-widest">New Event</h1>
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
                <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} className="w-full p-4 text-xl font-bold border border-gray-100 rounded-2xl focus:ring-1 focus:ring-blue-600 outline-none transition-all placeholder:font-normal" placeholder="イベントのタイトル" required />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
                <div className="space-y-4">
                  <label className="block text-[10px] font-bold text-gray-400 uppercase tracking-widest">Start Date & Time</label>
                  <div className="flex gap-2">
                    <input type="date" min={today} value={startDate} onChange={(e) => setStartDate(e.target.value)} className="flex-1 p-3 border border-gray-100 rounded-xl focus:ring-1 focus:ring-blue-600 outline-none" required />
                    <input type="time" value={startTime} onChange={(e) => setStartTime(e.target.value)} className="w-32 p-3 border border-gray-100 rounded-xl focus:ring-1 focus:ring-blue-600 outline-none" required />
                  </div>
                </div>

                <div className="space-y-4">
                  <label className="block text-[10px] font-bold text-gray-400 uppercase tracking-widest">End Date & Time</label>
                  <div className="flex gap-2">
                    <input type="date" min={startDate || today} value={endDate} onChange={(e) => { setEndDate(e.target.value); setIsEndTimeManual(true); }} className="flex-1 p-3 border border-gray-100 rounded-xl focus:ring-1 focus:ring-blue-600 outline-none" />
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
                <textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={10} className="w-full p-4 border border-gray-100 rounded-2xl focus:ring-1 focus:ring-blue-600 outline-none resize-y min-h-[200px] font-light leading-relaxed" placeholder="イベントの詳細" />
              </div>
            </div>

            <div className="flex gap-4">
              <button
                type="button"
                onClick={(e) => handleCreateEvent(e, 'draft')}
                disabled={isLoading}
                className="flex-1 py-5 bg-gray-50 text-gray-400 text-xs font-bold uppercase tracking-[0.2em] rounded-2xl hover:bg-gray-100 transition-all"
              >
                下書きとして保存
              </button>
              <button
                type="button"
                onClick={(e) => handleCreateEvent(e, 'published')}
                disabled={isLoading}
                className="flex-[2] py-5 bg-blue-600 text-white text-xs font-bold uppercase tracking-[0.2em] rounded-2xl hover:bg-blue-700 shadow-xl shadow-blue-50 transition-all disabled:bg-blue-300"
              >
                {isLoading ? 'Creating...' : 'イベントを公開する'}
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
