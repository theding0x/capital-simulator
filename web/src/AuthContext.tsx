import { createContext, useContext, useEffect, useState } from "react";
import type { ReactNode } from "react";

import { setCanWrite } from "./api";

export interface AuthState {
  status: "loading" | "guest" | "authed";
  login: string;
  isOwner: boolean;
  oauthEnabled: boolean;
}

interface AuthContextShape extends AuthState {
  canWrite: boolean;
  signIn: () => void;
  signOut: () => Promise<void>;
}

const initial: AuthState = {
  status: "loading",
  login: "",
  isOwner: false,
  oauthEnabled: false,
};

const AuthContext = createContext<AuthContextShape>({
  ...initial,
  canWrite: false,
  signIn: () => {},
  signOut: async () => {},
});

interface MeResponse {
  authenticated: boolean;
  login: string;
  is_owner: boolean;
  oauth_enabled: boolean;
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>(initial);

  async function refresh() {
    try {
      const res = await fetch("/api/v1/auth/me");
      const body = (await res.json()) as MeResponse;
      setState({
        status: body.authenticated ? "authed" : "guest",
        login: body.login ?? "",
        isOwner: !!body.is_owner,
        oauthEnabled: !!body.oauth_enabled,
      });
      setCanWrite(!!body.is_owner);
    } catch {
      setState({ status: "guest", login: "", isOwner: false, oauthEnabled: false });
      setCanWrite(false);
    }
  }

  useEffect(() => {
    void refresh();
  }, []);

  const signIn = () => {
    window.location.href = "/api/v1/auth/github/login";
  };

  const signOut = async () => {
    try {
      await fetch("/api/v1/auth/logout", { method: "POST" });
    } catch {
      /* ignore */
    }
    await refresh();
  };

  return (
    <AuthContext.Provider value={{ ...state, canWrite: state.isOwner, signIn, signOut }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextShape {
  return useContext(AuthContext);
}
