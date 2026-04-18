import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/router";
import Link from "next/link";
import TopBar from "../../components/TopBar";
import NameGate from "../../components/NameGate";
import { api, writeSocket } from "../../lib/api";
import { diff } from "../../lib/diff";
import { loadName } from "../../lib/name";

export default function WritePage() {
  const router = useRouter();
  if (!router.isReady) return null;
  const { id } = router.query;
  if (!id) return null;
  return (
    <NameGate prompt="pick a pen name before you write">
      <Writer id={id} />
    </NameGate>
  );
}

function Writer({ id }) {
  const [meta, setMeta] = useState(null);
  const [text, setText] = useState("");
  const [viewers, setViewers] = useState(0);
  const [status, setStatus] = useState("connecting");
  const [error, setError] = useState(null);
  const textareaRef = useRef(null);
  const wsRef = useRef(null);
  const prevTextRef = useRef("");
  const cursorSendTimer = useRef(null);
  const lastCursorSent = useRef(-1);

  useEffect(() => {
    let cancelled = false;
    let ws;

    (async () => {
      try {
        const nb = await api.getNotebook(id);
        if (cancelled) return;
        if (nb.live && !nb.you_write) {
          // no-op; the server will reject the writer claim and we'll show the error.
        }
        setMeta(nb);
        setText(nb.text || "");
        prevTextRef.current = nb.text || "";
      } catch (e) {
        if (!cancelled) setError(e.message);
        return;
      }

      const name = loadName();
      ws = writeSocket(id, name);
      wsRef.current = ws;
      ws.onopen = () => setStatus("live");
      ws.onclose = () => setStatus("disconnected");
      ws.onerror = () => setStatus("error");
      ws.onmessage = (ev) => {
        try {
          const msg = JSON.parse(ev.data);
          if (msg.t === "init") {
            setMeta((m) => ({ ...(m || {}), ...msg }));
            setText(msg.text || "");
            prevTextRef.current = msg.text || "";
            setViewers(msg.viewers || 0);
          } else if (msg.t === "presence") {
            setViewers(msg.viewers || 0);
          } else if (msg.t === "err") {
            setError(msg.msg || "someone else is already writing");
            ws.close();
          }
        } catch {}
      };
    })();

    return () => {
      cancelled = true;
      if (wsRef.current) {
        try {
          wsRef.current.send(JSON.stringify({ t: "end" }));
        } catch {}
        wsRef.current.close();
      }
    };
  }, [id]);

  const sendOp = (op) => {
    const ws = wsRef.current;
    if (!ws || ws.readyState !== WebSocket.OPEN) return;
    ws.send(JSON.stringify({ t: "op", ...op }));
  };

  const sendCursor = (p) => {
    const ws = wsRef.current;
    if (!ws || ws.readyState !== WebSocket.OPEN) return;
    if (p === lastCursorSent.current) return;
    lastCursorSent.current = p;
    ws.send(JSON.stringify({ t: "cur", p }));
  };

  const queueCursor = (p) => {
    if (cursorSendTimer.current) clearTimeout(cursorSendTimer.current);
    cursorSendTimer.current = setTimeout(() => sendCursor(p), 30);
  };

  const handleChange = (e) => {
    const next = e.target.value;
    const prev = prevTextRef.current;
    const op = diff(prev, next);
    setText(next);
    prevTextRef.current = next;
    if (op) sendOp(op);
    queueCursor(e.target.selectionStart);
  };

  const handleSelect = (e) => queueCursor(e.target.selectionStart);

  if (error) {
    return (
      <div className="min-h-screen">
        <TopBar />
        <main className="max-w-2xl mx-auto px-5 pt-20 text-center">
          <p className="italic text-muted mb-4">{error}</p>
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
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-semibold">{meta?.title || "…"}</h1>
            <p className="mono text-xs text-muted mt-1">
              by {meta?.author || "…"}
            </p>
          </div>
          <div className="mono text-xs flex items-center gap-3 text-muted">
            <span className="flex items-center gap-2">
              <span
                className={`live-dot ${status === "live" ? "" : "opacity-30"}`}
                style={
                  status === "live" ? undefined : { animation: "none" }
                }
              />
              {status === "live" ? "broadcasting" : status}
            </span>
            <span>·</span>
            <span>
              {viewers} {viewers === 1 ? "viewer" : "viewers"}
            </span>
          </div>
        </div>

        <div className="rule mb-6" />

        <textarea
          ref={textareaRef}
          value={text}
          onChange={handleChange}
          onKeyUp={handleSelect}
          onClick={handleSelect}
          onSelect={handleSelect}
          placeholder="begin. they're listening."
          className="script min-h-[60vh]"
          autoFocus
        />

        <div className="rule mt-8 mb-4" />
        <div className="flex items-center justify-between mono text-xs text-muted">
          <span>{text.length} chars</span>
          <Link href="/" className="hover:text-ink">
            close (autosaved)
          </Link>
        </div>
      </main>
    </div>
  );
}
