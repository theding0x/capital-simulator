# market-service

Models exchange and circulation — where commodities meet money. Owns the C-M-C and M-C-M' circuits, price formation, and trade matching.

- Default port: `8083`
- Persistence: MongoDB; Redis for hot order-book / price caches

Run locally:

```bash
make run-market-service
```
