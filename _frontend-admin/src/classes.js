import {
  ArrayInput,
  BooleanInput,
  Create,
  Datagrid,
  DateField,
  DateInput,
  DateTimeInput,
  Edit,
  EditButton,
  FormDataConsumer,
  List,
  NumberInput,
  ReferenceInput,
  SelectInput,
  SimpleForm,
  SimpleFormIterator,
  TextField,
  TextInput,
  required,
  useRedirect,
} from 'react-admin';
import './App.css';

const fieldChoices = [
  { id: 'date',          name: 'Date' },
  { id: 'datetime',      name: 'Date & Time' },
  { id: 'tiny',          name: 'Editor' },
  { id: 'email',         name: 'Email' },
  { id: 'number',        name: 'Number' },
  { id: 'select-class',  name: 'Select (Class)' },
  { id: 'multi-class',   name: 'Mutli-Select (Class)'},
  { id: 'select-static', name: 'Select (Static)' },
  { id: 'text',          name: 'Text' },
  { id: 'textarea',      name: 'Textarea' },
  { id: 'time',          name: 'Time' },
  { id: 'any-upload',    name: 'Upload (Any)' },
  { id: 'image-upload',  name: 'Upload (Image)' },
];

export const ClassCreate = (props) => {
  const { update, ...rest } = props;
  const redirect = useRedirect();
  const onSuccess = () => {
    update((new Date()).getTime());
    redirect('list', 'classes');
  };

  const ensureFields = data => ({
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

export const ClassEdit = (props) => {
  const { update, ...rest } = props;
  const redirect = useRedirect();
  const onSuccess = () => {
    update((new Date()).getTime());
    redirect('list', 'classes');
  };

  return (
    <Edit {...rest} mutationOptions={{ onSuccess }} mutationMode="pessimistic">
      <SimpleForm>
        <TextInput source="name" validate={[required()]} fullWidth />
        <ArrayInput source="fields">
          <SimpleFormIterator className="field-row">
            <TextInput source="label" />
            <TextInput source="name" />
            <BooleanInput source="sort" />
            <NumberInput source="column" />
            <SelectInput source="type" choices={fieldChoices} defaultValue="text" />
            <FormDataConsumer>
              {({
                formData,
                scopedFormData,
                getSource,
                ...rest
              }) => {
                switch (scopedFormData.type) {
                case 'date':
                  return (
                    <>
                      <DateInput source={getSource('min')} {...rest} />
                      <DateInput source={getSource('max')} {...rest} />
                      <TextInput source={getSource('step')} label="Step (days)" {...rest} />
                      <TextInput source={getSource('format')} label="Format (Jan 2, 2006 3:04pm)" {...rest} />
                    </>
                  );
                case 'datetime':
                  return (
                    <>
                      <DateTimeInput source={getSource('min')} {...rest} />
                      <DateTimeInput source={getSource('max')} {...rest} />
                      <TextInput source={getSource('step')} label="Step (days)" {...rest} />
                      <TextInput source={getSource('format')} label="Format (Jan 2, 2006 3:04pm)" {...rest} />
                    </>
                  );
                case 'time':
                  return (
                    <>
                      <TextInput source={getSource('format')} label="Format (3:04pm)" {...rest} />
                    </>
                  );
                case 'number':
                  return (
                    <>
                      <TextInput source={getSource('min')} {...rest} />
                      <TextInput source={getSource('max')} {...rest} />
                      <TextInput source={getSource('step')} {...rest} />
                    </>
                  );
                case 'select-static':
                  return (
                    <>
                      <TextInput source={getSource('options')} label="Options (one per line, key | value or just value" multiline {...rest} />
                    </>
                  );
                case 'select-class':
                case 'multi-class':
                  return (
                    <>
                      <ReferenceInput source={getSource('class_id')} reference="classes">
                        <SelectInput optionText="name" />
                      </ReferenceInput>
                      <TextInput source={getSource('field')} />
                    </>
                  );
                default:
                  return null;
                }
              }}
            </FormDataConsumer>
          </SimpleFormIterator>
        </ArrayInput>
      </SimpleForm>
    </Edit>
  );
};

export const ClassList = () => (
  <List>
    <Datagrid rowClick="edit">
      <TextField source="name" />
      <DateField source="created" />
      <DateField source="updated" />
      <EditButton />
    </Datagrid>
  </List>
);
