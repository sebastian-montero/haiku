import { useEffect, useState } from "react";
import { loadName, saveName } from "../lib/name";

// Blocks children until a pen name is set. Used by write/watch pages.
export default function NameGate({ children, prompt }) {
  const [name, setName] = useState(null); // null = loading
  const [draft, setDraft] = useState("");

  useEffect(() => {
    const existing = loadName();
    setName(existing);
    setDraft(existing);
  }, []);

  if (name === null) return null;
  if (name) return children;

  const submit = (e) => {
    e?.preventDefault();
    const trimmed = draft.trim();
    if (!trimmed) return;
    saveName(trimmed);
    setName(trimmed);
  };

  return (
    <div className="min-h-[60vh] flex items-center justify-center px-5">
      <form
        onSubmit={submit}
        className="max-w-sm w-full text-center space-y-4"
      >
        <p className="italic text-muted">{prompt || "what do they call you?"}</p>
        <input
          autoFocus
          value={draft}
          onChange={(e) => setDraft(e.target.value)}
          className="w-full bg-transparent border-b border-[#c9c2b2] focus:border-ink outline-none text-center py-2 text-lg"
          placeholder="pen name"
        />
        <button
          type="submit"
          disabled={!draft.trim()}
          className="mono text-xs border border-ink px-4 py-2 hover:bg-ink hover:text-paper disabled:opacity-30"
        >
          enter
        </button>
      </form>
    </div>
  );
}
