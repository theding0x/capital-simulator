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

/** Minimal two-view router: Atlas at "/", the chapter dashboard at "/chapters". */
export default function Root() {
  const route = useHashRoute();
  if (route.startsWith("/chapters")) return <App />;
  return <Atlas />;
}
