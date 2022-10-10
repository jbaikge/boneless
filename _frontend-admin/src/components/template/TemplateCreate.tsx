import {
  Create,
  CreateProps,
  SimpleForm,
  TextInput,
} from 'react-admin';
import { CodeInput } from '../input';

const TemplateCreate = (props: CreateProps) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="name" fullWidth />
      <CodeInput source="body" />
    </SimpleForm>
  </Create>
);

export default TemplateCreate;
