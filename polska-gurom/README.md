# Polska górą — a tiny GitHub Pages site

A small static site that explains the Polish phrase **“Polska górą”** (often
written phonetically as **“Polska gurom”**) — what it means, how to pronounce
it, and where it's used.

Built with [SvelteKit](https://svelte.dev/) and
[`adapter-static`](https://kit.svelte.dev/docs/adapter-static), so it compiles
to plain HTML/CSS/JS that GitHub Pages can serve.

## Develop locally

```bash
cd polska-gurom
npm install
npm run dev      # http://localhost:5173
```

## Build

```bash
npm run build    # outputs to ./build
npm run preview  # serve the production build locally
```

## Deploy

Deployment uses a GitHub Actions workflow that, on a push to the default branch
touching `polska-gurom/**` (or via manual run), builds the site and publishes it
to GitHub Pages. It sets `BASE_PATH=/<repo-name>` so links resolve on a project
Pages site (`https://<user>.github.io/<repo>/`).

The workflow file is provided here as
[`deploy-workflow.yml`](deploy-workflow.yml) **but not installed automatically**:
the automation account that pushed this branch does not have the GitHub
`workflows` permission, so it cannot write into `.github/workflows/`. A repo
admin adds it once:

```bash
mkdir -p .github/workflows
cp polska-gurom/deploy-workflow.yml .github/workflows/deploy-polska-gurom.yml
git add .github/workflows/deploy-polska-gurom.yml
git commit -m "Add Polska górą Pages deploy workflow"
git push
```

(Or, in the GitHub UI: **Add file → Create new file**, path
`.github/workflows/deploy-polska-gurom.yml`, and paste the contents of
`deploy-workflow.yml`.)

**Then, one time:** in the repository go to **Settings → Pages → Build and
deployment → Source** and choose **GitHub Actions**.

## Editing content

All copy lives in [`src/lib/content.js`](src/lib/content.js); layout and styles
live in [`src/routes/+page.svelte`](src/routes/+page.svelte).
