import { type ReactNode, useCallback, useLayoutEffect, useState } from "react";
import { Chart as ChartJS } from "chart.js";
import type { ThemeType } from "../types/app";
import { ThemeContext } from "./themeContextValue";

const STORAGE_KEY = "uq-theme";
const DEFAULT_THEME: ThemeType = "featherine";

const THEME_CHART_COLOURS: Record<ThemeType, string> = {
    featherine: "#a89bb8",
    bernkastel: "#9bb5d0",
    lambdadelta: "#d09bb8",
};

function getStoredTheme(): ThemeType {
    try {
        const stored = localStorage.getItem(STORAGE_KEY);
        if (stored === "featherine" || stored === "bernkastel" || stored === "lambdadelta") {
            return stored;
        }
    } catch {
        // localStorage unavailable
    }
    return DEFAULT_THEME;
}

export function ThemeProvider({ children }: { children: ReactNode }) {
    const [theme, setThemeState] = useState<ThemeType>(getStoredTheme);

    useLayoutEffect(() => {
        if (theme === DEFAULT_THEME) {
            document.documentElement.removeAttribute("data-theme");
        } else {
            document.documentElement.setAttribute("data-theme", theme);
        }
        ChartJS.defaults.color = THEME_CHART_COLOURS[theme];
    }, [theme]);

    const setTheme = useCallback((newTheme: ThemeType) => {
        setThemeState(newTheme);
        try {
            localStorage.setItem(STORAGE_KEY, newTheme);
        } catch {
            // localStorage unavailable
        }
    }, []);

    return <ThemeContext.Provider value={{ theme, setTheme }}>{children}</ThemeContext.Provider>;
}
