import { useState } from "react";
import { Link } from "react-router-dom";
import { Bell, ChevronDown, Menu, Search } from "lucide-react";
import styles from "./Topbar.module.css";

type TopbarProps = { onMenuClick: () => void };

export function Topbar({ onMenuClick }: TopbarProps) {
  const [profileOpen, setProfileOpen] = useState(false);
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
        <div className={styles.profileWrap}>
          <button
            className={styles.profileButton}
            onClick={() => setProfileOpen((open) => !open)}
            aria-expanded={profileOpen}
            aria-label="Open profile menu"
          >
            <span className={styles.avatar}>JD</span>
            <span className={styles.profileName}>Jordan Davis</span>
            <ChevronDown size={14} />
          </button>
          {profileOpen && (
            <div className={styles.profileMenu}>
              <span>
                Signed in as <strong>Jordan Davis</strong>
              </span>
              <Link to="/auth/login" onClick={() => setProfileOpen(false)}>
                Log out
              </Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
