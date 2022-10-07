import {
  Create,
  SimpleForm,
  TextInput,
  TransformData,
  required,
  useRedirect,
} from 'react-admin';
import ClassProps from './ClassProps';
import CreateUpdateProps from './CreateUpdateProps';

const ClassCreate = (props: CreateUpdateProps) => {
  const { update, ...rest } = props;
  const redirect = useRedirect();
  const onSuccess = (data: ClassProps) => {
    update((new Date()).getTime());
    redirect('edit', 'classes', data.id);
  };

  const ensureFields = (data: TransformData) => ({
    ...data,
    fields: [],
  });

  return (
    <Create {...rest} mutationOptions={{ onSuccess }} transform={ensureFields}>
      <SimpleForm>
        <TextInput source="name" validate={[required()]} fullWidth />
      </SimpleForm>
    </Create>
  );
};

export default ClassCreate;
