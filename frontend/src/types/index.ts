export interface Task {
  id: string;
  name: string;
  desc?: string;
  ts?: string;
  te?: string;
  alertType?: "none" | "1_day_before" | "custom_date";
  alertDates?: string[];
}

export interface ActiveDate {
  year: number;
  month: number;
  day: number;
}

export type TasksState = Record<string, Task[]>;
