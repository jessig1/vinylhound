<script>
  import { createEventDispatcher } from "svelte";

  export let loading = false;

  const dispatch = createEventDispatcher();
  const signupPlaceholder = "Welcome to Vinyhound!\nMake your first playlist.";

  let signupUsername = "";
  let signupPassword = "";
  let signupContent = "";
  let loginUsername = "demo";
  let loginPassword = "demo123";

  function submitSignup(event) {
    event.preventDefault();
    dispatch("signup", {
      username: signupUsername.trim(),
      password: signupPassword,
      content: signupContent,
    });
  }

  function submitLogin(event) {
    event.preventDefault();
    dispatch("login", {
      username: loginUsername.trim(),
      password: loginPassword,
    });
  }

  export function resetSignup() {
    signupUsername = "";
    signupPassword = "";
    signupContent = "";
  }

  export function resetLogin() {
    loginPassword = "";
  }
</script>

<div class="auth-grid">
  <section class="panel">
    <h2>Create an account</h2>
    <form on:submit={submitSignup}>
      <label>
        Username
        <input
          type="text"
          bind:value={signupUsername}
          placeholder="alice"
          required
          autocomplete="username"
        />
      </label>
      <label>
        Password
        <input
          type="password"
          bind:value={signupPassword}
          required
          autocomplete="new-password"
        />
      </label>
      <label>
        Initial content <span class="hint">(one per line)</span>
        <textarea
          rows="4"
          bind:value={signupContent}
          placeholder={signupPlaceholder}
        ></textarea>
      </label>
      <button type="submit" disabled={loading}>Sign up</button>
    </form>
  </section>

  <section class="panel">
    <h2>Log in</h2>
    <form on:submit={submitLogin}>
      <label>
        Username
        <input
          type="text"
          bind:value={loginUsername}
          placeholder="demo"
          required
          autocomplete="username"
        />
      </label>
      <label>
        Password
        <input
          type="password"
          bind:value={loginPassword}
          placeholder="demo123"
          required
          autocomplete="current-password"
        />
      </label>
      <button type="submit" disabled={loading}>Log in</button>
    </form>
    <p class="hint">A demo account (demo / demo123) is available.</p>
  </section>
</div>

<style>
  .auth-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1.5rem;
  }
</style>
