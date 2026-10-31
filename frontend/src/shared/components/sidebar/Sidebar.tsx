import {
  CreditCard,
  LayoutDashboard,
  RefreshCw,
  Boxes,
  CalendarClock,
  FileText,
  Percent,
  ReceiptText,
  UsersRound,
  X,
} from "lucide-react";
import { Show } from "@clerk/react";
import { NavLink } from "react-router-dom";
import styles from "./Sidebar.module.css";

const navigation = [
  { label: "Dashboard", to: "/dashboard", Icon: LayoutDashboard },
  { label: "Subscriptions", to: "/subscriptions", Icon: RefreshCw },
  { label: "Billing & Invoices", to: "/billing", Icon: CreditCard },
];
const adminNavigation = [
  { label: "Products", to: "/products", Icon: Boxes },
  { label: "Plans", to: "/plans", Icon: CalendarClock },
  { label: "Taxes", to: "/taxes", Icon: ReceiptText },
  { label: "Discounts", to: "/discounts", Icon: Percent },
  { label: "Staff", to: "/staff", Icon: UsersRound },
  { label: "Quotation templates", to: "/quotation-templates", Icon: FileText },
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
        <Show when={{ role: "org:admin" }}>
          <p className={styles.sectionLabel}>Admin</p>
          <nav className={styles.navigation} aria-label="Admin navigation">
            {adminNavigation.map(({ label, to, Icon }) => (
              <NavLink key={to} to={to} onClick={onClose}
                className={({ isActive }) => `${styles.navLink} ${isActive ? styles.active : ""}`}>
                <Icon size={17} strokeWidth={1.8} /><span>{label}</span>
              </NavLink>
            ))}
          </nav>
        </Show>
      </aside>
    </>
  );
}
