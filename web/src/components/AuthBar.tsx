import { useAuth } from "../AuthContext";

/**
 * Fixed-position auth control (top-right). Mounted once globally in main.tsx
 * so it overlays both Atlas and the chapter dashboard without touching their
 * layouts. Read-only state is conveyed by the control itself (guest / sign-in)
 * and by the inline message shown when a non-owner attempts a write.
 */
export function AuthBar() {
  const { status, login, isOwner, oauthEnabled, signIn, signOut } = useAuth();

  if (status === "loading") return null;

  return (
    <div className="auth-bar">
      {isOwner ? (
        <>
          <span className="auth-bar-user">{login}</span>
          <button className="auth-bar-btn" onClick={() => void signOut()}>
            Sign out
          </button>
        </>
      ) : status === "authed" ? (
        <>
          <span className="auth-bar-user">{login} (guest)</span>
          <button className="auth-bar-btn" onClick={() => void signOut()}>
            Sign out
          </button>
        </>
      ) : oauthEnabled ? (
        <button className="auth-bar-btn" onClick={signIn}>
          Sign in with GitHub
        </button>
      ) : (
        <span className="auth-bar-user">guest</span>
      )}
    </div>
  );
}
