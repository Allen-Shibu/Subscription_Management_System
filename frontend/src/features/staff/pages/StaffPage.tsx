import { useState } from "react";
import type { FormEvent } from "react";
import { StaffHeader } from "../components/StaffHeader";
import { StaffModal } from "../components/StaffModal";
import { StaffOverview } from "../components/StaffOverview";
import { StaffTable } from "../components/StaffTable";
import type { StaffForm, StaffMember } from "../types";
import styles from "../Staff.module.css";

const initialStaff: StaffMember[] = [
  { id: "1", name: "Jordan Davis", email: "jordan@recur.app", role: "admin" },
  { id: "2", name: "Taylor Morgan", email: "taylor@recur.app", role: "staff" },
];

const emptyForm: StaffForm = {
  name: "",
  email: "",
  password: "",
  isAdmin: false,
};

export function StaffPage() {
  const [staff, setStaff] = useState<StaffMember[]>(initialStaff);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [mode, setMode] = useState<"create" | "edit">("create");
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState({
    name: "",
    email: "",
    isAdmin: false,
  });
  const adminCount = staff.filter((member) => member.role === "admin").length;

  const openCreate = () => {
    setMode("create");
    setEditingId(null);
    setForm({ name: "", email: "", isAdmin: false });
    setIsModalOpen(true);
  };

  const openEdit = (member: StaffMember) => {
    setMode("edit");
    setEditingId(member.id);
    setForm({
      name: member.name,
      email: member.email,
      isAdmin: member.role === "admin",
    });
    setIsModalOpen(true);
  };

  const closeModal = () => setIsModalOpen(false);

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (mode === "create") {
      const newMember: StaffMember = {
        id: Date.now().toString(),
        name: form.name.trim(),
        email: form.email.trim(),
        role: form.isAdmin ? "admin" : "staff",
      };
      setStaff((current) => [newMember, ...current]);
    } else if (editingId) {
      setStaff((current) =>
        current.map((member) =>
          member.id === editingId
            ? {
                ...member,
                name: form.name.trim(),
                email: form.email.trim(),
                role: form.isAdmin ? "admin" : "staff",
              }
            : member,
        ),
      );
    }

    closeModal();
  };

  const handleDelete = (id: string) => {
    setStaff((current) => current.filter((member) => member.id !== id));
  };

  return (
    <div className={styles.staff}>
      <header className={styles.header}>
        <div>
          <p className={styles.eyebrow}>Staff Management</p>
          <h1>Team Members</h1>
          <p className={styles.subtitle}>
            Preview staff members and roles. Changes stay in this tab.
          </p>
        </div>
        <div className={styles.headerActions}>
          <button className={styles.primaryButton} onClick={openCreate}>
            <Plus size={15} /> Add Staff
          </button>
        </div>
      </header>
      <p className={styles.subtitle}>Sample data only. Adding or editing a member here does not create a Clerk account or grant a role.</p>
      <section className={styles.summary} aria-label="Staff overview">
        <article className={`${styles.summaryCard} ${styles.summaryTotal}`}>
          <span className={styles.summaryIcon}>
            <Users size={18} />
          </span>
          <div>
            <span>Team members</span>
            <strong>{staff.length}</strong>
          </div>
        </article>
        <article className={`${styles.summaryCard} ${styles.summaryAdmins}`}>
          <span className={styles.summaryIcon}>
            <ShieldCheck size={18} />
          </span>
          <div>
            <span>Admins</span>
            <strong>{adminCount}</strong>
          </div>
        </article>
        <article className={`${styles.summaryCard} ${styles.summaryStaff}`}>
          <span className={styles.summaryIcon}>
            <UserRound size={18} />
          </span>
          <div>
            <span>Staff</span>
            <strong>{teamCount}</strong>
          </div>
        </article>
      </section>
      <section className={styles.content}>
        <div className={styles.panel}>
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <colgroup>
                <col className={styles.nameColumn} />
                <col className={styles.emailColumn} />
                <col className={styles.roleColumn} />
                <col className={styles.actionsColumn} />
              </colgroup>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th className={styles.actionsHeader}>Actions</th>
                </tr>
              </thead>
              <tbody>
                {staff.map((member) => (
                  <tr key={member.id}>
                    <td data-label="Name">
                      <div className={styles.nameCell}>
                        <span className={styles.avatar}>
                          {member.name
                            .split(" ")
                            .map((n) => n[0])
                            .join("")
                            .slice(0, 2)}
                        </span>
                        <span className={styles.nameText}>{member.name}</span>
                      </div>
                    </td>
                    <td data-label="Email">{member.email}</td>
                    <td data-label="Role">
                      <span
                        className={`${styles.roleBadge} ${
                          member.role === "admin"
                            ? styles.admin
                            : styles.regular
                        }`}
                      >
                        {member.role === "admin" ? "Admin" : "Staff"}
                      </span>
                    </td>
                    <td className={styles.actionsCell} data-label="Actions">
                      <button
                        className={styles.iconButton}
                        onClick={() => openEdit(member)}
                        aria-label="Edit staff"
                      >
                        <Pencil size={16} />
                      </button>
                      <button
                        className={`${styles.iconButton} ${styles.danger}`}
                        onClick={() => handleDelete(member.id)}
                        aria-label="Delete staff"
                      >
                        <Trash2 size={16} />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </section>

      {isModalOpen && (
        <div className={styles.modalBackdrop} onClick={closeModal}>
          <div
            className={styles.modal}
            role="dialog"
            aria-modal="true"
            onClick={(e) => e.stopPropagation()}
          >
            <header className={styles.modalHeader}>
              <h2>{mode === "create" ? "Add New Staff" : "Edit Staff"}</h2>
              <button
                className={styles.closeButton}
                onClick={closeModal}
                aria-label="Close"
              >
                ×
              </button>
            </header>
            <form className={styles.modalForm} onSubmit={handleSubmit}>
              <label className={styles.field}>
                <span>Staff Name</span>
                <input
                  type="text"
                  placeholder="John Doe"
                  value={form.name}
                  onChange={(e) =>
                    setForm((prev) => ({ ...prev, name: e.target.value }))
                  }
                  required
                />
              </label>
              <label className={styles.field}>
                <span>Email</span>
                <input
                  type="email"
                  placeholder="johndoe@odoo.com"
                  value={form.email}
                  onChange={(e) =>
                    setForm((prev) => ({ ...prev, email: e.target.value }))
                  }
                  required
                />
              </label>
              <label className={styles.checkboxField}>
                <input
                  type="checkbox"
                  checked={form.isAdmin}
                  onChange={(e) =>
                    setForm((prev) => ({ ...prev, isAdmin: e.target.checked }))
                  }
                />
                <span>Create as Admin</span>
              </label>
              <div className={styles.modalActions}>
                <button
                  type="button"
                  className={styles.secondaryButton}
                  onClick={closeModal}
                >
                  Cancel
                </button>
                <button type="submit" className={styles.primaryButton}>
                  {mode === "create" ? "Add demo staff" : "Save demo changes"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </main>
  );
}
