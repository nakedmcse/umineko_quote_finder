function css(name: string): string {
    return getComputedStyle(document.documentElement).getPropertyValue(name).trim();
}

export function getThemeColours() {
    return {
        gold: css("--gold"),
        goldDark: css("--gold-dark"),
        goldLight: css("--gold-light"),
        purple: css("--purple"),
        purpleLight: css("--purple-light"),
        purpleMuted: css("--purple-muted"),
        text: css("--text"),
        textMuted: css("--text-muted"),
        rose: css("--rose"),
    };
}

export function getPalette(): string[] {
    const c = getThemeColours();
    return [
        c.gold,
        c.purpleLight,
        "#ff3333",
        "#3399ff",
        c.purple,
        c.goldLight,
        c.rose,
        c.goldDark,
        c.purpleMuted,
        c.text,
        "#c97bb4",
        "#7bc9a3",
    ];
}

export function getGridColour(): string {
    return `rgba(${css("--purple-rgb")}, 0.4)`;
}

export const zoomConfig = {
    zoom: {
        wheel: { enabled: true },
        pinch: { enabled: true },
        mode: "xy" as const,
    },
    pan: {
        enabled: true,
        mode: "xy" as const,
    },
};
