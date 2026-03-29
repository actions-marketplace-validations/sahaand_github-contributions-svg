# github-contributions-svg

A GitHub composite action that generates SVG assets for your GitHub profile README:

- **Stats card** — stars, commits, PRs, issues
- **Streak card** — current and longest contribution streak
- **Trophies card** — GitHub achievement trophies
- **Top languages card** — language breakdown by code size
- **Snake animation** — animated snake eating your contribution grid
- **3D contribution graph** — isometric 3D view with language pie and contribution radar charts

---

<div align="center">

[![Snake animation](https://raw.githubusercontent.com/saeedata/saeedata/master/assets/snake.svg)]()

</div>

<div align="center">

[![3D Contribution Graph](https://raw.githubusercontent.com/saeedata/saeedata/master/assets/contrib-3d.svg)]()

</div>

<div align="center">

[![Trophies](https://raw.githubusercontent.com/saeedata/saeedata/master/assets/trophies.svg)]()

</div>

<div align="center">

[![Streak](https://raw.githubusercontent.com/saeedata/saeedata/master/assets/streak.svg)]()

</div>

<div align="center">

[![GitHub Stats](https://raw.githubusercontent.com/saeedata/saeedata/master/assets/stats.svg)]()

</div>

<div align="center">

[![Top Languages](https://raw.githubusercontent.com/saeedata/saeedata/master/assets/top-langs.svg)]()

</div>

---

## Usage

Add this step to your workflow **after** `actions/checkout`:

```yaml
- uses: saeedata/github-contributions-svg@v1
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
```

The action writes the following files into the `assets/` directory of your workspace:

| File | Description |
|------|-------------|
| `assets/stats.svg` | Stats card |
| `assets/streak.svg` | Streak card |
| `assets/trophies.svg` | Trophies card |
| `assets/top-langs.svg` | Top languages card |
| `assets/snake.svg` | Snake animation |
| `assets/contrib-3d.svg` | 3D contribution graph |

### Full workflow example

```yaml
name: Generate README Assets

on:
  schedule:
    - cron: "0 5 * * *"   # daily at 05:00 UTC
  workflow_dispatch:

jobs:
  generate:
    runs-on: ubuntu-latest
    permissions:
      contents: write

    steps:
      - uses: actions/checkout@v4

      - uses: saeedata/github-contributions-svg@v1
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}

      - name: Commit and push updated assets
        run: |
          git config user.name  "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git add assets/
          git diff --staged --quiet \
            && echo "No changes." \
            || git commit -m "chore: regenerate README assets [skip ci]" && git push
```

Then reference the generated SVGs in your README:

```markdown
![Snake](https://raw.githubusercontent.com/<user>/<repo>/main/assets/snake.svg)
![Stats](https://raw.githubusercontent.com/<user>/<repo>/main/assets/stats.svg)
![Streak](https://raw.githubusercontent.com/<user>/<repo>/main/assets/streak.svg)
![3D Graph](https://raw.githubusercontent.com/<user>/<repo>/main/assets/contrib-3d.svg)
```

---

## Inputs

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `github-token` | Yes | — | GitHub token for API access. Use a PAT with `repo` scope to include private contributions. |
| `username` | No | `${{ github.repository_owner }}` | GitHub username to generate assets for. |
| `output-dir` | No | `assets` | Directory (relative to workspace root) where SVG files are written. |

### Private contributions

`GITHUB_TOKEN` only has access to the current repository. To count contributions from private
repos, create a [Personal Access Token](https://github.com/settings/tokens) with the `repo` scope,
add it as a repository secret (e.g. `GH_PAT`), and pass it instead:

```yaml
- uses: saeedata/github-contributions-svg@v1
  with:
    github-token: ${{ secrets.GH_PAT }}
```

---

## Running locally

```bash
export GITHUB_TOKEN=your_token
export GITHUB_USERNAME=your_username

go run ./scripts/stats
go run ./scripts/snake --output=assets/snake.svg
go run ./scripts/contrib3d --output=assets/contrib-3d.svg
```

Requires Go 1.22+. No external dependencies.
