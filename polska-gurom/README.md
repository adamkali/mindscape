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

Deployment is automated by
[`.github/workflows/deploy-polska-gurom.yml`](../.github/workflows/deploy-polska-gurom.yml):
on a push to the default branch that touches `polska-gurom/**` (or via manual
run), it builds the site and publishes it to GitHub Pages.

The workflow sets `BASE_PATH=/<repo-name>` so links resolve correctly on a
project Pages site (`https://<user>.github.io/<repo>/`).

**One-time setup:** in the repository, go to **Settings → Pages → Build and
deployment → Source** and choose **GitHub Actions**.

## Editing content

All copy lives in [`src/lib/content.js`](src/lib/content.js); layout and styles
live in [`src/routes/+page.svelte`](src/routes/+page.svelte).
