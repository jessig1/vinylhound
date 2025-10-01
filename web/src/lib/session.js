const TOKEN_KEY = "vinyhound:token";
const USER_KEY = "vinyhound:user";

function hasWindow() {
  return typeof window !== "undefined" && typeof window.localStorage !== "undefined";
}

export function readSession() {
  if (!hasWindow()) {
    return null;
  }
  const token = window.localStorage.getItem(TOKEN_KEY);
  const username = window.localStorage.getItem(USER_KEY);
  if (!token || !username) {
    return null;
  }
  return { token, username };
}

export function writeSession({ token, username }) {
  if (!hasWindow()) {
    return;
  }
  window.localStorage.setItem(TOKEN_KEY, token);
  window.localStorage.setItem(USER_KEY, username);
}

export function clearSession() {
  if (!hasWindow()) {
    return;
  }
  window.localStorage.removeItem(TOKEN_KEY);
  window.localStorage.removeItem(USER_KEY);
}
