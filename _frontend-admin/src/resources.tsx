import * as React from 'react';
import { Route } from 'react-router-dom';
import { AdminUI, CustomRoutes, Loading, Resource, useDataProvider } from 'react-admin';
import { ClassCreate, ClassEdit, ClassImport, ClassList } from './classes';
import { DocumentCreate, DocumentEdit, DocumentList, DocumentShow } from './documents';
import { TemplateCreate, TemplateEdit, TemplateImport, TemplateList } from './templates';
import { FormCreate, FormEdit, FormList } from './forms';

interface ClassData {
  id: string;
  name: string;
}

interface ClassResponse {
  data: Array<ClassData>;
}

export const AsyncResources = () => {
    const [resources, setResources] = React.useState<ClassData[]>([]);
    const [updateResources, setUpdateResources] = React.useState(0);
    const dataProvider = useDataProvider();

    React.useEffect(() => {
      dataProvider.getList('classes', {
        filter: '',
        pagination: {page: 1, perPage: 50},
        sort: {field: 'name', order: 'ASC'},
      }).then((list: ClassResponse) => setResources(list.data));
    }, [updateResources, dataProvider]);

    return (
      <AdminUI ready={Loading}>
        {resources.map(resource => (
          <Resource
            options={{ label: resource.name }}
            key={resource.id}
            name={`classes/${resource.id}/documents`}
            create={DocumentCreate}
            edit={DocumentEdit}
            list={DocumentList}
            show={DocumentShow}
          />
        ))}
        <Resource
          name="classes"
          options={{ label: 'Manage Classes' }}
          create={<ClassCreate update={setUpdateResources} />}
          edit={<ClassEdit update={setUpdateResources} />}
          list={ClassList}
        />
        <Resource
          name="forms"
          options={{ label: 'Manage Forms' }}
          create={FormCreate}
          edit={FormEdit}
          list={FormList}
        />
        <Resource
          name="templates"
          create={TemplateCreate}
          edit={TemplateEdit}
          list={TemplateList}
        />
        <CustomRoutes>
          <Route path="/class-import" element={<ClassImport update={setUpdateResources} />} />
          <Route path="/template-import" element={<TemplateImport />} />
        </CustomRoutes>
      </AdminUI>
    )
  }
