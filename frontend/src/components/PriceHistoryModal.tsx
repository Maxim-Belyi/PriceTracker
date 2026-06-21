import { useEffect, useState } from "react";
import { Item, PriceHistory, fetchItemHistory, proxyImageUrl } from "../api/client";
import { X, Loader2, TrendingDown, TrendingUp, Minus } from "lucide-react";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";
import { cn } from "../lib/utils";

interface PriceHistoryModalProps {
  item: Item | null;
  onClose: () => void;
}

export function PriceHistoryModal({ item, onClose }: PriceHistoryModalProps) {
  const [history, setHistory] = useState<PriceHistory[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (!item) return;
    
    let isMounted = true;
    const loadHistory = async () => {
      setIsLoading(true);
      try {
        const data = await fetchItemHistory(item.id);
        if (isMounted) setHistory(data);
      } catch (e) {
        console.error(e);
      } finally {
        if (isMounted) setIsLoading(false);
      }
    };
    
    loadHistory();
    return () => { isMounted = false };
  }, [item]);

  if (!item) return null;

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("ru-RU", { style: "currency", currency: "RUB", maximumFractionDigits: 0 }).format(price);
  };
  
  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString("ru-RU", { day: 'numeric', month: 'short' });
  };

  // Calculate trend
  const sortedHistory = [...history].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());
  const firstPrice = sortedHistory.length > 0 ? sortedHistory[0].price : item.current_price;
  const currentPrice = item.current_price;
  const diff = currentPrice - firstPrice;
  const percentChange = firstPrice ? Math.abs((diff / firstPrice) * 100).toFixed(1) : "0";

  return (
    <>
      <div 
        className="fixed inset-0 bg-slate-900/20 backdrop-blur-sm z-40 transition-opacity"
        onClick={onClose}
      />
      
      <div className="fixed inset-y-0 right-0 w-full max-w-xl bg-white shadow-2xl z-50 animate-in slide-in-from-right duration-300 sm:rounded-l-3xl border-l border-slate-200 flex flex-col">
        <div className="flex items-center justify-between px-6 py-4 border-b border-slate-100">
          <h2 className="font-semibold text-lg text-slate-800">История цены</h2>
          <button 
            onClick={onClose}
            className="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-full transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="p-6 flex flex-col flex-1 overflow-y-auto">
          {/* Header Info */}
          <div className="flex gap-4 mb-8">
            <img src={proxyImageUrl(item.image_url)} alt={item.title ?? undefined} className="w-20 h-20 object-contain rounded-xl border border-slate-100 shadow-sm" />
            <div className="flex flex-col justify-center">
              <h3 className="font-medium text-slate-700 leading-snug mb-1">{item.title ?? "Без названия"}</h3>
              <div className="flex items-baseline gap-3">
                <span className="text-2xl font-bold text-slate-900 tracking-tight">{formatPrice(item.current_price)}</span>
                
                {diff !== 0 && (
                  <div className={cn(
                    "flex items-center text-sm font-medium",
                    diff > 0 ? "text-red-500" : "text-emerald-500"
                  )}>
                    {diff > 0 ? <TrendingUp className="w-4 h-4 mr-1" /> : <TrendingDown className="w-4 h-4 mr-1" />}
                    {percentChange}%
                  </div>
                )}
                {diff === 0 && history.length > 1 && (
                  <div className="flex items-center text-sm font-medium text-slate-400">
                    <Minus className="w-4 h-4 mr-1" />
                    0%
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Chart Area */}
          <div className="flex-1 min-h-[300px] w-full bg-slate-50/50 rounded-2xl border border-slate-100 p-4 relative">
            {isLoading ? (
              <div className="absolute inset-0 flex items-center justify-center">
                <Loader2 className="w-8 h-8 text-blue-500 animate-spin" />
              </div>
            ) : history.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={sortedHistory} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#e2e8f0" />
                  <XAxis 
                    dataKey="date" 
                    tickFormatter={formatDate}
                    axisLine={false}
                    tickLine={false}
                    tick={{ fill: '#64748b', fontSize: 12 }}
                    dy={10}
                  />
                  <YAxis 
                    domain={['auto', 'auto']}
                    tickFormatter={(val) => new Intl.NumberFormat("ru-RU", { notation: "compact" }).format(val)}
                    axisLine={false}
                    tickLine={false}
                    tick={{ fill: '#64748b', fontSize: 12 }}
                  />
                  <Tooltip 
                    formatter={(value: any) => [formatPrice(Number(value)), "Цена"]}
                    labelFormatter={(label) => formatDate(label as string)}
                    contentStyle={{ borderRadius: '12px', border: 'none', boxShadow: '0 10px 15px -3px rgb(0 0 0 / 0.1)' }}
                  />
                  <Line 
                    type="monotone" 
                    dataKey="price" 
                    stroke="#3b82f6" 
                    strokeWidth={3}
                    dot={{ fill: '#3b82f6', strokeWidth: 2, r: 4 }}
                    activeDot={{ r: 6, strokeWidth: 0 }}
                  />
                </LineChart>
              </ResponsiveContainer>
            ) : (
              <div className="h-full flex items-center justify-center text-slate-400 text-sm">
                Нет данных для отображения
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
}
