import * as React from 'react';
import { Admin, Resource } from 'react-admin';
import simpleRestProvider from 'ra-data-simple-rest';
import { ClassCreate, ClassEdit, ClassList } from './classes';
// import dataProvider from './dataProvider';
import { darkTheme } from './theme';

const App = () => (
  <Admin dataProvider={simpleRestProvider(window._env_.API_URL)} theme={darkTheme}>
    <Resource name="classes" create={ClassCreate} edit={ClassEdit} list={ClassList} />
  </Admin>
);

export default App;
