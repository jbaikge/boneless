import {
  Create,
  Datagrid,
  DateTimeInput,
  Edit,
  EditButton,
  List,
  Loading,
  SimpleForm,
  TextField,
  TextInput,
  useGetOne,
  useResourceContext,
} from 'react-admin';
import { RichTextInput } from 'ra-input-rich-text';

export const DocumentCreate = (props) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="path" fullWidth />
      <DocumentInputs />
    </SimpleForm>
  </Create>
);

export const DocumentEdit = (props) => (
  <Edit {...props}>
    <SimpleForm>
      <TextInput source="path" fullWidth />
      <DocumentInputs />
    </SimpleForm>
  </Edit>
);

export const DocumentInputs = () => {
  const resourceContext = useResourceContext();
  // resourceContext should be "classes/<id>/documents"
  const [ , resource, id ] = /([^/]+)\/([^/]+)\/.*/.exec(resourceContext);
  const { data, isLoading } = useGetOne(resource, { id });

  if (isLoading) {
    return <Loading />;
  }

  return data.fields.map(field => {
    const source = `values.${field.name}`;
    switch (field.type) {
      case 'datetime':
        return <DateTimeInput key={field.name} source={source} label={field.label} inputProps={{ min: field.min, max: field.max, step: field.step }} />
      case 'richtext':
        return <RichTextInput key={field.name} source={source} label={field.label} fullWidth />
      case 'text':
        return <TextInput key={field.name} source={source} label={field.label} fullWidth />
      default:
        return <TextInput key={field.name} source={source} label={`Unknown type (${field.type}) - ${field.label}`} fullWidth />
    }
  });
}

export const DocumentList = (props) => (
  <List {...props}>
    <Datagrid rowClick="edit">
      <DocumentFields />
      <EditButton />
    </Datagrid>
  </List>
);

export const DocumentFields = () => {
  const resourceContext = useResourceContext();
  const [ , resource, id ] = /([^/]+)\/([^/]+)\/.*/.exec(resourceContext);
  const { data, isLoading } = useGetOne(resource, { id });

  if (isLoading) {
    return <Loading />;
  }

  return data.fields.filter(field => field.column > 0).sort((a, b) => a.column - b.column).map(field => {
    const source = `values.${field.name}`;
    switch (field.type) {
      default:
        return <TextField source={source} label={field.label} />
    }
  });
}
