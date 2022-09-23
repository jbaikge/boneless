import {
  Create,
  CreateProps,
  Datagrid,
  Edit,
  EditButton,
  EditProps,
  List,
  ListProps,
  ShowButton,
  SimpleForm,
  TextField,
  TextInput,
} from 'react-admin';
import { FormBuilderInput } from './formBuilder';
import { GlobalPagination } from './pagination';

export const FormCreate = (props: CreateProps) => {
  return (
    <Create {...props}>
      <FormForm />
    </Create>
  );
};

export const FormEdit = (props: EditProps) => {
  return (
    <Edit {...props}>
      <FormForm />
    </Edit>
  );
};

const FormForm = () => {
  return (
    <SimpleForm>
      <TextInput source="id" label="Form ID" helperText="Available after creation" disabled fullWidth />
      <TextInput source="name" fullWidth />
      <FormBuilderInput source="schema" />
    </SimpleForm>
  );
};

export const FormList = (props: ListProps) => {
  return (
    <List {...props} pagination={<GlobalPagination />}>
      <Datagrid sx={{
        '& td:last-child': { width: '5em' },
        '& td:nth-last-child(2)': { width: '5em' },
      }}>
        <TextField source="name" label="Form Name" />
        <ShowButton />
        <EditButton />
      </Datagrid>
    </List>
  )
}
