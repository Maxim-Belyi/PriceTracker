import { LayoutGrid, MonitorSmartphone, Package, ShoppingBag, Store } from "lucide-react";
import { cn } from "../lib/utils";

const tabs = [
  { id: "all", label: "Все магазины", icon: LayoutGrid },
  { id: "citilink", label: "Ситилинк", icon: MonitorSmartphone },
  { id: "ozon", label: "Ozon", icon: Package },
  { id: "wb", label: "Wildberries", icon: Store },
];

interface SidebarProps {
  activeTab: string;
  setActiveTab: (id: string) => void;
}

export function Sidebar({ activeTab, setActiveTab }: SidebarProps) {
  return (
    <aside className="w-64 h-screen border-r border-slate-200 bg-white flex flex-col fixed left-0 top-0">
      <div className="h-16 flex items-center px-6 border-b border-slate-100">
        <div className="flex items-center gap-2 text-slate-900">
          <ShoppingBag className="w-6 h-6 text-blue-600" />
          <span className="font-semibold text-lg tracking-tight">PriceTracker</span>
        </div>
      </div>
      <nav className="flex-1 p-4 space-y-1">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          const isActive = activeTab === tab.id;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={cn(
                "w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-200",
                isActive
                  ? "bg-blue-50 text-blue-700"
                  : "text-slate-600 hover:bg-slate-50 hover:text-slate-900"
              )}
            >
              <Icon className={cn("w-5 h-5", isActive ? "text-blue-600" : "text-slate-400")} />
              {tab.label}
            </button>
          );
        })}
      </nav>
      <div className="p-4 border-t border-slate-100">
        <div className="text-xs text-slate-400 text-center">
          &copy; 2026 PriceTracker
        </div>
      </div>
    </aside>
  );
}
