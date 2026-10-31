import { Pencil, Trash2 } from "lucide-react";
import type { StaffMember } from "../types";
import styles from "../Staff.module.css";

type StaffTableProps = {
  staff: StaffMember[];
  onEdit: (member: StaffMember) => void;
  onDelete: (id: string) => void;
};

export function StaffTable({ staff, onEdit, onDelete }: StaffTableProps) {
  return (
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
                          .map((namePart) => namePart[0])
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
                        member.role === "admin" ? styles.admin : styles.regular
                      }`}
                    >
                      {member.role === "admin" ? "Admin" : "Staff"}
                    </span>
                  </td>
                  <td className={styles.actionsCell} data-label="Actions">
                    <button
                      className={styles.iconButton}
                      onClick={() => onEdit(member)}
                      aria-label={`Edit ${member.name}`}
                    >
                      <Pencil size={16} />
                    </button>
                    <button
                      className={`${styles.iconButton} ${styles.danger}`}
                      onClick={() => onDelete(member.id)}
                      aria-label={`Delete ${member.name}`}
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
  );
}
