import {
  Edit,
  EditProps,
  SimpleForm,
  TextInput,
} from 'react-admin';
import { CodeInput } from '../input';

const TemplateEdit = (props: EditProps) => (
  <Edit {...props}>
    <SimpleForm>
      <TextInput source="name" fullWidth />
      <CodeInput source="body" />
    </SimpleForm>
  </Edit>
);

export default TemplateEdit;
