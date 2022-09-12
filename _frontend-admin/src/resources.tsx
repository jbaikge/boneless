import * as React from 'react';
import { AdminUI, Loading, Resource, useDataProvider } from 'react-admin';
import { ClassCreate, ClassEdit, ClassList } from './classes';
import { DocumentCreate, DocumentEdit, DocumentList } from './documents';
import { TemplateCreate, TemplateEdit, TemplateList } from './templates';

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
            name={"classes/" + resource.id + "/documents"}
            create={DocumentCreate}
            edit={DocumentEdit}
            list={DocumentList}
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
          name="templates"
          create={TemplateCreate}
          edit={TemplateEdit}
          list={TemplateList}
        />
      </AdminUI>
    )
  }