import {
  Create,
  CreateProps,
  Datagrid,
  DateField,
  Edit,
  EditButton,
  EditProps,
  List,
  ListProps,
  SimpleForm,
  TextField,
  TextInput,
} from 'react-admin';
import { CodeInput } from './codeInput';

export const TemplateCreate = (props: CreateProps) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="name" fullWidth />
      <CodeInput source="body" />
    </SimpleForm>
  </Create>
);

export const TemplateEdit = (props: EditProps) => (
  <Edit {...props}>
    <SimpleForm>
      <TextInput source="name" fullWidth />
      <CodeInput source="body" />
    </SimpleForm>
  </Edit>
);

export const TemplateList = (props: ListProps) => (
  <List {...props}>
    <Datagrid>
      <TextField source="name" />
      <DateField source="created" />
      <DateField source="updated" />
      <EditButton />
    </Datagrid>
  </List>
);
