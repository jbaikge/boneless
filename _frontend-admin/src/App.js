import * as React from 'react';
import { AdminContext, defaultI18nProvider, localStorageStore } from 'react-admin';
import simpleRestProvider from 'ra-data-simple-rest';
import { darkTheme, lightTheme } from './theme';
import { AsyncResources } from './resources';

const API_URL = process.env.REACT_APP_API_URL;
const store = localStorageStore();

const App = () => {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

  return (
    <AdminContext
      dataProvider={simpleRestProvider(API_URL)}
      i18nProvider={defaultI18nProvider}
      store={store}
      theme={prefersDark ? darkTheme : lightTheme}
    >
      <AsyncResources />
    </AdminContext>
  );
};

export default App;
