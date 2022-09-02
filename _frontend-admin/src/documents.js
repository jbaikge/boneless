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
  ReferenceInput,
  SelectInput,
  SimpleForm,
  TextField,
  TextInput,
  useGetOne,
  useResourceContext,
} from 'react-admin';
import { TinyInput } from './tinyInput';

const resourceRE = /([^/]+)\/([^/]+)\/.*/;

export const DocumentCreate = (props) => {
  return (
    <Create {...props}>
      <DocumentForm />
    </Create>
  );
};

export const DocumentEdit = (props) => {
  return (
    <Edit {...props}>
      <DocumentForm />
    </Edit>
  );
};

export const DocumentForm = () => {
  const resourceContext = useResourceContext();
  // resourceContext should be "classes/<id>/documents"
  const [ , resource, id ] = resourceRE.exec(resourceContext);
  const { data, isLoading } = useGetOne(resource, { id });

  if (isLoading) {
    return <Loading />;
  }

  let parentInput;
  if (data.parent_id != "") {
    parentInput = <ReferenceInput source="parent_id" reference={'classes/' + data.parent_id + '/documents'} perPage={25}>
      <SelectInput optionText="values.title" fullWidth />
    </ReferenceInput>;
  }

  return (
    <SimpleForm>
      <TextInput source="path" fullWidth />
      {parentInput}
      <ReferenceInput source="template_id" reference="templates" perPage={100}>
        <SelectInput fullWidth />
      </ReferenceInput>
      {data.fields.map(field => {
        const source = `values.${field.name}`;
        switch (field.type) {
          case 'datetime':
            return <DateTimeInput key={field.name} source={source} label={field.label} inputProps={{ min: field.min, max: field.max, step: field.step }} />
          case 'richtext':
            return <RichTextInput key={field.name} source={source} label={field.label} fullWidth />
          case 'select-class':
            return (
              <ReferenceInput reference={'/classes/' + field.class_id + '/documents'} source={source}>
                <SelectInput optionText={'values.' + field.field} label={field.label} fullWidth />
              </ReferenceInput>
            );
          case 'text':
            return <TextInput key={field.name} source={source} label={field.label} fullWidth />
          case 'time':
            return <TextInput type="time" key={field.name} source={source} label={field.label} fullWidth />
          case 'tiny':
            return <TinyInput key={field.name} source={source} label={field.label} fullWidth />
          default:
            return <TextInput key={field.name} source={source} label={`Unknown type (${field.type}) - ${field.label}`} fullWidth />
        }
      })}
    </SimpleForm>
  );
};

export const DocumentList = (props) => {
  const resourceContext = useResourceContext();
  const [ , resource, id ] = resourceRE.exec(resourceContext);
  const { data, isLoading } = useGetOne(resource, { id });

  if (isLoading) {
    return <Loading />;
  }

  return (
    <List {...props}>
      <Datagrid>
        {data.fields.filter(field => field.column > 0).sort((a, b) => a.column - b.column).map(field => {
          const source = `values.${field.name}`;
          switch (field.type) {
            case 'date':
              return <DateField key={field.label} source={source} label={field.label} />
            case 'datetime':
              return <DateField key={field.label} source={source} label={field.label} showTime />
            default:
              return <TextField key={field.label} source={source} label={field.label} />
          }
        })}
        <EditButton />
      </Datagrid>
    </List>
  );
};
