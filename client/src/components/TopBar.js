import Link from "next/link";
import { useEffect, useState } from "react";
import { loadName, saveName } from "../lib/name";

export default function TopBar() {
  const [name, setName] = useState("");
  const [editing, setEditing] = useState(false);

  useEffect(() => setName(loadName()), []);

  const commit = () => {
    const trimmed = name.trim();
    if (trimmed) saveName(trimmed);
    setEditing(false);
  };

  return (
    <header className="w-full border-b border-[#e5e0d3]">
      <div className="max-w-3xl mx-auto flex items-center justify-between px-5 py-4">
        <Link href="/" className="text-xl font-semibold tracking-tight">
          haiku<span className="text-accent">·</span>
        </Link>
        <div className="mono text-xs text-muted flex items-center gap-3">
          {editing ? (
            <input
              autoFocus
              value={name}
              onChange={(e) => setName(e.target.value)}
              onBlur={commit}
              onKeyDown={(e) => e.key === "Enter" && commit()}
              className="bg-transparent border-b border-[#c9c2b2] focus:border-ink outline-none w-32"
              placeholder="pen name"
            />
          ) : (
            <button
              onClick={() => setEditing(true)}
              className="hover:text-ink underline-offset-4 hover:underline"
              title="change pen name"
            >
              {name || "set pen name"}
            </button>
          )}
        </div>
      </div>
    </header>
  );
}
