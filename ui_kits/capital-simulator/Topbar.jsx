function Topbar({ activeChapter }) {
  return (
    <header className="topbar">
      <span className="topbar-logo">Capital Simulator</span>
      {activeChapter && (
        <span className="topbar-chapter">
          Vol. I &middot; Ch.{String(activeChapter.number).padStart(2, "0")} &mdash;{" "}
          {activeChapter.title}
        </span>
      )}
    </header>
  );
}
window.Topbar = Topbar;
