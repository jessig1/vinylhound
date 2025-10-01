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

  let token = "";
  let activeUser = "";
  let content = [];
  let contentDraft = "";
  let loading = false;
  let message = "";
  let messageKind = "info";

  let authPanel;

  onMount(() => {
    if (typeof window === "undefined") {
      return;
    }
    const storedToken = window.localStorage.getItem("vinyhound:token");
    const storedUser = window.localStorage.getItem("vinyhound:user");
    if (storedToken && storedUser) {
      token = storedToken;
      activeUser = storedUser;
      loadContent(false);
    }
  });

  function normalizeContent(value) {
    if (!value) {
      return [];
    }
    return value
      .split(/\r?\n/)
      .map(function (entry) {
        return entry.trim();
      })
      .filter(Boolean);
  }

  function persistSession() {
    if (typeof window === "undefined") {
      return;
    }
    window.localStorage.setItem("vinyhound:token", token);
    window.localStorage.setItem("vinyhound:user", activeUser);
  }

  function clearSession() {
    if (typeof window === "undefined") {
      return;
    }
    window.localStorage.removeItem("vinyhound:token");
    window.localStorage.removeItem("vinyhound:user");
  }

  function logout() {
    token = "";
    activeUser = "";
    content = [];
    contentDraft = "";
    clearSession();
  }

  function clearMessage() {
    message = "";
    messageKind = "info";
  }

  function setMessage(value, kind) {
    message = value;
    messageKind = kind || "info";
  }

  async function handleSignup(event) {
    const detail = event.detail || {};
    const username = (detail.username || "").trim();
    const password = detail.password || "";

    if (!username || !password) {
      setMessage("Username and password are required.", "error");
      return;
    }

    loading = true;
    clearMessage();
    try {
      const contentPayload = normalizeContent(detail.content || "");
      await apiSignup({ username, password, content: contentPayload });
      setMessage("Account created for " + username + ". You can log in now.", "success");
      if (authPanel && authPanel.resetSignup) {
        authPanel.resetSignup();
      }
    } catch (err) {
      setMessage(err.message || "Unable to sign up.", "error");
    } finally {
      loading = false;
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

    loading = true;
    clearMessage();
    try {
      const data = await apiLogin({ username, password });
      if (authPanel && authPanel.resetLogin) {
        authPanel.resetLogin();
      }
      token = data && data.token ? data.token : "";
      activeUser = username;
      persistSession();
      await loadContent(false);
      setMessage("Welcome back, " + activeUser + "!", "success");
    } catch (err) {
      setMessage(err.message || "Unable to log in.", "error");
    } finally {
      loading = false;
    }
  }

  async function loadContent(showSuccess) {
    if (!token) {
      return;
    }

    loading = true;
    try {
      const data = await fetchContent(token);
      const items =
        data && typeof data === "object" && Array.isArray(data.content)
          ? data.content
          : [];
      content = items;
      contentDraft = items.join("\n");
      if (showSuccess) {
        setMessage("Content refreshed.", "success");
      }
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) {
        logout();
        setMessage("Session expired. Please log in again.", "error");
        return;
      }
      setMessage(err.message || "Unable to load content.", "error");
    } finally {
      loading = false;
    }
  }

  async function handleSaveContent(event) {
    if (!token) {
      setMessage("You need to log in before updating content.", "error");
      return;
    }

    loading = true;
    clearMessage();
    try {
      const detail = event.detail || {};
      const draftValue = detail.draft !== undefined ? detail.draft : contentDraft;
      const entries = normalizeContent(draftValue);
      await updateContent(token, entries);
      content = entries;
      contentDraft = entries.join("\n");
      setMessage("Content updated.", "success");
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) {
        logout();
        setMessage("Session expired. Please log in again.", "error");
        return;
      }
      setMessage(err.message || "Unable to save content.", "error");
    } finally {
      loading = false;
    }
  }

  async function handleRefreshContent() {
    clearMessage();
    await loadContent(true);
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
