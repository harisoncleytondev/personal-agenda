import { useState } from "react";
import CalendarPage from "./pages/CalendarPage";
import DayPage from "./pages/DayPage";
import { getDateKey } from "./utils/constants";
import type { ActiveDate, TasksState, Task } from "./types";
import AuthPage from "./pages/AuthPage";

const getInitialMocks = (): TasksState => {
  const t = new Date();
  const todayKey = getDateKey(t.getFullYear(), t.getMonth(), t.getDate());

  const amanhã = new Date(t);
  amanhã.setDate(amanhã.getDate() + 1);
  const amanhaStr = `${amanhã.getFullYear()}-${String(amanhã.getMonth() + 1).padStart(2, "0")}-${String(amanhã.getDate()).padStart(2, "0")}`;

  const daquiA3Dias = new Date(t);
  daquiA3Dias.setDate(daquiA3Dias.getDate() + 3);
  const tresDiasStr = `${daquiA3Dias.getFullYear()}-${String(daquiA3Dias.getMonth() + 1).padStart(2, "0")}-${String(daquiA3Dias.getDate()).padStart(2, "0")}`;

  return {
    [todayKey]: [
      {
        id: "mock-1",
        name: "Reunião de Alinhamento",
        ts: "10:00",
        te: "11:00",
        alertType: "1_day_before",
      },
      {
        id: "mock-2",
        name: "Estudar React",
        desc: "Fazer o módulo de hooks e componentes",
        alertType: "custom_date",
        alertDates: [amanhaStr, tresDiasStr],
      },
    ],
  };
};

export default function App() {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [today] = useState<Date>(new Date());
  const [viewDate, setViewDate] = useState<Date>(new Date());
  const [activeDate, setActiveDate] = useState<ActiveDate | null>(null);
  const [tasks, setTasks] = useState<TasksState>(getInitialMocks());

  const handleOpenDay = (year: number, month: number, day: number) => {
    setActiveDate({ year, month, day });
  };

  const handleAddTask = (dateObj: ActiveDate, taskObj: Task) => {
    const key = getDateKey(dateObj.year, dateObj.month, dateObj.day);
    setTasks((prev) => ({
      ...prev,
      [key]: [...(prev[key] || []), taskObj],
    }));
  };

  const handleUpdateTask = (dateObj: ActiveDate, updatedTask: Task) => {
    const key = getDateKey(dateObj.year, dateObj.month, dateObj.day);
    setTasks((prev) => ({
      ...prev,
      [key]: prev[key].map((t) => (t.id === updatedTask.id ? updatedTask : t)),
    }));
  };

  const handleDeleteTask = (dateKey: string, taskId: string) => {
    setTasks((prev) => ({
      ...prev,
      [dateKey]: prev[dateKey].filter((t) => t.id !== taskId),
    }));
  };

  return (
    <div className="app-shell">
      {!isAuthenticated ? (
        <AuthPage onLogin={() => setIsAuthenticated(true)} />
      ) : !activeDate ? (
        <CalendarPage
          today={today}
          viewYear={viewDate.getFullYear()}
          viewMonth={viewDate.getMonth()}
          setViewMonth={setViewDate}
          tasks={tasks}
          onOpenDay={handleOpenDay}
        />
      ) : (
        <DayPage
          activeDate={activeDate}
          tasks={tasks}
          onBack={() => setActiveDate(null)}
          onAddTask={handleAddTask}
          onUpdateTask={handleUpdateTask}
          onDeleteTask={handleDeleteTask}
        />
      )}
    </div>
  );
}
