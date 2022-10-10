import {
  Create,
  CreateProps,
} from 'react-admin';
import DocumentForm from './DocumentForm';

const DocumentCreate = (props: CreateProps) => {
  return (
    <Create {...props} redirect="list">
      <DocumentForm />
    </Create>
  );
};

export default DocumentCreate;
