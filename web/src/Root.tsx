import { useEffect, useState } from "react";
import App from "./App";
import Atlas from "./atlas/Atlas";

/** Reads the current hash route, e.g. "#/chapters" → "/chapters". */
function useHashRoute(): string {
  const read = () => window.location.hash.replace(/^#/, "") || "/";
  const [route, setRoute] = useState(read);
  useEffect(() => {
    const onHash = () => setRoute(read());
    window.addEventListener("hashchange", onHash);
    return () => window.removeEventListener("hashchange", onHash);
  }, []);
  return route;
}

/**
 * Minimal router:
 *   "/"                        → Atlas (observatory)
 *   "/chapters"                → App (chapter dashboard)
 *   "/atlas/chapter/:id"       → Atlas with SpineMatrix open + drawer pre-opened on :id
 */
export default function Root() {
  const route = useHashRoute();
  if (route.startsWith("/chapters")) return <App />;

  // Deep-link: /atlas/chapter/v1-ch07 → open SpineMatrix with that chapter
  const chMatch = route.match(/^\/atlas\/chapter\/([^/]+)$/);
  if (chMatch) return <Atlas initialChapterId={chMatch[1]} />;

  return <Atlas />;
}
