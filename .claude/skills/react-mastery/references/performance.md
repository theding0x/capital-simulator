# Performance

> Load when: Profiling, optimization, handling large lists, reducing bundle size.

## Profile First, Optimize Never (Without Data)

The #1 performance mistake is optimizing without measuring. React is fast by default. Most "slow" React apps have one of three problems:
1. An expensive computation running on every render
2. A parent re-rendering thousands of children unnecessarily
3. A waterfall of sequential data fetches

## React DevTools Profiler

```
1. Open React DevTools → Profiler tab
2. Click Record → interact with the app → Stop
3. Look at the flame graph — wide bars = slow renders
4. Check "Why did this render?" for each component
```

Key signals:
- Component rendering > 16ms (blocks a frame at 60fps)
- Component rendering when its props haven't changed
- Hundreds of components rendering for a single state change

## Optimization Techniques (In Priority Order)

### 1. Move State Down

The most impactful optimization. State in a parent re-renders all children:

```tsx
// ❌ Bad: typing in search re-renders the entire product list
function ProductPage() {
  const [search, setSearch] = useState('');
  const [products] = useState(expensiveProductList);
  return (
    <div>
      <input value={search} onChange={e => setSearch(e.target.value)} />
      <ProductList products={products} />
    </div>
  );
}

// ✅ Good: isolate the search input into its own component
function SearchInput({ onSearch }: { onSearch: (q: string) => void }) {
  const [search, setSearch] = useState('');
  return <input value={search} onChange={e => {
    setSearch(e.target.value);
    onSearch(e.target.value);
  }} />;
}
```

### 2. Composition (Children Pattern)

```tsx
// ❌ Scrolling re-renders everything inside ScrollTracker
function ScrollTracker() {
  const [scrollY, setScrollY] = useState(0);
  useEffect(() => {
    const handler = () => setScrollY(window.scrollY);
    window.addEventListener('scroll', handler);
    return () => window.removeEventListener('scroll', handler);
  }, []);
  return (
    <div>
      <ScrollIndicator y={scrollY} />
      <ExpensiveContent /> {/* Re-renders on every scroll! */}
    </div>
  );
}

// ✅ Pass children — they're already rendered, won't re-render
function ScrollTracker({ children }: { children: ReactNode }) {
  const [scrollY, setScrollY] = useState(0);
  // ... same effect
  return (
    <div>
      <ScrollIndicator y={scrollY} />
      {children} {/* Stable reference, no re-render */}
    </div>
  );
}
```

### 3. useMemo / useCallback (After Profiling)

Only use these when profiling shows a real problem:

```tsx
// useMemo for expensive computations
const sortedItems = useMemo(
  () => items.slice().sort((a, b) => a.price - b.price),
  [items]
);

// useCallback for stable references passed to memoized children
const handleDelete = useCallback((id: string) => {
  setItems(prev => prev.filter(item => item.id !== id));
}, []);
```

### 4. Virtualization for Large Lists

Don't render 10,000 DOM nodes. Use `@tanstack/react-virtual`:

```tsx
import { useVirtualizer } from '@tanstack/react-virtual';

function VirtualList({ items }: { items: Item[] }) {
  const parentRef = useRef<HTMLDivElement>(null);
  const virtualizer = useVirtualizer({
    count: items.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 60,
    overscan: 5,
  });

  return (
    <div ref={parentRef} style={{ height: '600px', overflow: 'auto' }}>
      <div style={{ height: `${virtualizer.getTotalSize()}px`, position: 'relative' }}>
        {virtualizer.getVirtualItems().map(virtualRow => (
          <div
            key={virtualRow.key}
            style={{
              position: 'absolute',
              top: 0,
              transform: `translateY(${virtualRow.start}px)`,
              height: `${virtualRow.size}px`,
              width: '100%',
            }}
          >
            <ItemRow item={items[virtualRow.index]} />
          </div>
        ))}
      </div>
    </div>
  );
}
```

### 5. Code Splitting with lazy()

```tsx
const AdminPanel = lazy(() => import('./admin-panel'));
const Settings = lazy(() => import('./settings'));

function App() {
  return (
    <Suspense fallback={<Spinner />}>
      <Routes>
        <Route path="/admin" element={<AdminPanel />} />
        <Route path="/settings" element={<Settings />} />
      </Routes>
    </Suspense>
  );
}
```

## Bundle Analysis

```bash
# Next.js
ANALYZE=true next build

# Vite
npx vite-bundle-visualizer

# Generic
npx source-map-explorer 'build/static/js/*.js'
```

Look for: large dependencies used in few places, duplicate packages, polyfills you don't need.
