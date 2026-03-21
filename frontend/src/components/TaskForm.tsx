import { useState, useEffect } from "react";
import type { ActiveDate, Task } from "../types";

interface TaskFormProps {
  activeDate: ActiveDate;
  editingTask?: Task;
  onAddTask: (date: ActiveDate, task: Task) => void;
  onUpdateTask: (date: ActiveDate, task: Task) => void;
  onCancelEdit: () => void;
}

export default function TaskForm({
  activeDate,
  editingTask,
  onAddTask,
  onUpdateTask,
  onCancelEdit,
}: TaskFormProps) {
  const [name, setName] = useState<string>("");
  const [desc, setDesc] = useState<string>("");
  const [hasTime, setHasTime] = useState<boolean>(false);
  const [timeStart, setTimeStart] = useState<string>("");
  const [timeEnd, setTimeEnd] = useState<string>("");
  const [alertType, setAlertType] = useState<
    "none" | "1_day_before" | "custom_date"
  >("none");
  const [alertDates, setAlertDates] = useState<string[]>([]);
  const [tempDate, setTempDate] = useState<string>("");
  const [isAlertMenuOpen, setIsAlertMenuOpen] = useState<boolean>(false);
  const [error, setError] = useState<boolean>(false);

  const alertOptions = [
    { value: "none", label: "Sem alerta" },
    { value: "1_day_before", label: "1 dia antes" },
    { value: "custom_date", label: "Data(s) específica(s)" },
  ];

  useEffect(() => {
    if (editingTask) {
      setName(editingTask.name);
      setDesc(editingTask.desc || "");
      setHasTime(!!editingTask.ts);
      setTimeStart(editingTask.ts || "");
      setTimeEnd(editingTask.te || "");
      setAlertType(editingTask.alertType || "none");
      setAlertDates(editingTask.alertDates || []);
    } else {
      resetForm();
    }
  }, [editingTask]);

  const resetForm = () => {
    setName("");
    setDesc("");
    setHasTime(false);
    setTimeStart("");
    setTimeEnd("");
    setAlertType("none");
    setAlertDates([]);
    setTempDate("");
    setIsAlertMenuOpen(false);
  };

  const handleAddDate = () => {
    if (tempDate && !alertDates.includes(tempDate)) {
      setAlertDates([...alertDates, tempDate].sort());
      setTempDate("");
    }
  };

  const handleRemoveDate = (dateToRemove: string) => {
    setAlertDates(alertDates.filter((d) => d !== dateToRemove));
  };

  const handleSubmit = () => {
    if (!name.trim()) {
      setError(true);
      setTimeout(() => setError(false), 400);
      return;
    }

    const taskData: Task = {
      id: editingTask ? editingTask.id : Date.now().toString(),
      name: name.trim(),
      desc: desc.trim(),
      ts: hasTime ? timeStart : "",
      te: hasTime ? timeEnd : "",
      alertType,
      alertDates: alertType === "custom_date" ? alertDates : undefined,
    };

    if (editingTask) {
      onUpdateTask(activeDate, taskData);
    } else {
      onAddTask(activeDate, taskData);
    }

    resetForm();
  };

  const handleSelectOption = (value: any) => {
    setAlertType(value);
    setIsAlertMenuOpen(false);
  };

  return (
    <div className="form-zone">
      <div className="bs-field">
        <label className="bs-label">
          Nome <span className="req">*</span>
        </label>
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          className={`bs-input ${error ? "err" : ""}`}
          type="text"
          placeholder="Ex: Reunião de equipe"
          maxLength={80}
          autoComplete="off"
        />
      </div>

      <div className="bs-field">
        <label className="bs-label">Descrição</label>
        <textarea
          value={desc}
          onChange={(e) => setDesc(e.target.value)}
          className="bs-input"
          placeholder="Detalhes opcionais..."
        />
      </div>

      <div className="bs-field">
        <label className="bs-label">Alerta de Lembrete</label>

        <div className="custom-select-wrapper">
          <div
            className={`bs-input custom-select-trigger ${isAlertMenuOpen ? "open" : ""}`}
            onClick={() => setIsAlertMenuOpen(!isAlertMenuOpen)}
          >
            <span>
              {alertOptions.find((opt) => opt.value === alertType)?.label}
            </span>
            <svg
              viewBox="0 0 24 24"
              fill="none"
              strokeWidth="2.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
          </div>

          {isAlertMenuOpen && (
            <div className="custom-select-dropdown">
              {alertOptions.map((opt) => (
                <div
                  key={opt.value}
                  className={`custom-select-option ${alertType === opt.value ? "selected" : ""}`}
                  onClick={() => handleSelectOption(opt.value)}
                >
                  {opt.label}
                </div>
              ))}
            </div>
          )}
        </div>

        {alertType === "custom_date" && (
          <div style={{ marginTop: "12px", animation: "slideDown 0.3s ease" }}>
            <div style={{ display: "flex", gap: "8px" }}>
              <input
                type="date"
                value={tempDate}
                onChange={(e) => setTempDate(e.target.value)}
                className="bs-input"
                style={{ flex: 1 }}
              />
              <button
                type="button"
                onClick={handleAddDate}
                className="btn-secondary-custom"
                style={{ marginTop: 0, padding: "0 20px", width: "auto" }}
              >
                +
              </button>
            </div>

            {alertDates.length > 0 && (
              <div className="date-badges-container">
                {alertDates.map((dateStr) => {
                  const [y, m, d] = dateStr.split("-");
                  return (
                    <div key={dateStr} className="date-badge">
                      {`${d}/${m}/${y}`}
                      <button
                        type="button"
                        onClick={() => handleRemoveDate(dateStr)}
                        aria-label="Remover data"
                      >
                        &times;
                      </button>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        )}
      </div>

      <label className="toggle-row" htmlFor="tgl-time">
        <div className="tgl">
          <input
            type="checkbox"
            id="tgl-time"
            checked={hasTime}
            onChange={(e) => setHasTime(e.target.checked)}
          />
          <div className="tgl-track"></div>
          <div className="tgl-thumb"></div>
        </div>
        <span className="toggle-label">Definir horário</span>
      </label>

      <div className={`time-row ${hasTime ? "show" : ""}`} id="time-row">
        <div className="flex-fill">
          <label className="bs-label">Início</label>
          <input
            value={timeStart}
            onChange={(e) => setTimeStart(e.target.value)}
            className="bs-input"
            type="time"
          />
        </div>
        <div className="flex-fill">
          <label className="bs-label">Fim</label>
          <input
            value={timeEnd}
            onChange={(e) => setTimeEnd(e.target.value)}
            className="bs-input"
            type="time"
          />
        </div>
      </div>

      <button className="btn-primary-custom" onClick={handleSubmit}>
        {editingTask ? "Salvar alterações" : "Adicionar afazer"}
      </button>

      {editingTask && (
        <button className="btn-secondary-custom" onClick={onCancelEdit}>
          Cancelar
        </button>
      )}
    </div>
  );
}
