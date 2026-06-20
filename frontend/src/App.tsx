import { useEffect, useState } from "react";
import { Layout } from "./components/Layout";
import { TopBar } from "./components/TopBar";
import { ItemCard } from "./components/ItemCard";
import { PriceHistoryModal } from "./components/PriceHistoryModal";
import { fetchItems, Item } from "./api/client";
import { Loader2 } from "lucide-react";

function App() {
  const [activeTab, setActiveTab] = useState("all");
  const [items, setItems] = useState<Item[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedItem, setSelectedItem] = useState<Item | null>(null);

  const loadItems = async () => {
    setIsLoading(true);
    try {
      const data = await fetchItems();
      setItems(data);
    } catch (e) {
      console.error(e);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadItems();
  }, []);

  const handleTrackSuccess = () => {
    loadItems();
  };

  return (
    <Layout activeTab={activeTab} setActiveTab={setActiveTab}>
      <TopBar onTrackSuccess={handleTrackSuccess} />
      
      <main className="p-8 flex-1">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">
            {activeTab === "all" ? "Все товары" : 
             activeTab === "citilink" ? "Ситилинк" : 
             activeTab === "ozon" ? "Ozon" : "Wildberries"}
          </h1>
          <p className="text-slate-500 mt-1">Отслеживайте изменение цен на выбранные товары</p>
        </div>

        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <Loader2 className="w-8 h-8 text-blue-500 animate-spin" />
          </div>
        ) : items.length > 0 ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-6">
            {items.map((item) => (
              <ItemCard 
                key={item.id} 
                item={item} 
                onClick={(item) => setSelectedItem(item)} 
              />
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center h-64 text-slate-400 bg-white rounded-3xl border border-slate-200 border-dashed">
            <p>Нет отслеживаемых товаров</p>
            <p className="text-sm mt-1">Добавьте ссылку сверху, чтобы начать</p>
          </div>
        )}
      </main>

      {selectedItem && (
        <PriceHistoryModal 
          item={selectedItem} 
          onClose={() => setSelectedItem(null)} 
        />
      )}
    </Layout>
  );
}

export default App;
