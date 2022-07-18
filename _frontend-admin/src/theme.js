import { Button, defaultTheme, useTheme } from "react-admin";

export const lightTheme = defaultTheme;
export const darkTheme = {
    ...defaultTheme,
    palette: {
        mode: 'dark',
    },
};

export const ThemeToggler = () => {
    const [theme, setTheme] = useTheme();

    return (
        <Button onClick={() => setTheme(theme.palette.mode === 'dark' ? lightTheme : darkTheme)}>
            {theme.palette.mode === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}
        </Button>
    );
};
