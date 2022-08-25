import { RichTextInput } from 'ra-input-rich-text';
import {
  Create,
  Datagrid,
  DateField,
  Edit,
  EditButton,
  List,
  SimpleForm,
  TextField,
  TextInput,
} from 'react-admin';

export const TemplateCreate = (props) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="name" fullWidth />
      <RichTextInput source="body" fullWidth />
    </SimpleForm>
  </Create>
);

export const TemplateEdit = (props) => (
  <Edit {...props}>
    <SimpleForm>
      <TextInput source="name" fullWidth />
      <RichTextInput source="body" fullWidth />
    </SimpleForm>
  </Edit>
);

export const TemplateList = (props) => (
  <List>
    <Datagrid rowClick="edit">
      <TextField source="name" />
      <DateField source="created" />
      <DateField source="updated" />
      <EditButton />
    </Datagrid>
  </List>
);