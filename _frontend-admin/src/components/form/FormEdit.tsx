import {
  Edit,
  EditProps,
} from 'react-admin';
import FormForm from './FormForm';

export const FormEdit = (props: EditProps) => {
  return (
    <Edit {...props}>
      <FormForm />
    </Edit>
  );
};

export default FormEdit;
