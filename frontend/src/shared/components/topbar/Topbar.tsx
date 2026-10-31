import { UserButton } from "@clerk/react";
import { Bell, ChevronDown, Menu, Search } from "lucide-react";
import styles from "./Topbar.module.css";

type TopbarProps = { onMenuClick: () => void };

export function Topbar({ onMenuClick }: TopbarProps) {
  return (
    <header className={styles.topbar}>
      <button
        className={styles.menuButton}
        onClick={onMenuClick}
        aria-label="Open navigation"
      >
        <Menu size={20} />
      </button>
      <label className={styles.search}>
        <Search size={16} />
        <input
          aria-label="Search"
          placeholder="Search subscriptions, invoices, customers..."
        />
        <kbd>⌘ K</kbd>
      </label>
      <div className={styles.actions}>
        <button className={styles.rangeButton}>
          Last 30 Days <ChevronDown size={14} />
        </button>
        <button className={styles.iconButton} aria-label="Notifications">
          <Bell size={17} />
        </button>
        <UserButton />
      </div>
    </header>
  );
}
