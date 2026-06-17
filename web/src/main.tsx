import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import Root from "./Root";
import { AuthProvider } from "./AuthContext";
import { AuthBar } from "./components/AuthBar";
import "./index.css";
import "./chapters/ch-panel.css";

const rootEl = document.getElementById("root");
if (!rootEl) {
  throw new Error("Root element #root not found");
}

createRoot(rootEl).render(
  <StrictMode>
    <AuthProvider>
      <AuthBar />
      <Root />
    </AuthProvider>
  </StrictMode>
);
