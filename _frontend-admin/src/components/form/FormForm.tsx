import {
  SimpleForm,
  TextInput,
} from 'react-admin';

const FormForm = () => {
  return (
    <SimpleForm>
      <TextInput source="id" label="Form ID" helperText="Available after creation" disabled fullWidth />
      <TextInput source="name" fullWidth />
    </SimpleForm>
  );
};

export default FormForm;
