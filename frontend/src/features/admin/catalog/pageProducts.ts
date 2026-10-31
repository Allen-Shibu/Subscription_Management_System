import type { ListParams, Page } from "./types";

export function pageProducts<T extends { name: string }>(items: T[], { q, page, pageSize }: ListParams): Page<T> {
  const filtered = items.filter(item => item.name.toLocaleLowerCase().includes(q.toLocaleLowerCase()));
  return { items: filtered.slice((page - 1) * pageSize, page * pageSize), total: filtered.length, page, pageSize };
}
