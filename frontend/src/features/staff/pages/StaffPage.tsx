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
  const [form, setForm] = useState<StaffForm>(emptyForm);

  const adminCount = staff.filter((member) => member.role === "admin").length;

  const openCreate = () => {
    setMode("create");
    setEditingId(null);
    setForm(emptyForm);
    setIsModalOpen(true);
  };

  const openEdit = (member: StaffMember) => {
    setMode("edit");
    setEditingId(member.id);
    setForm({
      name: member.name,
      email: member.email,
      password: "",
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
    <main className={styles.staff}>
      <StaffHeader onAddStaff={openCreate} />
      <StaffOverview
        total={staff.length}
        admins={adminCount}
        staff={staff.length - adminCount}
      />
      <StaffTable staff={staff} onEdit={openEdit} onDelete={handleDelete} />
      {isModalOpen && (
        <StaffModal
          mode={mode}
          form={form}
          setForm={setForm}
          onClose={closeModal}
          onSubmit={handleSubmit}
        />
      )}
    </main>
  );
}
