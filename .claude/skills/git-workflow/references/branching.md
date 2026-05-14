# Branching Strategies

> Load when: Choosing a branching strategy, setting up workflows.

## Trunk-Based (Recommended)

```
main ─────●────●────●────●────●────●
           \  /      \  /      \  /
feature-a ─●      feature-b─●    feature-c─●
```

- Short-lived feature branches (1-3 days max)
- Merge to main frequently
- Use feature flags for incomplete features
- CI runs on every push to main

## GitHub Flow (Simple)

```
main ─────●────────●────────●
           \        |        |
feature ───●──●──●─PR──merge
```

- Branch from main, PR back to main
- Good for continuous deployment
- No release branches needed

## Commit Message Convention

```
feat(auth): add OAuth2 login with Google
fix(api): handle null response from payment gateway
refactor(db): extract query builder into separate module
docs(readme): add deployment instructions
chore(deps): upgrade express to 4.19.2

# Breaking changes
feat(api)!: change response format for /users endpoint

BREAKING CHANGE: response.data is now an array instead of object
```