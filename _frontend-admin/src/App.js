import * as React from 'react';
import { Admin, Resource } from 'react-admin';
import simpleRestProvider from 'ra-data-simple-rest';
import { ClassCreate, ClassEdit, ClassList } from './classes';
// import dataProvider from './dataProvider';
import { darkTheme } from './theme';
import { DocumentCreate, DocumentEdit, DocumentList } from './documents';

const API_URL = process.env.REACT_APP_API_URL;

const App = () => {
  const [resources, setResources] = React.useState([]);
  const [updateResources, setUpdateResources] = React.useState(0);

  React.useEffect(() => {
    const fetchClasses = async () => {
      console.log('Setting up resources');
      const response = await fetch(API_URL + '/classes?_start=0&_end=50');
      const data = await response.json();
      const classes = data.map(c => <Resource options={{ label: c.name }} key={c.id} name={"classes/" + c.id + "/documents"} create={DocumentCreate} edit={DocumentEdit} list={DocumentList} />);
      setResources(classes);
      console.log('Finished setting up resources');
    };
    fetchClasses();
  }, [updateResources]);

  return (
    <Admin dataProvider={simpleRestProvider(API_URL)} theme={darkTheme}>
      {resources}
      <Resource name="classes" options={{ label: 'Manage Classes' }} create={<ClassCreate update={setUpdateResources} />} edit={<ClassEdit update={setUpdateResources} />} list={ClassList} />
      <Resource name="templates" list={ClassList} />
    </Admin>
  );
};

export default App;
