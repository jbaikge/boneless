import {
  Create,
  CreateProps,
} from 'react-admin';
import FormForm from './FormForm';

const FormCreate = (props: CreateProps) => {
  return (
    <Create {...props}>
      <FormForm />
    </Create>
  );
};

export default FormCreate;
