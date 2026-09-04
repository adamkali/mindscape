<script>
  import { phrase, meanings, pronunciation, usage, examples } from '$lib/content.js';

  const year = new Date().getFullYear();
</script>

<svelte:head>
  <title>Polska górą — what "Polska gurom" means</title>
</svelte:head>

<div class="page">
  <header class="hero">
    <div class="flag" aria-hidden="true"></div>
    <p class="eyebrow">A Polish phrase, explained</p>
    <h1>Polska&nbsp;górą!</h1>
    <p class="tagline">
      also spelled <strong>“{phrase.colloquial}”</strong> — literally
      <em>“Poland on top!”</em>
    </p>
    <p class="respell">{phrase.ipa} · <span>{phrase.respell}</span></p>
  </header>

  <main>
    <section class="meanings" aria-label="Meaning">
      {#each meanings as m}
        <article class="card">
          <h2>{m.heading}</h2>
          <p>{m.body}</p>
        </article>
      {/each}
    </section>

    <section class="block" aria-label="Pronunciation">
      <h2>How to say it</h2>
      <dl class="pron">
        {#each pronunciation as p}
          <div class="pron-row">
            <dt>{p.label}</dt>
            <dd>{p.value}</dd>
          </div>
        {/each}
      </dl>
      <p class="note">
        The tricky part: Polish <b>ó</b> sounds like English “oo”, and a final
        <b>-ą</b> is often pronounced “-om”. Say the sounds and “górą” comes out
        as “gu-rom”.
      </p>
    </section>

    <section class="block" aria-label="Where it is used">
      <h2>Where you'll hear it</h2>
      <ul class="usage">
        {#each usage as u}
          <li>
            <span class="tag">{u.where}</span>
            <span>{u.detail}</span>
          </li>
        {/each}
      </ul>
    </section>

    <section class="block" aria-label="Examples">
      <h2>In a sentence</h2>
      <ul class="examples">
        {#each examples as ex}
          <li>
            <span class="pl" lang="pl">{ex.pl}</span>
            <span class="en">{ex.en}</span>
          </li>
        {/each}
      </ul>
    </section>
  </main>

  <footer>
    <p>
      “Polska górą” · a small static site built with
      <a href="https://svelte.dev" rel="noopener">Svelte</a> and served on
      GitHub&nbsp;Pages.
    </p>
    <p class="fine">© {year} · Content is informational.</p>
  </footer>
</div>

<style>
  :global(:root) {
    --red: #dc143c;
    --red-dark: #b01030;
    --ink: #1a1a1f;
    --muted: #5a5a66;
    --paper: #ffffff;
    --wash: #fbf4f5;
    --line: #eadfe1;
  }

  :global(html) {
    scroll-behavior: smooth;
  }

  :global(body) {
    margin: 0;
    background: var(--wash);
    color: var(--ink);
    font-family: system-ui, -apple-system, "Segoe UI", Roboto, Helvetica, Arial,
      sans-serif;
    line-height: 1.6;
  }

  .page {
    max-width: 780px;
    margin: 0 auto;
    padding: 0 1.25rem 4rem;
  }

  /* Hero */
  .hero {
    text-align: center;
    padding: 3.5rem 0 2.5rem;
  }

  .flag {
    width: 88px;
    height: 58px;
    margin: 0 auto 1.5rem;
    border-radius: 6px;
    overflow: hidden;
    box-shadow: 0 6px 20px rgba(176, 16, 48, 0.25);
    background: linear-gradient(to bottom, #fff 0 50%, var(--red) 50% 100%);
    border: 1px solid var(--line);
  }

  .eyebrow {
    text-transform: uppercase;
    letter-spacing: 0.14em;
    font-size: 0.72rem;
    font-weight: 600;
    color: var(--red-dark);
    margin: 0 0 0.5rem;
  }

  h1 {
    font-size: clamp(2.6rem, 9vw, 4.2rem);
    line-height: 1.05;
    margin: 0;
    letter-spacing: -0.02em;
  }

  .tagline {
    font-size: 1.15rem;
    color: var(--muted);
    margin: 0.75rem 0 0.35rem;
  }

  .respell {
    font-family: ui-monospace, "SFMono-Regular", Menlo, Consolas, monospace;
    font-size: 0.9rem;
    color: var(--muted);
  }
  .respell span {
    color: var(--ink);
  }

  h2 {
    font-size: 1.15rem;
    margin: 0 0 0.5rem;
  }

  /* Meaning cards */
  .meanings {
    display: grid;
    gap: 1rem;
    margin-bottom: 2rem;
  }

  .card {
    background: var(--paper);
    border: 1px solid var(--line);
    border-left: 4px solid var(--red);
    border-radius: 10px;
    padding: 1.25rem 1.4rem;
  }
  .card p {
    margin: 0;
    color: #33333b;
  }

  .block {
    background: var(--paper);
    border: 1px solid var(--line);
    border-radius: 10px;
    padding: 1.4rem 1.5rem;
    margin-bottom: 1.5rem;
  }

  /* Pronunciation */
  .pron {
    margin: 0;
  }
  .pron-row {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.55rem 0;
    border-bottom: 1px dashed var(--line);
  }
  .pron-row:last-child {
    border-bottom: none;
  }
  .pron dt {
    color: var(--muted);
    font-size: 0.9rem;
  }
  .pron dd {
    margin: 0;
    font-weight: 600;
    text-align: right;
  }
  .note {
    margin: 1rem 0 0;
    font-size: 0.95rem;
    color: var(--muted);
  }

  /* Usage */
  .usage {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    gap: 0.9rem;
  }
  .usage li {
    display: grid;
    grid-template-columns: 120px 1fr;
    gap: 0.9rem;
    align-items: start;
  }
  .tag {
    display: inline-block;
    background: var(--wash);
    border: 1px solid var(--line);
    color: var(--red-dark);
    font-size: 0.8rem;
    font-weight: 600;
    padding: 0.15rem 0.6rem;
    border-radius: 999px;
    text-align: center;
  }

  /* Examples */
  .examples {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    gap: 0.9rem;
  }
  .examples li {
    padding: 0.75rem 1rem;
    background: var(--wash);
    border-radius: 8px;
  }
  .examples .pl {
    display: block;
    font-weight: 600;
  }
  .examples .en {
    display: block;
    color: var(--muted);
    font-size: 0.95rem;
  }

  footer {
    text-align: center;
    color: var(--muted);
    font-size: 0.9rem;
    margin-top: 2.5rem;
  }
  footer a {
    color: var(--red-dark);
  }
  .fine {
    font-size: 0.8rem;
    opacity: 0.8;
  }

  @media (max-width: 520px) {
    .usage li {
      grid-template-columns: 1fr;
      gap: 0.35rem;
    }
    .pron-row {
      flex-direction: column;
      gap: 0.1rem;
    }
    .pron dd {
      text-align: left;
    }
  }
</style>
