import { useState } from "react";
import { useAuth } from "@clerk/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { Eye, Pencil, Plus, Search, Trash2 } from "lucide-react";
import { ConfirmDelete } from "./ListControls";
import { useListControls } from "./useListControls";
import { pageProducts } from "./pageProducts";
import { serverProducts, type ServerProduct } from "./serverProducts";
import styles from "./Catalog.module.css";

export function ServerProductsPage({ admin = false }: { admin?: boolean }) {
  const { getToken } = useAuth();
  const queryClient = useQueryClient();
  const { params, search, setSearch, setPage, setPageSize } = useListControls();
  const [deleting, setDeleting] = useState<ServerProduct | null>(null);
  const [error, setError] = useState("");
  const products = useQuery({ queryKey: ["server-products"], queryFn: ({ signal }) => serverProducts.list(getToken, signal), retry: false });
  const remove = useMutation({
    mutationFn: (id: number) => serverProducts.delete(id, getToken),
    onSuccess: async () => { setDeleting(null); if (page.items.length === 1 && params.page > 1) setPage(params.page - 1); await queryClient.invalidateQueries({ queryKey: ["server-products"] }); },
    onError: (failure: Error) => { setDeleting(null); setError(failure.message); },
  });
  const page = pageProducts(admin ? products.data ?? [] : (products.data ?? []).filter(item => item.active), params);
  const pages = Math.max(1, Math.ceil(page.total / params.pageSize));

  return <div className={styles.page}>
    <header className={styles.header}><div><p className={styles.eyebrow}>Live products</p><h1>Products</h1>
      <p>{admin ? "Manage products and billing tiers." : "Browse available products and billing tiers."}</p></div>
      {admin && <Link className={styles.primary} to="/admin/products/new"><Plus size={16} /> Add product</Link>}</header>
    {error && <div className={styles.error} role="alert">{error}</div>}
    <div className={styles.toolbar}>
      <label className={styles.search}><Search size={16} /><input aria-label="Search products" placeholder="Search products..." value={search} onChange={event => setSearch(event.target.value)} /></label>
      <label className={styles.pageSize}>Rows per page <select value={params.pageSize} onChange={event => setPageSize(Number(event.target.value))}>{[5, 10, 20].map(size => <option key={size} value={size}>{size}</option>)}</select></label>
    </div>
    <div className={styles.tableWrap}>
      {products.isPending ? <div className={styles.skeleton} aria-label="Loading" role="status">{Array.from({ length: 5 }, (_, index) => <span key={index} />)}</div> :
        products.isError ? <div className={styles.empty} role="alert"><h2>Could not load products</h2><p>{products.error.message}</p><button onClick={() => products.refetch()}>Retry</button></div> :
        !page.items.length ? <div className={styles.empty}><h2>{params.q ? "No products match your search" : "No products yet"}</h2></div> :
          <table className={styles.table}><thead><tr><th>Name</th><th>Description</th><th>Status</th><th>{admin ? "Actions" : "Details"}</th></tr></thead>
            <tbody>{page.items.map(product => <tr key={product.id}><td className={styles.name}>{product.name}</td><td>{product.description || "—"}</td>
              <td><span className={styles.badge}>{product.active ? "Active" : "Inactive"}</span></td>
              <td><div className={styles.actions}>{!admin && <Link aria-label={`View ${product.name}`} to={`/products/${product.id}`}><Eye size={16} /></Link>}
                {admin && <><Link aria-label={`Edit ${product.name}`} to={`/admin/products/${product.id}/edit`}><Pencil size={16} /></Link>
                  <button aria-label={`Delete ${product.name}`} onClick={() => setDeleting(product)}><Trash2 size={16} /></button></>}</div></td></tr>)}</tbody></table>}
    </div>
    <div className={styles.pagination}><span>{page.total} total · Search and pagination run in this browser.</span><div>
      <button disabled={params.page <= 1} onClick={() => setPage(params.page - 1)}>Previous</button><span>Page {params.page} of {pages}</span>
      <button disabled={params.page >= pages} onClick={() => setPage(params.page + 1)}>Next</button>
    </div></div>
    {deleting && <ConfirmDelete name={deleting.name} pending={remove.isPending} onCancel={() => setDeleting(null)} onConfirm={() => remove.mutate(deleting.id)} />}
  </div>;
}
