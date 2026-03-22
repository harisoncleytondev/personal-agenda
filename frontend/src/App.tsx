import { useState, useEffect, useCallback } from "react";
import CalendarPage from "./pages/CalendarPage";
import DayPage from "./pages/DayPage";
import { GET_BASE_URL, getDateKey } from "./utils/constants";
import type { ActiveDate, TasksState, Task } from "./types";
import AuthPage from "./pages/AuthPage";

export default function App() {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [today] = useState<Date>(new Date());
  const [viewDate, setViewDate] = useState<Date>(new Date());
  const [activeDate, setActiveDate] = useState<ActiveDate | null>(null);
  const [tasks, setTasks] = useState<TasksState>({});

  const fetchWithAuth = useCallback(
    async (url: string, options: RequestInit = {}) => {
      let res = await fetch(url, options);

      if (res.status === 401) {
        const refreshRes = await fetch(`${GET_BASE_URL}/auth/refresh`, {
          credentials: "include",
        });

        if (refreshRes.ok) {
          res = await fetch(url, options);
        } else {
          setIsAuthenticated(false);
        }
      }

      return res;
    },
    [],
  );

  const checkSession = useCallback(async () => {
    try {
      const res = await fetchWithAuth(`${GET_BASE_URL}/logged/appointment/`, {
        credentials: "include",
      });
      if (res.ok) {
        setIsAuthenticated(true);
      }
    } catch (err) {
    } finally {
      setIsLoading(false);
    }
  }, [fetchWithAuth]);

  useEffect(() => {
    checkSession();
  }, [checkSession]);

  const loadTasks = useCallback(async () => {
    try {
      const res = await fetchWithAuth(
        `${GET_BASE_URL}/logged/appointment/getall`,
        {
          credentials: "include",
        },
      );
      if (res.ok) {
        const data = await res.json();
        const newTasks: TasksState = {};
        if (Array.isArray(data)) {
          data.forEach((ap: any) => {
            const dateObj = new Date(ap.task_date);
            const key = getDateKey(
              dateObj.getUTCFullYear(),
              dateObj.getUTCMonth(),
              dateObj.getUTCDate(),
            );
            if (!newTasks[key]) newTasks[key] = [];
            const alertDates = ap.alert_dates
              ? ap.alert_dates.map((d: string) => d.split("T")[0])
              : undefined;
            newTasks[key].push({
              id: ap.id,
              name: ap.name,
              desc: ap.description || "",
              ts: ap.time_start ? ap.time_start.substring(0, 5) : "",
              te: ap.time_end ? ap.time_end.substring(0, 5) : "",
              alertType: ap.alert_type,
              alertDates: alertDates,
            });
          });
        }
        setTasks(newTasks);
      }
    } catch (err) {}
  }, [fetchWithAuth]);

  useEffect(() => {
    if (isAuthenticated) loadTasks();
  }, [isAuthenticated, loadTasks]);

  const handleOpenDay = (year: number, month: number, day: number) => {
    setActiveDate({ year, month, day });
  };

  const handleAddTask = async (dateObj: ActiveDate, taskObj: Task) => {
    const taskDateStr = `${dateObj.year}-${String(dateObj.month + 1).padStart(2, "0")}-${String(dateObj.day).padStart(2, "0")}`;
    const payload = {
      task_date: taskDateStr,
      name: taskObj.name,
      description: taskObj.desc || "",
      time_start: taskObj.ts || "",
      time_end: taskObj.te || "",
      alert_type: taskObj.alertType || "none",
      alert_dates: taskObj.alertDates || [],
    };
    try {
      const res = await fetchWithAuth(
        `${GET_BASE_URL}/logged/appointment/create`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
          credentials: "include",
        },
      );
      if (res.ok) loadTasks();
    } catch (err) {}
  };

  const handleUpdateTask = async (dateObj: ActiveDate, taskObj: Task) => {
    const taskDateStr = `${dateObj.year}-${String(dateObj.month + 1).padStart(2, "0")}-${String(dateObj.day).padStart(2, "0")}`;
    const payload = {
      task_date: taskDateStr,
      name: taskObj.name,
      description: taskObj.desc || "",
      time_start: taskObj.ts || "",
      time_end: taskObj.te || "",
      alert_type: taskObj.alertType || "none",
      alert_dates: taskObj.alertDates || [],
    };
    try {
      const res = await fetchWithAuth(
        `${GET_BASE_URL}/logged/appointment/update/${taskObj.id}`,
        {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
          credentials: "include",
        },
      );
      if (res.ok) loadTasks();
    } catch (err) {}
  };

  const handleDeleteTask = async (dateKey: string, taskId: string) => {
    setTasks((prev) => ({
      ...prev,
      [dateKey]: prev[dateKey].filter((t) => t.id !== taskId),
    }));
    try {
      await fetchWithAuth(
        `${GET_BASE_URL}/logged/appointment/delete/${taskId}`,
        {
          method: "DELETE",
          credentials: "include",
        },
      );
    } catch (err) {}
  };

  if (isLoading) {
    return null;
  }

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
