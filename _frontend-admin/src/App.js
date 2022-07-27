import * as React from 'react';
import { Admin, Resource } from 'react-admin';
import simpleRestProvider from 'ra-data-simple-rest';
import { ClassCreate, ClassEdit, ClassList } from './classes';
// import dataProvider from './dataProvider';
import { darkTheme } from './theme';

const fetchResources = () =>
  fetch(window._env_.API_URL + '/classes?_start=0&_end=50')
    .then(response => response.json())
    .then(classes => classes.map(c => <Resource options={{ label: c.name }} name={"classes/" + c.id + "/documents"} key={c.id} list={ClassList} />));

const App = () => (
  <Admin dataProvider={simpleRestProvider(window._env_.API_URL)} theme={darkTheme}>
    <Resource name="classes" create={ClassCreate} edit={ClassEdit} list={ClassList} />
    {fetchResources}
  </Admin>
);

export default App;
