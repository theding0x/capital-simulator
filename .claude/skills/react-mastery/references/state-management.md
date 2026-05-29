# State Management

> Load when: Choosing a state management solution, implementing global state, managing server state.

## Decision Matrix

| Need | Solution | Why |
|------|----------|-----|
| Local UI state | `useState` / `useReducer` | Simplest. No overhead. |
| Shared between siblings | Lift state up or composition | Avoids unnecessary libraries |
| Theme/locale/auth | React Context | Changes rarely, affects many components |
| Server data (REST) | TanStack Query | Caching, deduplication, background refresh |
| Server data (GraphQL) | Apollo Client or urql | Protocol-specific caching |
| Complex client state | Zustand | Tiny, no boilerplate, works outside React |
| Atomic/granular | Jotai | Bottom-up, per-atom subscriptions |
| Form state | React Hook Form | Uncontrolled by default, minimal re-renders |

## Zustand — The Recommended Default

Zustand wins for most cases because it's simple, tiny (1KB), works outside React, and has no boilerplate:

```tsx
import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';

interface AppState {
  user: User | null;
  notifications: Notification[];
  setUser: (user: User | null) => void;
  addNotification: (n: Notification) => void;
  dismissNotification: (id: string) => void;
}

const useAppStore = create<AppState>()(
  devtools(
    persist(
      (set) => ({
        user: null,
        notifications: [],
        setUser: (user) => set({ user }, false, 'setUser'),
        addNotification: (n) =>
          set((s) => ({ notifications: [...s.notifications, n] }), false, 'addNotification'),
        dismissNotification: (id) =>
          set(
            (s) => ({ notifications: s.notifications.filter((n) => n.id !== id) }),
            false,
            'dismissNotification'
          ),
      }),
      { name: 'app-store' }
    )
  )
);

// Selector pattern — component only re-renders when selected value changes
function UserAvatar() {
  const user = useAppStore((s) => s.user);
  if (!user) return null;
  return <Avatar src={user.avatar} />;
}

// Use shallow compare for object selections
import { shallow } from 'zustand/shallow';
function NotificationBadge() {
  const { count, hasUnread } = useAppStore(
    (s) => ({ count: s.notifications.length, hasUnread: s.notifications.some(n => !n.read) }),
    shallow
  );
  return hasUnread ? <Badge count={count} /> : null;
}
```

## TanStack Query — Server State

Server state is fundamentally different from client state. It's async, has a source of truth you don't control, and can become stale:

```tsx
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

// Query with proper stale/cache times
function useUser(userId: string) {
  return useQuery({
    queryKey: ['user', userId],
    queryFn: () => api.getUser(userId),
    staleTime: 5 * 60 * 1000,     // Data fresh for 5 minutes
    gcTime: 30 * 60 * 1000,       // Keep in cache for 30 minutes
    retry: 2,
  });
}

// Mutation with optimistic update
function useUpdateUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: UpdateUserInput) => api.updateUser(data),
    onMutate: async (newData) => {
      await queryClient.cancelQueries({ queryKey: ['user', newData.id] });
      const previous = queryClient.getQueryData(['user', newData.id]);
      queryClient.setQueryData(['user', newData.id], (old: User) => ({
        ...old,
        ...newData,
      }));
      return { previous };
    },
    onError: (_err, newData, context) => {
      queryClient.setQueryData(['user', newData.id], context?.previous);
    },
    onSettled: (_data, _err, variables) => {
      queryClient.invalidateQueries({ queryKey: ['user', variables.id] });
    },
  });
}
```

## Anti-Patterns

```tsx
// ❌ Storing server data in useState — goes stale, no caching
const [users, setUsers] = useState([]);
useEffect(() => {
  fetch('/api/users').then(r => r.json()).then(setUsers);
}, []);

// ❌ Putting everything in a single context — all consumers re-render
const AppContext = createContext({ user, theme, locale, cart, notifications });

// ❌ Redux for simple apps — enormous boilerplate for minimal benefit
// Use Redux Toolkit only when you truly need time-travel debugging,
// middleware chains, or have 50+ interconnected state slices

// ❌ Zustand without selectors — defeats the purpose
const state = useAppStore(); // Re-renders on ANY state change
```

## When Context Is Enough

Context works well for infrequently-changing values that many components need:

```tsx
// Good use: theme that changes on user action, not on every render
const ThemeContext = createContext<Theme>(defaultTheme);

function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setTheme] = useState<Theme>(() => {
    return localStorage.getItem('theme') === 'dark' ? darkTheme : lightTheme;
  });
  const value = useMemo(() => ({ theme, setTheme }), [theme]);
  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}
```
