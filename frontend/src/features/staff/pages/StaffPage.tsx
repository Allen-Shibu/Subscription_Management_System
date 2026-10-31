import { useState, type FormEvent } from "react";
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
const emptyForm: StaffForm = { name: "", email: "", isAdmin: false };

export function StaffPage() {
  const [staff, setStaff] = useState<StaffMember[]>(initialStaff);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<StaffForm>(emptyForm);
  const adminCount = staff.filter(member => member.role === "admin").length;

  const openCreate = () => {
    setEditingId(null);
    setForm(emptyForm);
    setIsModalOpen(true);
  };
  const openEdit = (member: StaffMember) => {
    setEditingId(member.id);
    setForm({ name: member.name, email: member.email, isAdmin: member.role === "admin" });
    setIsModalOpen(true);
  };
  const save = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const member = { name: form.name.trim(), email: form.email.trim(), role: form.isAdmin ? "admin" as const : "staff" as const };
    if (editingId) setStaff(current => current.map(item => item.id === editingId ? { ...item, ...member } : item));
    else setStaff(current => [{ id: crypto.randomUUID(), ...member }, ...current]);
    setIsModalOpen(false);
  };

  return <div className={styles.staff}>
    <StaffHeader onAddStaff={openCreate} />
    <p className={styles.subtitle}>Sample data only. Changes stay in this tab and do not create Clerk accounts or grant roles.</p>
    <StaffOverview total={staff.length} admins={adminCount} staff={staff.length - adminCount} />
    <StaffTable staff={staff} onEdit={openEdit} onDelete={id => setStaff(current => current.filter(member => member.id !== id))} />
    {isModalOpen && <StaffModal mode={editingId ? "edit" : "create"} form={form} setForm={setForm}
      onClose={() => setIsModalOpen(false)} onSubmit={save} />}
  </div>;
}
