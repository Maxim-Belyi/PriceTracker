import { Item } from "../api/client";
import { cn } from "../lib/utils";
import { AlertCircle, CheckCircle2, Clock, CloudOff } from "lucide-react";

interface ItemCardProps {
  item: Item;
  onClick: (item: Item) => void;
}

export function ItemCard({ item, onClick }: ItemCardProps) {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("ru-RU", {
      style: "currency",
      currency: "RUB",
      maximumFractionDigits: 0,
    }).format(price);
  };

  const statusConfig = {
    pending: { icon: Clock, className: "text-amber-500 bg-amber-50", label: "В ожидании" },
    success: { icon: CheckCircle2, className: "text-emerald-500 bg-emerald-50", label: "Обновлено" },
    error: { icon: AlertCircle, className: "text-red-500 bg-red-50", label: "Ошибка" },
    processed: { icon: CloudOff, className: "text-blue-500 bg-blue-50", label: "В процессе" },
  };

  const StatusIcon = statusConfig[item.status].icon;

  return (
    <div
      onClick={() => onClick(item)}
      className="group bg-white rounded-2xl p-4 border border-slate-200 shadow-sm hover:shadow-md hover:border-slate-300 transition-all duration-200 cursor-pointer flex flex-col"
    >
      <div className="relative aspect-square mb-4 bg-slate-50 rounded-xl overflow-hidden">
        <img
          src={item.image_url}
          alt={item.title}
          className="object-cover w-full h-full group-hover:scale-105 transition-transform duration-300"
        />
        <div className={cn(
          "absolute top-3 left-3 flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-semibold tracking-wide backdrop-blur-md shadow-sm border border-white/20",
          statusConfig[item.status].className
        )}>
          <StatusIcon className="w-3.5 h-3.5" />
          {statusConfig[item.status].label}
        </div>
      </div>
      
      <div className="flex-1 flex flex-col">
        <h3 className="text-sm font-medium text-slate-700 line-clamp-2 leading-snug mb-2" title={item.title}>
          {item.title}
        </h3>
        <div className="mt-auto">
          <span className="text-lg font-bold text-slate-900 tracking-tight">
            {formatPrice(item.current_price)}
          </span>
        </div>
      </div>
    </div>
  );
}
