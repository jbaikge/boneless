import { RichTextInput } from 'ra-input-rich-text';
import {
  Create,
  Datagrid,
  DateField,
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
import { TinyInput } from './tinyInput';

const resourceRE = /([^/]+)\/([^/]+)\/.*/;

export const DocumentCreate = (props) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="path" fullWidth />
      <TextInput source="template_id" fullWidth />
      <DocumentInputs />
    </SimpleForm>
  </Create>
);

export const DocumentEdit = (props) => (
  <Edit {...props}>
    <SimpleForm>
      <TextInput source="path" fullWidth />
      <TextInput source="template_id" fullWidth />
      <DocumentInputs />
    </SimpleForm>
  </Edit>
);

export const DocumentInputs = () => {
  const resourceContext = useResourceContext();
  // resourceContext should be "classes/<id>/documents"
  const [ , resource, id ] = resourceRE.exec(resourceContext);
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
      case 'tiny':
        return <TinyInput key={field.name} source={source} label={field.label} fullWidth />
      default:
        return <TextInput key={field.name} source={source} label={`Unknown type (${field.type}) - ${field.label}`} fullWidth />
    }
  });
}

export const DocumentList = (props) => {
  const resourceContext = useResourceContext();
  const [ , resource, id ] = resourceRE.exec(resourceContext);
  const { data, isLoading } = useGetOne(resource, { id });

  if (isLoading) {
    return <Loading />;
  }

  return (
    <List {...props}>
      <Datagrid rowClick="edit">
        {data.fields.filter(field => field.column > 0).sort((a, b) => a.column - b.column).map(field => {
          const source = `values.${field.name}`;
          switch (field.type) {
            case 'date':
              return <DateField source={source} label={field.label} />
            case 'datetime':
              return <DateField source={source} label={field.label} showTime />
            default:
              return <TextField source={source} label={field.label} />
          }
        })}
        <EditButton />
      </Datagrid>
    </List>
  );
};
