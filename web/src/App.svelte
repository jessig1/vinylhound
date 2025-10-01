<script>
  import { onMount } from "svelte";
  import Header from "./components/Header.svelte";
  import FlashMessage from "./components/FlashMessage.svelte";
  import AuthPanel from "./components/AuthPanel.svelte";
  import ContentPanel from "./components/ContentPanel.svelte";
  import {
    signup as apiSignup,
    login as apiLogin,
    fetchContent,
    updateContent,
    ApiError,
  } from "./lib/api";
  import { normalizeContent } from "./lib/content";
  import { readSession, writeSession, clearSession } from "./lib/session";

  let token = "";
  let activeUser = "";
  let content = [];
  let contentDraft = "";
  let loading = false;
  let message = "";
  let messageKind = "info";

  let authPanel;

  onMount(() => {
    const stored = readSession();
    if (!stored) {
      return;
    }
    token = stored.token;
    activeUser = stored.username;
    loadContent({ silent: true });
  });

  function setMessage(value = "", kind = "info") {
    message = value;
    messageKind = kind;
  }

  function clearMessage() {
    setMessage();
  }

  async function execute(task, fallbackMessage) {
    loading = true;
    try {
      await task();
      return true;
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) {
        logout();
        setMessage("Session expired. Please log in again.", "error");
      } else {
        setMessage(err?.message || fallbackMessage, "error");
      }
      return false;
    } finally {
      loading = false;
    }
  }

  function applySession(newToken, username) {
    token = newToken;
    activeUser = username;
    if (newToken && username) {
      writeSession({ token: newToken, username });
    } else {
      clearSession();
    }
  }

  function logout() {
    applySession("", "");
    content = [];
    contentDraft = "";
  }

  async function handleSignup(event) {
    const detail = event.detail || {};
    const username = (detail.username || "").trim();
    const password = detail.password || "";

    if (!username || !password) {
      setMessage("Username and password are required.", "error");
      return;
    }

    clearMessage();
    const succeeded = await execute(async () => {
      const initialContent = normalizeContent(detail.content || "");
      await apiSignup({ username, password, content: initialContent });
    }, "Unable to sign up.");

    if (succeeded) {
      setMessage("Account created for " + username + ". You can log in now.", "success");
      authPanel?.resetSignup?.();
    }
  }

  async function handleLogin(event) {
    const detail = event.detail || {};
    const username = (detail.username || "").trim();
    const password = detail.password || "";

    if (!username || !password) {
      setMessage("Username and password are required.", "error");
      return;
    }

    clearMessage();
    const succeeded = await execute(async () => {
      const data = await apiLogin({ username, password });
      const tokenValue = data?.token || "";
      if (!tokenValue) {
        throw new Error("Login response missing token.");
      }
      applySession(tokenValue, username);
      authPanel?.resetLogin?.();
    }, "Unable to log in.");

    if (succeeded) {
      await loadContent({ silent: true });
      setMessage("Welcome back, " + username + "!", "success");
    }
  }

  async function loadContent({ silent = true } = {}) {
    if (!token) {
      return;
    }
    if (!silent) {
      clearMessage();
    }

    const succeeded = await execute(async () => {
      const data = await fetchContent(token);
      const items = Array.isArray(data?.content) ? data.content : [];
      content = items;
      contentDraft = items.join("\n");
    }, "Unable to load content.");

    if (succeeded && !silent) {
      setMessage("Content refreshed.", "success");
    }
  }

  async function handleSaveContent(event) {
    if (!token) {
      setMessage("You need to log in before updating content.", "error");
      return;
    }

    clearMessage();
    const draftValue = event.detail?.draft ?? contentDraft;
    const entries = normalizeContent(draftValue);

    const succeeded = await execute(async () => {
      await updateContent(token, entries);
    }, "Unable to save content.");

    if (succeeded) {
      content = entries;
      contentDraft = entries.join("\n");
      setMessage("Content updated.", "success");
    }
  }

  async function handleRefreshContent() {
    await loadContent({ silent: false });
  }
</script>

<main>
  <Header {token} {activeUser} on:logout={logout} />

  <FlashMessage {message} kind={messageKind} />

  {#if token}
    <ContentPanel
      {content}
      bind:draft={contentDraft}
      {loading}
      on:save={handleSaveContent}
      on:refresh={handleRefreshContent}
    />
  {:else}
    <AuthPanel
      bind:this={authPanel}
      {loading}
      on:signup={handleSignup}
      on:login={handleLogin}
    />
  {/if}
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: "Segoe UI", Roboto, sans-serif;
    background: radial-gradient(circle at top left, #f5f5ff, #eef2f7 45%, #dde4ee);
    min-height: 100vh;
    color: #1f2933;
  }

  main {
    max-width: 960px;
    margin: 0 auto;
    padding: 2.5rem 1.5rem 4rem;
  }

  :global(button) {
    background: #4f46e5;
    border: none;
    color: #ffffff;
    padding: 0.6rem 1.2rem;
    border-radius: 0.5rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s ease, transform 0.2s ease;
  }

  :global(button:hover:not(:disabled)) {
    background: #4338ca;
    transform: translateY(-1px);
  }

  :global(button:disabled) {
    background: #94a3b8;
    cursor: wait;
  }

  :global(.panel) {
    background: #ffffff;
    border-radius: 1rem;
    box-shadow: 0 12px 30px rgba(79, 70, 229, 0.08);
    padding: 2rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  :global(form) {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  :global(label) {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    font-weight: 600;
    font-size: 0.95rem;
  }

  :global(input),
  :global(textarea) {
    border: 1px solid #cbd5e1;
    border-radius: 0.6rem;
    padding: 0.7rem 0.85rem;
    font-size: 1rem;
    font-family: inherit;
    resize: vertical;
  }

  :global(input:focus),
  :global(textarea:focus) {
    outline: none;
    border-color: #4f46e5;
    box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.15);
  }

  :global(.hint) {
    font-size: 0.85rem;
    color: #64748b;
    font-weight: 400;
  }
</style>
