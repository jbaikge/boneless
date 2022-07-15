import * as React from 'react';
import { Admin, Resource, ListGuesser } from 'react-admin';
import { ClassCreate, ClassList } from './classes';
import dataProvider from './dataProvider';

// import jsonServerProvider from 'ra-data-json-server';
// const dataProvider = jsonServerProvider('https://jsonplaceholder.typicode.com');

const App = () => (
  <Admin dataProvider={dataProvider}>
    <Resource name="classes" create={ClassCreate} list={ClassList} />
  </Admin>
);

export default App;
