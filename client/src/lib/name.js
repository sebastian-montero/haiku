// Per-browser pen name. No auth; just a pseudonym.
const KEY = "haiku:name";

export function loadName() {
  if (typeof window === "undefined") return "";
  return localStorage.getItem(KEY) || "";
}

export function saveName(name) {
  if (typeof window === "undefined") return;
  localStorage.setItem(KEY, name);
}
