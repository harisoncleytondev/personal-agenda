import type { ReactNode } from "react";

interface StatCardProps {
  number: number;
  label: ReactNode;
}

export default function StatCard({ number, label }: StatCardProps) {
  return (
    <div className="stat-card">
      <div className="stat-num">{number}</div>
      <div className="stat-lbl">{label}</div>
    </div>
  );
}
