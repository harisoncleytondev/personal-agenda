import React, { useMemo } from "react";
import { MONTHS, MONTHS_S, WDS_S, getDateKey } from "../utils/constants";
import type { TasksState } from "../types";
import StatCard from "../components/ui/StatCard";

interface CalendarPageProps {
  today: Date;
  viewYear: number;
  viewMonth: number;
  setViewMonth: React.Dispatch<React.SetStateAction<Date>>;
  tasks: TasksState;
  onOpenDay: (year: number, month: number, day: number) => void;
}

export default function CalendarPage({
  today,
  viewYear,
  viewMonth,
  setViewMonth,
  tasks,
  onOpenDay,
}: CalendarPageProps) {
  const handlePrevMonth = () =>
    setViewMonth((prev) => new Date(viewYear, prev.getMonth() - 1));
  const handleNextMonth = () =>
    setViewMonth((prev) => new Date(viewYear, prev.getMonth() + 1));

  const stats = useMemo(() => {
    let daysWithTasks = 0;
    let totalTasks = 0;
    const todayKey = getDateKey(
      today.getFullYear(),
      today.getMonth(),
      today.getDate(),
    );
    const todayTasksCount = (tasks[todayKey] || []).length;

    Object.keys(tasks).forEach((key) => {
      const [y, m] = key.split("-");
      if (parseInt(y) === viewYear && parseInt(m) === viewMonth + 1) {
        if (tasks[key].length > 0) {
          daysWithTasks++;
          totalTasks += tasks[key].length;
        }
      }
    });

    return { daysWithTasks, totalTasks, todayTasksCount };
  }, [tasks, viewYear, viewMonth, today]);

  const calendarDays = useMemo(() => {
    const firstDay = new Date(viewYear, viewMonth, 1).getDay();
    const daysInMonth = new Date(viewYear, viewMonth + 1, 0).getDate();
    const daysInPrevMonth = new Date(viewYear, viewMonth, 0).getDate();
    const days: {
      day: number;
      isOtherMonth: boolean;
      isToday?: boolean;
      hasTasks?: boolean;
    }[] = [];

    for (let i = 0; i < firstDay; i++) {
      days.push({
        day: daysInPrevMonth - firstDay + 1 + i,
        isOtherMonth: true,
      });
    }

    for (let d = 1; d <= daysInMonth; d++) {
      const isToday =
        d === today.getDate() &&
        viewMonth === today.getMonth() &&
        viewYear === today.getFullYear();
      const key = getDateKey(viewYear, viewMonth, d);
      const hasTasks = tasks[key] && tasks[key].length > 0;
      days.push({ day: d, isOtherMonth: false, isToday, hasTasks });
    }

    const rem =
      (firstDay + daysInMonth) % 7 === 0
        ? 0
        : 7 - ((firstDay + daysInMonth) % 7);
    for (let d = 1; d <= rem; d++) {
      days.push({ day: d, isOtherMonth: true });
    }

    return days;
  }, [viewYear, viewMonth, today, tasks]);

  const todayBadgeText = `${today.getDate()} ${MONTHS_S[today.getMonth()]}. ${today.getFullYear()}`;

  return (
    <div id="page-cal" className="page active">
      <div className="cal-topbar">
        <div className="greeting">
          Olá,
          <br />
          <em>Harison.</em>
        </div>
        <div className="today-badge">{todayBadgeText}</div>
      </div>

      <div className="month-strip">
        <div>
          <span className="month-name">{MONTHS[viewMonth]}</span>
          <span className="month-year">{viewYear}</span>
        </div>
        <div className="nav-arrows">
          <button
            className="arrow-btn"
            onClick={handlePrevMonth}
            aria-label="Anterior"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              strokeWidth="2.2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <polyline points="15 18 9 12 15 6" />
            </svg>
          </button>
          <button
            className="arrow-btn"
            onClick={handleNextMonth}
            aria-label="Próximo"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              strokeWidth="2.2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <polyline points="9 18 15 12 9 6" />
            </svg>
          </button>
        </div>
      </div>

      <div className="cal-grid">
        <div className="wd-row">
          {WDS_S.map((wd) => (
            <div key={wd} className="wd">
              {wd}
            </div>
          ))}
        </div>

        <div className="days-grid">
          {calendarDays.map((d, i) => {
            if (d.isOtherMonth) {
              return (
                <div key={i} className="dc oth">
                  <span className="dnum">{d.day}</span>
                </div>
              );
            }

            return (
              <div
                key={i}
                className={`dc ${d.isToday ? "today" : ""} ${d.hasTasks ? "has-tasks" : ""}`}
                onClick={() => onOpenDay(viewYear, viewMonth, d.day)}
              >
                <span className="dnum">{d.day}</span>
                {d.hasTasks && <span className="edot"></span>}
              </div>
            );
          })}
        </div>
      </div>

      <div className="section-divider">
        <div className="line"></div>
        <span>Resumo do mês</span>
        <div className="line"></div>
      </div>

      <div className="stats-row">
        <StatCard
          number={stats.daysWithTasks}
          label={
            <>
              Dias com
              <br />
              afazeres
            </>
          }
        />
        <StatCard
          number={stats.totalTasks}
          label={
            <>
              Total de
              <br />
              afazeres
            </>
          }
        />
        <StatCard
          number={stats.todayTasksCount}
          label={
            <>
              Hoje
              <br />
              agendado
            </>
          }
        />
      </div>
    </div>
  );
}
