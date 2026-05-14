# REST Conventions

> Load when: URL design, HTTP methods, pagination, filtering.

## URL Design

```
GET    /api/v1/users              # List users
POST   /api/v1/users              # Create user
GET    /api/v1/users/:id          # Get user
PUT    /api/v1/users/:id          # Replace user
PATCH  /api/v1/users/:id          # Update user fields
DELETE /api/v1/users/:id          # Delete user

# Nested resources (when relationship is strong)
GET    /api/v1/users/:id/orders   # User's orders
POST   /api/v1/users/:id/orders   # Create order for user

# Actions (when CRUD doesn't fit)
POST   /api/v1/orders/:id/cancel
POST   /api/v1/users/:id/verify-email
```

## Pagination

```
# Cursor-based (recommended for large datasets)
GET /api/v1/users?limit=20&after=eyJpZCI6MTAwfQ

Response:
{
  "data": [...],
  "pagination": {
    "has_more": true,
    "next_cursor": "eyJpZCI6MTIwfQ"
  }
}

# Offset-based (simpler, OK for small datasets)
GET /api/v1/users?page=2&per_page=20

Response:
{
  "data": [...],
  "pagination": {
    "page": 2,
    "per_page": 20,
    "total": 156,
    "pages": 8
  }
}
```

## Filtering & Sorting

```
GET /api/v1/products?status=active&min_price=10&max_price=100
GET /api/v1/products?sort=-created_at,name  # - prefix = descending
GET /api/v1/products?fields=id,name,price   # Sparse fieldsets
```