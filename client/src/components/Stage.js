// Stage renders text with an inline caret. Used on the watch page (and for
// replay) — the writer is shown as a textarea, not a Stage.
export default function Stage({ text, cursor, placeholder, ghost = false }) {
  const safeCursor = Math.max(0, Math.min(cursor ?? text.length, text.length));
  const before = text.slice(0, safeCursor);
  const after = text.slice(safeCursor);
  const isEmpty = text.length === 0;
  return (
    <div className="stage">
      {isEmpty ? (
        <span className="italic text-[#a59d8d]">{placeholder || ""}</span>
      ) : null}
      <span>{before}</span>
      <span className={`caret ${ghost ? "caret-ghost" : ""}`} />
      <span>{after}</span>
    </div>
  );
}
