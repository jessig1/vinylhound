const API_BASE = "/api";

export class ApiError extends Error {
  constructor(message, status) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

async function request(path, options = {}) {
  const headers = options.headers || {};
  const rest = { ...options, headers };

  if (options.body !== undefined) {
    rest.body = options.body;
  }

  const response = await fetch(API_BASE + path, rest);
  const text = await response.text();
  let data = null;

  if (text) {
    try {
      data = JSON.parse(text);
    } catch (err) {
      data = text;
    }
  }

  if (!response.ok) {
    const message =
      data && typeof data === "object" && "error" in data
        ? data.error
        : typeof data === "string" && data
        ? data
        : response.statusText;
    throw new ApiError(message, response.status);
  }

  return data;
}

function jsonHeaders(options = {}) {
  const headers = new Headers();
  if (options.includeContentType) {
    headers.set("Content-Type", "application/json");
  }
  if (options.token) {
    headers.set("Authorization", "Bearer " + options.token);
  }
  return headers;
}

export async function signup(payload) {
  const headers = jsonHeaders({ includeContentType: true });
  await request("/signup", {
    method: "POST",
    headers,
    body: JSON.stringify(payload),
  });
}

export async function login(credentials) {
  const headers = jsonHeaders({ includeContentType: true });
  const data = await request("/login", {
    method: "POST",
    headers,
    body: JSON.stringify(credentials),
  });
  return data;
}

export async function fetchContent(token) {
  const headers = jsonHeaders({ token });
  return await request("/me/content", {
    method: "GET",
    headers,
  });
}

export async function updateContent(token, content) {
  const headers = jsonHeaders({ token, includeContentType: true });
  await request("/me/content", {
    method: "PUT",
    headers,
    body: JSON.stringify({ content }),
  });
}
