import * as React from 'react';
import { AdminContext, defaultI18nProvider, defaultTheme, localStorageStore } from 'react-admin';
import dataProvider from './dataProvider';
import { AsyncResources } from './resources';

const store = localStorageStore();

const darkTheme = {
    ...defaultTheme,
    palette: {
        mode: 'dark',
    },
};

const App = () => {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

  return (
    <AdminContext
      dataProvider={dataProvider}
      i18nProvider={defaultI18nProvider}
      store={store}
      theme={prefersDark ? darkTheme : defaultTheme}
    >
      <AsyncResources />
    </AdminContext>
  );
};

export default App;
