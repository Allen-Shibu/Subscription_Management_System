import {
  BarChart3,
  CreditCard,
  LayoutDashboard,
  RefreshCw,
  Settings,
  UsersRound,
  X,
} from "lucide-react";
import { NavLink } from "react-router-dom";
import styles from "./Sidebar.module.css";

const navigation = [
  { label: "Dashboard", to: "/dashboard", Icon: LayoutDashboard },
  { label: "Subscriptions", to: "/subscriptions", Icon: RefreshCw },
  { label: "Customers", to: "/customers", Icon: UsersRound },
  { label: "Billing & Invoices", to: "/billing", Icon: CreditCard },
  { label: "Analytics", to: "/analytics", Icon: BarChart3 },
  { label: "Settings", to: "/settings", Icon: Settings },
];

type SidebarProps = { open: boolean; onClose: () => void };

export function Sidebar({ open, onClose }: SidebarProps) {
  return (
    <>
      <button
        aria-label="Close navigation"
        className={`${styles.backdrop} ${open ? styles.backdropVisible : ""}`}
        onClick={onClose}
        tabIndex={open ? 0 : -1}
      />
      <aside className={`${styles.sidebar} ${open ? styles.sidebarOpen : ""}`}>
        <div className={styles.brandRow}>
          <NavLink to="/dashboard" className={styles.brand} onClick={onClose}>
            <span className={styles.brandMark}>R</span>
            <span>Recur</span>
          </NavLink>
          <button
            className={styles.closeButton}
            onClick={onClose}
            aria-label="Close navigation"
          >
            <X size={19} />
          </button>
        </div>
        <p className={styles.sectionLabel}>Workspace</p>
        <nav className={styles.navigation} aria-label="Main navigation">
          {navigation.map(({ label, to, Icon }) => (
            <NavLink
              key={to}
              to={to}
              onClick={onClose}
              className={({ isActive }) =>
                `${styles.navLink} ${isActive ? styles.active : ""}`
              }
            >
              <Icon size={17} strokeWidth={1.8} />
              <span>{label}</span>
            </NavLink>
          ))}
        </nav>
      </aside>
    </>
  );
}
