import { useState } from "react";
import { Link, Loader2 } from "lucide-react";
import { cn } from "../lib/utils";
import { trackUrl } from "../api/client";

interface TopBarProps {
  onTrackSuccess: () => void;
}

export function TopBar({ onTrackSuccess }: TopBarProps) {
  const [url, setUrl] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState<{ text: string, type: "error" | "success" } | null>(null);

  const handleTrack = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url) return;

    setIsLoading(true);
    setMessage(null);
    try {
      const res = await trackUrl(url);
      setMessage({ text: `Успешно: ${res.status} (ID: ${res.id})`, type: "success" });
      setUrl("");
      onTrackSuccess();
      setTimeout(() => setMessage(null), 3000);
    } catch (err) {
      setMessage({ text: "Ошибка при добавлении ссылки. Сервер недоступен?", type: "error" });
      setTimeout(() => setMessage(null), 3000);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <header className="h-20 flex items-center px-8 bg-white border-b border-slate-200 sticky top-0 z-10">
      <form onSubmit={handleTrack} className="flex-1 max-w-2xl relative flex items-center">
        <div className="absolute left-4 text-slate-400">
          <Link className="w-5 h-5" />
        </div>
        <input
          type="url"
          placeholder="Вставьте ссылку на товар..."
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          className="w-full pl-12 pr-32 py-3 bg-slate-50 border border-slate-200 rounded-2xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all placeholder:text-slate-400"
          required
        />
        <button
          type="submit"
          disabled={isLoading || !url}
          className="absolute right-2 px-4 py-1.5 bg-slate-900 text-white rounded-xl text-sm font-medium hover:bg-slate-800 focus:outline-none focus:ring-2 focus:ring-slate-900/20 transition-all disabled:opacity-50 flex items-center gap-2"
        >
          {isLoading && <Loader2 className="w-4 h-4 animate-spin" />}
          Отследить
        </button>
      </form>
      
      {/* Toast Notification Container */}
      {message && (
        <div className={cn(
          "fixed top-6 right-6 px-4 py-3 rounded-xl shadow-lg border text-sm font-medium flex items-center gap-2 animate-in slide-in-from-top-2 fade-in duration-300 z-50",
          message.type === "error" ? "bg-red-50 text-red-700 border-red-100" : "bg-green-50 text-green-700 border-green-100"
        )}>
          {message.text}
        </div>
      )}
    </header>
  );
}
