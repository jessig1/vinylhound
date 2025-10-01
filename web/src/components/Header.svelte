<script>
  import { createEventDispatcher } from "svelte";

  export let token = "";
  export let activeUser = "";

  const dispatch = createEventDispatcher();

  function handleLogout() {
    dispatch("logout");
  }
</script>

<header class="app-header">
  <h1>Vinyhound</h1>
  {#if token}
    <div class="session">
      <span>Signed in as <strong>{activeUser}</strong></span>
      <button type="button" on:click={handleLogout}>Log out</button>
    </div>
  {/if}
</header>

<style>
  .app-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  h1 {
    font-size: 2.5rem;
    margin: 0;
  }

  .session {
    display: flex;
    gap: 1rem;
    align-items: center;
    font-size: 0.95rem;
  }

  @media (max-width: 640px) {
    .app-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 1rem;
    }
  }
</style>
