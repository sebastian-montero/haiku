// Compute a minimal single-range edit between two strings.
// Returns { p, n, s } suitable for the server Op type, or null if unchanged.
// Positions and lengths are JS UTF-16 code-unit offsets, matching the Go side.
export function diff(prev, next) {
  if (prev === next) return null;
  const lenA = prev.length;
  const lenB = next.length;
  let start = 0;
  const max = Math.min(lenA, lenB);
  while (start < max && prev.charCodeAt(start) === next.charCodeAt(start)) {
    start++;
  }
  let endA = lenA;
  let endB = lenB;
  while (
    endA > start &&
    endB > start &&
    prev.charCodeAt(endA - 1) === next.charCodeAt(endB - 1)
  ) {
    endA--;
    endB--;
  }
  return {
    p: start,
    n: endA - start,
    s: next.slice(start, endB),
  };
}

// Apply a single op (server-shape) to a string.
export function applyOp(text, op) {
  const p = Math.max(0, Math.min(op.p, text.length));
  const end = Math.min(p + (op.n || 0), text.length);
  return text.slice(0, p) + (op.s || "") + text.slice(end);
}
