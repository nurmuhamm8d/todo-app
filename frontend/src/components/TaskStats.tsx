import React from "react";
import type { main as NS } from "../wailsjs/go/models";

type StatsDTO = NS.StatsDTO;

type TaskStatsProps = {
  stats: StatsDTO | null;
};

export default function TaskStats({ stats }: TaskStatsProps) {
  const total = stats?.total ?? 0;
  const active = stats?.active ?? 0;
  const completed = stats?.completed ?? 0;
  const overdue = stats?.overdue ?? 0;
  const progress = total > 0 ? Math.round((completed / total) * 100) : 0;

  return (
    <div className="p-4 grid grid-cols-2 md:grid-cols-4 gap-4">
      <div className="rounded-xl bg-gray-800 text-white p-4">
        <div className="text-sm opacity-80">Total</div>
        <div className="text-3xl font-semibold">{total}</div>
      </div>
      <div className="rounded-xl bg-gray-800 text-white p-4">
        <div className="text-sm opacity-80">Active</div>
        <div className="text-3xl font-semibold">{active}</div>
      </div>
      <div className="rounded-xl bg-gray-800 text-white p-4">
        <div className="text-sm opacity-80">Completed</div>
        <div className="text-3xl font-semibold">{completed}</div>
      </div>
      <div className="rounded-xl bg-gray-800 text-white p-4">
        <div className="text-sm opacity-80">Overdue</div>
        <div className="text-3xl font-semibold">{overdue}</div>
      </div>
      <div className="col-span-2 md:col-span-4">
        <div className="h-3 w-full bg-gray-700 rounded-xl overflow-hidden">
          <div className="h-3 bg-green-500" style={{ width: `${progress}%` }} />
        </div>
        <div className="mt-2 text-sm text-white opacity-80">{progress}%</div>
      </div>
    </div>
  );
}
