import React, { useEffect, useMemo, useRef, useState } from "react";

type Option = { value: string; label: string };
type Props = {
  value: string;
  onChange: (v: string) => void;
  options: Option[];
  className?: string;
  placeholder?: string;
};

export default function Select({ value, onChange, options, className, placeholder }: Props) {
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement | null>(null);
  const btnRef = useRef<HTMLButtonElement | null>(null);
  const listRef = useRef<HTMLUListElement | null>(null);

  const current = useMemo(
    () => options.find(o => o.value === value)?.label ?? placeholder ?? "Select",
    [value, options, placeholder]
  );

  useEffect(() => {
    const onDoc = (e: MouseEvent) => {
      if (!rootRef.current) return;
      if (!rootRef.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onDoc);
    return () => document.removeEventListener("mousedown", onDoc);
  }, []);

  useEffect(() => {
    if (open && listRef.current) {
      const idx = Math.max(0, options.findIndex(o => o.value === value));
      const el = listRef.current.children[idx] as HTMLElement | undefined;
      el?.scrollIntoView({ block: "nearest" });
    }
  }, [open, options, value]);

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (!open && (e.key === "ArrowDown" || e.key === "ArrowUp" || e.key === " " || e.key === "Enter")) {
      e.preventDefault();
      setOpen(true);
      return;
    }
    if (!open) return;
    const idx = Math.max(0, options.findIndex(o => o.value === value));
    if (e.key === "ArrowDown") {
      e.preventDefault();
      onChange(options[Math.min(idx + 1, options.length - 1)].value);
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      onChange(options[Math.max(idx - 1, 0)].value);
    } else if (e.key === "Enter") {
      e.preventDefault();
      setOpen(false);
      btnRef.current?.focus();
    } else if (e.key === "Escape" || e.key === "Tab") {
      setOpen(false);
    }
  };

  return (
    <div className={`ui-select ${className ?? ""}`} ref={rootRef}>
      <button
        ref={btnRef}
        type="button"
        className={`ui-select-trigger ${open ? "open" : ""}`}
        aria-haspopup="listbox"
        aria-expanded={open}
        onClick={() => setOpen(v => !v)}
        onKeyDown={onKeyDown}
      >
        <span className="ui-select-value">{current}</span>
        <span className="ui-select-caret" aria-hidden>▾</span>
      </button>
      {open && (
        <ul className="ui-select-menu" role="listbox" ref={listRef} tabIndex={-1} onKeyDown={onKeyDown}>
          {options.map(o => (
            <li
              key={o.value}
              role="option"
              aria-selected={o.value === value}
              className={`ui-select-option ${o.value === value ? "selected" : ""}`}
              data-value={o.value}
              onMouseDown={e => e.preventDefault()}
              onClick={() => {
                onChange(o.value);
                setOpen(false);
                btnRef.current?.focus();
              }}
            >
              <span className="dot" />
              <span>{o.label}</span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
