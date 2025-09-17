import React from "react";
import { useTheme } from "../contexts/ThemeContext";

export default function ThemeToggle() {
  const { theme, toggle } = useTheme();
  return (
    <button className="theme-toggle" aria-label="Toggle theme" onClick={toggle}>
      {theme === "light" ? "🌞" : "🌙"}
    </button>
  );
}
