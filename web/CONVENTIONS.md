# Frontend Conventions

These are the load-bearing conventions for the SolidJS frontend. **Read this before
adding a feature.** The goal is to keep logic centralized and DRY — most new work
should be a method on an existing context or a use of an existing primitive, not a
new copy of the same pattern.

> Rule of thumb: if you're about to write a `fetch`/mutation/auth-header in a
> component, stop — it almost certainly belongs in a context.

---

## 1. Contexts own shared state + mutations

Any state or API mutation used by more than one component lives in a **context**
(`createContext` + `ParentComponent` provider + `useX()` hook that throws if used
outside its provider). Components are presentational: they collect input and call
context methods.

Existing contexts and their domains — **extend these, don't reinvent them**:

| Context | Owns |
|---|---|
| `AuthContext` | session, token, login/logout |
| `TreeContext` | folders + bookmarks: focus/selection, create/update/delete/move, targeted refresh |
| `ViewContext` | agenda/tasks CRUD + filtering |
| `WidgetContext` | user widgets list, schemas, add |
| `BackgroundContext` | background image |
| `KeyboardNavContext` | nav mode + the tree node registry |

- Adding a create/update/delete/move? → **add a method to the relevant context.**
  Never re-implement a mutation inside a component.
- Keep **view-only state local**: form field signals, search text, `isDragging`,
  modal open/close. Those are not shared state and do not belong in a context.
- Refresh-after-mutation belongs with the mutation (in the context), not in the
  component that triggered it.

## 2. API access — one configured client, auth is automatic

- **Authenticated calls:** use `useAuthenticatedApi()` (returns ready-configured api
  instances), or `getAuthenticatedApiConfig()` when you must construct
  `new XxxApi(getAuthenticatedApiConfig())` inside a non-reactive function.
- **Never** `new XxxApi()` bare, and **never** `new Configuration({...})` by hand for
  authenticated calls — those instances miss the auth token *and* the 401 interceptor.
- **Do not pass an `authorization` argument.** The Bearer token is injected globally by
  the `apiKey` provider in `utils/apiConfig.ts` (the `BearerAuth` security scheme).
  Generated methods no longer accept it.
- Public endpoints only (login/signup): a bare `new UsersApi()` is fine.

## 3. Async data → `createResource`

Fetching that has a loading state uses `createResource` (it integrates with
`<Suspense>`); call `refetch()` after a mutation. Don't hand-roll
`createSignal` + `onMount(fetch)` — and never capture a value like a token once at
component init (it goes stale); read it inside the fetcher.

## 4. Reusable behavior → primitives + shared contracts

- Element-local reactive behavior reused by several components → a `createX` **primitive**
  in `hooks/` (e.g. `useDragNode`'s `useDraggableNode` / `useDropTarget`).
- Cross-component data contracts (serialized payloads, shared types) → `utils/`
  (e.g. `dragNode.ts`'s `DragNodePayload` + `serializeNode`/`parseNode`).
- The second time you write the same derivation, extract one helper — don't copy it
  (see `TreeContext`'s `parentOf` / `refreshFolderById`).

## 5. Auto-generated client — never edit

- `web/src/api/**` is generated from the OpenAPI spec. **Never hand-edit it.**
  Regenerate with `go-task swag` (spec + TS client together) after backend changes.
- Backend auth contract: annotate endpoints with `@Security BearerAuth` (not an
  `Authorization` `@Param`). The scheme is declared once in `main.go`.

## 6. SolidJS idioms

- `<For each>` over `.map` for lists.
- Derive with `createMemo`; don't duplicate computed state.
- Match the existing context shape (provider + `useX()` hook that throws when unbound).

---

## Pre-flight checklist before adding code

1. New mutation? → a method on a context, not in a component.
2. New shared state? → a context, not prop-drilling.
3. New API call? → configured client, **no** `authorization` arg.
4. Writing the same logic a second time? → extract a helper/primitive.
5. About to touch `web/src/api/**`? → stop; change the backend and regenerate.
