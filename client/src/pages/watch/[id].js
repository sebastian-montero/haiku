import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/router";
import Link from "next/link";
import TopBar from "../../components/TopBar";
import Stage from "../../components/Stage";
import { api, watchSocket } from "../../lib/api";
import { applyOp } from "../../lib/diff";
import { loadName } from "../../lib/name";

export default function WatchPage() {
  const router = useRouter();
  if (!router.isReady) return null;
  const { id, replay } = router.query;
  if (!id) return null;
  if (replay) return <Replayer id={id} />;
  return <LiveWatcher id={id} />;
}

function LiveWatcher({ id }) {
  const [meta, setMeta] = useState(null);
  const [text, setText] = useState("");
  const [cursor, setCursor] = useState(0);
  const [viewers, setViewers] = useState(0);
  const [writerName, setWriterName] = useState("");
  const [live, setLive] = useState(false);
  const [status, setStatus] = useState("connecting");
  const [error, setError] = useState(null);
  const wsRef = useRef(null);

  useEffect(() => {
    const name = loadName();
    const ws = watchSocket(id, name);
    wsRef.current = ws;
    ws.onopen = () => setStatus("open");
    ws.onclose = () => setStatus("closed");
    ws.onerror = () => setStatus("error");
    ws.onmessage = (ev) => {
      try {
        const msg = JSON.parse(ev.data);
        if (msg.t === "init") {
          setMeta(msg);
          setText(msg.text || "");
          setCursor(msg.cursor ?? (msg.text || "").length);
          setViewers(msg.viewers || 0);
          setLive(!!msg.live);
          setWriterName(msg.writer_name || "");
        } else if (msg.t === "op") {
          setText((prev) => {
            const next = applyOp(prev, msg.op);
            // Nudge cursor to land just after the inserted text.
            setCursor(msg.op.p + (msg.op.s?.length || 0));
            return next;
          });
        } else if (msg.t === "cur") {
          setCursor(msg.p);
        } else if (msg.t === "presence") {
          setViewers(msg.viewers || 0);
          setLive(!!msg.live);
          setWriterName(msg.writer_name || "");
        }
      } catch {}
    };
    return () => {
      try {
        ws.close();
      } catch {}
    };
  }, [id]);

  useEffect(() => {
    api
      .getNotebook(id)
      .then((nb) => setMeta((m) => ({ ...(m || {}), ...nb })))
      .catch((e) => setError(e.message));
  }, [id]);

  if (error && !meta) {
    return (
      <div className="min-h-screen">
        <TopBar />
        <main className="max-w-2xl mx-auto px-5 pt-20 text-center">
          <p className="italic text-muted">{error}</p>
          <Link href="/" className="mono text-xs underline">
            back
          </Link>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <TopBar />
      <main className="max-w-2xl mx-auto px-5 pt-10 pb-24">
        <div className="flex items-start justify-between mb-6">
          <div>
            <h1 className="text-2xl font-semibold">{meta?.title || "…"}</h1>
            <p className="mono text-xs text-muted mt-1">
              by {meta?.author || "…"}
            </p>
          </div>
          <div className="mono text-xs text-muted text-right">
            <div className="flex items-center gap-2 justify-end">
              {live ? <span className="live-dot" /> : null}
              <span>
                {live
                  ? writerName
                    ? `${writerName} is writing`
                    : "live"
                  : "not live"}
              </span>
            </div>
            <div>
              {viewers} {viewers === 1 ? "watching" : "watching"}
            </div>
            <div className="opacity-60">{status}</div>
          </div>
        </div>

        <div className="rule mb-6" />

        {text || live ? (
          <Stage
            text={text}
            cursor={cursor}
            ghost={!live}
            placeholder={live ? "they're thinking…" : ""}
          />
        ) : (
          <p className="italic text-muted">(nothing has been written yet)</p>
        )}

        <div className="rule mt-8 mb-4" />
        <div className="flex items-center justify-between mono text-xs text-muted">
          <Link href={`/watch/${id}?replay=1`} className="hover:text-ink">
            ⟲ replay how it was written
          </Link>
          <Link href="/" className="hover:text-ink">
            close
          </Link>
        </div>
      </main>
    </div>
  );
}

function Replayer({ id }) {
  const [meta, setMeta] = useState(null);
  const [text, setText] = useState("");
  const [cursor, setCursor] = useState(0);
  const [idx, setIdx] = useState(0); // next op to apply
  const [playing, setPlaying] = useState(false);
  const [speed, setSpeed] = useState(2);
  const [error, setError] = useState(null);
  const timerRef = useRef(null);

  useEffect(() => {
    api
      .getNotebook(id)
      .then((nb) => {
        setMeta(nb);
        setText("");
        setCursor(0);
        setIdx(0);
      })
      .catch((e) => setError(e.message));
    return () => clearTimeout(timerRef.current);
  }, [id]);

  // Drive playback.
  useEffect(() => {
    if (!playing || !meta) return;
    const ops = meta.ops || [];
    if (idx >= ops.length) {
      setPlaying(false);
      return;
    }
    const prevTs = idx === 0 ? 0 : ops[idx - 1].ts;
    const delay = Math.max(8, (ops[idx].ts - prevTs) / speed);
    timerRef.current = setTimeout(() => {
      setText((t) => applyOp(t, ops[idx]));
      setCursor(ops[idx].p + (ops[idx].s?.length || 0));
      setIdx((i) => i + 1);
    }, Math.min(delay, 2000));
    return () => clearTimeout(timerRef.current);
  }, [playing, idx, meta, speed]);

  const restart = () => {
    clearTimeout(timerRef.current);
    setText("");
    setCursor(0);
    setIdx(0);
    setPlaying(true);
  };

  const jumpToEnd = () => {
    clearTimeout(timerRef.current);
    setPlaying(false);
    setText(meta?.text || "");
    setCursor((meta?.text || "").length);
    setIdx((meta?.ops || []).length);
  };

  if (error) {
    return (
      <div className="min-h-screen">
        <TopBar />
        <main className="max-w-2xl mx-auto px-5 pt-20 text-center">
          <p className="italic text-muted">{error}</p>
        </main>
      </div>
    );
  }

  const ops = meta?.ops || [];
  const pct =
    ops.length === 0 ? 100 : Math.round((idx / ops.length) * 100);
  const done = idx >= ops.length;

  return (
    <div className="min-h-screen">
      <TopBar />
      <main className="max-w-2xl mx-auto px-5 pt-10 pb-24">
        <div className="flex items-start justify-between mb-6">
          <div>
            <h1 className="text-2xl font-semibold">{meta?.title || "…"}</h1>
            <p className="mono text-xs text-muted mt-1">
              replay · by {meta?.author || "…"}
            </p>
          </div>
          <div className="mono text-xs text-muted text-right">
            {ops.length} edits · {pct}%
          </div>
        </div>

        <div className="rule mb-6" />

        <Stage
          text={text}
          cursor={cursor}
          ghost={!playing || done}
          placeholder="press play"
        />

        <div className="rule mt-8 mb-4" />
        <div className="flex items-center justify-between mono text-xs">
          <div className="flex items-center gap-3">
            <button
              onClick={() => setPlaying((p) => !p)}
              className="border border-ink px-3 py-1 hover:bg-ink hover:text-paper"
            >
              {done ? "done" : playing ? "pause" : "play"}
            </button>
            <button
              onClick={restart}
              className="text-muted hover:text-ink"
            >
              restart
            </button>
            <button
              onClick={jumpToEnd}
              className="text-muted hover:text-ink"
            >
              skip
            </button>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-muted">speed</span>
            {[1, 2, 4, 8].map((s) => (
              <button
                key={s}
                onClick={() => setSpeed(s)}
                className={`px-2 py-1 border ${
                  s === speed
                    ? "border-ink bg-ink text-paper"
                    : "border-transparent text-muted hover:text-ink"
                }`}
              >
                {s}×
              </button>
            ))}
            <Link href={`/watch/${id}`} className="ml-3 text-muted hover:text-ink">
              live
            </Link>
          </div>
        </div>
      </main>
    </div>
  );
}
