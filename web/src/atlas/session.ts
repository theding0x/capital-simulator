// A per-page-load Atlas session id. Held in memory only — never persisted — so a
// reload mints a new id and the server starts that session's run fresh from seed.
function newSessionId(): string {
  const c = globalThis.crypto;
  if (c && typeof c.randomUUID === "function") return c.randomUUID();
  if (c && typeof c.getRandomValues === "function") {
    const b = new Uint8Array(16);
    c.getRandomValues(b);
    return Array.from(b, (x) => x.toString(16).padStart(2, "0")).join("");
  }
  return `s-${Date.now()}-${Math.random().toString(16).slice(2)}`;
}

export const atlasSessionId: string = newSessionId();

/** Header to attach to every Atlas observatory request. */
export const atlasSessionHeader: Record<string, string> = {
  "X-Atlas-Session": atlasSessionId,
};
