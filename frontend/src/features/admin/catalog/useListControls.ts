import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import type { ListParams } from "./types";

export function useListControls(): {
  params: ListParams; search: string; setSearch: (value: string) => void;
  setPage: (value: number) => void; setPageSize: (value: number) => void;
} {
  const [url, setUrl] = useSearchParams();
  const q = url.get("q") ?? "";
  const [draft, setDraft] = useState<{ base: string; value: string } | null>(null);
  const search = draft?.base === q ? draft.value : q;
  const page = Math.max(1, Number(url.get("page")) || 1);
  const pageSize = [5, 10, 20].includes(Number(url.get("pageSize"))) ? Number(url.get("pageSize")) : 10;
  useEffect(() => {
    if (search === q) return;
    const timer = setTimeout(() => setUrl(current => {
      const next = new URLSearchParams(current);
      if (search) next.set("q", search); else next.delete("q");
      next.set("page", "1");
      return next;
    }, { replace: true }), 300);
    return () => clearTimeout(timer);
  }, [search, q, setUrl]);
  const set = (key: string, value: number) => setUrl(current => {
    const next = new URLSearchParams(current);
    next.set(key, String(value));
    if (key === "pageSize") next.set("page", "1");
    return next;
  });
  return { params: { q, page, pageSize }, search, setSearch: value => setDraft({ base: q, value }),
    setPage: value => set("page", value), setPageSize: value => set("pageSize", value) };
}
