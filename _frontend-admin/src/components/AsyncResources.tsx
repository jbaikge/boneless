import * as React from 'react';
import { Route } from 'react-router-dom';
import { AdminUI, CustomRoutes, Loading, Resource, useDataProvider } from 'react-admin';
import { ClassCreate, ClassEdit, ClassImport, ClassList } from './class';
import { DocumentCreate, DocumentEdit, DocumentList, DocumentShow } from './document';
import { TemplateCreate, TemplateEdit, TemplateImport, TemplateList } from './template';
import { FormCreate, FormEdit, FormList } from './form';

// const ClassCreate = React.lazy(() => import('./class/ClassCreate'));
// const ClassEdit = React.lazy(() => import('./class/ClassEdit'));
// const ClassImport = React.lazy(() => import('./class/ClassImport'));
// const ClassList = React.lazy(() => import('./class/ClassList'));
// const DocumentCreate = React.lazy(() => import('./document/DocumentCreate'));
// const DocumentEdit = React.lazy(() => import('./document/DocumentEdit'));
// const DocumentList = React.lazy(() => import('./document/DocumentList'));
// const DocumentShow = React.lazy(() => import('./document/DocumentShow'));
// const TemplateCreate = React.lazy(() => import('./template/TemplateCreate'));
// const TemplateEdit = React.lazy(() => import('./template/TemplateEdit'));
// const TemplateImport = React.lazy(() => import('./template/TemplateImport'));
// const TemplateList = React.lazy(() => import('./template/TemplateList'));
// const FormCreate = React.lazy(() => import('./form/FormCreate'));
// const FormEdit = React.lazy(() => import('./form/FormEdit'));
// const FormList = React.lazy(() => import('./form/FormList'));

interface ClassData {
  id: string;
  name: string;
}

interface ClassResponse {
  data: Array<ClassData>;
}

const AsyncResources = () => {
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
    <React.Suspense fallback={<Loading />}>
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
    </React.Suspense>
  )
}

export default AsyncResources;
