import { defaultTheme } from "react-admin";

export const lightTheme = defaultTheme;
export const darkTheme = {
    ...defaultTheme,
    palette: {
        mode: 'dark',
    },
};
