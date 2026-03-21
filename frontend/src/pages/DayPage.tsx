import { useState } from "react";
import { MONTHS, WDS, getDateKey } from "../utils/constants";
import type { ActiveDate, TasksState, Task } from "../types";
import TaskForm from "../components/TaskForm";
import TaskCard from "../components/TaskCard";

interface DayPageProps {
  activeDate: ActiveDate;
  tasks: TasksState;
  onBack: () => void;
  onAddTask: (date: ActiveDate, task: Task) => void;
  onUpdateTask: (date: ActiveDate, task: Task) => void;
  onDeleteTask: (dateKey: string, taskId: string) => void;
}

export default function DayPage({
  activeDate,
  tasks,
  onBack,
  onAddTask,
  onUpdateTask,
  onDeleteTask,
}: DayPageProps) {
  const [editingTaskId, setEditingTaskId] = useState<string | null>(null);

  const { year, month, day } = activeDate;
  const dateKey = getDateKey(year, month, day);
  const dayTasks = tasks[dateKey] || [];

  const dt = new Date(year, month, day);
  const wd = WDS[dt.getDay()];

  const editingTask = dayTasks.find((t) => t.id === editingTaskId);

  const handleEditClick = (taskId: string) => {
    setEditingTaskId(taskId);
    document
      .querySelector(".scroll-content")
      ?.scrollTo({ top: 0, behavior: "smooth" });
  };

  const handleUpdate = (date: ActiveDate, task: Task) => {
    onUpdateTask(date, task);
    setEditingTaskId(null);
  };

  return (
    <div id="page-day" className="page active">
      <div className="day-hero">
        <button className="back-btn" onClick={onBack}>
          <svg
            viewBox="0 0 24 24"
            fill="none"
            strokeWidth="2.2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <polyline points="15 18 9 12 15 6" />
          </svg>
          Calendário
        </button>
        <div className="day-hero-date">
          {day} de <em>{MONTHS[month].toLowerCase()}</em>
        </div>
        <div className="day-hero-wd">
          {wd} · {year}
        </div>
      </div>

      <div className="scroll-content">
        <TaskForm
          activeDate={activeDate}
          editingTask={editingTask}
          onAddTask={onAddTask}
          onUpdateTask={handleUpdate}
          onCancelEdit={() => setEditingTaskId(null)}
        />

        <div className="tasks-zone">
          <div className="zone-header">
            <div className="zone-title">Afazeres do dia</div>
            <div className="task-count">{dayTasks.length}</div>
          </div>

          <div id="tasks-list">
            {dayTasks.length === 0 ? (
              <div className="no-tasks">
                <span>∅</span>Nenhum afazer ainda.
              </div>
            ) : (
              dayTasks.map((task) => (
                <TaskCard
                  key={task.id}
                  task={task}
                  onEdit={handleEditClick}
                  onDelete={(taskId) => onDeleteTask(dateKey, taskId)}
                />
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
