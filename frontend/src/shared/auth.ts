import { useAuth } from "@clerk/react";

export function useIsAdmin() {
  const { orgRole, sessionClaims } = useAuth();
  const customRole = (sessionClaims as { role?: string } | null)?.role;
  return orgRole === "org:admin" || customRole === "admin" || customRole === "org:admin";
}
