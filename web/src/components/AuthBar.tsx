import { useAuth } from "../AuthContext";

/**
 * Fixed-position auth control (top-right) plus a full-width read-only banner.
 * Mounted once globally in main.tsx so it overlays both Atlas and the
 * chapter dashboard without touching their layouts.
 */
export function AuthBar() {
  const { status, login, isOwner, oauthEnabled, canWrite, signIn, signOut } = useAuth();

  if (status === "loading") return null;

  return (
    <>
      {!canWrite && (
        <div className="auth-banner" role="status">
          Read-only — sign in as the owner to make changes.
        </div>
      )}
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
    </>
  );
}
