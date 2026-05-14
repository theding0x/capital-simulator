# Hooks Patterns

> Load when: Building custom hooks, composing hooks, optimizing hook usage.

## Custom Hook Design

A good custom hook extracts **a single concern** with a clean return interface:

```tsx
// Good: Single concern, clear interface
function useDebounce<T>(value: T, delay: number): T {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const timer = setTimeout(() => setDebounced(value), delay);
    return () => clearTimeout(timer);
  }, [value, delay]);
  return debounced;
}

// Good: Composing hooks for complex behavior
function useSearch(query: string) {
  const debouncedQuery = useDebounce(query, 300);
  const { data, isLoading, error } = useSWR(
    debouncedQuery ? `/api/search?q=${debouncedQuery}` : null,
    fetcher
  );
  return { results: data, isLoading, error, isDebouncing: query !== debouncedQuery };
}
```

## Rules of Hooks — The Why

Hooks must be called in the same order every render because React tracks them by call index, not by name. This means:

- No hooks inside conditions, loops, or early returns
- No hooks after a conditional return
- Custom hooks must start with `use` (this enables the linter to enforce the rules)

## useEffect — When You Actually Need It

Most things people put in `useEffect` belong somewhere else:

```tsx
// ❌ Wrong: Deriving state in effect
const [fullName, setFullName] = useState('');
useEffect(() => {
  setFullName(`${firstName} ${lastName}`);
}, [firstName, lastName]);

// ✅ Right: Derive during render
const fullName = `${firstName} ${lastName}`;

// ❌ Wrong: Resetting state on prop change
useEffect(() => {
  setSelection(null);
}, [items]);

// ✅ Right: Use key to reset component
<ItemList key={categoryId} items={items} />

// ✅ Legitimate useEffect: Synchronizing with external system
useEffect(() => {
  const connection = createConnection(serverUrl);
  connection.connect();
  return () => connection.disconnect();
}, [serverUrl]);
```

## useRef Patterns

Refs are for values that don't affect rendering:

```tsx
// Previous value tracking
function usePrevious<T>(value: T): T | undefined {
  const ref = useRef<T | undefined>(undefined);
  useEffect(() => { ref.current = value; });
  return ref.current;
}

// Stable callback (avoids re-renders from changing callbacks)
function useStableCallback<T extends (...args: any[]) => any>(callback: T): T {
  const ref = useRef(callback);
  useLayoutEffect(() => { ref.current = callback; });
  return useCallback((...args: any[]) => ref.current(...args), []) as T;
}

// Interval that respects latest closure values
function useInterval(callback: () => void, delay: number | null) {
  const savedCallback = useRef(callback);
  useEffect(() => { savedCallback.current = callback; });
  useEffect(() => {
    if (delay === null) return;
    const id = setInterval(() => savedCallback.current(), delay);
    return () => clearInterval(id);
  }, [delay]);
}
```

## React 19: `use()` Hook

The `use` hook can read promises and context conditionally (unlike other hooks):

```tsx
// Reading a promise — works with Suspense
function Comments({ commentsPromise }: { commentsPromise: Promise<Comment[]> }) {
  const comments = use(commentsPromise);
  return comments.map(c => <Comment key={c.id} comment={c} />);
}

// Conditional context reading
function StatusIcon({ isAdmin }: { isAdmin: boolean }) {
  if (isAdmin) {
    const theme = use(ThemeContext); // This is allowed with use()
    return <Icon color={theme.primary} />;
  }
  return <DefaultIcon />;
}
```

## Composition Pattern: Reducer + Context

For complex state shared across a component tree:

```tsx
type Action = 
  | { type: 'ADD_ITEM'; payload: Item }
  | { type: 'REMOVE_ITEM'; payload: string }
  | { type: 'UPDATE_ITEM'; payload: { id: string; changes: Partial<Item> } };

function cartReducer(state: CartState, action: Action): CartState {
  switch (action.type) {
    case 'ADD_ITEM':
      return { ...state, items: [...state.items, action.payload] };
    case 'REMOVE_ITEM':
      return { ...state, items: state.items.filter(i => i.id !== action.payload) };
    case 'UPDATE_ITEM':
      return {
        ...state,
        items: state.items.map(i =>
          i.id === action.payload.id ? { ...i, ...action.payload.changes } : i
        ),
      };
  }
}

// Split context to prevent unnecessary re-renders
const CartStateContext = createContext<CartState>(initialState);
const CartDispatchContext = createContext<Dispatch<Action>>(() => {});

function CartProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(cartReducer, initialState);
  return (
    <CartStateContext.Provider value={state}>
      <CartDispatchContext.Provider value={dispatch}>
        {children}
      </CartDispatchContext.Provider>
    </CartStateContext.Provider>
  );
}

// Consumer hooks
const useCartState = () => useContext(CartStateContext);
const useCartDispatch = () => useContext(CartDispatchContext);
```
