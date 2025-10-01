export function normalizeContent(value) {
  if (!value) {
    return [];
  }

  return value
    .split(/\r?\n/)
    .map((entry) => entry.trim())
    .filter(Boolean);
}
