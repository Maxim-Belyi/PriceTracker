export interface Item {
  id: number;
  title: string;
  image_url: string;
  current_price: number;
  status: "pending" | "success" | "error";
}

export interface PriceHistory {
  price: number;
  date: string;
}

const API_BASE_URL = 'http://localhost:8080';

export async function fetchItems(): Promise<Item[]> {
  try {
    const res = await fetch(`${API_BASE_URL}/items`);
    if (!res.ok) throw new Error("Failed to fetch items");
    return await res.json();
  } catch (error) {
    console.warn("API unreachable, using dummy data for items", error);
    return [
      { id: 1, title: "Apple iPhone 15 Pro Max 256GB, Natural Titanium", image_url: "https://images.unsplash.com/photo-1695048133142-1a20484d2569?auto=format&fit=crop&q=80&w=200&h=200", current_price: 139990, status: "success" },
      { id: 2, title: "Sony PlayStation 5 Slim (Disk Edition)", image_url: "https://images.unsplash.com/photo-1606813907291-d86efa9b94db?auto=format&fit=crop&q=80&w=200&h=200", current_price: 54990, status: "success" },
      { id: 3, title: "Dyson V15 Detect Absolute", image_url: "https://images.unsplash.com/photo-1558317374-067fb5f30001?auto=format&fit=crop&q=80&w=200&h=200", current_price: 69990, status: "pending" },
      { id: 4, title: "Samsung Galaxy S24 Ultra 512GB", image_url: "https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?auto=format&fit=crop&q=80&w=200&h=200", current_price: 119990, status: "error" },
    ];
  }
}

export async function trackUrl(url: string): Promise<{ id: number, status: string }> {
  try {
    const res = await fetch(`${API_BASE_URL}/track`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url })
    });
    if (!res.ok) throw new Error("Failed to track URL");
    return await res.json();
  } catch (error) {
    console.warn("API unreachable, using dummy tracking response", error);
    await new Promise(resolve => setTimeout(resolve, 1000));
    return { id: Math.floor(Math.random() * 1000) + 10, status: "Сохранено" };
  }
}

export async function fetchItemHistory(id: number): Promise<PriceHistory[]> {
  try {
    const res = await fetch(`${API_BASE_URL}/items/${id}/history`);
    if (!res.ok) throw new Error("Failed to fetch history");
    return await res.json();
  } catch (error) {
    console.warn(`API unreachable, using dummy history for item ${id}`, error);
    await new Promise(resolve => setTimeout(resolve, 500));
    
    // Generate some random realistic looking dummy data based on ID
    const basePrice = 100000 + (id * 10000);
    return [
      { price: basePrice + 15000, date: "2026-05-01" },
      { price: basePrice + 10000, date: "2026-05-15" },
      { price: basePrice + 5000, date: "2026-06-01" },
      { price: basePrice, date: "2026-06-13" },
    ];
  }
}
