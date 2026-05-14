# Server Components

> Load when: RSC architecture, streaming, Suspense boundaries, Server Actions.

## Mental Model

Server Components run ONLY on the server. They can access databases, file systems, and secrets directly. They send rendered HTML (not JavaScript) to the client. Client Components are the traditional React you know — they run in the browser and handle interactivity.

**Default to Server Components.** Only add `'use client'` when you need:
- Event handlers (onClick, onChange, etc.)
- State or lifecycle hooks (useState, useEffect)
- Browser-only APIs (localStorage, window, IntersectionObserver)
- Custom hooks that use any of the above

## Architecture Patterns

### The "Donut" Pattern

Server Component wraps Client Component, passing server-fetched data as props:

```tsx
// page.tsx (Server Component — no 'use client')
import { ProductGrid } from './product-grid';

export default async function ProductsPage() {
  const products = await db.products.findMany({ take: 50 });
  return (
    <main>
      <h1>Products</h1>
      <ProductGrid initialProducts={products} />
    </main>
  );
}

// product-grid.tsx (Client Component — needs interactivity)
'use client';
import { useState } from 'react';

export function ProductGrid({ initialProducts }: { initialProducts: Product[] }) {
  const [sortBy, setSortBy] = useState<'price' | 'name'>('name');
  const sorted = [...initialProducts].sort((a, b) =>
    sortBy === 'price' ? a.price - b.price : a.name.localeCompare(b.name)
  );
  return (
    <>
      <select onChange={(e) => setSortBy(e.target.value as any)}>
        <option value="name">Name</option>
        <option value="price">Price</option>
      </select>
      <div className="grid grid-cols-3 gap-4">
        {sorted.map(p => <ProductCard key={p.id} product={p} />)}
      </div>
    </>
  );
}
```

### Children as Server Content

Pass Server Components as children to Client Components:

```tsx
// layout.tsx (Server)
export default function Layout({ children }: { children: ReactNode }) {
  return (
    <Sidebar>           {/* Client Component */}
      <NavLinks />      {/* Server Component — rendered on server, passed as children */}
    </Sidebar>
  );
}
```

## Streaming with Suspense

Suspense boundaries let you stream parts of the page independently:

```tsx
export default async function Dashboard() {
  return (
    <div className="grid grid-cols-2 gap-6">
      <Suspense fallback={<MetricsSkeleton />}>
        <Metrics />     {/* Streams when ready */}
      </Suspense>
      <Suspense fallback={<ChartSkeleton />}>
        <RevenueChart /> {/* Streams independently */}
      </Suspense>
      <Suspense fallback={<TableSkeleton />}>
        <RecentOrders /> {/* Doesn't block above */}
      </Suspense>
    </div>
  );
}

async function Metrics() {
  const data = await analyticsApi.getMetrics(); // Takes 200ms
  return <MetricsDisplay data={data} />;
}

async function RevenueChart() {
  const data = await analyticsApi.getRevenue(); // Takes 2s — streams late
  return <Chart data={data} />;
}
```

## Server Actions

Server Actions are async functions that run on the server, invoked from the client:

```tsx
// actions.ts
'use server';

import { revalidatePath } from 'next/cache';
import { redirect } from 'next/navigation';

export async function createPost(formData: FormData) {
  const title = formData.get('title') as string;
  const content = formData.get('content') as string;

  // Validate
  if (!title || title.length < 3) {
    return { error: 'Title must be at least 3 characters' };
  }

  // Mutate
  const post = await db.posts.create({ data: { title, content } });

  // Revalidate and redirect
  revalidatePath('/posts');
  redirect(`/posts/${post.id}`);
}

// In a Client Component:
'use client';
import { useActionState } from 'react';
import { createPost } from './actions';

function CreatePostForm() {
  const [state, formAction, isPending] = useActionState(createPost, null);

  return (
    <form action={formAction}>
      <input name="title" required />
      <textarea name="content" required />
      {state?.error && <p className="text-red-500">{state.error}</p>}
      <button disabled={isPending}>
        {isPending ? 'Creating...' : 'Create Post'}
      </button>
    </form>
  );
}
```

## Anti-Patterns

- **Adding `'use client'` to everything** — defeats the purpose. You ship more JS and lose server-side data access.
- **Importing a Server Component into a Client Component** — impossible. The Server Component becomes a Client Component. Pass as children instead.
- **Large serializable props** — data crossing the server→client boundary must be JSON-serializable. Don't pass entire database rows; select only needed fields.
- **Fetching in Client Components when Server Components would work** — if the component doesn't need interactivity, keep it on the server.
