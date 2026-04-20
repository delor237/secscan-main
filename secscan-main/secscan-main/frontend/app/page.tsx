"use client";

import { useState } from 'react';

export default function Home() {
  const [url, setUrl] = useState('');
  const [scanId, setScanId] = useState<string | null>(null);
  const [progress, setProgress] = useState(0);
  const [messages, setMessages] = useState<string[]>([]);
  const [status, setStatus] = useState<'idle' | 'running' | 'completed' | 'error'>('idle');

  const startScan = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url) return;

    setStatus('running');
    setProgress(0);
    setMessages(['Tarama başlatılıyor...']);

    try {
      const res = await fetch('http://localhost:8080/api/scan', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url, modules: ["SBOM", "Nmap", "Header"] }),
      });
      const data = await res.json();
      
      if (data.scan_id) {
        setScanId(data.scan_id);
        connectSSE(data.scan_id);
      }
    } catch (err) {
      setStatus('error');
      setMessages([...messages, 'Hata oluştu! Sunucuya bağlanılamadı.']);
    }
  };

  const connectSSE = (id: string) => {
    const eventSource = new EventSource(`http://localhost:8080/api/scan/${id}/stream`);
    
    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setProgress(data.progress || 0);
      setMessages(prev => [...prev, data.message]);
      
      if (data.status === 'completed') {
        setStatus('completed');
        eventSource.close();
      }
    };

    eventSource.onerror = (err) => {
      console.error('SSE Error:', err);
      eventSource.close();
      setStatus('error');
    };
  };

  return (
    <main className="min-h-screen bg-gray-900 text-white p-10 flex flex-col items-center">
      <div className="max-w-3xl w-full">
        <h1 className="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-emerald-400 mb-8">
          SecScan Finale Projesi
        </h1>
        
        <form onSubmit={startScan} className="bg-gray-800 p-6 rounded-xl shadow-lg border border-gray-700">
          <div className="flex gap-4">
            <input 
              type="url" 
              placeholder="https://example.com" 
              className="flex-1 bg-gray-900 border border-gray-600 rounded-lg px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              required
            />
            <button 
              type="submit" 
              disabled={status === 'running'}
              className="bg-emerald-500 hover:bg-emerald-600 font-semibold px-8 py-3 rounded-lg transition-colors disabled:opacity-50"
            >
              Tara
            </button>
          </div>
        </form>

        {(status === 'running' || status === 'completed') && (
          <div className="mt-8 bg-gray-800 p-6 rounded-xl shadow-lg border border-gray-700">
            <div className="flex justify-between mb-2">
              <span className="font-semibold text-gray-300">İlerleme Durumu</span>
              <span className="text-emerald-400 font-bold">{progress}%</span>
            </div>
            <div className="w-full bg-gray-700 rounded-full h-4 mb-6">
              <div 
                className="bg-emerald-500 h-4 rounded-full transition-all duration-500"
                style={{ width: \`\${progress}%\` }}
              ></div>
            </div>

            <div className="bg-gray-900 p-4 rounded-lg h-64 overflow-y-auto font-mono text-sm border border-gray-700">
              {messages.map((m, i) => (
                <div key={i} className="mb-2 text-green-400">
                  <span className="text-gray-500">[{new Date().toLocaleTimeString()}]</span> {m}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </main>
  );
}
