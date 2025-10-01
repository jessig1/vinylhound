<script>
  import { createEventDispatcher } from "svelte";

  export let content = [];
  export let draft = "";
  export let loading = false;

  const dispatch = createEventDispatcher();

  function submit(event) {
    event.preventDefault();
    dispatch("save", { draft });
  }

  function refresh() {
    dispatch("refresh");
  }
</script>

<section class="panel">
  <h2>Your content</h2>
  <p class="hint">Each line will be saved as a separate item.</p>
  <form on:submit={submit} class="content-form">
    <textarea
      bind:value={draft}
      rows="8"
      placeholder="Add one entry per line"
      aria-label="Content entries"
    ></textarea>
    <div class="actions">
      <button type="submit" disabled={loading}>Save</button>
      <button type="button" on:click={refresh} disabled={loading}>Refresh</button>
    </div>
  </form>
  {#if content.length}
    <ul class="content-list">
      {#each content as item}
        <li>{item}</li>
      {/each}
    </ul>
  {:else}
    <p class="empty">No content yet. Add some above!</p>
  {/if}
</section>

<style>
  .content-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .actions {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .content-list {
    list-style: none;
    padding: 0;
    margin: 1rem 0 0;
    display: grid;
    gap: 0.75rem;
  }

  .content-list li {
    padding: 0.7rem 0.9rem;
    background: #eef2ff;
    border-radius: 0.7rem;
  }

  .empty {
    margin-top: 1rem;
    color: #64748b;
  }
</style>
