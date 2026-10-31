import { useState } from "react";
import { Outlet } from "react-router-dom";
import { Sidebar } from "../sidebar/Sidebar";
import { Topbar } from "../topbar/Topbar";
import styles from "./MainLayout.module.css";

export function MainLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  return (
    <div className={styles.appLayout}>
      <Sidebar open={sidebarOpen} onClose={() => setSidebarOpen(false)} />
      <div className={styles.mainColumn}>
        <Topbar onMenuClick={() => setSidebarOpen(true)} />
        <main className={styles.content}>
          <Outlet />
        </main>
      </div>
    </div>
  );
}
