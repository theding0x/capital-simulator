# Testing React

> Load when: Testing components, hooks, async flows, integration tests.

## Philosophy

Test behavior, not implementation. If you refactor a component's internals and the tests break, the tests were wrong.

## Setup: Vitest + Testing Library

```tsx
// vitest.config.ts
import { defineConfig } from 'vitest/config';
export default defineConfig({
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    globals: true,
  },
});

// src/test/setup.ts
import '@testing-library/jest-dom';
```

## Component Testing

```tsx
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

describe('TodoList', () => {
  it('adds a new todo when form is submitted', async () => {
    const user = userEvent.setup();
    render(<TodoList />);

    const input = screen.getByPlaceholderText('Add a todo');
    await user.type(input, 'Buy groceries');
    await user.click(screen.getByRole('button', { name: /add/i }));

    expect(screen.getByText('Buy groceries')).toBeInTheDocument();
    expect(input).toHaveValue(''); // Input cleared after submit
  });

  it('marks a todo as complete', async () => {
    const user = userEvent.setup();
    render(<TodoList initialTodos={[{ id: '1', text: 'Test', done: false }]} />);

    await user.click(screen.getByRole('checkbox'));
    expect(screen.getByText('Test')).toHaveClass('line-through');
  });
});
```

## Hook Testing

```tsx
import { renderHook, act } from '@testing-library/react';

describe('useCounter', () => {
  it('increments count', () => {
    const { result } = renderHook(() => useCounter(0));

    act(() => { result.current.increment(); });
    expect(result.current.count).toBe(1);

    act(() => { result.current.increment(); });
    expect(result.current.count).toBe(2);
  });
});

// Hooks with context need wrappers
describe('useAuth', () => {
  const wrapper = ({ children }: { children: ReactNode }) => (
    <AuthProvider>{children}</AuthProvider>
  );

  it('returns user after login', async () => {
    const { result } = renderHook(() => useAuth(), { wrapper });
    await act(async () => {
      await result.current.login('user@test.com', 'password');
    });
    expect(result.current.user?.email).toBe('user@test.com');
  });
});
```

## Async Testing

```tsx
// Waiting for data to load
it('displays user data after loading', async () => {
  render(<UserProfile userId="123" />);

  // Shows loading state
  expect(screen.getByText(/loading/i)).toBeInTheDocument();

  // Wait for data to appear
  const userName = await screen.findByText('John Doe');
  expect(userName).toBeInTheDocument();

  // Loading state is gone
  expect(screen.queryByText(/loading/i)).not.toBeInTheDocument();
});
```

## API Mocking with MSW

```tsx
import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';

const server = setupServer(
  http.get('/api/users/:id', ({ params }) => {
    return HttpResponse.json({ id: params.id, name: 'John Doe' });
  }),
  http.post('/api/users', async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json({ id: '123', ...body }, { status: 201 });
  })
);

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

// Override for specific test
it('shows error when API fails', async () => {
  server.use(
    http.get('/api/users/:id', () => {
      return new HttpResponse(null, { status: 500 });
    })
  );
  render(<UserProfile userId="123" />);
  expect(await screen.findByText(/error/i)).toBeInTheDocument();
});
```

## What NOT to Test

- Implementation details (state variable names, internal methods)
- Third-party library behavior (test YOUR code, not React's)
- Snapshot tests of large components (brittle, low signal)
- CSS styling (use visual regression tools instead)
