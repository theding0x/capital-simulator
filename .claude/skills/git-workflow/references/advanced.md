# Advanced Git Operations

> Load when: Interactive rebase, bisect, cherry-pick, hooks.

## Interactive Rebase

```bash
# Rewrite last 5 commits
git rebase -i HEAD~5

# In the editor:
pick abc1234 feat: add login
squash def5678 fix: login typo        # Merge into previous
reword ghi9012 feat: add dashboard    # Change message
drop jkl3456 WIP: temporary           # Remove commit
pick mno7890 feat: add settings

# Rebase onto main (before PR)
git fetch origin
git rebase origin/main
```

## Git Bisect

Find which commit introduced a bug:

```bash
git bisect start
git bisect bad                    # Current commit is broken
git bisect good v1.0.0            # This version was working
# Git checks out middle commit
# Test it, then:
git bisect good   # or
git bisect bad
# Repeat until Git finds the exact commit
git bisect reset  # Done, go back to original
```

## Git Hooks (with Husky)

```json
// package.json
{
  "scripts": {
    "prepare": "husky"
  }
}
```

```bash
# .husky/pre-commit
npm run lint-staged

# .husky/commit-msg
npx commitlint --edit $1
```