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

const API_BASE_URL = 'http://localhost:9090';

export async function fetchItems(): Promise<Item[]> {
  try {
    const res = await fetch(`${API_BASE_URL}/items`);
    if (!res.ok) throw new Error("Failed to fetch items");
    return await res.json();
  } catch (error) {
    console.warn("API unreachable, using dummy data for items", error);
    throw error;
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
    throw error;
  }
}

export async function fetchItemHistory(id: number): Promise<PriceHistory[]> {
  try {
    const res = await fetch(`${API_BASE_URL}/history/${id}`);
    if (!res.ok) throw new Error("Failed to fetch history");
    return await res.json();
  } catch (error) {
    console.warn("failed to fetch history", error);
    await new Promise(resolve => setTimeout(resolve, 500));

    const basePrice = 100000 + (id * 10000);
    return [
      { price: basePrice + 15000, date: "2026-05-01" },
      { price: basePrice + 10000, date: "2026-05-15" },
      { price: basePrice + 5000, date: "2026-06-01" },
      { price: basePrice, date: "2026-06-13" },
    ];
  }
}
