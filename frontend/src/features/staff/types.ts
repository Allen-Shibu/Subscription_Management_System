export type StaffRole = "admin" | "staff";

export type StaffMember = {
  id: string;
  name: string;
  email: string;
  role: StaffRole;
};

export type StaffForm = {
  name: string;
  email: string;
  password: string;
  isAdmin: boolean;
};
