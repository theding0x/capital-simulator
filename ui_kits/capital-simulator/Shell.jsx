function Topbar({ chapter }) {
  return (
    <header className="topbar">
      <span className="topbar-logo">Capital Simulator</span>
      {chapter && (
        <span className="topbar-chapter">
          Vol. I &middot; Ch.{String(chapter.number).padStart(2, "0")} &mdash; {chapter.title}
        </span>
      )}
    </header>
  );
}

function Sidebar({ activeId, onSelect }) {
  const parts = [];
  for (const ch of window.CHAPTERS) {
    let p = parts.find((x) => x.name === ch.part);
    if (!p) { p = { name: ch.part, chapters: [] }; parts.push(p); }
    p.chapters.push(ch);
  }
  return (
    <nav className="sidebar">
      {parts.map((part) => (
        <div key={part.name}>
          <div className="sidebar-part">{part.name}</div>
          {part.chapters.map((ch) => {
            const isActive = ch.id === activeId;
            const cls = ["sidebar-item",
              isActive ? "sidebar-item--active"
                : ch.status === "done" ? "sidebar-item--done"
                : "sidebar-item--pending"].join(" ");
            return (
              <button key={ch.id} className={cls} onClick={() => onSelect(ch.id)}>
                <span className="sidebar-num">{String(ch.number).padStart(2, "0")}</span>
                <span className="sidebar-title">{ch.title}</span>
                {ch.status === "done" && !isActive && <span className="sidebar-icon">✓</span>}
                {isActive && <span className="sidebar-icon">›</span>}
              </button>
            );
          })}
        </div>
      ))}
    </nav>
  );
}

function ChapterHeader({ chapter, quote }) {
  return (
    <header className="chapter-header">
      <div className="chapter-header-num">
        Vol. I &middot; Chapter {String(chapter.number).padStart(2, "0")}
      </div>
      <h1 className="chapter-header-title">{chapter.title}</h1>
      {quote && <p className="chapter-header-quote">{quote}</p>}
    </header>
  );
}

function Card({ title, description, children }) {
  return (
    <section className="card">
      {title && <h2>{title}</h2>}
      {description && <p className="description">{description}</p>}
      {children}
    </section>
  );
}

function Placeholder({ children }) {
  return <div className="chapter-placeholder">{children}</div>;
}

Object.assign(window, { Topbar, Sidebar, ChapterHeader, Card, Placeholder });
