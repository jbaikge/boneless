import * as React from 'react';
import { AdminContext, defaultI18nProvider, localStorageStore } from 'react-admin';
import dataProvider from './dataProvider';
import { darkTheme, lightTheme } from './theme';
import { AsyncResources } from './resources';

const store = localStorageStore();

const App = () => {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

  return (
    <AdminContext
      dataProvider={dataProvider}
      i18nProvider={defaultI18nProvider}
      store={store}
      theme={prefersDark ? darkTheme : lightTheme}
    >
      <AsyncResources />
    </AdminContext>
  );
};

export default App;
