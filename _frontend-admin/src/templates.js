import React from 'react';
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
import { CodeInput } from './codeInput';

export const TemplateCreate = (props) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="name" fullWidth />
      <CodeInput source="body" />
    </SimpleForm>
  </Create>
);

export const TemplateEdit = (props) => (
  <Edit {...props}>
    <SimpleForm>
      <TextInput source="name" fullWidth />
      <CodeInput source="body" />
    </SimpleForm>
  </Edit>
);

export const TemplateList = (props) => (
  <List {...props}>
    <Datagrid rowClick="edit">
      <TextField source="name" />
      <DateField source="created" />
      <DateField source="updated" />
      <EditButton />
    </Datagrid>
  </List>
);
