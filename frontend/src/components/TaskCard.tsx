import type { Task } from "../types";

interface TaskCardProps {
  task: Task;
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
}

export default function TaskCard({ task, onEdit, onDelete }: TaskCardProps) {
  const renderAlertText = () => {
    if (task.alertType === "1_day_before") return "1 dia antes";
    if (
      task.alertType === "custom_date" &&
      task.alertDates &&
      task.alertDates.length > 0
    ) {
      return task.alertDates
        .map((dateStr) => {
          const [year, month, day] = dateStr.split("-");
          return `${day}/${month}/${year}`;
        })
        .join(", ");
    }
    return null;
  };

  const alertText = renderAlertText();

  return (
    <div className="task-card">
      <div className="tc-actions">
        <button
          className="tc-btn edit"
          onClick={() => onEdit(task.id)}
          aria-label="Editar"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            strokeWidth="2.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M12 20h9"></path>
            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path>
          </svg>
        </button>
        <button
          className="tc-btn del"
          onClick={() => onDelete(task.id)}
          aria-label="Remover"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            strokeWidth="2.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </div>

      <div className="tc-name">{task.name}</div>

      <div className="tc-info-row">
        {task.ts && (
          <div className="tc-badge">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              strokeWidth="2.2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <circle cx="12" cy="12" r="10" />
              <polyline points="12 6 12 12 16 14" />
            </svg>
            {task.ts}
            {task.te ? ` – ${task.te}` : ""}
          </div>
        )}

        {alertText && (
          <div className="tc-badge alert">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              strokeWidth="2.2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
              <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
            </svg>
            {alertText}
          </div>
        )}
      </div>

      {task.desc && <div className="tc-desc">{task.desc}</div>}
    </div>
  );
}
