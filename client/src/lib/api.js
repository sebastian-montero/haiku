// Thin fetch wrapper + WebSocket URL helper. No auth, all endpoints are open.

export const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

export const WS_BASE = API_BASE.replace(/^http/, "ws");

async function j(path, init = {}) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(text || `${res.status} ${res.statusText}`);
  }
  if (res.status === 204) return null;
  return res.json();
}

export const api = {
  listNotebooks: (filter) =>
    j(`/notebooks${filter ? `?filter=${filter}` : ""}`),
  getNotebook: (id) => j(`/notebooks/${id}`),
  createNotebook: (title, author) =>
    j(`/notebooks`, {
      method: "POST",
      body: JSON.stringify({ title, author }),
    }),
  renameNotebook: (id, title) =>
    j(`/notebooks/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ title }),
    }),
  deleteNotebook: (id) => j(`/notebooks/${id}`, { method: "DELETE" }),
};

export function writeSocket(id, name) {
  return new WebSocket(
    `${WS_BASE}/ws/write/${id}?name=${encodeURIComponent(name || "anon")}`,
  );
}

export function watchSocket(id, name) {
  return new WebSocket(
    `${WS_BASE}/ws/watch/${id}?name=${encodeURIComponent(name || "anon")}`,
  );
}
