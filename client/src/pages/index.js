import { useEffect, useState } from "react";
import { useRouter } from "next/router";
import Link from "next/link";
import TopBar from "../components/TopBar";
import { api } from "../lib/api";
import { loadName, saveName } from "../lib/name";

export default function Home() {
  const router = useRouter();
  const [notebooks, setNotebooks] = useState([]);
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState(null);
  const [title, setTitle] = useState("");
  const [creating, setCreating] = useState(false);

  const refresh = async () => {
    try {
      const data = await api.listNotebooks();
      setNotebooks(Array.isArray(data) ? data : []);
      setError(null);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoaded(true);
    }
  };

  useEffect(() => {
    refresh();
    const id = setInterval(refresh, 3000);
    return () => clearInterval(id);
  }, []);

  const handleCreate = async (e) => {
    e.preventDefault();
    const t = title.trim();
    if (!t) return;
    let author = loadName().trim();
    if (!author) {
      author = window.prompt("pen name?") || "";
      author = author.trim();
      if (!author) return;
      saveName(author);
    }
    setCreating(true);
    try {
      const nb = await api.createNotebook(t, author);
      router.push(`/write/${nb.id}`);
    } catch (e) {
      setError(e.message);
    } finally {
      setCreating(false);
    }
  };

  const live = notebooks.filter((n) => n.live);
  const recent = notebooks.filter((n) => !n.live);

  return (
    <div className="min-h-screen">
      <TopBar />
      <main className="max-w-3xl mx-auto px-5 pt-10 pb-24">
        <section className="mb-12">
          <h1 className="text-3xl font-semibold mb-2">watch them write.</h1>
          <p className="text-muted italic mb-6">
            open a notebook and whoever's tuned in sees every keystroke —
            the pauses, the reversals, the cross-outs. writing as performance.
          </p>
          <form
            onSubmit={handleCreate}
            className="flex items-center gap-3 border-b border-[#e5e0d3] pb-4"
          >
            <input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="title of a new notebook…"
              className="flex-1 bg-transparent outline-none py-2 text-lg placeholder:italic placeholder:text-[#a59d8d]"
            />
            <button
              type="submit"
              disabled={!title.trim() || creating}
              className="mono text-xs border border-ink px-4 py-2 hover:bg-ink hover:text-paper disabled:opacity-30"
            >
              {creating ? "…" : "begin"}
            </button>
          </form>
        </section>

        {error && (
          <p className="mono text-xs text-accent mb-6">
            couldn't reach haiku: {error}
          </p>
        )}

        <Section
          title="writing now"
          empty="nobody is writing at this moment."
          items={live}
          render={(n) => <LiveRow key={n.id} n={n} />}
        />

        <Section
          title="recent notebooks"
          empty={loaded ? "no notebooks yet — start one above." : "…"}
          items={recent}
          render={(n) => <RecentRow key={n.id} n={n} onDeleted={refresh} />}
        />
      </main>
    </div>
  );
}

function Section({ title, items, empty, render }) {
  return (
    <section className="mb-12">
      <h2 className="mono text-xs uppercase tracking-[0.2em] text-muted mb-4">
        {title}
      </h2>
      {items.length === 0 ? (
        <p className="italic text-muted">{empty}</p>
      ) : (
        <ul className="divide-y divide-[#e5e0d3]">{items.map(render)}</ul>
      )}
    </section>
  );
}

function LiveRow({ n }) {
  return (
    <li>
      <Link
        href={`/watch/${n.id}`}
        className="py-4 flex items-start justify-between gap-4 hover:bg-[#f4f0e4] -mx-2 px-2 rounded-sm"
      >
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-3">
            <span className="live-dot" />
            <span className="font-medium truncate">{n.title}</span>
          </div>
          <p className="mono text-xs text-muted mt-1">
            by {n.author}
            {n.writer_name && n.writer_name !== n.author
              ? ` · ${n.writer_name} is writing`
              : ""}
            {" · "}
            {n.viewers} watching
          </p>
          {n.preview && (
            <p className="mt-2 text-sm text-muted line-clamp-2">{n.preview}</p>
          )}
        </div>
        <span className="mono text-xs text-accent whitespace-nowrap self-center">
          watch →
        </span>
      </Link>
    </li>
  );
}

function RecentRow({ n, onDeleted }) {
  const [confirming, setConfirming] = useState(false);
  const del = async (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (!confirming) {
      setConfirming(true);
      setTimeout(() => setConfirming(false), 2500);
      return;
    }
    await api.deleteNotebook(n.id);
    onDeleted?.();
  };
  return (
    <li className="py-4 flex items-start justify-between gap-4">
      <Link
        href={`/watch/${n.id}`}
        className="min-w-0 flex-1 hover:bg-[#f4f0e4] -mx-2 px-2 py-1 rounded-sm"
      >
        <div className="font-medium truncate">{n.title}</div>
        <p className="mono text-xs text-muted mt-1">
          by {n.author} · {n.ops_count} edits ·{" "}
          {new Date(n.updated_at).toLocaleString()}
        </p>
        {n.preview && (
          <p className="mt-2 text-sm text-muted line-clamp-2">{n.preview}</p>
        )}
      </Link>
      <div className="mono text-xs flex items-center gap-3 whitespace-nowrap pt-1">
        <Link href={`/write/${n.id}`} className="text-muted hover:text-ink">
          continue
        </Link>
        <Link href={`/watch/${n.id}?replay=1`} className="text-accent">
          replay
        </Link>
        <button
          onClick={del}
          className="text-muted hover:text-accent"
          title="delete notebook"
        >
          {confirming ? "sure?" : "delete"}
        </button>
      </div>
    </li>
  );
}
